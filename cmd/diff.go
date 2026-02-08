package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bilal-sayyed/repocheck-go/internal/checker"
	"github.com/bilal-sayyed/repocheck-go/internal/license"
	"github.com/bilal-sayyed/repocheck-go/pkg/report"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff [old.json] [new.json]",
	Short: "Compare two scan reports (Pro only)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Check License
		if !license.IsPro() {
			fmt.Println("❌ 'diff' is a Pro feature. Please add a license key.")
			os.Exit(1)
		}

		oldFile := args[0]
		newFile := args[1]

		// 2. Load JSONs
		oldRes, err := loadScanResult(oldFile)
		if err != nil {
			fmt.Printf("Error loading %s: %v\n", oldFile, err)
			os.Exit(1)
		}
		newRes, err := loadScanResult(newFile)
		if err != nil {
			fmt.Printf("Error loading %s: %v\n", newFile, err)
			os.Exit(1)
		}

		// 3. Compare
		diff := checker.CompareScans(oldRes, newRes)

		// 4. Output
		fmt.Println("=== Repocheck Diff Report ===")

		fmt.Printf("Score Change: %d -> %d (%+d)\n",
			getScore(oldRes), getScore(newRes), diff.ScoreChange)

		if len(diff.NewIssues) > 0 {
			fmt.Println("\n🔴 New Issues (Regressions):")
			for _, issue := range diff.NewIssues {
				fmt.Println("  - " + issue)
			}
		} else {
			fmt.Println("\n✅ No new issues found.")
		}

		if len(diff.FixedIssues) > 0 {
			fmt.Println("\n🟢 Fixed Issues:")
			for _, issue := range diff.FixedIssues {
				fmt.Println("  - " + issue)
			}
		}
	},
}

func loadScanResult(path string) (*report.RepoScanResult, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Strip BOMs
	// UTF-8 BOM
	if len(content) >= 3 {
		if content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
			// Found UTF-8 BOM
			content = content[3:]
		}
	}

	var res report.RepoScanResult
	if err := json.Unmarshal(content, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func getScore(r *report.RepoScanResult) int {
	if r.ReadinessScore != nil {
		return r.ReadinessScore.Score
	}
	return 0
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
