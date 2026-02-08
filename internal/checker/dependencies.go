package checker

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/bilal-sayyed/repocheck-go/pkg/report"
)

func CheckDependencies(path string) report.DependencyStats {
	stats := report.DependencyStats{
		Manager: "Unknown",
	}

	// Check for go.mod
	if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
		stats.Manager = "go.mod"
		stats.DirectCount, stats.TotalCount = countGoDependencies(path)
		return stats
	}

	// Check for package.json
	if _, err := os.Stat(filepath.Join(path, "package.json")); err == nil {
		stats.Manager = "npm/yarn"
		stats.DirectCount, stats.TotalCount = countNpmDependencies(path) // Placeholder implementation
		return stats
	}

	// Check for requirements.txt
	if _, err := os.Stat(filepath.Join(path, "requirements.txt")); err == nil {
		stats.Manager = "pip"
		stats.DirectCount = countLines(filepath.Join(path, "requirements.txt"))
		stats.TotalCount = stats.DirectCount
		return stats
	}

	return stats
}

func countGoDependencies(root string) (int, int) {
	// Simple parsing of go.mod
	f, err := os.Open(filepath.Join(root, "go.mod"))
	if err != nil {
		return 0, 0
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	direct := 0
	total := 0
	inRequire := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "require (") {
			inRequire = true
			continue
		}
		if line == ")" {
			inRequire = false
			continue
		}

		if strings.HasPrefix(line, "require ") || inRequire {
			// Basic line check: has at least 2 fields (module version)
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				total++
				if !strings.Contains(line, "// indirect") {
					direct++
				}
			}
		}
	}
	return direct, total
}

func countNpmDependencies(root string) (int, int) {
	// Parsing package.json is better with a struct, but for now simple line count or 0
	// TODO: Implement real JSON parsing if needed for accuracy
	return 0, 0
}

func countLines(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) != "" {
			count++
		}
	}
	return count
}
