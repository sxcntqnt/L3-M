package report

import (
	"encoding/json"
	"fmt"
	"os"
)

// SelectorResult = result for a single selector check
type SelectorResult struct {
	Label  string `json:"label"`
	Status string `json:"status"`
}

// BookieReport = detailed verification for one bookie
type BookieReport struct {
	Name    string           `json:"name"`
	URL     string           `json:"url"`
	Results []SelectorResult `json:"results"`
	AllPass bool             `json:"all_pass"`
}

// FullReport = JSON structure with summary + details
type FullReport struct {
	Summary []BookieReport `json:"summary"`
	Details []BookieReport `json:"details"`
}

// SaveJSON writes the full report to JSON
func SaveJSON(report FullReport, filename string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}
	fmt.Printf("📄 Saved JSON report: %s\n", filename)
	return nil
}

// SaveMarkdown writes the full report to Markdown
func SaveMarkdown(report FullReport, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer f.Close()

	// Summary Table
	fmt.Fprintf(f, "# Verification Report\n\n")
	fmt.Fprintf(f, "## 📊 Summary\n")
	fmt.Fprintf(f, "| Bookie | URL | Status |\n")
	fmt.Fprintf(f, "|--------|-----|--------|\n")
	for _, s := range report.Summary {
		status := "✅"
		if !s.AllPass {
			status = "❌"
		}
		fmt.Fprintf(f, "| %s | %s | %s |\n", s.Name, s.URL, status)
	}

	// Details
	fmt.Fprintf(f, "\n---\n\n")
	for _, d := range report.Details {
		overall := "✅ Passed"
		if !d.AllPass {
			overall = "❌ Failed"
		}
		fmt.Fprintf(f, "## %s (%s)\n", d.Name, d.URL)
		for _, res := range d.Results {
			fmt.Fprintf(f, "- %s: %s\n", res.Label, res.Status)
		}
		fmt.Fprintf(f, "Overall: %s\n\n", overall)
	}

	fmt.Printf("📄 Saved Markdown report: %s\n", filename)
	return nil
}

