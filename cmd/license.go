package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/bilal-sayyed/repocheck-go/internal/license"
	"github.com/spf13/cobra"
)

var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Manage repocheck license",
	Long:  `View status or add a license key to unlock Pro features.`,
}

var licenseAddCmd = &cobra.Command{
	Use:   "add [key]",
	Short: "Add a license key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		_, err := license.VerifyLicense(key)
		if err == nil {
			err := license.SaveLicense(key)
			if err != nil {
				fmt.Printf("Error saving license: %v\n", err)
				return
			}
			fmt.Println("✅ License validated and saved! Pro features unlocked.")
		} else {
			fmt.Println("❌ Invalid license key.")
		}
	},
}

var licenseStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check license status",
	Run: func(cmd *cobra.Command, args []string) {
		key, err := license.LoadLicense()
		if err != nil || key == "" {
			fmt.Println("Status: Free Tier (No active license)")
			return
		}

		lic, err := license.VerifyLicense(key)
		if err == nil {
			fmt.Println("Status: ✅ Pro Tier (Active)")
			fmt.Printf("Plan: %s %s\n", lic.Plan, lic.Variant)

			expTime := time.Unix(lic.Exp, 0)
			daysRemaining := int(time.Until(expTime).Hours() / 24)

			fmt.Printf("Expires: %s\n", expTime.Format("02 Jan 2006"))
			fmt.Printf("Days remaining: %d\n", daysRemaining)
		} else {
			fmt.Printf("Status: ❌ Invalid License (%v)\n", err)
		}
	},
}

var licenseRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove license key",
	Run: func(cmd *cobra.Command, args []string) {
		path := license.GetLicensePath()
		err := os.Remove(path)
		if err != nil {
			// Ignore error if file doesn't exist
			if !os.IsNotExist(err) {
				fmt.Printf("Error removing license: %v\n", err)
				return
			}
		}
		fmt.Println("License removed. Reverted to Free Tier.")
	},
}

func init() {
	rootCmd.AddCommand(licenseCmd)
	licenseCmd.AddCommand(licenseAddCmd)
	licenseCmd.AddCommand(licenseRemoveCmd)
	licenseCmd.AddCommand(licenseStatusCmd)
}
