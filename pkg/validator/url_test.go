package validator

import (
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid http URL",
			input:   "http://example.com",
			wantErr: false,
		},
		{
			name:    "valid https URL",
			input:   "https://example.com",
			wantErr: false,
		},
		{
			name:    "valid https URL with path",
			input:   "https://example.com/path/to/page",
			wantErr: false,
		},
		{
			name:    "valid https URL with query",
			input:   "https://example.com?param=value",
			wantErr: false,
		},
		{
			name:    "valid ftp URL",
			input:   "ftp://ftp.example.com",
			wantErr: false,
		},
		{
			name:    "valid ftps URL",
			input:   "ftps://ftp.example.com",
			wantErr: false,
		},
		{
			name:    "empty URL",
			input:   "",
			wantErr: true,
			errMsg:  "URL cannot be empty",
		},
		{
			name:    "URL without scheme",
			input:   "example.com",
			wantErr: true,
			errMsg:  "URL must include a scheme",
		},
		{
			name:    "URL with invalid scheme",
			input:   "gopher://example.com",
			wantErr: true,
			errMsg:  "unsupported URL scheme: gopher",
		},
		{
			name:    "URL without host",
			input:   "https://",
			wantErr: true,
			errMsg:  "URL must include a host",
		},
		{
			name:    "malformed URL",
			input:   "https://[invalid",
			wantErr: true,
			errMsg:  "invalid URL format",
		},
		{
			name:    "URL with uppercase scheme",
			input:   "HTTPS://example.com",
			wantErr: false,
		},
		{
			name:    "URL with port",
			input:   "https://example.com:8080",
			wantErr: false,
		},
		{
			name:    "URL with fragment",
			input:   "https://example.com#section",
			wantErr: false,
		},
		{
			name:    "URL with auth",
			input:   "https://user:pass@example.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				// Check if error message contains expected text for wrapped errors
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateURL() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "URL with scheme",
			input:   "https://example.com",
			want:    "https://example.com",
			wantErr: false,
		},
		{
			name:    "URL without scheme",
			input:   "example.com",
			want:    "https://example.com",
			wantErr: false,
		},
		{
			name:    "URL with whitespace",
			input:   "  https://example.com  ",
			want:    "https://example.com",
			wantErr: false,
		},
		{
			name:    "URL without scheme and with path",
			input:   "example.com/path",
			want:    "https://example.com/path",
			wantErr: false,
		},
		{
			name:    "URL with http scheme",
			input:   "http://example.com",
			want:    "http://example.com",
			wantErr: false,
		},
		{
			name:    "empty URL",
			input:   "",
			want:    "",
			wantErr: true,
			errMsg:  "URL cannot be empty",
		},
		{
			name:    "URL with only whitespace",
			input:   "   ",
			want:    "",
			wantErr: true,
			errMsg:  "URL cannot be empty",
		},
		{
			name:    "invalid URL after normalization",
			input:   "://invalid",
			want:    "",
			wantErr: true,
			errMsg:  "invalid URL format",
		},
		{
			name:    "URL with query parameters",
			input:   "example.com?q=test&page=1",
			want:    "https://example.com?q=test&page=1",
			wantErr: false,
		},
		{
			name:    "URL with fragment",
			input:   "example.com#section",
			want:    "https://example.com#section",
			wantErr: false,
		},
		{
			name:    "URL with port",
			input:   "example.com:8080",
			want:    "https://example.com:8080",
			wantErr: false,
		},
		{
			name:    "URL with path and trailing slash",
			input:   "example.com/path/",
			want:    "https://example.com/path/",
			wantErr: false,
		},
		{
			name:    "URL with special characters in path",
			input:   "example.com/path%20with%20spaces",
			want:    "https://example.com/path%20with%20spaces",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("NormalizeURL() error message = %v, want containing %v", err.Error(), tt.errMsg)
				}
			}
			if got != tt.want {
				t.Errorf("NormalizeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizeURL_IdempotentBehavior(t *testing.T) {
	// Test that normalizing an already normalized URL doesn't change it
	tests := []string{
		"https://example.com",
		"http://example.com/path",
		"https://example.com:8080/path?query=value#fragment",
	}

	for _, url := range tests {
		t.Run(url, func(t *testing.T) {
			normalized1, err := NormalizeURL(url)
			if err != nil {
				t.Fatalf("First normalization failed: %v", err)
			}

			normalized2, err := NormalizeURL(normalized1)
			if err != nil {
				t.Fatalf("Second normalization failed: %v", err)
			}

			if normalized1 != normalized2 {
				t.Errorf("Normalization is not idempotent: first=%v, second=%v", normalized1, normalized2)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}