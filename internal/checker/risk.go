package checker

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bilal-sayyed/repocheck-go/pkg/report"
)

func CheckLicenseRisk(path string) *report.LicenseRisk {
	risk := &report.LicenseRisk{
		RiskLevel:       "LOW",
		RiskyLicenses:   []string{},
		OldDependencies: []string{},
		Reasons:         []string{},
	}

	// 1. Check Root License
	rootLicenseBytes, err := os.ReadFile(filepath.Join(path, "LICENSE"))
	if err == nil {
		content := string(rootLicenseBytes)
		if strings.Contains(content, "GPL") && !strings.Contains(content, "Lesser General Public License") {
			// If key is AGPL or GPL, it might be risk depending on usage, but usually safe for the Repo itself.
			// However, checking dependencies is the real deal.
		}
	} else {
		risk.RiskyLicenses = append(risk.RiskyLicenses, "Missing Root LICENSE file")
		risk.RiskLevel = "MEDIUM"
		risk.Reasons = append(risk.Reasons, "Missing LICENSE file")
	}

	// 2. Check Dependency Abandonment (Go Version)
	// Just checking go.mod version as a proxy for "Modernity"
	goModPath := filepath.Join(path, "go.mod")
	if data, err := os.ReadFile(goModPath); err == nil {
		re := regexp.MustCompile(`go\s+(\d+\.\d+)`)
		matches := re.FindStringSubmatch(string(data))
		if len(matches) > 1 {
			verStr := matches[1]
			ver, _ := strconv.ParseFloat(verStr, 64)
			if ver < 1.18 {
				risk.OldDependencies = append(risk.OldDependencies, "Go version < 1.18 (Critical EOL)")
				risk.RiskLevel = "HIGH"
				risk.Reasons = append(risk.Reasons, "Go version < 1.18")
			} else if ver < 1.21 {
				risk.OldDependencies = append(risk.OldDependencies, "Go version < 1.21 (Recommended Upgrade)")
				if risk.RiskLevel != "HIGH" {
					risk.RiskLevel = "MEDIUM"
					risk.Reasons = append(risk.Reasons, "Go version < 1.21")
				}
			}
		}
	}

	// 3. Heuristic Scan for GPL in vendor (if exists)
	// Shallow scan
	vendorPath := filepath.Join(path, "vendor")
	if info, err := os.Stat(vendorPath); err == nil && info.IsDir() {
		filepath.Walk(vendorPath, func(p string, info os.FileInfo, err error) error {
			if !info.IsDir() && strings.Contains(strings.ToUpper(info.Name()), "LICENSE") {
				// Read file
				b, _ := os.ReadFile(p)
				c := string(b)
				if strings.Contains(c, "GNU GENERAL PUBLIC LICENSE") {
					rel, _ := filepath.Rel(path, p)
					risk.RiskyLicenses = append(risk.RiskyLicenses, rel+" (GPL Detected)")
					risk.RiskLevel = "HIGH"
					risk.Reasons = append(risk.Reasons, "GPL License detected in vendor")
				}
			}
			return nil
		})
	}

	// 4. Check CI Configuration (Quick check for Risk Analysis)
	hasCI := false
	possibleCIPaths := []string{
		".github",
		".gitlab-ci.yml",
		"circleci",
		"jenkins",
		".travis.yml",
		"azure-pipelines.yml",
	}
	for _, p := range possibleCIPaths {
		if _, err := os.Stat(filepath.Join(path, p)); err == nil {
			hasCI = true
			break
		}
	}
	if !hasCI {
		if risk.RiskLevel == "LOW" {
			risk.RiskLevel = "MEDIUM"
		}
		risk.Reasons = append(risk.Reasons, "No CI configuration")
	}

	return risk
}
