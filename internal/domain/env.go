package domain

// DiffType represents the result of comparing two environment values.
type DiffType int

const (
	Unchanged DiffType = iota
	Added
	Removed
	Modified
)

// EnvDiff represents a single difference between two environments.
type EnvDiff struct {
	Key       string
	DiffType  DiffType
	ValueA    string
	ValueB    string
}

// EnvContext represents a named set of environment variables.
type EnvContext struct {
	Name  string
	Vars  map[string]string
}
