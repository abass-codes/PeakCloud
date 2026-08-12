package config

import (
	"testing"
	"time"
)

func TestLoadReliabilityConfigDefaults(t *testing.T) {
	keys := []string{
		"HTTP_READ_HEADER_TIMEOUT_SECONDS",
		"HTTP_READ_TIMEOUT_SECONDS",
		"HTTP_WRITE_TIMEOUT_SECONDS",
		"HTTP_IDLE_TIMEOUT_SECONDS",
		"SHUTDOWN_TIMEOUT_SECONDS",
	}

	for _, key := range keys {
		t.Setenv(key, "")
	}

	cfg, err := loadReliabilityConfig()
	if err != nil {
		t.Fatalf("load reliability config: %v", err)
	}

	if cfg.ReadHeaderTimeout != 5*time.Second {
		t.Fatalf(
			"expected read header timeout 5s, got %s",
			cfg.ReadHeaderTimeout,
		)
	}

	if cfg.ReadTimeout != 15*time.Second {
		t.Fatalf(
			"expected read timeout 15s, got %s",
			cfg.ReadTimeout,
		)
	}

	if cfg.WriteTimeout != 15*time.Second {
		t.Fatalf(
			"expected write timeout 15s, got %s",
			cfg.WriteTimeout,
		)
	}

	if cfg.IdleTimeout != 60*time.Second {
		t.Fatalf(
			"expected idle timeout 60s, got %s",
			cfg.IdleTimeout,
		)
	}

	if cfg.ShutdownTimeout != 10*time.Second {
		t.Fatalf(
			"expected shutdown timeout 10s, got %s",
			cfg.ShutdownTimeout,
		)
	}
}

func TestLoadReliabilityConfigCustomValues(t *testing.T) {
	t.Setenv("HTTP_READ_HEADER_TIMEOUT_SECONDS", "7")
	t.Setenv("HTTP_READ_TIMEOUT_SECONDS", "20")
	t.Setenv("HTTP_WRITE_TIMEOUT_SECONDS", "25")
	t.Setenv("HTTP_IDLE_TIMEOUT_SECONDS", "90")
	t.Setenv("SHUTDOWN_TIMEOUT_SECONDS", "30")

	cfg, err := loadReliabilityConfig()
	if err != nil {
		t.Fatalf("load reliability config: %v", err)
	}

	if cfg.ReadHeaderTimeout != 7*time.Second {
		t.Fatalf("unexpected read header timeout")
	}

	if cfg.ReadTimeout != 20*time.Second {
		t.Fatalf("unexpected read timeout")
	}

	if cfg.WriteTimeout != 25*time.Second {
		t.Fatalf("unexpected write timeout")
	}

	if cfg.IdleTimeout != 90*time.Second {
		t.Fatalf("unexpected idle timeout")
	}

	if cfg.ShutdownTimeout != 30*time.Second {
		t.Fatalf("unexpected shutdown timeout")
	}
}

func TestLoadReliabilityConfigRejectsInvalidValue(t *testing.T) {
	t.Setenv("HTTP_READ_TIMEOUT_SECONDS", "invalid")

	_, err := loadReliabilityConfig()
	if err == nil {
		t.Fatal("expected invalid timeout error")
	}
}

func TestLoadReliabilityConfigRejectsNonPositiveValue(t *testing.T) {
	t.Setenv("SHUTDOWN_TIMEOUT_SECONDS", "0")

	_, err := loadReliabilityConfig()
	if err == nil {
		t.Fatal("expected non-positive timeout error")
	}
}
