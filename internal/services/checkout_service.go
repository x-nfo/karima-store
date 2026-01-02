package services

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/karima-store/internal/database"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"gorm.io/gorm"
)

// CheckoutService handles checkout operations
type CheckoutService struct {
	db                  *database.PostgreSQL
	orderRepo           repository.OrderRepository
	productRepo         repository.ProductRepository
	variantRepo         repository.VariantRepository
	stockLogRepo        repository.StockLogRepository
	pricingService      *PricingService
	notificationService *NotificationService
	midtransConfig      *MidtransConfig
}

type MidtransConfig struct {
	ServerKey           string
	ClientKey           string
	APIBaseURL          string
	IsProduction        bool
	SandboxServerKey    string
	ProductionServerKey string
}

// NewCheckoutService creates a new checkout service instance
func NewCheckoutService(
	db *database.PostgreSQL,
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	variantRepo repository.VariantRepository,
	stockLogRepo repository.StockLogRepository,
	pricingService *PricingService,
	notificationService *NotificationService,
	midtransConfig *MidtransConfig,
) *CheckoutService {
	return &CheckoutService{
		db:                  db,
		orderRepo:           orderRepo,
		productRepo:         productRepo,
		variantRepo:         variantRepo,
		stockLogRepo:        stockLogRepo,
		pricingService:      pricingService,
		notificationService: notificationService,
		midtransConfig:      midtransConfig,
	}
}

