package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"diago/config"
	"diago/fetch"
	"diago/report"
	"diago/utils"

	_ "diago/bookies"
	"gopkg.in/yaml.v3"
)

func main() {
	mode := flag.String("mode", "fetch", "Mode: generate, fetch, or auto")
	bookiesFile := flag.String("bookies-file", "bookies.txt", "Bookies file")
	outputDir := flag.String("output-dir", "EMC", "Output directory")
	bakeOverrides := flag.Bool("bake-overrides", false, "Apply overrides and persist them to config.yaml, then delete overrides.yaml")
	flag.Parse()

	// Load enabled bookies from bookies.txt
	enabledBookies, err := utils.EnabledBookies(*bookiesFile)
	if err != nil {
		fmt.Printf("❌ Failed to load bookies: %v\n", err)
		os.Exit(1)
	}

	overridesPath := filepath.Join(*outputDir, "overrides.yaml")
	var overrides config.OverrideMap
	usingOverrides := false

	if _, err := os.Stat(overridesPath); err == nil {
		overrides, _ = config.LoadOverrides(overridesPath)
		usingOverrides = true
	} else {
		overrides = make(config.OverrideMap)
	}

	switch *mode {
	case "generate":
		generateConfigs(enabledBookies, overrides, *outputDir)
		if *bakeOverrides && usingOverrides {
			bakeOverridesFile(*outputDir, overridesPath)
		}

	case "fetch":
		fullReport := fetchConfigs(enabledBookies, *outputDir)
		createLatestSnippet(fullReport, *outputDir)

	case "auto":
		if configsMissing(enabledBookies, *outputDir) {
			fmt.Println("🛠️ Configs missing – generating first...")
			generateConfigs(enabledBookies, overrides, *outputDir)
			if *bakeOverrides && usingOverrides {
				bakeOverridesFile(*outputDir, overridesPath)
			}
		}
		fullReport := fetchConfigs(enabledBookies, *outputDir)
		createLatestSnippet(fullReport, *outputDir)

	default:
		fmt.Printf("❌ Unknown mode: %s\n", *mode)
		os.Exit(1)
	}
}

// configsMissing checks if any bookie's config.yaml is missing
func configsMissing(bookies []utils.Bookie, outputDir string) bool {
	for _, b := range bookies {
		cfgPath := filepath.Join(outputDir, strings.ToLower(b.Name()), "config.yaml")
		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			return true
		}
	}
	return false
}

// generateConfigs writes config files for all bookies
func generateConfigs(bookies []utils.Bookie, overrides config.OverrideMap, outputDir string) {
    fmt.Println("🛠️ Generating configs...")

    for _, b := range bookies {
        folderName := strings.ToLower(b.Name())
        folderPath := filepath.Join(outputDir, folderName)
        if err := os.MkdirAll(folderPath, 0755); err != nil {
            fmt.Printf("❌ Failed to create folder %s: %v\n", folderPath, err)
            continue
        }

        // Ensure we are passing a map[string]interface{} as expected
        overrideForBookie, exists := overrides[b.Name()]
        if !exists {
            overrideForBookie = make(map[string]interface{})
        }

        // Wrap the override in a map of OverrideMap
        err := config.GenerateConfig(b.Name(), config.OverrideMap{b.Name(): overrideForBookie}, outputDir, b.URL(), "/usr/bin/chrome", 0)
        if err != nil {
            fmt.Printf("❌ Error generating config for %s: %v\n", b.Name(), err)
            continue
        }

        fmt.Printf("✅ Updated config for %s at %s/%s/config.yaml\n", b.Name(), outputDir, folderName)
    }
}

// fetchConfigs fetches all bookies and returns the full report
func fetchConfigs(bookies []utils.Bookie, outputDir string) report.FullReport {
	fmt.Println("🌐 Fetching and verifying bookies...")

	var summary []report.BookieReport
	var details []report.BookieReport

	for _, b := range bookies {
		cfgPath := filepath.Join(outputDir, strings.ToLower(b.Name()), "config.yaml")
		cfg, err := loadConfig(cfgPath)
		if err != nil {
			fmt.Printf("⚠️ Skipping %s, failed to load config: %v\n", b.Name(), err)
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

// createLatestSnippet generates latest_report.md from full report
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

// loadConfig reads a YAML config for a single bookie
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

// bakeOverridesFile merges overrides into the base config and renames the original overrides.yaml
func bakeOverridesFile(outputDir, overridesPath string) {
    // Load the overrides
    overrideMap, err := config.LoadOverrides(overridesPath)
    if err != nil {
        fmt.Printf("⚠️ Failed to load overrides file: %v\n", err)
        return
    }

    // Iterate over all bookies in overrideMap
    for bookieName, overridesForBookie := range overrideMap {
        bookieDir := filepath.Join(outputDir, strings.ToLower(bookieName))
        configPath := filepath.Join(bookieDir, "config.yaml")

        // Read existing config
        baseConfigData, err := os.ReadFile(configPath)
        if err != nil {
            fmt.Printf("⚠️ Failed to read config for %s: %v\n", bookieName, err)
            continue
        }

        // Unmarshal YAML into Sportsbook struct
        var baseConfig config.Sportsbook
        if err := yaml.Unmarshal(baseConfigData, &baseConfig); err != nil {
            fmt.Printf("⚠️ Failed to unmarshal config for %s: %v\n", bookieName, err)
            continue
        }

        // Apply overrides
        baseConfig.ApplyOverrides(overridesForBookie)

        // Marshal updated config back to YAML
        updatedConfigData, err := yaml.Marshal(baseConfig)
        if err != nil {
            fmt.Printf("⚠️ Failed to marshal updated config for %s: %v\n", bookieName, err)
            continue
        }

        // Save updated config
        if err := os.WriteFile(configPath, updatedConfigData, 0644); err != nil {
            fmt.Printf("⚠️ Failed to write updated config for %s: %v\n", bookieName, err)
            continue
        }

        fmt.Printf("✅ Baked overrides into config for %s\n", bookieName)
    }

    // Rename overrides.yaml → overrides.baked.yaml
    bakedPath := filepath.Join(outputDir, "overrides.baked.yaml")
    if err := os.Rename(overridesPath, bakedPath); err != nil {
        fmt.Printf("⚠️ Failed to rename overrides.yaml: %v\n", err)
    } else {
        fmt.Printf("🍞 All overrides baked. Original overrides.yaml renamed to %s\n", bakedPath)
    }
}

