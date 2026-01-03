package repository

import (
	"github.com/karima-store/internal/models"

	"gorm.io/gorm"
)

type FlashSaleRepository interface {
	GetByID(id uint) (*models.FlashSale, error)
	GetAll() ([]models.FlashSale, error)
	GetActiveFlashSales() ([]models.FlashSale, error)
	GetUpcomingFlashSales() ([]models.FlashSale, error)
	Create(flashSale *models.FlashSale) error
	Update(flashSale *models.FlashSale) error
	Delete(id uint) error
	AddProductToFlashSale(flashSaleProduct *models.FlashSaleProduct) error
	RemoveProductFromFlashSale(flashSaleID, productID uint) error
	GetFlashSaleProducts(flashSaleID uint) ([]models.FlashSaleProduct, error)
	UpdateFlashSaleProduct(flashSaleProduct *models.FlashSaleProduct) error
}

type flashSaleRepository struct {
	db *gorm.DB
}

func NewFlashSaleRepository(db *gorm.DB) FlashSaleRepository {
	return &flashSaleRepository{db: db}
}

// GetByID retrieves a flash sale by ID
func (r *flashSaleRepository) GetByID(id uint) (*models.FlashSale, error) {
	var flashSale models.FlashSale
	err := r.db.Preload("Products").First(&flashSale, id).Error
	if err != nil {
		return nil, err
	}
	return &flashSale, nil
}

// GetAll retrieves all flash sales
func (r *flashSaleRepository) GetAll() ([]models.FlashSale, error) {
	var flashSales []models.FlashSale
	err := r.db.Preload("Products").Find(&flashSales).Error
	if err != nil {
		return nil, err
	}
	return flashSales, nil
}

// GetActiveFlashSales retrieves all currently active flash sales
func (r *flashSaleRepository) GetActiveFlashSales() ([]models.FlashSale, error) {
	var flashSales []models.FlashSale
	err := r.db.Preload("Products").
		Where("status = ?", models.FlashSaleActive).
		Where("start_time <= ?", gorm.Expr("NOW()")).
		Where("end_time >= ?", gorm.Expr("NOW()")).
		Find(&flashSales).Error
	if err != nil {
		return nil, err
	}
	return flashSales, nil
}

// GetUpcomingFlashSales retrieves all upcoming flash sales
func (r *flashSaleRepository) GetUpcomingFlashSales() ([]models.FlashSale, error) {
	var flashSales []models.FlashSale
	err := r.db.Preload("Products").
		Where("status = ?", models.FlashSaleUpcoming).
		Where("start_time > ?", gorm.Expr("NOW()")).
		Find(&flashSales).Error
	if err != nil {
		return nil, err
	}
	return flashSales, nil
}

// Create creates a new flash sale
func (r *flashSaleRepository) Create(flashSale *models.FlashSale) error {
	return r.db.Create(flashSale).Error
}

// Update updates an existing flash sale
func (r *flashSaleRepository) Update(flashSale *models.FlashSale) error {
	return r.db.Save(flashSale).Error
}

// Delete soft deletes a flash sale
func (r *flashSaleRepository) Delete(id uint) error {
	return r.db.Delete(&models.FlashSale{}, id).Error
}

// AddProductToFlashSale adds a product to a flash sale
func (r *flashSaleRepository) AddProductToFlashSale(flashSaleProduct *models.FlashSaleProduct) error {
	return r.db.Create(flashSaleProduct).Error
}

// RemoveProductFromFlashSale removes a product from a flash sale
func (r *flashSaleRepository) RemoveProductFromFlashSale(flashSaleID, productID uint) error {
	return r.db.Where("flash_sale_id = ? AND product_id = ?", flashSaleID, productID).
		Delete(&models.FlashSaleProduct{}).Error
}

// GetFlashSaleProducts retrieves all products in a flash sale
func (r *flashSaleRepository) GetFlashSaleProducts(flashSaleID uint) ([]models.FlashSaleProduct, error) {
	var flashSaleProducts []models.FlashSaleProduct
	err := r.db.Preload("FlashSale").
		Preload("Product").
		Where("flash_sale_id = ?", flashSaleID).
		Find(&flashSaleProducts).Error
	if err != nil {
		return nil, err
	}
	return flashSaleProducts, nil
}

// UpdateFlashSaleProduct updates a product in a flash sale
func (r *flashSaleRepository) UpdateFlashSaleProduct(flashSaleProduct *models.FlashSaleProduct) error {
	return r.db.Save(flashSaleProduct).Error
}
