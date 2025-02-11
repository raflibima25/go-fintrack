package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name       string `gorm:"type:varchar(255);not null"`
	Email      string `gorm:"type:varchar(255);unique;not null"`
	Username   string `gorm:"type:varchar(50);unique;not null"`
	Password   string `gorm:"type:varchar(255);omitempty"`
	IsAdmin    bool   `gorm:"type:boolean;default:false"`
	Provider   string `gorm:"type:varchar(50);omitempty"`
	ProfilePic string `gorm:"type:varchar(255);omitempty"`
}
