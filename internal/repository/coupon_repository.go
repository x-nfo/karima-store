package repository

import (
	"time"

	"github.com/karima-store/internal/models"
	"gorm.io/gorm"
)

type CouponRepository interface {
	Create(coupon *models.Coupon) error
	GetByID(id uint) (*models.Coupon, error)
	GetByCode(code string) (*models.Coupon, error)
	GetAll(limit, offset int) ([]models.Coupon, int64, error)
	GetActive() ([]models.Coupon, error)
	Update(coupon *models.Coupon) error
	Delete(id uint) error
	ValidateCoupon(code string, userID uint, purchaseAmount float64, customerType string) (*models.Coupon, error)
	RecordUsage(couponID, userID, orderID uint, discountAmount float64) error
	GetUserUsageCount(couponID, userID uint) (int, error)
}

type couponRepository struct {
	db *gorm.DB
}

func NewCouponRepository(db *gorm.DB) CouponRepository {
	return &couponRepository{db: db}
}

func (r *couponRepository) Create(coupon *models.Coupon) error {
	// Select("*") ensures that zero values (like bool false) are also saved
	return r.db.Select("*").Create(coupon).Error
}

func (r *couponRepository) GetByID(id uint) (*models.Coupon, error) {
	var coupon models.Coupon
	err := r.db.First(&coupon, id).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (r *couponRepository) GetByCode(code string) (*models.Coupon, error) {
	var coupon models.Coupon
	err := r.db.Where("code = ?", code).First(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (r *couponRepository) GetAll(limit, offset int) ([]models.Coupon, int64, error) {
	var coupons []models.Coupon
	var total int64

	query := r.db.Model(&models.Coupon{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&coupons).Error

	return coupons, total, err
}

func (r *couponRepository) GetActive() ([]models.Coupon, error) {
	var coupons []models.Coupon
	now := time.Now()

	err := r.db.Where("status = ?", models.CouponStatusActive).
		Where("(valid_from IS NULL OR valid_from <= ?)", now).
		Where("(valid_until IS NULL OR valid_until >= ?)", now).
		Order("created_at DESC").
		Find(&coupons).Error

	return coupons, err
}

func (r *couponRepository) Update(coupon *models.Coupon) error {
	return r.db.Save(coupon).Error
}

func (r *couponRepository) Delete(id uint) error {
	return r.db.Delete(&models.Coupon{}, id).Error
}

// ValidateCoupon checks if a coupon is valid for use
func (r *couponRepository) ValidateCoupon(code string, userID uint, purchaseAmount float64, customerType string) (*models.Coupon, error) {
	coupon, err := r.GetByCode(code)
	if err != nil {
		return nil, err
	}

	// Check if coupon is active
	if coupon.Status != models.CouponStatusActive {
		return nil, gorm.ErrRecordNotFound
	}

	// Check validity dates
	now := time.Now()
	if coupon.ValidFrom != nil && now.Before(*coupon.ValidFrom) {
		return nil, gorm.ErrRecordNotFound
	}
	if coupon.ValidUntil != nil && now.After(*coupon.ValidUntil) {
		return nil, gorm.ErrRecordNotFound
	}

	// Check minimum purchase amount
	if coupon.MinPurchaseAmount > 0 && purchaseAmount < coupon.MinPurchaseAmount {
		return nil, gorm.ErrRecordNotFound
	}

	// Check customer type restrictions
	if customerType == "retail" && !coupon.ForRetail {
		return nil, gorm.ErrRecordNotFound
	}
	if customerType == "reseller" && !coupon.ForReseller {
		return nil, gorm.ErrRecordNotFound
	}

	// Check usage limits
	if coupon.MaxUsageCount > 0 && coupon.UsageCount >= coupon.MaxUsageCount {
		return nil, gorm.ErrRecordNotFound
	}

	// Check per-user usage limit
	if coupon.MaxUsagePerUser > 0 {
		userUsageCount, err := r.GetUserUsageCount(coupon.ID, userID)
		if err != nil {
			return nil, err
		}
		if userUsageCount >= coupon.MaxUsagePerUser {
			return nil, gorm.ErrRecordNotFound
		}
	}

	return coupon, nil
}

// RecordUsage records a coupon usage
func (r *couponRepository) RecordUsage(couponID, userID, orderID uint, discountAmount float64) error {
	// Create coupon usage record
	usage := &models.CouponUsage{
		CouponID:       couponID,
		UserID:         userID,
		OrderID:        orderID,
		DiscountAmount: discountAmount,
	}

	if err := r.db.Create(usage).Error; err != nil {
		return err
	}

	// Update coupon usage count
	return r.db.Model(&models.Coupon{}).
		Where("id = ?", couponID).
		Update("usage_count", gorm.Expr("usage_count + 1")).Error
}

// GetUserUsageCount gets the number of times a user has used a coupon
func (r *couponRepository) GetUserUsageCount(couponID, userID uint) (int, error) {
	var count int64
	err := r.db.Model(&models.CouponUsage{}).
		Where("coupon_id = ? AND user_id = ?", couponID, userID).
		Count(&count).Error

	return int(count), err
}
