package cmd

import (
	"fmt"
	"os"

	"diago/fetch"
	"diago/report"
	"diago/utils"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "diago",
	Short: "Bookie verification CLI",
	Run: func(cmd *cobra.Command, args []string) {
		// Load enabled bookies from bookies.txt
		enabled, err := utils.EnabledBookies("bookies.txt")
		if err != nil {
			fmt.Println("❌ Failed to load bookies:", err)
			os.Exit(1)
		}

		if len(enabled) == 0 {
			fmt.Println("⚠️ No enabled bookies found in bookies.txt")
			os.Exit(0)
		}

		// Run verification for each bookie
		var reports []report.BookieReport
		for _, b := range enabled {
			r := fetch.VerifyBookieWithConfig(b)
			reports = append(reports, r)
		}

		// Save using the report package
		if err := report.SaveReports(".", reports); err != nil {
			fmt.Println("❌ Error saving reports:", err)
			os.Exit(1)
		}

		fmt.Println("✅ Verification complete. Reports saved in report.json & report.md")
	},
}

// Execute runs the root cobra command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

