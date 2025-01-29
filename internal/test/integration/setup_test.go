// test/integration/setup_test.go
package integration

import (
	"fmt"
	"go-fintrack/internal/controller"
	"go-fintrack/internal/payload/entity"
	"go-fintrack/internal/service"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestServer struct {
	DB          *gorm.DB
	Router      *gin.Engine
	UserService *service.UserService
}

// setupTestDatabase menginisialisasi koneksi database untuk testing
func setupTestDatabase() (*gorm.DB, error) {
	// Load env test
	if err := godotenv.Load("../../../.env.test"); err != nil {
		return nil, fmt.Errorf("error loading .env.test file: %v", err)
	}

	// Baca environment variables
	var (
		host     = os.Getenv("DB_HOST")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
		port     = os.Getenv("DB_PORT")
	)

	// Buat connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		host, user, password, dbname, port)

	// Buka koneksi database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %v", err)
	}

	// Auto migrate schemas
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Category{},
		&entity.Transaction{},
	); err != nil {
		return nil, fmt.Errorf("failed to run auto migration: %v", err)
	}

	return db, nil
}

// setupTestServer menginisialisasi server untuk testing
func setupTestServer(t *testing.T) *TestServer {
	// Setup database
	db, err := setupTestDatabase()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	// Setup services
	userService := &service.UserService{DB: db}

	// Setup router dengan mode test
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup controllers
	userController := &controller.UserController{
		UserService: userService,
	}

	// Setup routes
	router.POST("/api/v1/user/register", userController.RegisterHandler)
	router.POST("/api/v1/user/login", userController.LoginHandler)

	return &TestServer{
		DB:          db,
		Router:      router,
		UserService: userService,
	}
}

// cleanupDatabase membersihkan database setelah testing
func cleanupDatabase(db *gorm.DB) error {
	// Hapus semua data dari tabel-tabel
	if err := db.Exec(`
        TRUNCATE TABLE users, categories, transactions CASCADE;
    `).Error; err != nil {
		return fmt.Errorf("failed to cleanup database: %v", err)
	}
	return nil
}

// cleanupTest helper function untuk membersihkan test
func cleanupTest(t *testing.T, ts *TestServer) {
	if err := cleanupDatabase(ts.DB); err != nil {
		t.Logf("Warning: Failed to cleanup test database: %v", err)
	}
}
