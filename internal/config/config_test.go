package config

import (
	"os"
	"testing"
)

func Test_getEnv(t *testing.T) {
	t.Run("returns value when env exists", func(t *testing.T) {
		t.Setenv("TEST_ENV_KEY", "test-value")

		result := getEnv("TEST_ENV_KEY")

		if result != "test-value" {
			t.Errorf("getEnv should return env value, got %q", result)
		}
	})

	t.Run("returns empty string when env missing", func(t *testing.T) {
		key := "TEST_ENV_MISSING"

		oldValue, existed := os.LookupEnv(key)
		_ = os.Unsetenv(key)
		t.Cleanup(func() {
			if existed {
				_ = os.Setenv(key, oldValue)
				return
			}
			_ = os.Unsetenv(key)
		})

		result := getEnv(key)

		if result != "" {
			t.Errorf("getEnv should return empty string for missing key, got %q", result)
		}
	})
}

func Test_NewConfig(t *testing.T) {
	tests := []struct {
		name     string
		envs     map[string]string
		expected Config
	}{
		{
			name: "all env vars set",
			envs: map[string]string{
				"APP_PORT":     "8080",
				"DATABASE_URL": "postgres://localhost/db",
				"DEBUG":        "true",
				"JWT_SECRET":   "mysecret",
			},
			expected: Config{
				AppPort: "8080",
				DBUrl:   "postgres://localhost/db",
				Debug:   true,
				Secret:  "mysecret",
			},
		},
		{
			name: "debug false when not set",
			envs: map[string]string{
				"APP_PORT":   "3000",
				"JWT_SECRET": "s3cr3t",
			},
			expected: Config{
				AppPort: "3000",
				Debug:   false,
				Secret:  "s3cr3t",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ustaw zmienne środowiskowe tylko dla tego podtestu
			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			cfg := NewConfig()

			if cfg.AppPort != tt.expected.AppPort {
				t.Errorf("AppPort: expected %q, got %q", tt.expected.AppPort, cfg.AppPort)
			}
			if cfg.DBUrl != tt.expected.DBUrl {
				t.Errorf("DBUrl: expected %q, got %q", tt.expected.DBUrl, cfg.DBUrl)
			}
			if cfg.Debug != tt.expected.Debug {
				t.Errorf("Debug: expected %v, got %v", tt.expected.Debug, cfg.Debug)
			}
			if cfg.Secret != tt.expected.Secret {
				t.Errorf("Secret: expected %q, got %q", tt.expected.Secret, cfg.Secret)
			}
		})
	}
}

func Test_getEnvAsBool(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{name: "true string", value: "true", expected: true},
		{name: "false string", value: "false", expected: false},
		{name: "1", value: "1", expected: true},
		{name: "0", value: "0", expected: false},
		{name: "invalid string", value: "yes", expected: false},
		{name: "empty string", value: "", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TEST_BOOL_KEY", tt.value)

			result := getEnvAsBool("TEST_BOOL_KEY")

			if result != tt.expected {
				t.Errorf("getEnvAsBool(%q): expected %v, got %v", tt.value, tt.expected, result)
			}
		})
	}
}
