package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	entity2 "go-manajemen-keuangan/internal/payload/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func ConnectDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
		&entity2.User{},
		&entity2.Category{},
		&entity2.Transaction{},
	); err != nil {
		logrus.Fatal("Auto migration failed:", err)
	}

	return db
}
