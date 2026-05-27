// env-diff compares two .env files or environments and shows the differences.
//
// Usage:
//
//	env-diff <file1> <file2>
//	env-diff --env1 <env1> --env2 <env2>
//
// Examples:
//
//	env-diff .env.production .env.staging
//	env-diff .env.example .env
//	env-diff --env1 PROD --env2 STAGING
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type EnvPair struct {
	Key   string
	Value string
}

type DiffResult struct {
	Key      string
	Status   string // "added", "removed", "modified", "unchanged"
	LeftVal  string
	RightVal string
}

func parseEnvFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	env := make(map[string]string)
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE
		eqIdx := strings.Index(line, "=")
		if eqIdx == -1 {
			continue
		}

		key := strings.TrimSpace(line[:eqIdx])
		value := strings.TrimSpace(line[eqIdx+1:])

		// Remove surrounding quotes
		value = stripQuotes(value)

		if key != "" {
			env[key] = value
		}
	}

	return env, scanner.Err()
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func compareEnvs(left, right map[string]string) []DiffResult {
	allKeys := make(map[string]bool)
	for k := range left {
		allKeys[k] = true
	}
	for k := range right {
		allKeys[k] = true
	}

	var results []DiffResult
	for k := range allKeys {
		lVal, lOk := left[k]
		rVal, rOk := right[k]

		var status string
		if lOk && !rOk {
			status = "removed"
		} else if !lOk && rOk {
			status = "added"
		} else if lVal != rVal {
			status = "modified"
		} else {
			status = "unchanged"
		}

		results = append(results, DiffResult{
			Key:      k,
			Status:   status,
			LeftVal:  lVal,
			RightVal: rVal,
		})
	}

	// Sort by key
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	return results
}

func printReport(results []DiffResult, leftName, rightName string) {
	// Count by status
	summary := map[string]int{}
	for _, r := range results {
		summary[r.Status]++
	}

	fmt.Println()
	fmt.Printf("  Comparing: %s vs %s\n", leftName, rightName)
	fmt.Printf("  Total keys: %d | Added: %d | Removed: %d | Modified: %d | Unchanged: %d\n",
		len(results), summary["added"], summary["removed"], summary["modified"], summary["unchanged"])
	fmt.Println()

	// Print differences only by default
	hasDiffs := false
	for _, r := range results {
		if r.Status == "unchanged" {
			continue
		}
		hasDiffs = true

		switch r.Status {
		case "added":
			fmt.Printf("  + %-30s %s\n", r.Key, r.RightVal)
		case "removed":
			fmt.Printf("  - %-30s %s\n", r.Key, r.LeftVal)
		case "modified":
			fmt.Printf("  ~ %-30s\n", r.Key)
			fmt.Printf("    %s: %s\n", leftName, r.LeftVal)
			fmt.Printf("    %s: %s\n", rightName, r.RightVal)
		}
	}

	if !hasDiffs {
		fmt.Println("  No differences found.")
	}

	fmt.Println()
}

func printHelp() {
	fmt.Println("env-diff — Compare two .env files or environments")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  env-diff <file1> <file2>")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  file1    First .env file")
	fmt.Println("  file2    Second .env file")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  env-diff .env.production .env.staging")
	fmt.Println("  env-diff .env.example .env")
	fmt.Println()
	fmt.Println("Output:")
	fmt.Println("  - Summary of differences (added, removed, modified, unchanged)")
	fmt.Println("  - Detailed view of each difference")
}

func main() {
	if len(os.Args) < 3 {
		if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
			printHelp()
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Usage: env-diff <file1> <file2>\n")
		fmt.Fprintf(os.Stderr, "Try 'env-diff --help' for more information.\n")
		os.Exit(1)
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		printHelp()
		os.Exit(0)
	}

	file1 := os.Args[1]
	file2 := os.Args[2]

	// Check files exist
	if _, err := os.Stat(file1); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: file not found: %s\n", file1)
		os.Exit(1)
	}
	if _, err := os.Stat(file2); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: file not found: %s\n", file2)
		os.Exit(1)
	}

	left, err := parseEnvFile(file1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", file1, err)
		os.Exit(1)
	}

	right, err := parseEnvFile(file2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", file2, err)
		os.Exit(1)
	}

	results := compareEnvs(left, right)
	printReport(results, file1, file2)
}
