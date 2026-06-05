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

	corsOrigins := []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
	}
	if frontendURL := strings.TrimSpace(os.Getenv("FRONTEND_URL")); frontendURL != "" {
		corsOrigins = append(corsOrigins, frontendURL)
	}
	if raw := os.Getenv("CORS_ORIGINS"); raw != "" {
		for _, o := range strings.Split(raw, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				corsOrigins = append(corsOrigins, o)
			}
		}
	}
	corsOrigins = uniqueStrings(corsOrigins)

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

func uniqueStrings(items []string) []string {
	seen := make(map[string]bool, len(items))
	out := make([]string, 0, len(items))
	for _, item := range items {
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}
	return out
}

// IsAllowedOrigin returns true for configured origins, Railway deploys, and proxied requests (no Origin).
func (c *Config) IsAllowedOrigin(origin string) bool {
	if origin == "" {
		return true
	}
	for _, allowed := range c.CORSOrigins {
		if origin == allowed {
			return true
		}
	}
	// Railway production/staging frontends
	if strings.HasSuffix(origin, ".railway.app") || strings.HasSuffix(origin, ".up.railway.app") {
		return true
	}
	return false
}
