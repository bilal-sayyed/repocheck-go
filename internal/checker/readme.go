package checker

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bilal-sayyed/repocheck-go/pkg/report"
)

func CheckOnboarding(path string) report.OnboardingStatus {
	status := report.OnboardingStatus{
		MissingElements: []string{},
	}

	// Find README
	readmePath := findReadme(path)
	if readmePath == "" {
		status.HasReadme = false
		status.MissingElements = append(status.MissingElements, "README.md")
		return status
	}
	status.HasReadme = true

	// Read Content
	contentBytes, err := os.ReadFile(readmePath)
	if err != nil {
		return status
	}
	content := string(contentBytes)

	// Check for sections
	// Regex for "Installation", "Setup", "Getting Started"
	setupRegex := regexp.MustCompile(`(?i)(install|setup|build|dependencies)`)
	if setupRegex.MatchString(content) {
		status.HasSetupInstructions = true
	} else {
		status.MissingElements = append(status.MissingElements, "Setup Instructions")
	}

	// Regex for "Run", "Usage", "Execute"
	runRegex := regexp.MustCompile(`(?i)(run|start|usage|execute|cmd)`)
	if runRegex.MatchString(content) {
		status.HasRunInstructions = true
	} else {
		status.MissingElements = append(status.MissingElements, "Run Instructions")
	}

	return status
}

func findReadme(root string) string {
	// Look for README.md, readme.txt, etc.
	matches, _ := filepath.Glob(filepath.Join(root, "README*"))
	if len(matches) > 0 {
		return matches[0]
	}
	// Case insensitive fallback (Glob is case sensitive on some OSs)
	entries, _ := os.ReadDir(root)
	for _, e := range entries {
		if strings.EqualFold(e.Name(), "readme.md") || strings.EqualFold(e.Name(), "readme.txt") {
			return filepath.Join(root, e.Name())
		}
	}
	return ""
}
