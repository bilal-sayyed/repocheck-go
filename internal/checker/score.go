package checker

import (
	"fmt"
	"strings"

	"github.com/bilal-sayyed/repocheck-go/pkg/report"
)

func CalculateScore(result *report.RepoScanResult) *report.ReadinessScore {
	docScore := 0
	depScore := 0
	ciScore := 0
	secScore := 0

	// 1. Documentation (Max 25)
	if result.Readiness.HasReadme {
		docScore += 10
	}
	if result.Readiness.HasRunInstructions {
		docScore += 10
	}
	if result.Readiness.HasSetupInstructions {
		docScore += 5
	}
	// Bonus: Env Example can count towards docs or hygiene. Putting it in Security/Hygiene.

	// 2. Dependencies (Max 25)
	if result.Dependencies.Manager != "Unknown" {
		depScore += 25
	}

	// 3. CI / DevOps (Max 25)
	if result.CIDevOps != nil {
		if result.CIDevOps.HasDockerfile {
			ciScore += 10
		}
		if result.CIDevOps.HasCI {
			ciScore += 15
		}
	}

	// 4. Security (Max 25)
	// Secrets (10)
	if len(result.Hygiene.HardcodedSecrets) == 0 {
		secScore += 10
	}
	// License Risk (10)
	if result.LicenseRisk != nil {
		if result.LicenseRisk.RiskLevel == "LOW" {
			secScore += 10
		} else if result.LicenseRisk.RiskLevel == "MEDIUM" {
			secScore += 5
		}
	} else {
		// If no risk analysis (Free?), assume safe or 0?
		// CalculateScore is currently only called for Pro.
		// If LicenseRisk is missing, we can default to 0 or 5.
		secScore += 5
	}
	// Hygiene (5)
	if result.Hygiene.HasEnvExample {
		secScore += 3
	}
	if result.Hygiene.HasGitignore {
		secScore += 2
	}

	totalScore := docScore + depScore + ciScore + secScore

	// Grade
	grade := "F (Not production ready)"
	if totalScore >= 90 {
		grade = "A (Excellent)"
	} else if totalScore >= 80 {
		grade = "B (Good)"
	} else if totalScore >= 70 {
		grade = "C (Acceptable)"
	} else if totalScore >= 60 {
		grade = "D (Needs improvement)"
	}

	// Generate Executive Summary
	riskStr := "Unknown risk"
	if result.LicenseRisk != nil {
		riskStr = result.LicenseRisk.RiskLevel + " risk"
	}

	summary := fmt.Sprintf("%s repository", riskStr)

	// Strengths
	var strengths []string
	if docScore >= 20 {
		strengths = append(strengths, "documentation")
	}
	if depScore >= 20 {
		strengths = append(strengths, "dependencies")
	}
	if ciScore >= 20 {
		strengths = append(strengths, "CI/CD")
	}
	if secScore >= 20 {
		strengths = append(strengths, "security")
	}

	if len(strengths) > 0 {
		summary += " with strong " + strings.Join(strengths, " and ")
	}

	// Weaknesses
	var weaknesses []string
	if ciScore < 10 {
		weaknesses = append(weaknesses, "CI pipeline")
	}
	if secScore < 10 {
		weaknesses = append(weaknesses, "security compliance")
	}
	if result.LicenseRisk != nil && len(result.LicenseRisk.RiskyLicenses) > 0 {
		weaknesses = append(weaknesses, "license compliance")
	}

	if len(weaknesses) > 0 {
		if len(strengths) > 0 {
			summary += ", but missing "
		} else {
			summary += " but missing "
		}
		summary += strings.Join(weaknesses, " and ")
	}

	summary += "."
	result.ExecutiveSummary = summary

	return &report.ReadinessScore{
		Score: totalScore,
		Grade: grade,
		Breakdown: &report.ScoreBreakdown{
			Documentation: docScore,
			Dependencies:  depScore,
			CIDevOps:      ciScore,
			Security:      secScore,
		},
	}
}
