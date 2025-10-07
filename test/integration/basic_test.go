package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestBasicStackList tests that comet can list stacks
func TestBasicStackList(t *testing.T) {
	// Build comet binary
	buildCmd := exec.Command("go", "build", "-o", "comet-test", ".")
	buildCmd.Dir = filepath.Join("..", "..")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build comet: %v", err)
	}
	defer os.Remove(filepath.Join("..", "..", "comet-test"))

	// Run comet list
	cmd := exec.Command(filepath.Join("..", "..", "comet-test"), "list")
	cmd.Dir = filepath.Join("..", "..")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("comet list failed: %v\nOutput: %s", err, output)
	}

	// Check output contains "dev"
	if !strings.Contains(string(output), "dev") {
		t.Errorf("Expected output to contain 'dev', got: %s", output)
	}
}

// TestStackParse tests that all example stacks parse correctly
func TestStackParse(t *testing.T) {
	stacks := []string{"dev"}

	buildCmd := exec.Command("go", "build", "-o", "comet-test", ".")
	buildCmd.Dir = filepath.Join("..", "..")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build comet: %v", err)
	}
	defer os.Remove(filepath.Join("..", "..", "comet-test"))

	for _, stack := range stacks {
		t.Run(stack, func(t *testing.T) {
			cmd := exec.Command(filepath.Join("..", "..", "comet-test"), "list", stack)
			cmd.Dir = filepath.Join("..", "..")
			output, err := cmd.CombinedOutput()

			if err != nil {
				t.Errorf("Failed to list stack %s: %v\nOutput: %s", stack, err, output)
			}
		})
	}
}

// TestVersionCommand tests that comet version command works
func TestVersionCommand(t *testing.T) {
	// Build comet binary
	buildCmd := exec.Command("go", "build", "-o", "comet-test", ".")
	buildCmd.Dir = filepath.Join("..", "..")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build comet: %v", err)
	}
	defer os.Remove(filepath.Join("..", "..", "comet-test"))

	// Run comet version
	cmd := exec.Command(filepath.Join("..", "..", "comet-test"), "version")
	cmd.Dir = filepath.Join("..", "..")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("comet version failed: %v\nOutput: %s", err, output)
	}

	// Check output is not empty
	if len(output) == 0 {
		t.Error("Expected version output, got empty string")
	}
}
