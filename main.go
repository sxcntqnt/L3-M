package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"diago/config"
	"diago/fetch"
	"diago/report"

        _ "diago/bookies" 
	"gopkg.in/yaml.v3"
)

func main() {
	mode := flag.String("mode", "fetch", "Mode: generate, fetch, or auto")
	bookiesFile := flag.String("bookies-file", "bookies.txt", "Bookies file")
	outputDir := flag.String("output-dir", "EMC", "Output directory")
	flag.Parse()

	bookiesList, err := loadBookies(*bookiesFile)
	if err != nil {
		fmt.Printf("❌ Failed to read bookies file: %v\n", err)
		os.Exit(1)
	}

	overridesPath := filepath.Join(*outputDir, "overrides.yaml")
	var overrides config.OverrideMap
	if _, err := os.Stat(overridesPath); err == nil {
		overrides, _ = config.LoadOverrides(overridesPath)
	} else {
		overrides = make(config.OverrideMap)
	}

	switch *mode {
	case "generate":
		generateConfigs(bookiesList, overrides, *outputDir)
	case "fetch":
		fullReport := fetchConfigs(bookiesList, *outputDir)
		createLatestSnippet(fullReport, *outputDir)
	case "auto":
		if configsMissing(bookiesList, *outputDir) {
			fmt.Println("🛠️ Configs missing – generating first...")
			generateConfigs(bookiesList, overrides, *outputDir)
		}
		fullReport := fetchConfigs(bookiesList, *outputDir)
		createLatestSnippet(fullReport, *outputDir)
	default:
		fmt.Printf("❌ Unknown mode: %s\n", *mode)
		os.Exit(1)
	}
}

func loadBookies(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var bookies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		bookies = append(bookies, parts[0])
	}
	return bookies, scanner.Err()
}

func configsMissing(bookies []string, outputDir string) bool {
	for _, b := range bookies {
		cfgPath := filepath.Join(outputDir, strings.ToLower(b), "config.yaml")
		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			return true
		}
	}
	return false
}

func generateConfigs(bookies []string, overrides config.OverrideMap, outputDir string) {
	fmt.Println("🛠️ Generating configs...")
	if err := config.GenerateConfigs(bookies, overrides, outputDir, "https://%s.example.com", "/usr/bin/chrome"); err != nil {
		fmt.Printf("❌ Error generating configs: %v\n", err)
		os.Exit(1)
	}
}

// fetchConfigs now returns the full report
func fetchConfigs(bookies []string, outputDir string) report.FullReport {
	fmt.Println("🌐 Fetching and verifying bookies...")

	var summary []report.BookieReport
	var details []report.BookieReport

	for _, b := range bookies {
		cfgPath := filepath.Join(outputDir, strings.ToLower(b), "config.yaml")
		cfg, err := loadConfig(cfgPath)
		if err != nil {
			fmt.Printf("⚠️ Skipping %s, failed to load config: %v\n", b, err)
			continue
		}

		r := fetch.VerifyBookieWithConfig(cfg.Name, cfg.BaseURL, cfg)
		details = append(details, r)
		summary = append(summary, r)
	}

	fullReport := report.FullReport{
		Summary: summary,
		Details: details,
	}

	jsonPath := filepath.Join(outputDir, "report.json")
	mdPath := filepath.Join(outputDir, "report.md")

	if err := report.SaveJSON(fullReport, jsonPath); err != nil {
		fmt.Printf("❌ Failed to save JSON report: %v\n", err)
		os.Exit(1)
	}
	if err := report.SaveMarkdown(fullReport, mdPath); err != nil {
		fmt.Printf("❌ Failed to save Markdown report: %v\n", err)
		os.Exit(1)
	}

	return fullReport
}

// createLatestSnippet now builds it from FullReport in memory
func createLatestSnippet(fullReport report.FullReport, outputDir string) {
	latestMD := filepath.Join(outputDir, "latest_report.md")

	f, err := os.Create(latestMD)
	if err != nil {
		fmt.Printf("❌ Failed to create latest_report.md: %v\n", err)
		return
	}
	defer f.Close()

	fmt.Fprintf(f, "## 📊 Summary\n")
	fmt.Fprintf(f, "| Bookie | URL | Status |\n")
	fmt.Fprintf(f, "|--------|-----|--------|\n")
	for _, s := range fullReport.Summary {
		status := "✅"
		if !s.AllPass {
			status = "❌"
		}
		fmt.Fprintf(f, "| %s | %s | %s |\n", s.Name, s.URL, status)
	}

	fmt.Fprintf(f, "\n_Updated automatically via GitHub Actions_\n")
	fmt.Printf("✅ Created latest report snippet: %s\n", latestMD)
}

func loadConfig(path string) (*config.Sportsbook, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var sb config.Sportsbook
	if err := yaml.Unmarshal(data, &sb); err != nil {
		return nil, err
	}
	return &sb, nil
}

