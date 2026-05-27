package presentation

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/TataneSan/env-diff/internal/domain"
)

// Printer handles all CLI output.
type Printer struct {
	useColor bool
}

// NewPrinter creates a new Printer.
func NewPrinter() *Printer {
	return &Printer{useColor: supportsColor()}
}

// PrintDiff displays the differences between two environments.
func (p *Printer) PrintDiff(diffs []domain.EnvDiff, nameA, nameB string) {
	if len(diffs) == 0 {
		fmt.Println("No differences found.")
		return
	}

	fmt.Printf("Comparing %s vs %s\n", nameA, nameB)
	fmt.Printf("Found %d difference(s)\n\n", len(diffs))

	// Sort diffs by type then key
	sort.Slice(diffs, func(i, j int) bool {
		if diffs[i].DiffType != diffs[j].DiffType {
			return diffs[i].DiffType < diffs[j].DiffType
		}
		return diffs[i].Key < diffs[j].Key
	})

	for _, diff := range diffs {
		p.printDiff(diff)
	}

	p.printSummary(diffs)
}

// PrintStats displays summary statistics.
func (p *Printer) PrintStats(diffs []domain.EnvDiff) {
	added := 0
	removed := 0
	modified := 0
	for _, d := range diffs {
		switch d.DiffType {
		case domain.Added:
			added++
		case domain.Removed:
			removed++
		case domain.Modified:
			modified++
		}
	}
	fmt.Printf("Summary: %d added, %d removed, %d modified\n", added, removed, modified)
}

// PrintError displays an error message.
func (p *Printer) PrintError(msg string) {
	fmt.Fprintf(os.Stderr, "%sError: %s%s\n", red(), msg, reset())
}

// PrintUsage displays the help text.
func (p *Printer) PrintUsage() {
	fmt.Println("env-diff - Compare environment variables between two contexts")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  env-diff file <fileA> <fileB>     Compare two .env files")
	fmt.Println("  env-diff current <file>            Compare current env with a .env file")
	fmt.Println("  env-diff cmd <cmdA> <cmdB>         Compare output of two commands")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  env-diff file .env .env.production")
	fmt.Println("  env-diff current .env")
	fmt.Println("  env-diff cmd 'docker inspect --format=\"{{json .Config.Env}}\" mycontainer' 'cat .env | xargs -I{} echo {}'")
}

func (p *Printer) printDiff(diff domain.EnvDiff) {
	switch diff.DiffType {
	case domain.Added:
		fmt.Printf("%s+%s %-30s = %s\n", green(), reset(), diff.Key, diff.ValueB)
	case domain.Removed:
		fmt.Printf("%s-%s %-30s = %s\n", red(), reset(), diff.Key, diff.ValueA)
	case domain.Modified:
		fmt.Printf("%s~%s %-30s\n", yellow(), reset(), diff.Key)
		fmt.Printf("    %s-%s %s\n", red(), reset(), diff.ValueA)
		fmt.Printf("    %s+%s %s\n", green(), reset(), diff.ValueB)
	}
}

func (p *Printer) printSummary(diffs []domain.EnvDiff) {
	added := 0
	removed := 0
	modified := 0
	for _, d := range diffs {
		switch d.DiffType {
		case domain.Added:
			added++
		case domain.Removed:
			removed++
		case domain.Modified:
			modified++
		}
	}
	fmt.Printf("\nTotal: %d added, %d removed, %d modified\n", added, removed, modified)
}

func supportsColor() bool {
	return os.Getenv("TERM") != "dumb" && os.Getenv("NO_COLOR") == ""
}

func red() string   { return "\033[31m" }
func green() string { return "\033[32m" }
func yellow() string { return "\033[33m" }
func reset() string { return "\033[0m" }

// MaskValue replaces all but the last 4 characters with asterisks.
func MaskValue(value string) string {
	if len(value) <= 4 {
		return strings.Repeat("*", len(value))
	}
	return strings.Repeat("*", len(value)-4) + value[len(value)-4:]
}
