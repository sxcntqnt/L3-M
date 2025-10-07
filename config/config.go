package config

import "gopkg.in/yaml.v3"

// Sportsbook represents a single bookie's configuration
type Sportsbook struct {
	Name        string      `yaml:"name"`
	BaseURL     string      `yaml:"base_url"`
	BrowserPath string      `yaml:"browser_path"`
	Username    string      `yaml:"username"`
	Password    string      `yaml:"password"`
	Region      string      `yaml:"region"`
	Selectors   Selectors   `yaml:"selectors"`
	BetButton   string      `yaml:"bet_button"`
	BetHistory  string      `yaml:"bet_history"`
	LiveBetting LiveBetting `yaml:"live_betting"`
}

// Selectors holds CSS selectors for login, event search, and odds
type Selectors struct {
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
}

// LiveBetting holds selectors for live betting section
type LiveBetting struct {
	LiveEvent       string `yaml:"live_event"`
	InPlayBetButton string `yaml:"in_play_bet_button"`
}

// OverrideMap holds per-bookie overrides loaded from YAML
type OverrideMap map[string]map[string]interface{}

// ApplyOverrides merges overrides into a Sportsbook struct
func (sb *Sportsbook) ApplyOverrides(overrides map[string]interface{}) {
	data, _ := yaml.Marshal(overrides)
	_ = yaml.Unmarshal(data, sb)
}

