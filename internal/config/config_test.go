package config

import "testing"

func TestGetEnvReturnsFallback(t *testing.T) {
	const key = "PEAKCLOUD_TEST_MISSING_ENV"

	t.Setenv(key, "")

	got := getEnv(key, "fallback")

	if got != "fallback" {
		t.Fatalf("expected fallback, got %q", got)
	}
}

func TestGetEnvReturnsEnvironmentValue(t *testing.T) {
	const key = "PEAKCLOUD_TEST_ENV"

	t.Setenv(key, "configured")

	got := getEnv(key, "fallback")

	if got != "configured" {
		t.Fatalf("expected configured, got %q", got)
	}
}
