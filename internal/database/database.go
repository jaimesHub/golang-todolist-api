package database

import (
	"fmt"

	"github.com/jaimesHub/golang-todo-app/internal/config"
	"github.com/jaimesHub/golang-todo-app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a connection to the database using standard config
func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	return connectWithDSN(dsn)
}

// ConnectWithSupabase establishes a connection to the database using Supabase config
func ConnectWithSupabase() (*gorm.DB, error) {
	supabaseConfig, err := LoadSupabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load Supabase config: %w", err)
	}

	dbConfig := supabaseConfig.GetDatabaseConfig()
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode,
	)

	return connectWithDSN(dsn)
}

// connectWithDSN is a helper function to connect to the database with a DSN
func connectWithDSN(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	// Add models to migrate here
	return db.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.Activity{},
	)
}

// RunSQLMigration executes SQL migration files
func RunSQLMigration(db *gorm.DB, filePath string) error {
	// This is a placeholder for running SQL migrations
	// In a real implementation, you would read the SQL file and execute it
	// For now, we'll rely on GORM's AutoMigrate
	return nil
}
