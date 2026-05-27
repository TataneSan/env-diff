package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

const (
	ansiReset  = "\033[0m"
	ansiRed    = "\033[31m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiCyan   = "\033[36m"
	ansiGray   = "\033[90m"
	ansiBold   = "\033[1m"
)

// EnvDiff holds the comparison result between two env files.
type EnvDiff struct {
	OnlyInA     map[string]string `json:"only_in_first"`
	OnlyInB     map[string]string `json:"only_in_second"`
	Changed     map[string]Pair   `json:"changed"`
	Unchanged   map[string]string `json:"unchanged"`
}

// Pair holds old and new values for a changed key.
type Pair struct {
	Old string `json:"old"`
	New string `json:"new"`
}

var (
	rootCmd = &cobra.Command{
		Use:   "env-diff",
		Short: "Compare two .env files and show differences",
		Long: `env-diff compares two .env files and displays the differences:
  - Keys only in the first file (removed)
  - Keys only in the second file (added)
  - Keys with different values (changed)
  - Keys with identical values (unchanged)

Supports output in table (default) or JSON format.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("two file paths are required\n\nUsage: env-diff <file1> <file2>\n")
			}
			fileA, fileB := args[0], args[1]

			mapA, err := parseEnv(fileA)
			if err != nil {
				return fmt.Errorf("error reading %s: %w", fileA, err)
			}
			mapB, err := parseEnv(fileB)
			if err != nil {
				return fmt.Errorf("error reading %s: %w", fileB, err)
			}

			diff := computeDiff(mapA, mapB)

			switch outputFormat {
			case "json":
				return printJSON(diff, fileA, fileB)
			default:
				printTable(diff, fileA, fileB)
			}
			return nil
		},
	}

	outputFormat string
	quietMode    bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "table", "Output format: table or json")
	rootCmd.PersistentFlags().BoolVarP(&quietMode, "quiet", "q", false, "Only show differences (hide unchanged keys)")
}

// parseEnv reads a .env file and returns a map of key=value pairs.
// It skips blank lines and lines starting with #.
// It handles optional quotes around values.
func parseEnv(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		if key == "" {
			continue
		}

		value := strings.TrimSpace(parts[1])

		// Strip surrounding quotes
		value = stripQuotes(value)

		result[key] = value
	}

	return result, scanner.Err()
}

// stripQuotes removes leading/trailing single or double quotes from a value.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// computeDiff compares two env maps and categorises keys.
func computeDiff(a, b map[string]string) *EnvDiff {
	diff := &EnvDiff{
		OnlyInA:     make(map[string]string),
		OnlyInB:     make(map[string]string),
		Changed:     make(map[string]Pair),
		Unchanged:   make(map[string]string),
	}

	// Keys in A
	for k, v := range a {
		if vb, ok := b[k]; ok {
			if v == vb {
				diff.Unchanged[k] = v
			} else {
				diff.Changed[k] = Pair{Old: v, New: vb}
			}
		} else {
			diff.OnlyInA[k] = v
		}
	}

	// Keys only in B
	for k, v := range b {
		if _, ok := a[k]; !ok {
			diff.OnlyInB[k] = v
		}
	}

	return diff
}

// printTable outputs a colored table of differences.
func printTable(diff *EnvDiff, fileA, fileB string) {
	fmt.Println()
	fmt.Printf("  Comparing: %s  →  %s\n", bold(fileA), bold(fileB))
	fmt.Println(strings.Repeat("─", 56))
	fmt.Printf("  %s keys in %s (removed)\n", countLabel(len(diff.OnlyInA)), bold(fileA))
	fmt.Printf("  %s keys in %s (added)\n", countLabel(len(diff.OnlyInB)), bold(fileB))
	fmt.Printf("  %s keys changed\n", countLabel(len(diff.Changed)))
	fmt.Printf("  %s keys unchanged\n", countLabel(len(diff.Unchanged)))
	fmt.Println(strings.Repeat("─", 56))
	fmt.Println()

	if !quietMode && len(diff.OnlyInA) > 0 {
		fmt.Println(colorize("  REMOVED (only in first file)", ansiRed))
		for _, k := range sortedKeys(diff.OnlyInA) {
			fmt.Printf("    %s-%s %s = %s\n", ansiRed, ansiBold, colorize(k, ansiCyan), colorize(maskValue(diff.OnlyInA[k]), ansiGray))
		}
		fmt.Println()
	}

	if len(diff.OnlyInB) > 0 {
		fmt.Println(colorize("  ADDED (only in second file)", ansiGreen))
		for _, k := range sortedKeys(diff.OnlyInB) {
			fmt.Printf("    %s+%s %s = %s\n", ansiGreen, ansiBold, colorize(k, ansiCyan), colorize(maskValue(diff.OnlyInB[k]), ansiGray))
		}
		fmt.Println()
	}

	if len(diff.Changed) > 0 {
		fmt.Println(colorize("  CHANGED", ansiYellow))
		for _, k := range sortedKeys(diff.Changed) {
			p := diff.Changed[k]
			fmt.Printf("    %s~%s %s\n", ansiYellow, ansiBold, colorize(k, ansiCyan))
			fmt.Printf("      - %s\n", colorize(maskValue(p.Old), ansiRed))
			fmt.Printf("      + %s\n", colorize(maskValue(p.New), ansiGreen))
		}
		fmt.Println()
	}

	if !quietMode && len(diff.Unchanged) > 0 {
		fmt.Println(colorize(fmt.Sprintf("  UNCHANGED (%d keys)", len(diff.Unchanged)), ansiGray))
		for _, k := range sortedKeys(diff.Unchanged) {
			fmt.Printf("    %s %s = %s\n", ansiGray, colorize(k, ansiCyan), colorize(maskValue(diff.Unchanged[k]), ansiGray))
		}
		fmt.Println()
	}

	if len(diff.OnlyInA) == 0 && len(diff.OnlyInB) == 0 && len(diff.Changed) == 0 {
		fmt.Println("  Files are identical.")
		fmt.Println()
	}
}

// printJSON outputs the diff as JSON.
func printJSON(diff *EnvDiff, fileA, fileB string) error {
	output := map[string]interface{}{
		"first_file":  fileA,
		"second_file": fileB,
		"summary": map[string]int{
			"removed":   len(diff.OnlyInA),
			"added":     len(diff.OnlyInB),
			"changed":   len(diff.Changed),
			"unchanged": len(diff.Unchanged),
		},
	}

	if !quietMode {
		output["only_in_first"]  = diff.OnlyInA
		output["only_in_second"] = diff.OnlyInB
		output["unchanged"]      = diff.Unchanged
	}
	output["changed"] = diff.Changed

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}

func sortedKeys(m interface{}) []string {
	// Helper: extract keys from any map type we use.
	var keys []string
	switch v := m.(type) {
	case map[string]string:
		keys = make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
	case map[string]Pair:
		keys = make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}

func countLabel(n int) string {
	return fmt.Sprintf("%d", n)
}

func maskValue(v string) string {
	if v == "" {
		return ""
	}
	// Mask long values to avoid cluttering the output.
	if len(v) > 40 {
		return v[:37] + "..."
	}
	return v
}

func supportsColor() bool {
	return os.Getenv("NO_COLOR") == "" && os.Getenv("TERM") != "dumb"
}

func colorize(s string, color string) string {
	if !supportsColor() {
		return s
	}
	return color + s + ansiReset
}

func bold(s string) string {
	if !supportsColor() {
		return s
	}
	return ansiBold + s + ansiReset
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
