package config

import (
	"fmt"
	"net/url"
	"strings"
)

func (c *Config) ValidateProduction() error {
	if !strings.EqualFold(c.AppEnv, "production") {
		return nil
	}

	if !c.SessionSecure {
		return fmt.Errorf("SESSION_SECURE must be true in production")
	}

	if !c.S3UseSSL {
		return fmt.Errorf("S3_USE_SSL must be true in production")
	}

	webURL, err := url.Parse(c.WebURL)
	if err != nil {
		return fmt.Errorf("parse WEB_URL: %w", err)
	}

	if webURL.Scheme != "https" {
		return fmt.Errorf("WEB_URL must use https in production")
	}

	if webURL.Host == "" {
		return fmt.Errorf("WEB_URL must contain a host in production")
	}

	return nil
}
