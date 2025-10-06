package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Sportsbook structure representing each sportsbook configuration
type Sportsbook struct {
	Name       string `yaml:"name"`
	BaseURL    string `yaml:"base_url"`
	BrowserPath string `yaml:"browser_path"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Region     string `yaml:"region"`
	Selectors  struct {
		Login struct {
			UsernameInput string `yaml:"username_input"`
			PasswordInput string `yaml:"password_input"`
			LoginButton   string `yaml:"login_button"`
		} `yaml:"login"`
		EventSearch struct {
			SportDropdown string `yaml:"sport_dropdown"`
			DatePicker    string `yaml:"date_picker"`
			SearchButton  string `yaml:"search_button"`
		} `yaml:"event_search"`
		OddsSelector struct {
			Moneyline string `yaml:"moneyline"`
			Spread    string `yaml:"spread"`
			Totals    string `yaml:"totals"`
		} `yaml:"odds_selector"`
		BetButton   string `yaml:"bet_button"`
		BetHistory  string `yaml:"bet_history"`
		LiveBetting struct {
			LiveEvent        string `yaml:"live_event"`
			InPlayBetButton  string `yaml:"in_play_bet_button"`
		} `yaml:"live_betting"`
	} `yaml:"selectors"`
}

type Config struct {
	Sportsbooks []Sportsbook `yaml:"sportsbooks"`
}

func generateConfig(numSportsbooks int, names []string) {
	// Directory to store the config file
	directory := "EMC"

	// Check if directory exists, if not, create it
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, os.ModePerm); err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}

	// Default base URL, browser path, and credentials
	baseURL := "https://sportsbook%d.example.com"
	browserPath := "/path/to/%s"
	username := "user%d"
	password := "pass%d"
	region := "US"
	if numSportsbooks%2 != 0 {
		region = "EU" // Alternate region for fun
	}

	// If names are provided, use them; otherwise, generate default names
	if names == nil {
		names = make([]string, numSportsbooks)
		for i := range names {
			names[i] = fmt.Sprintf("Sportsbook %d", i+1)
		}
	}

	// List to hold the sportsbook configurations
	var sportsbooks []Sportsbook

	// Generate the configuration for each sportsbook
	for i := 0; i < numSportsbooks; i++ {
		sportbook := Sportsbook{
			Name:       names[i],
			BaseURL:    fmt.Sprintf(baseURL, i+1),
			BrowserPath: fmt.Sprintf(browserPath, "chrome"),
			Username:   fmt.Sprintf(username, i+1),
			Password:   fmt.Sprintf(password, i+1),
			Region:     region,
		}
		// Set selectors for this sportsbook
		sportbook.Selectors.Login.UsernameInput = fmt.Sprintf("input#username%d", i+1)
		sportbook.Selectors.Login.PasswordInput = fmt.Sprintf("input#password%d", i+1)
		sportbook.Selectors.Login.LoginButton = fmt.Sprintf("button#login%d", i+1)
		sportbook.Selectors.EventSearch.SportDropdown = fmt.Sprintf("select#sport%d", i+1)
		sportbook.Selectors.EventSearch.DatePicker = fmt.Sprintf("input#date%d", i+1)
		sportbook.Selectors.EventSearch.SearchButton = fmt.Sprintf("button#search%d", i+1)
		sportbook.Selectors.OddsSelector.Moneyline = fmt.Sprintf("div.moneyline%d", i+1)
		sportbook.Selectors.OddsSelector.Spread = fmt.Sprintf("div.spread%d", i+1)
		sportbook.Selectors.OddsSelector.Totals = fmt.Sprintf("div.totals%d", i+1)
		sportbook.Selectors.BetButton = fmt.Sprintf("button#placeBet%d", i+1)
		sportbook.Selectors.BetHistory = fmt.Sprintf("div#betHistory%d", i+1)
		sportbook.Selectors.LiveBetting.LiveEvent = fmt.Sprintf("div.live-event%d", i+1)
		sportbook.Selectors.LiveBetting.InPlayBetButton = fmt.Sprintf("button#inPlayBet%d", i+1)

		// Append the sportsbook configuration to the list
		sportsbooks = append(sportsbooks, sportbook)
	}

	// Create the full config structure
	configData := Config{
		Sportsbooks: sportsbooks,
	}

	// Define the path to save the config file
	configPath := fmt.Sprintf("%s/config.yaml", directory)

	// Write to YAML file in the EMC directory
	file, err := os.Create(configPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	if err := encoder.Encode(configData); err != nil {
		fmt.Println("Error writing YAML data:", err)
		return
	}

	// Inform the user of success
	fmt.Printf("%d sportsbook configurations have been generated successfully and saved to %s!\n", numSportsbooks, configPath)
}

func main() {
	// Example: generate 3 sportsbook configurations with custom names
	generateConfig(3, []string{"Bookie A", "Bookie B", "Bookie C"})
}

