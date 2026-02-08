package checker

import (
	"github.com/bilal-sayyed/repocheck-go/pkg/report"
)

type DiffReport struct {
	ScoreChange int      `json:"score_change"` // New - Old
	NewIssues   []string `json:"new_issues"`   // Issues present in New but not Old
	FixedIssues []string `json:"fixed_issues"` // Issues in Old but not New
	Regressions bool     `json:"regressions"`  // True if score dropped or new issues found
}

func CompareScans(old, new *report.RepoScanResult) DiffReport {
	diff := DiffReport{
		NewIssues:   []string{},
		FixedIssues: []string{},
	}

	// 1. Score Diff
	oldScore := 0
	if old.ReadinessScore != nil {
		oldScore = old.ReadinessScore.Score
	}
	newScore := 0
	if new.ReadinessScore != nil {
		newScore = new.ReadinessScore.Score
	}

	diff.ScoreChange = newScore - oldScore

	// 2. Compare Hygiene (Secrets)
	// Simple set comparison
	diff.NewIssues = append(diff.NewIssues, findNewStrings(old.Hygiene.HardcodedSecrets, new.Hygiene.HardcodedSecrets, "Secret found: ")...)
	diff.FixedIssues = append(diff.FixedIssues, findFixedStrings(old.Hygiene.HardcodedSecrets, new.Hygiene.HardcodedSecrets, "Secret removed: ")...)

	// 3. Compare Onboarding
	diff.NewIssues = append(diff.NewIssues, findNewStrings(old.Readiness.MissingElements, new.Readiness.MissingElements, "Missing Onboarding: ")...)
	diff.FixedIssues = append(diff.FixedIssues, findFixedStrings(old.Readiness.MissingElements, new.Readiness.MissingElements, "Fixed Onboarding: ")...)

	// 4. Compare Risk
	if old.LicenseRisk != nil && new.LicenseRisk != nil {
		diff.NewIssues = append(diff.NewIssues, findNewStrings(old.LicenseRisk.RiskyLicenses, new.LicenseRisk.RiskyLicenses, "New Risk License: ")...)
		diff.FixedIssues = append(diff.FixedIssues, findFixedStrings(old.LicenseRisk.RiskyLicenses, new.LicenseRisk.RiskyLicenses, "Removed Risk License: ")...)
	}

	if len(diff.NewIssues) > 0 || diff.ScoreChange < 0 {
		diff.Regressions = true
	}

	return diff
}

func findNewStrings(oldList, newList []string, prefix string) []string {
	res := []string{}
	oldMap := make(map[string]bool)
	for _, s := range oldList {
		oldMap[s] = true
	}

	for _, s := range newList {
		if !oldMap[s] {
			res = append(res, prefix+s)
		}
	}
	return res
}

func findFixedStrings(oldList, newList []string, prefix string) []string {
	res := []string{}
	newMap := make(map[string]bool)
	for _, s := range newList {
		newMap[s] = true
	}

	for _, s := range oldList {
		if !newMap[s] {
			res = append(res, prefix+s)
		}
	}
	return res
}
