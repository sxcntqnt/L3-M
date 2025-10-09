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


// buildSportsbook constructs the Sportsbook struct with base settings and optional override
func buildSportsbook(bookie string, index int, baseURL, browserPath string, overrides OverrideMap) Sportsbook {
	sb := Sportsbook{
		Name:      bookie,
		BaseURL:   baseURL,
		BrowserPath: browserPath,
		Username: fmt.Sprintf("user%d", index+1),
		Password: fmt.Sprintf("pass%d", index+1),
		Region:   "KE",
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

	if ovr, ok := overrides[bookie]; ok {
		sb.ApplyOverrides(ovr)
	}

	return sb
}

// GenerateConfig builds a YAML config for a single bookie with optional override
func GenerateConfig(bookie string, overrides OverrideMap, outputDir, baseURL, browserPath string, index int) error {
	bookieDir := filepath.Join(outputDir, strings.ToLower(bookie))

	if err := os.MkdirAll(bookieDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create bookie dir: %w", err)
	}

	sb := buildSportsbook(bookie, index, baseURL, browserPath, overrides)

	// Encode YAML
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	if err := encoder.Encode(sb); err != nil {
		return fmt.Errorf("error encoding YAML: %w", err)
	}

	// Check for changes
	configPath := filepath.Join(bookieDir, "config.yaml")
	existing, err := os.ReadFile(configPath)

	if err == nil && bytes.Equal(existing, buf.Bytes()) {
		fmt.Printf("⏭️ No changes for %s, skipped.\n", bookie)
		return nil
	}

	// Write config file
	if err := os.WriteFile(configPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing YAML: %w", err)
	}

	fmt.Printf("✅ Updated config for %s at %s\n", bookie, configPath)
	return nil
}

// WriteMetadata writes pipeline metadata to output directory
func WriteMetadata(outputDir string) error {
	meta := map[string]string{
		"generated_at": time.Now().Format(time.RFC3339),
		"pipeline_id":  os.Getenv("CI_PIPELINE_ID"),
		"run_id":       os.Getenv("GITHUB_RUN_ID"),
		"commit_sha":   os.Getenv("GITHUB_SHA"),
	}
	metaPath := filepath.Join(outputDir, "metadata.yaml")

	var metaBuf bytes.Buffer
	if err := yaml.NewEncoder(&metaBuf).Encode(meta); err != nil {
		return fmt.Errorf("failed to encode metadata: %w", err)
	}

	if err := os.WriteFile(metaPath, metaBuf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

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

