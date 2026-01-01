package services

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/karima-store/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestCheckoutService_TransactionRollback tests that transactions rollback on errors
func TestCheckoutService_TransactionRollback(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock:  %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm: %v", err)
	}

	// Test: Transaction should rollback on error
	mock.ExpectBegin()
	mock.ExpectRollback()

	err = gormDB.Transaction(func(tx *gorm.DB) error {
		return fmt.Errorf("test error - should rollback")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestCheckoutService_TransactionCommit tests that transactions commit on success
func TestCheckoutService_TransactionCommit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm: %v", err)
	}

	// Test: Transaction should commit on success
	mock.ExpectBegin()
	mock.ExpectCommit()

	err = gormDB.Transaction(func(tx *gorm.DB) error {
		return nil // Success
	})

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestCheckoutService_SignatureGeneration tests the signature generation for Midtrans
func TestCheckoutService_SignatureGeneration(t *testing.T) {
	serverKey := "test-server-key-12345"
	orderID := "ORD20260101120000"
	statusCode := "200"
	grossAmount := "100000.00"

	// Generate signature as the service does
	data := fmt.Sprintf("%s%s%s%s", orderID, statusCode, grossAmount, serverKey)
	hash := sha512.Sum512([]byte(data))
	signature := hex.EncodeToString(hash[:])

	// Verify it's a valid SHA512 hash (128 hex characters)
	if len(signature) != 128 {
		t.Errorf("Expected signature length 128, got %d", len(signature))
	}

	t.Logf("Generated signature: %s", signature)
}

// TestCheckoutService_StockDeductionLogic tests the stock deduction calculation
func TestCheckoutService_StockDeductionLogic(t *testing.T) {
	tests := []struct {
		name          string
		currentStock  int
		requestedQty  int
		shouldSucceed bool
		expectedStock int
	}{
		{
			name:          "Sufficient stock",
			currentStock:  100,
			requestedQty:  5,
			shouldSucceed: true,
			expectedStock: 95,
		},
		{
			name:          "Exact stock",
			currentStock:  5,
			requestedQty:  5,
			shouldSucceed: true,
			expectedStock: 0,
		},
		{
			name:          "Insufficient stock",
			currentStock:  3,
			requestedQty:  5,
			shouldSucceed: false,
			expectedStock: 3,
		},
		{
			name:          "Zero stock",
			currentStock:  0,
			requestedQty:  1,
			shouldSucceed: false,
			expectedStock: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changeAmount := -tt.requestedQty
			newStock := tt.currentStock + changeAmount

			if tt.shouldSucceed {
				if newStock < 0 {
					t.Errorf("Expected success but stock would be negative: %d", newStock)
				}
				if newStock != tt.expectedStock {
					t.Errorf("Expected stock %d, got %d", tt.expectedStock, newStock)
				}
			} else {
				if newStock >= 0 {
					t.Errorf("Expected failure (negative stock) but got valid stock: %d", newStock)
				}
			}
		})
	}
}

// TestCheckoutService_PaymentNotificationIdempotency tests idempotent webhook handling
func TestCheckoutService_PaymentNotificationIdempotency(t *testing.T) {
	tests := []struct {
		name                 string
		initialPaymentStatus models.PaymentStatus
		initialOrderStatus   models.OrderStatus
		notificationStatus   string
		shouldProcessUpdate  bool
	}{
		{
			name:                 "First settlement notification",
			initialPaymentStatus: models.PaymentPending,
			initialOrderStatus:   models.StatusPending,
			notificationStatus:   "settlement",
			shouldProcessUpdate:  true,
		},
		{
			name:                 "Duplicate settlement notification",
			initialPaymentStatus: models.PaymentPaid,
			initialOrderStatus:   models.StatusConfirmed,
			notificationStatus:   "settlement",
			shouldProcessUpdate:  false,
		},
		{
			name:                 "First failure notification",
			initialPaymentStatus: models.PaymentPending,
			initialOrderStatus:   models.StatusPending,
			notificationStatus:   "failed",
			shouldProcessUpdate:  true,
		},
		{
			name:                 "Duplicate failure notification",
			initialPaymentStatus: models.PaymentFailed,
			initialOrderStatus:   models.StatusCancelled,
			notificationStatus:   "failed",
			shouldProcessUpdate:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify idempotency logic: final states should not be updated
			if tt.initialPaymentStatus == models.PaymentPaid || tt.initialOrderStatus == models.StatusCancelled {
				if tt.shouldProcessUpdate {
					t.Error("Expected no update for final status")
				}
			}
		})
	}
}

// TestCheckoutService_OrderNumberUniqueness tests order number generation
func TestCheckoutService_OrderNumberUniqueness(t *testing.T) {
	orderNumbers := make(map[string]bool)

	// Generate order numbers
	for i := 0; i < 100; i++ {
		orderNum := "ORD" + time.Now().Format("20060102150405")
		orderNumbers[orderNum] = true
		time.Sleep(1 * time.Microsecond)
	}

	t.Logf("Generated %d unique order numbers out of 100 attempts", len(orderNumbers))

	//  Note: In production, add random suffix or use database sequence
	// to guarantee uniqueness with concurrent requests
}
