package repository

import (
	"github.com/karima-store/internal/models"
	"gorm.io/gorm"
)

type StockLogRepository interface {
	Create(log *models.StockLog) error
	WithTx(tx *gorm.DB) StockLogRepository
}

type stockLogRepository struct {
	db *gorm.DB
}

func NewStockLogRepository(db *gorm.DB) StockLogRepository {
	return &stockLogRepository{db: db}
}

func (r *stockLogRepository) Create(log *models.StockLog) error {
	return r.db.Create(log).Error
}

func (r *stockLogRepository) WithTx(tx *gorm.DB) StockLogRepository {
	return &stockLogRepository{db: tx}
}
