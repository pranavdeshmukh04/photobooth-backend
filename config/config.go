package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	// Server
	Port        string
	Environment string
	APIPrefix   string
	AllowedOrigins []string

	// JWT
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration

	// Couchbase
	CouchbaseURL      string
	CouchbaseUsername string
	CouchbasePassword string
	CouchbaseBucket   string
	CouchbaseTimeout  time.Duration

	// Cloudflare R2
	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2BucketName      string
	R2Endpoint        string
	R2PublicURL       string

	// Stripe
	StripeSecretKey      string
	StripeWebhookSecret  string
	StripePriceIDMonthly string
	StripePriceIDYearly  string

	// Frontend
	FrontendURL string

	// Rate Limiting
	RateLimitEnabled  bool
	RateLimitRequests int
	RateLimitWindow   time.Duration
}

// AppConfig is the global application configuration
var AppConfig *Config

// Load loads the configuration from the environment variables
func Load() error {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get environment variables
	port := os.Getenv("PORT")
	environment := os.Getenv("ENVIRONMENT")
	apiPrefix := os.Getenv("API_PREFIX")
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtAccessExpiry := os.Getenv("JWT_ACCESS_TOKEN_EXPIRY")
	jwtRefreshExpiry := os.Getenv("JWT_ACCESS_TOKEN_EXPIRY")
	couchbaseURL := os.Getenv("COUCHBASE_URL")
	couchbaseUsername := os.Getenv("COUCHBASE_USERNAME")
	couchbasePassword := os.Getenv("COUCHBASE_PASSWORD")
	couchbaseBucket := os.Getenv("COUCHBASE_BUCKET")
	couchbaseTimeout := os.Getenv("COUCHBASE_TIMEOUT")
	r2AccountID := os.Getenv("R2_ACCOUNT_ID")
	r2AccessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	r2SecretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	r2BucketName := os.Getenv("R2_BUCKET_NAME")
	r2Endpoint := os.Getenv("R2_ENDPOINT")
	r2PublicURL := os.Getenv("R2_PUBLIC_URL")
	stripeSecretKey := os.Getenv("STRIPE_SECRET_KEY")
	stripeWebhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	stripePriceIDMonthly := os.Getenv("STRIPE_PRICE_ID_MONTHLY")
	stripePriceIDYearly := os.Getenv("STRIPE_PRICE_ID_YEARLY")
	frontendURL := os.Getenv("FRONTEND_URL")
	rateLimitEnabled := os.Getenv("RATE_LIMIT_ENABLED")
	rateLimitRequests := os.Getenv("RATE_LIMIT_REQUESTS")
	rateLimitWindow := os.Getenv("RATE_LIMIT_WINDOW")

	// Set default values
	switch{
		case port == "":
			port = "8080"
		case environment == "":
			environment = "development"
		case apiPrefix == "":
			apiPrefix = "/api/v1"
	}

	// Parse duration values for JWT and Couchbase
	jwtAccessExpiryDuration, err := time.ParseDuration(jwtAccessExpiry)
	if err != nil {
		return fmt.Errorf("invalid JWT access expiry duration: %w", err)
	}
	jwtRefreshExpiryDuration, err := time.ParseDuration(jwtRefreshExpiry)
	if err != nil {
		return fmt.Errorf("invalid JWT refresh expiry duration: %w", err)
	}
	couchbaseTimeoutDuration, err := time.ParseDuration(couchbaseTimeout)
	if err != nil {
		return fmt.Errorf("invalid Couchbase timeout duration: %w", err)
	}

	// Parse boolean and integer values for rate limiting
	rateLimitEnabledBool, err := strconv.ParseBool(rateLimitEnabled)
	if err != nil {
		return fmt.Errorf("invalid rate limit enabled: %w", err)
	}
	rateLimitRequestsInt, err := strconv.Atoi(rateLimitRequests)
	if err != nil {
		return fmt.Errorf("invalid rate limit requests: %w", err)
	}
	rateLimitWindowDuration, err := time.ParseDuration(rateLimitWindow)
	if err != nil {
		return fmt.Errorf("invalid rate limit window duration: %w", err)
	}

	allowedOrigins := parseAllowedOrigins(allowedOriginsStr)

	// Create configuration
	cfg := &Config{
		Port: port,
		Environment: environment,
		APIPrefix: apiPrefix,
		AllowedOrigins: allowedOrigins,
		JWTSecret: jwtSecret,
		JWTAccessExpiry: jwtAccessExpiryDuration,
		JWTRefreshExpiry: jwtRefreshExpiryDuration,
		CouchbaseURL: couchbaseURL,
		CouchbaseUsername: couchbaseUsername,
		CouchbasePassword: couchbasePassword,
		CouchbaseBucket: couchbaseBucket,
		CouchbaseTimeout: couchbaseTimeoutDuration,
		R2AccountID: r2AccountID,
		R2AccessKeyID: r2AccessKeyID,
		R2SecretAccessKey: r2SecretAccessKey,
		R2BucketName: r2BucketName,
		R2Endpoint: r2Endpoint,
		R2PublicURL: r2PublicURL,
		StripeSecretKey: stripeSecretKey,
		StripeWebhookSecret: stripeWebhookSecret,
		StripePriceIDMonthly: stripePriceIDMonthly,
		StripePriceIDYearly: stripePriceIDYearly,
		FrontendURL: frontendURL,
		RateLimitEnabled: rateLimitEnabledBool,
		RateLimitRequests: rateLimitRequestsInt,
		RateLimitWindow: rateLimitWindowDuration,
	}

	// Validate configuration
	if err := validateConfig(cfg); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	AppConfig = cfg

	log.Println("Configuration loaded successfully")

	return nil
}

// GetConfig returns the global application configuration
func GetConfig() *Config {
	return AppConfig
}

// validateConfig validates the configuration
func validateConfig(cfg *Config) error {
	requiredFields := map[string]string{
		"JWT_SECRET":           cfg.JWTSecret,
		"COUCHBASE_URL":        cfg.CouchbaseURL,
		"COUCHBASE_USERNAME":   cfg.CouchbaseUsername,
		"COUCHBASE_PASSWORD":   cfg.CouchbasePassword,
		"COUCHBASE_BUCKET":     cfg.CouchbaseBucket,
		"R2_ACCOUNT_ID":        cfg.R2AccountID,
		"R2_ACCESS_KEY_ID":     cfg.R2AccessKeyID,
		"R2_SECRET_ACCESS_KEY": cfg.R2SecretAccessKey,
		"R2_BUCKET_NAME":       cfg.R2BucketName,
		"STRIPE_SECRET_KEY":    cfg.StripeSecretKey,
	}

	for field, value := range requiredFields {
		if value == "" {
			return fmt.Errorf("%s is required", field)
		}
	}

	log.Println("Configuration validated successfully")

	return nil
}

func parseAllowedOrigins(allowedOriginsStr string) []string {
    if allowedOriginsStr == "" {
        // Try FRONTEND_URL as fallback
        frontendURL := os.Getenv("FRONTEND_URL")
        if frontendURL != "" {
            return []string{frontendURL}
        }
        return []string{"http://localhost:3000"} // default
    }
    
    // Split and trim
    origins := strings.Split(allowedOriginsStr, ",")
    result := make([]string, 0, len(origins))
    
    for _, origin := range origins {
        trimmed := strings.TrimSpace(origin)
        if trimmed != "" {
            result = append(result, trimmed)
        }
    }
    
    return result
}