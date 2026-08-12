package config

import (
	"strings"
	"testing"
)

func TestValidateProductionAllowsDevelopment(t *testing.T) {
	cfg := &Config{
		AppEnv:        "development",
		WebURL:        "http://localhost:3000",
		SessionSecure: false,
		S3UseSSL:      false,
	}

	if err := cfg.ValidateProduction(); err != nil {
		t.Fatalf("expected development config to pass: %v", err)
	}
}

func TestValidateProductionAcceptsSecureProduction(t *testing.T) {
	cfg := &Config{
		AppEnv:        "production",
		WebURL:        "https://peakcloud.example.com",
		SessionSecure: true,
		S3UseSSL:      true,
	}

	if err := cfg.ValidateProduction(); err != nil {
		t.Fatalf("expected production config to pass: %v", err)
	}
}

func TestValidateProductionRequiresSecureSessionCookie(t *testing.T) {
	cfg := &Config{
		AppEnv:        "production",
		WebURL:        "https://peakcloud.example.com",
		SessionSecure: false,
		S3UseSSL:      true,
	}

	err := cfg.ValidateProduction()
	if err == nil {
		t.Fatal("expected production validation to fail")
	}

	if !strings.Contains(err.Error(), "SESSION_SECURE") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateProductionRequiresSecureObjectStorage(t *testing.T) {
	cfg := &Config{
		AppEnv:        "production",
		WebURL:        "https://peakcloud.example.com",
		SessionSecure: true,
		S3UseSSL:      false,
	}

	err := cfg.ValidateProduction()
	if err == nil {
		t.Fatal("expected production validation to fail")
	}

	if !strings.Contains(err.Error(), "S3_USE_SSL") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateProductionRequiresHTTPSWebURL(t *testing.T) {
	cfg := &Config{
		AppEnv:        "production",
		WebURL:        "http://peakcloud.example.com",
		SessionSecure: true,
		S3UseSSL:      true,
	}

	err := cfg.ValidateProduction()
	if err == nil {
		t.Fatal("expected production validation to fail")
	}

	if !strings.Contains(err.Error(), "https") {
		t.Fatalf("unexpected error: %v", err)
	}
}
