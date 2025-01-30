package entity

import (
	"errors"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	UserID    uint   `gorm:"not null;index"`
	Name      string `gorm:"type:varchar(100);not null"`
	Color     string `gorm:"type:varchar(50);default:'bg-blue-100'"`
	IconColor string `gorm:"type:varchar(50);default:'text-blue-500'"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.Name == "" {
		return errors.New("category name cannot be empty")
	}
	return nil
}
