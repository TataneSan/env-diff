package infrastructure

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/TataneSan/env-diff/internal/domain"
)

// EnvCollector collects environment variables from various sources.
type EnvCollector struct{}

// NewEnvCollector creates a new EnvCollector.
func NewEnvCollector() *EnvCollector {
	return &EnvCollector{}
}

// CollectCurrent returns the current process environment.
func (c *EnvCollector) CollectCurrent() *domain.EnvContext {
	vars := make(map[string]string)
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	return &domain.EnvContext{Name: "current", Vars: vars}
}

// CollectFromFile reads environment variables from a .env file.
func (c *EnvCollector) CollectFromFile(path string) (*domain.EnvContext, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	vars := make(map[string]string)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove surrounding quotes
			value = strings.Trim(value, "\"'")
			vars[key] = value
		}
	}

	return &domain.EnvContext{Name: filepath.Base(path), Vars: vars}, nil
}

// CollectFromCommand runs a command and captures its environment output.
func (c *EnvCollector) CollectFromCommand(cmd string) (*domain.EnvContext, error) {
	output, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return nil, err
	}

	vars := make(map[string]string)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}

	return &domain.EnvContext{Name: cmd, Vars: vars}, nil
}

// EnvDiffer compares two environment contexts.
type EnvDiffer struct{}

// NewEnvDiffer creates a new EnvDiffer.
func NewEnvDiffer() *EnvDiffer {
	return &EnvDiffer{}
}

// Compare returns the differences between two environments.
func (d *EnvDiffer) Compare(a, b *domain.EnvContext) []domain.EnvDiff {
	var diffs []domain.EnvDiff

	// Check all keys in A
	for key, valA := range a.Vars {
		if valB, ok := b.Vars[key]; ok {
			if valA != valB {
				diffs = append(diffs, domain.EnvDiff{
					Key:      key,
					DiffType: domain.Modified,
					ValueA:   valA,
					ValueB:   valB,
				})
			}
		} else {
			diffs = append(diffs, domain.EnvDiff{
				Key:      key,
				DiffType: domain.Removed,
				ValueA:   valA,
			})
		}
	}

	// Check keys only in B
	for key, valB := range b.Vars {
		if _, ok := a.Vars[key]; !ok {
			diffs = append(diffs, domain.EnvDiff{
				Key:      key,
				DiffType: domain.Added,
				ValueB:   valB,
			})
		}
	}

	return diffs
}
