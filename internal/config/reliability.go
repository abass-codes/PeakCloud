package config

import (
	"fmt"
	"strconv"
	"time"
)

type ReliabilityConfig struct {
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func loadReliabilityConfig() (ReliabilityConfig, error) {
	readHeaderSeconds, err := positiveIntEnv(
		"HTTP_READ_HEADER_TIMEOUT_SECONDS",
		5,
	)
	if err != nil {
		return ReliabilityConfig{}, err
	}

	readSeconds, err := positiveIntEnv(
		"HTTP_READ_TIMEOUT_SECONDS",
		15,
	)
	if err != nil {
		return ReliabilityConfig{}, err
	}

	writeSeconds, err := positiveIntEnv(
		"HTTP_WRITE_TIMEOUT_SECONDS",
		15,
	)
	if err != nil {
		return ReliabilityConfig{}, err
	}

	idleSeconds, err := positiveIntEnv(
		"HTTP_IDLE_TIMEOUT_SECONDS",
		60,
	)
	if err != nil {
		return ReliabilityConfig{}, err
	}

	shutdownSeconds, err := positiveIntEnv(
		"SHUTDOWN_TIMEOUT_SECONDS",
		10,
	)
	if err != nil {
		return ReliabilityConfig{}, err
	}

	return ReliabilityConfig{
		ReadHeaderTimeout: time.Duration(readHeaderSeconds) * time.Second,
		ReadTimeout:       time.Duration(readSeconds) * time.Second,
		WriteTimeout:      time.Duration(writeSeconds) * time.Second,
		IdleTimeout:       time.Duration(idleSeconds) * time.Second,
		ShutdownTimeout:   time.Duration(shutdownSeconds) * time.Second,
	}, nil
}

func positiveIntEnv(key string, fallback int) (int, error) {
	value, err := strconv.Atoi(
		getEnv(key, strconv.Itoa(fallback)),
	)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}

	if value <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", key)
	}

	return value, nil
}
