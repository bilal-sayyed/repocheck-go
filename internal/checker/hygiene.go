package checker

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bilal-sayyed/repocheck-go/pkg/report"
)

func CheckHygiene(path string) report.HygieneChecks {
	hygiene := report.HygieneChecks{
		HardcodedSecrets: []string{},
	}

	// 1. Check for .env.example
	if _, err := os.Stat(filepath.Join(path, ".env.example")); err == nil {
		hygiene.HasEnvExample = true
	}

	// 2. Config files
	configs := []string{".env", "config.yaml", "config.json", "settings.py", "repocheck.yaml"}
	for _, conf := range configs {
		if _, err := os.Stat(filepath.Join(path, conf)); err == nil {
			hygiene.ConfigPresent = true
			break // Found at least one config file
		}
	}

	// 3. Gitignore Check
	gitignorePath := filepath.Join(path, ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil {
		hygiene.HasGitignore = true

		// Check if .env is ignored
		content, _ := os.ReadFile(gitignorePath)
		sContent := string(content)
		if strings.Contains(sContent, ".env") {
			hygiene.EnvIgnored = true
		}
	}

	// 4. Scan for Hardcoded Secrets
	// Very basic regex for demonstration.
	// Warning: This can match false positives.
	secretRegex := regexp.MustCompile(`(?i)(api_key|password|secret|token)\s*=\s*['"][a-zA-Z0-9_\-]{8,}['"]`)

	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == "node_modules" || info.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		// Only check small text files
		if info.Size() > 1024*1024 { // Skip files > 1MB
			return nil
		}

		f, err := os.Open(filePath)
		if err != nil {
			return nil
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if len(line) > 1000 {
				continue
			} // Skip very long lines

			if secretRegex.MatchString(line) {
				relPath, _ := filepath.Rel(path, filePath)
				hygiene.HardcodedSecrets = append(hygiene.HardcodedSecrets, relPath)
				break // Report file once
			}
		}
		return nil
	})

	return hygiene
}
