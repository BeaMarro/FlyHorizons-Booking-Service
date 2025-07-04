package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type BaseRepository struct {
	DB *gorm.DB
}

func (dal *BaseRepository) LoadDBConfig() string {
	// Load environment variables from .env file (optional)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Get environment variables
	server := os.Getenv("DB_SERVER")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")

	if server == "" || port == "" || user == "" || password == "" || database == "" {
		log.Fatalf("Missing environment variables for database connection")
	}

	// Build and return database connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		server, user, password, port, database)
	return connString
}

func (dal *BaseRepository) CreateConnection() (*gorm.DB, error) {
	if dal.DB != nil {
		return dal.DB, nil
	}

	connString := dal.LoadDBConfig()
	if connString == "" {
		return nil, fmt.Errorf("failed to load database configuration: environment variables missing")
	}

	// Connect to database
	sqlDB, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	// Ensure a successful connection by pinging the database
	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging the database: %w", err)
	}

	// Initialize GORM
	db, err := gorm.Open(sqlserver.New(sqlserver.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error initializing GORM: %w", err)
	}

	dal.DB = db
	return dal.DB, nil
}

func (dal *BaseRepository) CloseConnection() {
	if dal.DB != nil {
		sqlDB, _ := dal.DB.DB()
		sqlDB.Close()
	}
}
