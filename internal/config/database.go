package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go-manajemen-keuangan/internal/payload/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func ConnectDB() *gorm.DB {
	var (
		host     = os.Getenv("DB_HOST")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
		port     = os.Getenv("DB_PORT")
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Errorf("Failed connect to the database: %v", err)
	}

	if err = db.AutoMigrate(
		&entity.User{},
		&entity.Category{},
		&entity.Transaction{},
	); err != nil {
		logrus.Fatal("Auto migration failed:", err)
	}

	return db
}
