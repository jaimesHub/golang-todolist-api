package database

import (
	"fmt"
	"os"

	"github.com/jaimesHub/golang-todo-app/internal/config"
	"github.com/joho/godotenv"
)

// SupabaseConfig holds Supabase specific configuration
type SupabaseConfig struct {
	URL       string
	APIKey    string
	ProjectID string
}

// LoadSupabaseConfig loads Supabase configuration from environment variables
func LoadSupabaseConfig() (*SupabaseConfig, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	url := os.Getenv("SUPABASE_URL")
	if url == "" {
		return nil, fmt.Errorf("SUPABASE_URL environment variable is required")
	}

	apiKey := os.Getenv("SUPABASE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SUPABASE_API_KEY environment variable is required")
	}

	projectID := os.Getenv("SUPABASE_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("SUPABASE_PROJECT_ID environment variable is required")
	}

	return &SupabaseConfig{
		URL:       url,
		APIKey:    apiKey,
		ProjectID: projectID,
	}, nil
}

// GetDatabaseConfig converts Supabase config to standard database config
func (s *SupabaseConfig) GetDatabaseConfig() config.DatabaseConfig {
	// Extract database connection details from Supabase URL
	// In a real implementation, you might need to parse the URL or use specific environment variables
	return config.DatabaseConfig{
		Host:     os.Getenv("SUPABASE_DB_HOST"),
		Port:     5432, // Default PostgreSQL port
		User:     os.Getenv("SUPABASE_DB_USER"),
		Password: os.Getenv("SUPABASE_DB_PASSWORD"),
		DBName:   os.Getenv("SUPABASE_DB_NAME"),
		SSLMode:  "require", // Supabase typically requires SSL
	}
}
