package checker

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bilal-sayyed/repocheck-go/pkg/report"
)

func CheckCIDevOps(path string) *report.CIDevOps {
	res := &report.CIDevOps{
		Tools: []string{},
	}

	// 1. Docker
	if _, err := os.Stat(filepath.Join(path, "Dockerfile")); err == nil {
		res.HasDockerfile = true
		res.Tools = append(res.Tools, "Docker")
	}

	// 2. CI Detection
	// GitHub Actions
	if _, err := os.Stat(filepath.Join(path, ".github", "workflows")); err == nil {
		res.HasCI = true
		res.Tools = append(res.Tools, "GitHub Actions")
	}
	// GitLab CI
	if _, err := os.Stat(filepath.Join(path, ".gitlab-ci.yml")); err == nil {
		res.HasCI = true
		res.Tools = append(res.Tools, "GitLab CI")
	}
	// CircleCI
	if _, err := os.Stat(filepath.Join(path, ".circleci")); err == nil {
		res.HasCI = true
		res.Tools = append(res.Tools, "CircleCI")
	}

	// 3. Test Detection (Heuristic)
	// Go tests
	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == "vendor" || info.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(info.Name(), "_test.go") {
			res.HasTests = true
			return filepath.SkipAll // Found one, good enough for HasTests
		}
		return nil
	})

	return res
}
