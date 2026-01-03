package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleCustomer UserRole = "customer"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Basic Information
	FullName    string     `json:"full_name" gorm:"not null;size:100"`
	Email       string     `json:"email" gorm:"uniqueIndex;not null;size:100"`
	KratosID    string     `json:"kratos_id" gorm:"uniqueIndex;size:36"`
	Phone       string     `json:"phone" gorm:"size:20"`
	Password    string     `json:"-" gorm:"not null;size:255"`
	Avatar      string     `json:"avatar" gorm:"size:255"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Gender      string     `json:"gender" gorm:"size:10"`

	// Role & Status
	Role       UserRole `json:"role" gorm:"not null;default:'customer'"`
	IsVerified bool     `json:"is_verified" gorm:"default:false"`
	IsActive   bool     `json:"is_active" gorm:"default:true"`

	// Address Information
	Address    string `json:"address" gorm:"size:255"`
	City       string `json:"city" gorm:"size:100"`
	Province   string `json:"province" gorm:"size:100"`
	PostalCode string `json:"postal_code" gorm:"size:10"`

	// Timestamps
	LastLoginAt *time.Time `json:"last_login_at"`

	// Relations
	Orders   []Order    `json:"orders,omitempty" gorm:"foreignKey:UserID"`
	Cart     *Cart      `json:"cart,omitempty" gorm:"foreignKey:UserID"`
	Reviews  []Review   `json:"reviews,omitempty" gorm:"foreignKey:UserID"`
	Wishlist []Wishlist `json:"wishlist,omitempty" gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}
