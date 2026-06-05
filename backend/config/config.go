package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const DefaultJWTSecret = "your-secret-key-change-in-production"

type Config struct {
	DatabaseURL   string
	JWTSecret     string
	Port          string
	GinMode       string
	CORSOrigins   []string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
	OTPTTL        time.Duration
	DevFixedOTP   bool
}

func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	ginMode := os.Getenv("GIN_MODE")
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = DefaultJWTSecret
	}
	if ginMode == "release" && jwtSecret == DefaultJWTSecret {
		return nil, fmt.Errorf("JWT_SECRET must be set in production")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	corsOrigins := []string{"http://localhost:3000"}
	if raw := os.Getenv("CORS_ORIGINS"); raw != "" {
		corsOrigins = strings.Split(raw, ",")
		for i := range corsOrigins {
			corsOrigins[i] = strings.TrimSpace(corsOrigins[i])
		}
	}

	return &Config{
		DatabaseURL: dbURL,
		JWTSecret:   jwtSecret,
		Port:        port,
		GinMode:     ginMode,
		CORSOrigins: corsOrigins,
		AccessTTL:   15 * time.Minute,
		RefreshTTL:  30 * 24 * time.Hour,
		OTPTTL:      5 * time.Minute,
		DevFixedOTP: ginMode != "release",
	}, nil
}
