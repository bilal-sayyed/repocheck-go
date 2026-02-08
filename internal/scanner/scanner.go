package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bilal-sayyed/repocheck-go/pkg/report"
)

// ScanRepository performs a full scan of the directory
func ScanRepository(path string) (*report.RepoScanResult, error) {
	result := &report.RepoScanResult{}

	// 1. Basic Traversal & Summary
	summary, fileList, err := analyzeStructure(path)
	if err != nil {
		return nil, err
	}
	result.Summary = summary

	// 2. Identify Language explicitly
	result.Summary.Language = detectLanguage(fileList)

	// 3. Detect Git
	if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
		result.Summary.GitPresent = true
	}

	return result, nil
}

func analyzeStructure(root string) (report.RepoSummary, []string, error) {
	summary := report.RepoSummary{
		FileCount: 0,
		Type:      "Unknown",
	}
	var fileList []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == "node_modules" || info.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		summary.FileCount++
		fileList = append(fileList, info.Name())

		// Heuristic for Repo Type
		name := strings.ToLower(info.Name())
		if name == "main.go" || name == "index.js" || name == "cmd" {
			summary.EntryPoints = append(summary.EntryPoints, info.Name())
		}

		return nil
	})

	// Heuristics for Type
	if hasFile(fileList, "main.go") && strings.Contains(root, "cmd") {
		summary.Type = "CLI/Backend"
	} else if hasFile(fileList, "package.json") && hasFile(fileList, "index.html") {
		summary.Type = "Frontend"
	} else if hasFile(fileList, "go.mod") {
		summary.Type = "Go Library/Service"
	}

	return summary, fileList, err
}

func detectLanguage(files []string) string {
	counts := make(map[string]int)
	for _, f := range files {
		ext := filepath.Ext(f)
		switch ext {
		case ".go":
			counts["Go"]++
		case ".js", ".ts", ".jsx", ".tsx":
			counts["JavaScript/TypeScript"]++
		case ".py":
			counts["Python"]++
		case ".java":
			counts["Java"]++
		case ".rs":
			counts["Rust"]++
		}
	}

	// Simple majority wins
	maxCount := 0
	lang := "Unknown"
	for l, c := range counts {
		if c > maxCount {
			maxCount = c
			lang = l
		}
	}
	return lang
}

func hasFile(files []string, name string) bool {
	for _, f := range files {
		if strings.EqualFold(f, name) {
			return true
		}
	}
	return false
}