// Checkout creates an order and generates Midtrans Snap token
func (s *CheckoutService) Checkout(req *models.CheckoutRequest) (*models.CheckoutResponse, error) {
	// 1. Prepare Data & Calculate Prices (Read-Only)
	var priceReqItems []PriceCalculationRequest
	for _, item := range req.Items {
		var variantID *uint
		if item.VariantID != 0 {
			vID := item.VariantID
			variantID = &vID
		}

		priceReqItems = append(priceReqItems, PriceCalculationRequest{
			ProductID: item.ProductID,
			VariantID: variantID,
			Quantity:  item.Quantity,
		})
	}

	var shippingItems []ShippingItem
	for _, priceReq := range priceReqItems {
		product, err := s.productRepo.GetByID(priceReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product %d: %w", priceReq.ProductID, err)
		}
		shippingItems = append(shippingItems, ShippingItem{
			Weight:   product.Weight,
			Quantity: priceReq.Quantity,
		})
	}

	shippingReq := ShippingCalculationRequest{
		Items:        shippingItems,
		Destination:  req.ShippingCity,
		ShippingType: "jne",
	}

	customerType := CustomerRetail // Default to retail
	orderSummary, err := s.pricingService.CalculateOrderSummary(priceReqItems, shippingReq, customerType)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate order summary: %w", err)
	}

	orderNumber := s.generateOrderNumber()
	order := &models.Order{
		OrderNumber:      orderNumber,
		UserID:           req.UserID,
		PaymentMethod:    models.PaymentMethod(req.PaymentMethod),
		Subtotal:         orderSummary.Subtotal,
		Discount:         orderSummary.TotalDiscount,
		ShippingCost:     orderSummary.ShippingCost,
		Tax:              orderSummary.TaxAmount,
		TotalAmount:      orderSummary.Total,
		ShippingName:     req.ShippingName,
		ShippingProvider: "JNE",
		Status:           models.StatusPending,
		PaymentStatus:    models.PaymentPending,
		Items:            s.createOrderItems(priceReqItems, orderSummary),
	}

	// 2. Execution Phase: DB Transaction (Write)
	// Wraps Stock Deduction, Order Creation, and Snap Token Generation in an atomic block.
	// If Snap Token generation fails, the entire transaction will be rolled back automatically.
	var snapToken *models.MidtransSnapResponse
	err = s.db.DB().Transaction(func(tx *gorm.DB) error {
		txProductRepo := s.productRepo.WithTx(tx)
		txStockLogRepo := s.stockLogRepo.WithTx(tx)
		txOrderRepo := s.orderRepo.WithTx(tx)

		// A. Deduct Stock (Reservation)
		// We use the same method 'reduceStockWithTx' but must ensure it checks for negative stock.
		if err := s.reduceStockWithTx(txProductRepo, txStockLogRepo, order); err != nil {
			return fmt.Errorf("stock reservation failed: %w", err)
		}

		// B. Create Order
		if err := txOrderRepo.Create(order); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// C. Generate Snap Token (External API Call)
		// If this fails, the entire transaction (stock deduction + order creation) will be rolled back
		snapToken, err = s.generateSnapToken(order, priceReqItems, req)
		if err != nil {
			return fmt.Errorf("failed to generate snap token: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 4. Notifications (Non-blocking)
	if s.notificationService != nil {
		go func() {
			if err := s.notificationService.SendOrderCreatedNotification(order); err != nil {
				log.Printf("Failed to send order created notification: %v", err)
			}
		}()
	}

	return &models.CheckoutResponse{
		OrderNumber: orderNumber,
		OrderID:     order.ID,
		SnapToken:   snapToken.Token,
		RedirectURL: snapToken.RedirectURL,
		Amount:      order.TotalAmount,
		ExpiryTime:  time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}, nil
}

// verifySignature verifies Midtrans webhook signature
func (s *CheckoutService) verifySignature(notification *models.MidtransPaymentNotification) bool {
	// Signature format: SHA512(order_id + status_code + gross_amount + server_key)
	data := fmt.Sprintf("%s%s%.2f%s",
		notification.OrderID,
		notification.StatusCode,
		notification.GrossAmount,
		s.midtransConfig.ServerKey,
	)

	hash := sha512.Sum512([]byte(data))
	signature := hex.EncodeToString(hash[:])
	return signature == notification.SignatureKey
}

// ProcessPaymentNotification processes Midtrans webhook notification
func (s *CheckoutService) ProcessPaymentNotification(notification *models.MidtransPaymentNotification) error {
	// Verify signature
	if !s.verifySignature(notification) {
		return fmt.Errorf("invalid signature")
	}

	// Get DB instance for transaction
	db := s.db.DB()

	return db.Transaction(func(tx *gorm.DB) error {
		// Create transaction-aware repositories
		txOrderRepo := s.orderRepo.WithTx(tx)
		// txProductRepo is only needed for restore
		// txStockLogRepo is only needed for restore
		txProductRepo := s.productRepo.WithTx(tx)
		txStockLogRepo := s.stockLogRepo.WithTx(tx)

		// Get order by order number with transaction (using FOR UPDATE if needed, but simple Get here is mostly fine unless high concurrency on same order)
		order, err := txOrderRepo.GetByOrderNumber(notification.OrderID)
		if err != nil {
			return fmt.Errorf("order not found: %s", notification.OrderID)
		}

		// Idempotency check: if status is already final, ignore
		if order.PaymentStatus == models.PaymentPaid || order.Status == models.StatusCancelled {
			return nil
		}

		// Process based on transaction status
		switch notification.TransactionStatus {
		case "capture", "settlement":
			// Payment successful
			if order.PaymentStatus != models.PaymentPaid {
				// Update order status
				order.PaymentStatus = models.PaymentPaid
				order.Status = models.StatusConfirmed
				now := time.Now()
				order.ConfirmedAt = &now

				if err := txOrderRepo.Update(order); err != nil {
					return err
				}

				// NOTE: Stock already deducted at Checkout. No need to deduct here.

				// Send payment success notification
				defer func() {
					if s.notificationService != nil {
						if err := s.notificationService.SendPaymentSuccessNotification(order); err != nil {
							log.Printf("Failed to send payment success notification: %v", err)
						}
					}
				}()
			}
		case "failed", "cancelled", "expire":
			// Payment failed or cancelled
			if order.PaymentStatus == models.PaymentPending {
				order.PaymentStatus = models.PaymentFailed
				order.Status = models.StatusCancelled
				order.CancelReason = "Payment " + notification.TransactionStatus
				now := time.Now()
				order.CancelledAt = &now

				if err := txOrderRepo.Update(order); err != nil {
					return err
				}

				// RESTORE stock since it was reserved at Checkout
				if err := s.restoreStockWithTx(txProductRepo, txStockLogRepo, order); err != nil {
					return err
				}
			}
		case "refund":
			// Payment refunded
			if order.PaymentStatus == models.PaymentPaid {
				order.PaymentStatus = models.PaymentRefunded
				order.Status = models.StatusRefunded

				if err := txOrderRepo.Update(order); err != nil {
					return err
				}

				// Restore stock and log changes
				if err := s.restoreStockWithTx(txProductRepo, txStockLogRepo, order); err != nil {
					return err
				}
			}
		default:
			// Just return, no error to avoid retry storm from webhook
			log.Printf("Unknown transaction status: %s", notification.TransactionStatus)
		}

		return nil
	})
}

// generateOrderNumber generates a unique order number
func (s *CheckoutService) generateOrderNumber() string {
	return "ORD" + time.Now().Format("20060102150405")
}

// reduceStockWithTx reduces stock and logs changes
func (s *CheckoutService) reduceStockWithTx(
	productRepo repository.ProductRepository,
	stockLogRepo repository.StockLogRepository,
	order *models.Order,
) error {
	for _, item := range order.Items {
		// Get latest stock
		product, err := productRepo.GetByID(item.ProductID)
		if err != nil {
			return err
		}

		changeAmount := -item.Quantity
		previousStock := product.Stock
		newStock := previousStock + changeAmount

		// Critical Check: Prevent negative stock
		if newStock < 0 {
			return fmt.Errorf("insufficient stock for product %s (ID: %d). Available: %d, Requested: %d",
				product.Name, item.ProductID, previousStock, item.Quantity)
		}

		// Update stock
		if err := productRepo.UpdateStock(item.ProductID, changeAmount); err != nil {
			return err
		}

		// Create log
		log := &models.StockLog{
			ProductID:     item.ProductID,
			VariantID:     nil, // OrderItem might need VariantID if we track variants
			ChangeAmount:  changeAmount,
			PreviousStock: previousStock,
			NewStock:      newStock,
			Reason:        fmt.Sprintf("Order %s Placed (Reserved)", order.OrderNumber),
			ReferenceID:   order.OrderNumber,
			CreatedAt:     time.Now(),
		}
		if err := stockLogRepo.Create(log); err != nil {
			return err
		}
	}
	return nil
}

// restoreStockWithTx restores stock and logs changes
func (s *CheckoutService) restoreStockWithTx(
	productRepo repository.ProductRepository,
	stockLogRepo repository.StockLogRepository,
	order *models.Order,
) error {
	for _, item := range order.Items {
		// Get current stock
		product, err := productRepo.GetByID(item.ProductID)
		if err != nil {
			return err
		}

		changeAmount := item.Quantity
		previousStock := product.Stock
		newStock := previousStock + changeAmount

		// Update stock
		if err := productRepo.UpdateStock(item.ProductID, changeAmount); err != nil {
			return err
		}

		// Create log
		log := &models.StockLog{
			ProductID:     item.ProductID,
			VariantID:     nil,
			ChangeAmount:  changeAmount,
			PreviousStock: previousStock,
			NewStock:      newStock,
			Reason:        fmt.Sprintf("Order %s Cancelled/Refunded (Restored)", order.OrderNumber),
			ReferenceID:   order.OrderNumber,
			CreatedAt:     time.Now(),
		}
		if err := stockLogRepo.Create(log); err != nil {
			return err
		}
	}
	return nil
}

// createOrderItems creates order items from checkout items
func (s *CheckoutService) createOrderItems(items []PriceCalculationRequest, orderSummary *OrderSummary) []models.OrderItem {
	var orderItems []models.OrderItem
	for _, item := range items {
		orderItem := models.OrderItem{
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  0, // Placeholder
			TotalPrice: 0, // Placeholder
		}
		orderItems = append(orderItems, orderItem)
	}
	return orderItems
}

// generateSnapToken generates Midtrans Snap token (Stubbed)
func (s *CheckoutService) generateSnapToken(order *models.Order, items []PriceCalculationRequest, req *models.CheckoutRequest) (*models.MidtransSnapResponse, error) {
	// Ensure server key is set
	if s.midtransConfig.ServerKey == "" {
		return nil, fmt.Errorf("midtrans server key is not set")
	}

	// Stub implementation since midtrans package is missing
	// TODO: Integrate actual Midtrans library

	return &models.MidtransSnapResponse{
		Token:       "dummy_snap_token_" + order.OrderNumber,
		RedirectURL: "https://app.sandbox.midtrans.com/snap/v2/vtweb/" + order.OrderNumber,
	}, nil
}
