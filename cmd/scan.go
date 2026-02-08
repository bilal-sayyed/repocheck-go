package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bilal-sayyed/repocheck-go/internal/checker"
	"github.com/bilal-sayyed/repocheck-go/internal/config"
	"github.com/bilal-sayyed/repocheck-go/internal/license"
	"github.com/bilal-sayyed/repocheck-go/internal/scanner"
	"github.com/bilal-sayyed/repocheck-go/pkg/report"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var jsonOutput bool
var summaryOutput bool
var ciMode bool
var noColor bool

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan a repository for hygiene and risks",
	Long: `scans the repository at the given path (or current directory) and 
reports on repository type, language, and potential issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		pathToScan := "."
		if len(args) > 0 {
			pathToScan = args[0]
		}

		// Configure Color
		if noColor || ciMode || jsonOutput {
			color.NoColor = true
		} else {
			// Force valid colors for Windows IDEs unless explicitly disabled
			color.NoColor = false
		}

		// Load config
		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if !ciMode && !jsonOutput {
			fmt.Printf("Scanning %s...\n", pathToScan)
		}
		if cfg.StrictMode && !ciMode && !jsonOutput {
			fmt.Println("Strict mode enabled")
		}

		// 1. Run Core Scanner
		scanResult, err := scanner.ScanRepository(pathToScan)
		if err != nil {
			if ciMode {
				fmt.Printf("STATUS=FAIL\nERROR=%v\n", err)
				os.Exit(2)
			}
			fmt.Printf("Error scanning repository: %v\n", err)
			os.Exit(1)
		}

		// 2. Run Onboarding Checks
		scanResult.Readiness = checker.CheckOnboarding(pathToScan)

		// 3. Run Dependency Checks
		scanResult.Dependencies = checker.CheckDependencies(pathToScan)

		// 4. Run Hygiene Checks
		scanResult.Hygiene = checker.CheckHygiene(pathToScan)

		// 5. Pro Checks & Action Items
		isPro := license.IsPro()

		if isPro {
			scanResult.LicenseRisk = checker.CheckLicenseRisk(pathToScan)
			scanResult.CIDevOps = checker.CheckCIDevOps(pathToScan)
			scanResult.Ownership = checker.CheckOwnership(pathToScan)
			scanResult.ReadinessScore = checker.CalculateScore(scanResult)

			// Populate Action Items
			// 1. Root LICENSE
			if scanResult.LicenseRisk != nil {
				for _, r := range scanResult.LicenseRisk.Reasons {
					if r == "Missing LICENSE file" {
						scanResult.ActionItems = append(scanResult.ActionItems, "Add a root LICENSE file (Required for compliance)")
					}
				}
			}
			// 2. CI Config
			if scanResult.CIDevOps != nil && !scanResult.CIDevOps.HasCI {
				scanResult.ActionItems = append(scanResult.ActionItems, "Add CI configuration (GitHub Actions recommended)")
			}
			// 3. Env Example
			if !scanResult.Hygiene.HasEnvExample {
				scanResult.ActionItems = append(scanResult.ActionItems, "Add .env.example to document environment variables")
			}
		}

		// 6. Output
		if ciMode {
			printCIReport(scanResult, isPro)

			// Exit Codes
			// 0 = OK, 1 = Warning, 2 = Fail
			code := 0
			// Fail conditions
			if scanResult.LicenseRisk != nil && scanResult.LicenseRisk.RiskLevel == "HIGH" {
				code = 2
			} else if len(scanResult.Hygiene.HardcodedSecrets) > 0 {
				code = 2
			} else if scanResult.ReadinessScore != nil && scanResult.ReadinessScore.Score < 50 {
				code = 2
			} else if scanResult.LicenseRisk != nil && scanResult.LicenseRisk.RiskLevel == "MEDIUM" {
				// Warn conditions
				code = 1
			}
			os.Exit(code)
		} else if jsonOutput {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(scanResult)
		} else if summaryOutput {
			printSummaryReport(scanResult, isPro)
		} else {
			printTextReport(scanResult, isPro)
		}
	},
}

func printTextReport(r *report.RepoScanResult, isPro bool) {
	header := color.New(color.FgCyan, color.Bold).SprintfFunc()
	// bold := color.New(color.Bold).SprintfFunc()

	fmt.Println("\n" + header("=== Repocheck Analysis ==="))

	if isPro && r.ExecutiveSummary != "" {
		fmt.Printf("\nSummary:\n%s\n", r.ExecutiveSummary)
	}

	tier := "Free"
	if isPro {
		tier = "Pro"
	}
	fmt.Printf("Tier:     %s\n", tier)
	fmt.Printf("Language: %s\n", r.Summary.Language)
	fmt.Printf("Type:     %s\n", r.Summary.Type)
	fmt.Printf("Files:    %d\n", r.Summary.FileCount)
	fmt.Printf("Git:      %v\n", r.Summary.GitPresent)

	fmt.Println("\n" + header("[Onboarding]"))
	fmt.Printf("README:   %v\n", r.Readiness.HasReadme)
	if len(r.Readiness.MissingElements) > 0 {
		fmt.Printf("Missing:  %v\n", r.Readiness.MissingElements)
	} else {
		statusColor := color.New(color.FgGreen).SprintfFunc()
		fmt.Printf("Status:   %s\n", statusColor("All Clean"))
	}

	fmt.Println("\n" + header("[Dependencies]"))
	fmt.Printf("Manager:  %s\n", r.Dependencies.Manager)
	fmt.Printf("Count:    %d (Direct: %d)\n", r.Dependencies.TotalCount, r.Dependencies.DirectCount)

	fmt.Println("\n" + header("[Hygiene]"))
	fmt.Printf(".env.example: %v\n", r.Hygiene.HasEnvExample)
	if len(r.Hygiene.HardcodedSecrets) > 0 {
		// New Cleaner Warning
		files := ""
		if len(r.Hygiene.HardcodedSecrets) == 1 {
			files = fmt.Sprintf("(%s)", r.Hygiene.HardcodedSecrets[0])
		} else {
			files = fmt.Sprintf("(%s, ...)", r.Hygiene.HardcodedSecrets[0])
		}
		secretColor := color.New(color.FgRed, color.Bold).SprintfFunc()
		fmt.Printf("Secrets:      %s potential secret detected %s\n", secretColor(fmt.Sprintf("%d", len(r.Hygiene.HardcodedSecrets))), files)
	} else {
		fmt.Println("Secrets:      None detected")
	}

	// Pro Sections
	if isPro && r.LicenseRisk != nil {
		fmt.Println("\n" + header("[Risk Analysis] (PRO)"))

		riskColor := color.New(color.FgGreen).SprintfFunc()
		if r.LicenseRisk.RiskLevel == "MEDIUM" {
			riskColor = color.New(color.FgYellow).SprintfFunc()
		} else if r.LicenseRisk.RiskLevel == "HIGH" {
			riskColor = color.New(color.FgRed).SprintfFunc()
		}

		fmt.Printf("Level:    %s\n", riskColor(r.LicenseRisk.RiskLevel))
		if len(r.LicenseRisk.Reasons) > 0 {
			fmt.Println("Reason:")
			for _, re := range r.LicenseRisk.Reasons {
				fmt.Printf(" - %s\n", re)
			}
		}
		if len(r.LicenseRisk.RiskyLicenses) > 0 {
			fmt.Println("Licenses:")
			for _, l := range r.LicenseRisk.RiskyLicenses {
				fmt.Printf(" - %s\n", l)
			}
		}
	} else if !isPro {
		fmt.Println("\n" + header("[Risk Analysis]"))
		fmt.Println("Upgrade to Pro to view license and compliance risks.")
	}

	if isPro && r.ReadinessScore != nil {
		fmt.Println("\n" + header("[Readiness Score] (PRO)"))

		scoreColor := color.New(color.FgRed).SprintfFunc()
		if r.ReadinessScore.Score >= 80 {
			scoreColor = color.New(color.FgGreen).SprintfFunc()
		} else if r.ReadinessScore.Score >= 50 {
			scoreColor = color.New(color.FgYellow).SprintfFunc()
		}

		fmt.Printf("Score:    %s/100\n", scoreColor(fmt.Sprintf("%d", r.ReadinessScore.Score)))
		fmt.Printf("Grade:    %s\n", r.ReadinessScore.Grade)
		if r.ReadinessScore.Breakdown != nil {
			fmt.Println("[Readiness Score Breakdown]")
			fmt.Printf("Documentation:   %d/25\n", r.ReadinessScore.Breakdown.Documentation)
			fmt.Printf("Dependencies:    %d/25\n", r.ReadinessScore.Breakdown.Dependencies)
			fmt.Printf("CI / DevOps:     %d/25\n", r.ReadinessScore.Breakdown.CIDevOps)
			fmt.Printf("Security:        %d/25\n", r.ReadinessScore.Breakdown.Security)
		}
	} else if !isPro {
		fmt.Println("\n" + header("[Readiness Score]"))
		fmt.Println("Upgrade to Pro to view readiness score and breakdown.")
	}

	if isPro && r.CIDevOps != nil {
		fmt.Println("\n" + header("[CI / DevOps] (PRO)"))
		fmt.Printf("Dockerfile: %v\n", r.CIDevOps.HasDockerfile)
		fmt.Printf("CI Tools:   %v\n", r.CIDevOps.Tools)
	} else if !isPro {
		fmt.Println("\n" + header("[CI / DevOps]"))
		fmt.Println("Upgrade to Pro to analyze CI readiness.")
	}

	if isPro && r.Ownership != nil {
		fmt.Println("\n" + header("[Ownership] (PRO)"))
		fmt.Printf("CODEOWNERS: %v\n", r.Ownership.HasCodeowners)
		status := "Undefined"
		statusColor := color.New(color.FgYellow).SprintfFunc()
		if r.Ownership.HasCodeowners || r.Ownership.HasMaintainers {
			status = "Defined"
			statusColor = color.New(color.FgGreen).SprintfFunc()
		}
		fmt.Printf("Status:     %s\n", statusColor(status))
	}

	if isPro {
		fmt.Printf("\n%s\n", header("[Extended Hygiene] (PRO)"))
		fmt.Printf(".gitignore:   %v\n", r.Hygiene.HasGitignore)
		fmt.Printf(".env Ignored: %v\n", r.Hygiene.EnvIgnored)
	}

	if isPro && len(r.ActionItems) > 0 {
		fmt.Println("\n" + header("[Next Actions]"))
		bold := color.New(color.Bold).SprintfFunc()
		for i, item := range r.ActionItems {
			fmt.Printf("%s. %s\n", bold(fmt.Sprintf("%d", i+1)), item)
		}
	} else if !isPro {
		fmt.Println("\n" + header("[Next Actions]"))
		fmt.Println("Upgrade to Pro to receive prioritized remediation steps.")
	}
}

func printCIReport(r *report.RepoScanResult, isPro bool) {
	// Status
	status := "PASS"
	if r.ReadinessScore != nil && r.ReadinessScore.Score < 50 {
		status = "FAIL"
	} else if len(r.Hygiene.HardcodedSecrets) > 0 {
		status = "FAIL"
	}
	fmt.Printf("STATUS=%s\n", status)

	// Risk (If known)
	risk := "UNKNOWN"
	if r.LicenseRisk != nil {
		risk = r.LicenseRisk.RiskLevel
	}
	fmt.Printf("RISK=%s\n", risk)

	// Score
	score := 0
	if r.ReadinessScore != nil {
		score = r.ReadinessScore.Score
	}
	fmt.Printf("SCORE=%d\n", score)

	// Critical Issues
	critical := 0
	if len(r.Hygiene.HardcodedSecrets) > 0 {
		critical += len(r.Hygiene.HardcodedSecrets)
	}
	if r.LicenseRisk != nil && r.LicenseRisk.RiskLevel == "HIGH" {
		critical += 1
	}
	fmt.Printf("CRITICAL_ISSUES=%d\n", critical)
}

func printSummaryReport(r *report.RepoScanResult, isPro bool) {
	tier := "Free"
	if isPro {
		tier = "Pro"
	}
	fmt.Printf("Repocheck Summary (%s)\n", tier)
	fmt.Printf("Language: %s\n", r.Summary.Language)
	fmt.Printf("Issues:   %d missing onboarding items\n", len(r.Readiness.MissingElements))
	if isPro && r.ReadinessScore != nil {
		scoreColor := color.New(color.FgRed).SprintfFunc()
		if r.ReadinessScore.Score >= 80 {
			scoreColor = color.New(color.FgGreen).SprintfFunc()
		} else if r.ReadinessScore.Score >= 50 {
			scoreColor = color.New(color.FgYellow).SprintfFunc()
		}
		fmt.Printf("Score:    %s/100 (%s)\n", scoreColor(fmt.Sprintf("%d", r.ReadinessScore.Score)), r.ReadinessScore.Grade)
	}
}

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Print a high-level summary of repository status (Pro only)",
	Run: func(cmd *cobra.Command, args []string) {
		pathToScan := "."
		if len(args) > 0 {
			pathToScan = args[0]
		}

		if noColor || ciMode || jsonOutput {
			color.NoColor = true
		} else {
			color.NoColor = false
		}

		if !license.IsPro() {
			fmt.Println("The 'summary' command is a Pro-only feature.")
			fmt.Println("Please upgrade to Pro to access this summary.")
			os.Exit(1)
		}

		// Initial Output
		header := color.New(color.FgCyan, color.Bold).SprintfFunc()
		fmt.Println(header("Repocheck Summary"))
		fmt.Println("-----------------")

		// Re-use logic (simplified scan)
		res, err := scanner.ScanRepository(pathToScan)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Run checkers
		res.Readiness = checker.CheckOnboarding(pathToScan)
		res.Dependencies = checker.CheckDependencies(pathToScan)
		res.Hygiene = checker.CheckHygiene(pathToScan)
		res.LicenseRisk = checker.CheckLicenseRisk(pathToScan)
		// We need full score for summary
		res.CIDevOps = checker.CheckCIDevOps(pathToScan)
		res.Ownership = checker.CheckOwnership(pathToScan)
		res.ReadinessScore = checker.CalculateScore(res)

		// Output
		riskColor := color.New(color.FgGreen).SprintfFunc()
		if res.LicenseRisk.RiskLevel == "MEDIUM" {
			riskColor = color.New(color.FgYellow).SprintfFunc()
		} else if res.LicenseRisk.RiskLevel == "HIGH" {
			riskColor = color.New(color.FgRed).SprintfFunc()
		}

		fmt.Printf("Risk Level: %s\n", riskColor(res.LicenseRisk.RiskLevel))

		scoreColor := color.New(color.FgRed).SprintfFunc()
		if res.ReadinessScore.Score >= 80 {
			scoreColor = color.New(color.FgGreen).SprintfFunc()
		} else if res.ReadinessScore.Score >= 50 {
			scoreColor = color.New(color.FgYellow).SprintfFunc()
		}
		fmt.Printf("Readiness:  %s/100 (%s)\n", scoreColor(fmt.Sprintf("%d", res.ReadinessScore.Score)), res.ReadinessScore.Grade)

		fmt.Println("\nTop Issues:")
		if len(res.Readiness.MissingElements) > 0 {
			for _, m := range res.Readiness.MissingElements {
				fmt.Printf(" - Missing %s\n", m)
			}
		}
		for _, r := range res.LicenseRisk.Reasons {
			fmt.Printf(" - %s\n", r)
		}
		if len(res.Hygiene.HardcodedSecrets) > 0 {
			fmt.Println(" - Secrets detected")
		}

		fmt.Println("\nRecommendation:")
		if res.ReadinessScore.Score < 70 {
			fmt.Println("Fix licensing and CI to reach C-grade readiness.")
		} else {
			fmt.Println("Great job! Maintain this hygiene.")
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(summaryCmd)
	scanCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
	scanCmd.Flags().BoolVar(&summaryOutput, "summary", false, "Output brief summary only")
	scanCmd.Flags().BoolVar(&ciMode, "ci", false, "Strict exit codes for CI/CD")
	scanCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable color output")
	summaryCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable color output")
}
