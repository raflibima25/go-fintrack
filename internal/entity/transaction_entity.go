package entity

import (
	"fmt"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID      uint
	CategoryID  uint
	Amount      float64
	Type        string `gorm:"size:20;not null"`
	Description string
}

func (t *Transaction) BeforeSave(tx *gorm.DB) (err error) {
	validTypes := map[string]bool{
		"income":  true,
		"expense": false,
	}
	if !validTypes[t.Type] {
		return fmt.Errorf("invalid transaction type: %s", t.Type)
	}
	return nil
}
