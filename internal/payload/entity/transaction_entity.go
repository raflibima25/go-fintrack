package entity

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	gorm.Model
	UserID      uint      `gorm:"not null"`
	CategoryID  uint      `gorm:"not null"`
	Amount      float64   `gorm:"not null"`
	Type        string    `gorm:"size:20;not null"` // income atau expense
	Description string    `gorm:"type:text"`
	Date        time.Time `gorm:"not null"`
	User        User      `gorm:"foreignKey:UserID"`
	Category    Category  `gorm:"foreignKey:CategoryID"`
}

func (t *Transaction) BeforeSave(tx *gorm.DB) (err error) {
	// validasi tipe transaksi
	validTypes := map[string]bool{
		"income":  true,
		"expense": true,
	}
	if !validTypes[t.Type] {
		return fmt.Errorf("invalid transaction type: %s", t.Type)
	}

	// validasi jumlah
	if t.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	return nil
}
