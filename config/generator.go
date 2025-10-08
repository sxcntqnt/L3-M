package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// GenerateConfigs builds YAML configs for each bookie with optional overrides
func GenerateConfigs(bookies []string, overrides OverrideMap, outputDir, baseURL, browserPath string) error {
	for i, bookie := range bookies {
		bookieDir := filepath.Join(outputDir, strings.ToLower(bookie))
		if err := os.MkdirAll(bookieDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create bookie dir: %w", err)
		}

		sb := Sportsbook{
			Name:        bookie,
			BaseURL:     baseURL,
			BrowserPath: browserPath,
			Username:    fmt.Sprintf("user%d", i+1),
			Password:    fmt.Sprintf("pass%d", i+1),
			Region:      "KE",
			Selectors: Selectors{
				Login: struct {
					UsernameInput string `yaml:"username_input"`
					PasswordInput string `yaml:"password_input"`
					LoginButton   string `yaml:"login_button"`
				}{
					UsernameInput: "input#username",
					PasswordInput: "input#password",
					LoginButton:   "button#login",
				},
				EventSearch: struct {
					SportDropdown string `yaml:"sport_dropdown"`
					DatePicker    string `yaml:"date_picker"`
					SearchButton  string `yaml:"search_button"`
				}{
					SportDropdown: "select#sport",
					DatePicker:    "input#date",
					SearchButton:  "button#search",
				},
				OddsSelector: struct {
					Moneyline string `yaml:"moneyline"`
					Spread    string `yaml:"spread"`
					Totals    string `yaml:"totals"`
				}{
					Moneyline: "div.moneyline",
					Spread:    "div.spread",
					Totals:    "div.totals",
				},
			},
			BetButton:   "button#placeBet",
			BetHistory:  "div#betHistory",
			LiveBetting: LiveBetting{LiveEvent: "div.live-event", InPlayBetButton: "button#inPlayBet"},
		}

		// Apply overrides if any
		if ovr, ok := overrides[bookie]; ok {
			sb.ApplyOverrides(ovr)
		}

		// Encode YAML into memory
		var buf bytes.Buffer
		encoder := yaml.NewEncoder(&buf)
		encoder.SetIndent(2)
		if err := encoder.Encode(sb); err != nil {
			fmt.Println("❌ Error encoding YAML:", err)
			continue
		}

		// Write YAML incrementally
		configPath := filepath.Join(bookieDir, "config.yaml")
		existing, err := os.ReadFile(configPath)
		if err == nil && bytes.Equal(existing, buf.Bytes()) {
			fmt.Printf("⏭️  No changes for %s, skipped.\n", bookie)
			continue
		}

		if err := os.WriteFile(configPath, buf.Bytes(), 0644); err != nil {
			fmt.Println("❌ Error writing YAML:", err)
			continue
		}

		fmt.Printf("✅ Updated config for %s at %s\n", bookie, configPath)
	}

	// Write metadata
	meta := map[string]string{
		"generated_at": time.Now().Format(time.RFC3339),
		"pipeline_id":  os.Getenv("CI_PIPELINE_ID"),
		"run_id":       os.Getenv("GITHUB_RUN_ID"),
		"commit_sha":   os.Getenv("GITHUB_SHA"),
	}
	metaPath := filepath.Join(outputDir, "metadata.yaml")
	var metaBuf bytes.Buffer
	_ = yaml.NewEncoder(&metaBuf).Encode(meta)
	_ = os.WriteFile(metaPath, metaBuf.Bytes(), 0644)
	fmt.Printf("📝 Metadata updated at %s\n", metaPath)

	return nil
}

// LoadOverrides loads YAML overrides from file
func LoadOverrides(path string) (OverrideMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ovr := OverrideMap{}
	if err := yaml.Unmarshal(data, &ovr); err != nil {
		return nil, err
	}
	return ovr, nil
}

