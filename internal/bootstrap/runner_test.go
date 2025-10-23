package bootstrap

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"filippo.io/age"
)

func TestExpandPath(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		name     string
		input    string
		xdgHome  string
		expected string
	}{
		{
			name:     "regular tilde expansion",
			input:    "~/test/file.txt",
			xdgHome:  "",
			expected: filepath.Join(home, "test", "file.txt"),
		},
		{
			name:     "sops age key on macOS without XDG_CONFIG_HOME",
			input:    "~/.config/sops/age/keys.txt",
			xdgHome:  "",
			expected: filepath.Join(home, "Library", "Application Support", "sops", "age", "keys.txt"),
		},
		{
			name:     "sops age key with XDG_CONFIG_HOME set",
			input:    "~/.config/sops/age/keys.txt",
			xdgHome:  filepath.Join(home, ".config"),
			expected: filepath.Join(home, ".config", "sops", "age", "keys.txt"),
		},
		{
			name:     "non-sops path is not affected",
			input:    "~/.config/other/file.txt",
			xdgHome:  "",
			expected: filepath.Join(home, ".config", "other", "file.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up XDG_CONFIG_HOME environment variable
			if tt.xdgHome != "" {
				os.Setenv("XDG_CONFIG_HOME", tt.xdgHome)
				defer os.Unsetenv("XDG_CONFIG_HOME")
			} else {
				os.Unsetenv("XDG_CONFIG_HOME")
			}

			result := expandPath(tt.input)

			// Skip macOS-specific test on non-macOS systems
			if runtime.GOOS != "darwin" && tt.name == "sops age key on macOS without XDG_CONFIG_HOME" {
				// On Linux, expect ~/.config path
				expectedLinux := filepath.Join(home, ".config", "sops", "age", "keys.txt")
				if result != expectedLinux {
					t.Errorf("expandPath() = %v, want %v (Linux)", result, expectedLinux)
				}
				return
			}

			if result != tt.expected {
				t.Errorf("expandPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestResolveSopsAgePath(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		name     string
		input    string
		xdgHome  string
		expected string
	}{
		{
			name:     "returns empty for non-sops path",
			input:    "~/.config/other/file.txt",
			xdgHome:  "",
			expected: "",
		},
		{
			name:     "resolves with XDG_CONFIG_HOME set",
			input:    "sops/age/keys.txt",
			xdgHome:  filepath.Join(home, "custom-config"),
			expected: filepath.Join(home, "custom-config", "sops", "age", "keys.txt"),
		},
		{
			name:     "resolves macOS default without XDG_CONFIG_HOME",
			input:    "sops/age/keys.txt",
			xdgHome:  "",
			expected: getPlatformSopsPath(home),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up XDG_CONFIG_HOME environment variable
			if tt.xdgHome != "" {
				os.Setenv("XDG_CONFIG_HOME", tt.xdgHome)
				defer os.Unsetenv("XDG_CONFIG_HOME")
			} else {
				os.Unsetenv("XDG_CONFIG_HOME")
			}

			result := resolveSopsAgePath(tt.input)
			if result != tt.expected {
				t.Errorf("resolveSopsAgePath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// getPlatformSopsPath returns the expected SOPS age key path for the current platform
func getPlatformSopsPath(home string) string {
	if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Library", "Application Support", "sops", "age", "keys.txt")
	}
	return filepath.Join(home, ".config", "sops", "age", "keys.txt")
}

func TestFormatAgeKey(t *testing.T) {
	// Generate a real test age key
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("Failed to generate test age key: %v", err)
	}
	testSecretKey := identity.String()

	formatted, err := formatAgeKey(testSecretKey)
	if err != nil {
		t.Fatalf("formatAgeKey() error = %v", err)
	}

	// Should contain the public key comment
	if !strings.Contains(formatted, "# public key: age1") {
		t.Errorf("formatAgeKey() should contain public key comment, got: %v", formatted)
	}

	// Should contain the original secret key
	if !strings.Contains(formatted, testSecretKey) {
		t.Errorf("formatAgeKey() should contain secret key, got: %v", formatted)
	}

	// Should have public key before secret key
	lines := strings.Split(formatted, "\n")
	if len(lines) < 2 {
		t.Errorf("formatAgeKey() should have at least 2 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "# public key:") {
		t.Errorf("formatAgeKey() first line should be public key comment, got: %v", lines[0])
	}
	if !strings.Contains(lines[1], "AGE-SECRET-KEY-") {
		t.Errorf("formatAgeKey() second line should be secret key, got: %v", lines[1])
	}
}

func TestExtractSecretKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "extracts from formatted content",
			input:    "# public key: age1abc123\nAGE-SECRET-KEY-1ABC123DEF456\n",
			expected: "AGE-SECRET-KEY-1ABC123DEF456",
		},
		{
			name:     "extracts from plain key",
			input:    "AGE-SECRET-KEY-1XYZ789\n",
			expected: "AGE-SECRET-KEY-1XYZ789",
		},
		{
			name:     "returns empty for no match",
			input:    "some random text\n",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSecretKey(tt.input)
			if result != tt.expected {
				t.Errorf("extractSecretKey() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShouldAppendAgeKey(t *testing.T) {
	// Generate test keys
	identity1, _ := age.GenerateX25519Identity()
	identity2, _ := age.GenerateX25519Identity()

	formatted1, _ := formatAgeKey(identity1.String())
	formatted2, _ := formatAgeKey(identity2.String())

	// Create temp directory for test files
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-keys.txt")

	t.Run("should append to non-existent file", func(t *testing.T) {
		shouldAppend, err := shouldAppendAgeKey(testFile, formatted1)
		if err != nil {
			t.Fatalf("shouldAppendAgeKey() error = %v", err)
		}
		if !shouldAppend {
			t.Error("shouldAppendAgeKey() should return true for non-existent file")
		}
	})

	t.Run("should append first key to empty file", func(t *testing.T) {
		// Create empty file
		os.WriteFile(testFile, []byte(""), 0600)

		shouldAppend, err := shouldAppendAgeKey(testFile, formatted1)
		if err != nil {
			t.Fatalf("shouldAppendAgeKey() error = %v", err)
		}
		if !shouldAppend {
			t.Error("shouldAppendAgeKey() should return true for empty file")
		}
	})

	t.Run("should not append duplicate key", func(t *testing.T) {
		// Write first key
		os.WriteFile(testFile, []byte(formatted1), 0600)

		shouldAppend, err := shouldAppendAgeKey(testFile, formatted1)
		if err != nil {
			t.Fatalf("shouldAppendAgeKey() error = %v", err)
		}
		if shouldAppend {
			t.Error("shouldAppendAgeKey() should return false for duplicate key")
		}
	})

	t.Run("should append different key", func(t *testing.T) {
		// File already has key1
		os.WriteFile(testFile, []byte(formatted1), 0600)

		shouldAppend, err := shouldAppendAgeKey(testFile, formatted2)
		if err != nil {
			t.Fatalf("shouldAppendAgeKey() error = %v", err)
		}
		if !shouldAppend {
			t.Error("shouldAppendAgeKey() should return true for different key")
		}
	})

	t.Run("should detect duplicate in multi-key file", func(t *testing.T) {
		// Write both keys
		content := formatted1 + "\n" + formatted2
		os.WriteFile(testFile, []byte(content), 0600)

		// Try to add key1 again
		shouldAppend, err := shouldAppendAgeKey(testFile, formatted1)
		if err != nil {
			t.Fatalf("shouldAppendAgeKey() error = %v", err)
		}
		if shouldAppend {
			t.Error("shouldAppendAgeKey() should return false when key exists in multi-key file")
		}
	})
}

// Helper functions are not needed - using strings package directly
