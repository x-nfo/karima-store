package services

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
)

type VariantService interface {
	CreateVariant(variant *models.ProductVariant) error
	GetVariantByID(id uint) (*models.ProductVariant, error)
	GetVariantBySKU(sku string) (*models.ProductVariant, error)
	GetVariantsByProductID(productID uint) ([]models.ProductVariant, error)
	UpdateVariant(id uint, variant *models.ProductVariant) error
	DeleteVariant(id uint) error
	UpdateVariantStock(id uint, quantity int) error
	GenerateSKU(productName, size, color string) string
}

type variantService struct {
	variantRepo repository.VariantRepository
	productRepo repository.ProductRepository
}

func NewVariantService(variantRepo repository.VariantRepository, productRepo repository.ProductRepository) VariantService {
	return &variantService{
		variantRepo: variantRepo,
		productRepo: productRepo,
	}
}

func (s *variantService) CreateVariant(variant *models.ProductVariant) error {
	// Validate required fields
	if variant.ProductID == 0 {
		return errors.New("product ID is required")
	}
	if variant.Name == "" {
		return errors.New("variant name is required")
	}
	if variant.Price <= 0 {
		return errors.New("variant price must be greater than 0")
	}

	// Check if product exists
	_, err := s.productRepo.GetByID(variant.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	// Generate SKU if not provided
	if variant.SKU == "" {
		product, err := s.productRepo.GetByID(variant.ProductID)
		if err != nil {
			return err
		}
		variant.SKU = s.GenerateSKU(product.Name, variant.Size, variant.Color)
	}

	// Check if SKU already exists
	existingVariant, err := s.variantRepo.GetBySKU(variant.SKU)
	if err == nil && existingVariant != nil {
		return errors.New("variant with this SKU already exists")
	}

	return s.variantRepo.Create(variant)
}

func (s *variantService) GetVariantByID(id uint) (*models.ProductVariant, error) {
	variant, err := s.variantRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("variant not found")
		}
		return nil, err
	}
	return variant, nil
}

func (s *variantService) GetVariantBySKU(sku string) (*models.ProductVariant, error) {
	variant, err := s.variantRepo.GetBySKU(sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("variant not found")
		}
		return nil, err
	}
	return variant, nil
}

func (s *variantService) GetVariantsByProductID(productID uint) ([]models.ProductVariant, error) {
	// Check if product exists
	_, err := s.productRepo.GetByID(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return s.variantRepo.GetByProductID(productID)
}

func (s *variantService) UpdateVariant(id uint, variant *models.ProductVariant) error {
	// Check if variant exists
	existingVariant, err := s.variantRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("variant not found")
		}
		return err
	}

	// Ensure ID is set
	variant.ID = id

	// If product ID is being changed, validate it exists
	if variant.ProductID != 0 && variant.ProductID != existingVariant.ProductID {
		_, err := s.productRepo.GetByID(variant.ProductID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("product not found")
			}
			return err
		}
	}

	return s.variantRepo.Update(variant)
}

func (s *variantService) DeleteVariant(id uint) error {
	// Check if variant exists
	_, err := s.variantRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("variant not found")
		}
		return err
	}

	return s.variantRepo.Delete(id)
}

func (s *variantService) UpdateVariantStock(id uint, quantity int) error {
	// Check if variant exists
	variant, err := s.variantRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("variant not found")
		}
		return err
	}

	// Validate stock update
	newStock := variant.Stock + quantity
	if newStock < 0 {
		return errors.New("insufficient stock")
	}

	return s.variantRepo.UpdateStock(id, quantity)
}

func (s *variantService) GenerateSKU(productName, size, color string) string {
	// Get first 3 letters of product name (uppercase)
	namePart := strings.ToUpper(productName)
	if len(namePart) > 3 {
		namePart = namePart[:3]
	}

	// Get size code (uppercase)
	sizePart := strings.ToUpper(size)
	if len(sizePart) > 2 {
		sizePart = sizePart[:2]
	}

	// Get color code (uppercase)
	colorPart := strings.ToUpper(color)
	if len(colorPart) > 3 {
		colorPart = colorPart[:3]
	}

	// Generate SKU format: NAME-SIZE-COLOR-XXXX
	sku := fmt.Sprintf("%s-%s-%s", namePart, sizePart, colorPart)

	// Clean up any special characters
	sku = strings.Map(func(r rune) rune {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, sku)

	// Remove multiple consecutive hyphens
	sku = strings.ReplaceAll(sku, "--", "-")

	// Trim hyphens from start and end
	sku = strings.Trim(sku, "-")

	return sku
}
