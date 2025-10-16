package bootstrap

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	stateDir  = ".comet"
	stateFile = "bootstrap.state"
)

// State represents the bootstrap state
type State struct {
	Version string               `json:"version"`
	LastRun time.Time            `json:"last_run"`
	Steps   map[string]*StepState `json:"steps"`
}

// StepState represents the state of a single bootstrap step
type StepState struct {
	Status      string    `json:"status"`       // completed, failed, skipped
	CompletedAt time.Time `json:"completed_at"`
	Target      string    `json:"target"`
	Error       string    `json:"error,omitempty"`
	LastAttempt time.Time `json:"last_attempt,omitempty"`
}

// LoadState loads the bootstrap state from disk
func LoadState() (*State, error) {
	statePath := getStatePath()

	// If state file doesn't exist, return empty state
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		return &State{
			Version: "1",
			Steps:   make(map[string]*StepState),
		}, nil
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	if state.Steps == nil {
		state.Steps = make(map[string]*StepState)
	}

	return &state, nil
}

// SaveState saves the bootstrap state to disk
func (s *State) Save() error {
	statePath := getStatePath()
	stateDir := filepath.Dir(statePath)

	// Create .comet directory if it doesn't exist
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	s.LastRun = time.Now()

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(statePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// UpdateStep updates the state of a bootstrap step
func (s *State) UpdateStep(name string, stepState *StepState) {
	if s.Steps == nil {
		s.Steps = make(map[string]*StepState)
	}
	s.Steps[name] = stepState
}

// GetStep retrieves the state of a bootstrap step
func (s *State) GetStep(name string) *StepState {
	return s.Steps[name]
}

// Clear removes the bootstrap state file
func Clear() error {
	statePath := getStatePath()
	if err := os.Remove(statePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove state file: %w", err)
	}
	return nil
}

// getStatePath returns the full path to the state file
func getStatePath() string {
	return filepath.Join(stateDir, stateFile)
}
