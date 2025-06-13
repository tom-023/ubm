package validator

import (
	"fmt"
	"net/url"
	"strings"
)

func ValidateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Parse URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Check if scheme is present
	if u.Scheme == "" {
		return fmt.Errorf("URL must include a scheme (e.g., http:// or https://)")
	}

	// Check if scheme is valid
	validSchemes := []string{"http", "https", "ftp", "ftps"}
	schemeValid := false
	for _, scheme := range validSchemes {
		if strings.ToLower(u.Scheme) == scheme {
			schemeValid = true
			break
		}
	}
	if !schemeValid {
		return fmt.Errorf("unsupported URL scheme: %s", u.Scheme)
	}

	// Check if host is present
	if u.Host == "" {
		return fmt.Errorf("URL must include a host")
	}

	return nil
}

func NormalizeURL(rawURL string) (string, error) {
	// Trim whitespace
	rawURL = strings.TrimSpace(rawURL)

	// If no scheme, add https://
	if !strings.Contains(rawURL, "://") {
		rawURL = "https://" + rawURL
	}

	// Validate the URL
	if err := ValidateURL(rawURL); err != nil {
		return "", err
	}

	// Parse and rebuild to normalize
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}