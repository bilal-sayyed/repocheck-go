package checker

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/bilal-sayyed/repocheck-go/pkg/report"
)

func CheckOwnership(path string) *report.Ownership {
	res := &report.Ownership{}

	// 1. CODEOWNERS
	// Locations: root, .github/, docs/
	locations := []string{"CODEOWNERS", ".github/CODEOWNERS", "docs/CODEOWNERS"}
	for _, loc := range locations {
		if _, err := os.Stat(filepath.Join(path, loc)); err == nil {
			res.HasCodeowners = true
			break
		}
	}

	// 2. Maintainers (package.json)
	packageJsonPath := filepath.Join(path, "package.json")
	if _, err := os.Stat(packageJsonPath); err == nil {
		if hasMaintainersKey(packageJsonPath) {
			res.HasMaintainers = true
		}
	}

	return res
}

func hasMaintainersKey(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return false
	}

	// Check "maintainers" (array) or "author" (string/object) - treating author as maintainer mostly
	if _, ok := data["maintainers"]; ok {
		return true
	}
	// Strict "maintainers" check based on spec, but author is good too.
	// Let's stick to strict "maintainers" or "contributors" for now.
	return false
}
