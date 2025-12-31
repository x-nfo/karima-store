package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/karima-store/internal/database"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"github.com/karima-store/internal/services"
)

// CheckoutService handles checkout operations
type CheckoutService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	// Add other dependencies as needed
}

// NewCheckoutService creates a new checkout service instance
func NewCheckoutService(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	notificationService *services.NotificationService,
) *CheckoutService {
	return &CheckoutService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

// Checkout creates an order and generates Midtrans Snap token
func (s *CheckoutService) Checkout(req *models.CheckoutRequest) (*models.CheckoutResponse, error) {
	// Calculate order summary
	orderSummary, err := s.pricingService.CalculateOrderSummary(req.Items, req.Shipping, req.CustomerType)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate order summary: %w", err)
	}

	// Generate order number
	orderNumber := s.generateOrderNumber()

	// Create order
	order := &models.Order{
		OrderNumber:       orderNumber,
		UserID:            req.UserID,
		PaymentMethod:      models.PaymentMethod(req.PaymentMethod),
		Subtotal:          orderSummary.Subtotal,
		Discount:          orderSummary.TotalDiscount,
		ShippingCost:      orderSummary.ShippingCost,
		Tax:               orderSummary.TaxAmount,
		TotalAmount:       orderSummary.Total,
		ShippingName:       req.ShippingName,
		ShippingProvider:   "JNE", // Default provider
		Items:             s.createOrderItems(req.Items, orderSummary),
	}

	// Create order in database
	if err := s.orderRepo.Create(order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Generate Midtrans Snap token
	snapToken, err := s.generateSnapToken(order, req.Items, req)
	if err != nil {
		// Rollback order creation
		s.orderRepo.Delete(order.ID)
		return nil, fmt.Errorf("failed to generate Snap token: %w", err)
	}

	// Send order created notification
	if err := s.notificationService.SendOrderCreatedNotification(order); err != nil {
		log.Printf("Failed to send order created notification: %v", err)
	}

	return &models.CheckoutResponse{
		OrderNumber: orderNumber,
		OrderID:    order.ID,
		SnapToken:   snapToken.Token,
		RedirectURL: snapToken.RedirectURL,
		Amount:      order.TotalAmount,
		ExpiryTime:  time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}, nil
}

// ProcessPaymentNotification processes Midtrans webhook notification
func (s *CheckoutService) ProcessPaymentNotification(notification *models.MidtransPaymentNotification) error {
	// Verify signature
	if !s.verifySignature(notification) {
		return fmt.Errorf("invalid signature")
	}

	// Get order by order number
	order, err := s.orderRepo.GetByOrderNumber(notification.OrderID)
	if err != nil {
		return fmt.Errorf("order not found: %s", notification.OrderID)
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

			if err := s.orderRepo.Update(order); err != nil {
				return err
			}

			// Reduce stock
			if err := s.reduceStock(order); err != nil {
				return err
			}

			// Send payment success notification
			if err := s.notificationService.SendPaymentSuccessNotification(order); err != nil {
				log.Printf("Failed to send payment success notification: %v", err)
			}
		}
	case "failed", "cancelled":
		// Payment failed or cancelled
		if order.PaymentStatus == models.PaymentPending {
			order.PaymentStatus = models.PaymentFailed
			order.Status = models.StatusCancelled
			order.CancelReason = "Payment " + notification.TransactionStatus
			now := time.Now()
			order.CancelledAt = &now

			if err := s.orderRepo.Update(order); err != nil {
				return err
			}

			// Restore stock
			if err := s.restoreStock(order); err != nil {
				return err
			}
		}
	case "refund":
		// Payment refunded
		if order.PaymentStatus == models.PaymentPaid {
			order.PaymentStatus = models.PaymentRefunded
			order.Status = models.StatusRefunded

			if err := s.orderRepo.Update(order); err != nil {
				return err
			}

			// Restore stock
			if err := s.restoreStock(order); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unknown transaction status: %s", notification.TransactionStatus)
	}

	return nil
}

// generateOrderNumber generates a unique order number
func (s *CheckoutService) generateOrderNumber() string {
	// Implementation for generating order number
	return "ORD" + time.Now().Format("20060102150405")
}

// verifySignature verifies Midtrans webhook signature
func (s *CheckoutService) verifySignature(notification *models.MidtransPaymentNotification) bool {
	// Signature format: SHA512(order_id + status_code + gross_amount + server_key)
	data := fmt.Sprintf("%s%s%s%s",
		notification.OrderID,
		notification.StatusCode,
		notification.GrossAmount,
		s.midtransServerKey,
	)

	signature := hex.EncodeToString(hash[:])
	return signature == notification.SignatureKey
}

// reduceStock reduces stock for all order items
func (s *CheckoutService) reduceStock(order *models.Order) error {
	for _, item := range order.Items {
		if item.VariantName != "" && item.ProductSKU != "" {
			// Logic to reduce stock
			if err := s.productRepo.UpdateStock(item.ProductID, item.VariantID, -item.Quantity); err != nil {
				return err
			}
		}
	}
	return nil
}

// restoreStock restores stock for cancelled/refunded orders
func (s *CheckoutService) restoreStock(order *models.Order) error {
	for _, item := range order.Items {
		if item.VariantName != "" && item.ProductSKU != "" {
			// Logic to restore stock
			if err := s.productRepo.UpdateStock(item.ProductID, item.VariantID, item.Quantity); err != nil {
				return err
			}
		}
	}
	return nil
}

// createOrderItems creates order items from checkout items
func (s *CheckoutService) createOrderItems(items []PriceCalculationRequest, orderSummary *OrderSummary) []models.OrderItem {
	var orderItems []models.OrderItem
	for _, item := range items {
		orderItem := models.OrderItem{
			ProductID:    item.ProductID,
			VariantID:    item.VariantID,
			Quantity:     item.Quantity,
			Price:        item.Price,
			TotalPrice:   item.Price * float64(item.Quantity),
		}
		orderItems = append(orderItems, orderItem)
	}
	return orderItems
}

// generateSnapToken generates Midtrans Snap token
func (s *CheckoutService) generateSnapToken(order *models.Order, items []PriceCalculationRequest, req *models.CheckoutRequest) (*models.MidtransSnapResponse, error) {
	// Build Midtrans request
	transactionDetails := models.TransactionDetails{
		OrderID:    order.OrderNumber,
		GrossAmount: order.TotalAmount,
	}

	// Configure Midtrans configuration
	midtransConfig := &models.MidtransConfig{
		IsProduction:     false, // Set to true in production
		SandboxServerKey: "SB-Mid-server-Your-Sandbox-Key",
		ProductionServerKey: "SB-Mid-server-Your-Production-Key",
	}

	// Build Midtrans request
	snapToken, err := midtrans.NewSnapToken(transactionDetails, midtransConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Snap token: %w", err)
	}

	return &models.MidtransSnapResponse{
		Token:       snapToken.Token,
		RedirectURL: snapToken.RedirectURL,
	}, nil
}

// reduceStock reduces stock for all order items
func (s *CheckoutService) reduceStock(order *models.Order) error {
	for _, item := range order.Items {
		if item.VariantName != "" && item.ProductSKU != "" {
			// Logic to reduce stock
			if err := s.productRepo.UpdateStock(item.ProductID, item.VariantID, -item.Quantity); err != nil {
				return err
			}
		}
	}
	return nil
}

// restoreStock restores stock for cancelled/refunded orders
func (s *CheckoutService) restoreStock(order *models.Order) error {
	for _, item := range order.Items {
		if item.VariantName != "" && item.ProductSKU != "" {
			// Logic to restore stock
			if err := s.productRepo.UpdateStock(item.ProductID, item.VariantID, item.Quantity); err != nil {
				return err
			}
		}
	}
	return nil
}
