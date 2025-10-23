package bootstrap

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
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
