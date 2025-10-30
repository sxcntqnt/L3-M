package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-yaml/yaml/v3"
)

// buildSportsbook constructs the Sportsbook struct with base settings and optional override
func buildSportsbook(bookie string, index int, baseURL, browserPath string, overrides OverrideMap) Sportsbook {
	sb := Sportsbook{
		Name:        bookie,
		BaseURL:     baseURL,
		BrowserPath: browserPath,
		Username:    fmt.Sprintf("user%d", index+1),
		Password:    fmt.Sprintf("pass%d", index+1),
		Region:      "KE", // Default region
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
			Session: struct {
				LogoutButton    string `yaml:"logout_button"`
				SessionUserInfo string `yaml:"session_user_info"`
			}{
				LogoutButton:    "button#logout",
				SessionUserInfo: "div#userInfo",
			},
			EventSearch: struct {
				SportDropdown string `yaml:"sport_dropdown"`
				DatePicker    string `yaml:"date_picker"`
				SearchButton  string `yaml:"search_button"`
				EventResults  string `yaml:"event_results"`
				EventItem     string `yaml:"event_item"`
				EventTitle    string `yaml:"event_title"`
				EventTeam     string `yaml:"event_team"`
			}{
				SportDropdown: "select#sport",
				DatePicker:    "input#date",
				SearchButton:  "button#search",
				EventResults:  "div#eventResults",
				EventItem:     "div.event-item",
				EventTitle:    "div.event-title",
				EventTeam:     "div.event-team",
			},
			OddsSelector: struct {
				Moneyline  string `yaml:"moneyline"`
				Spread    string `yaml:"spread"`
				Totals     string `yaml:"totals"`
				OddsDropdown string `yaml:"odds_dropdown"`
			}{
				Moneyline:  "div.match-result",
				Spread:     "div.over-under",
				Totals:      "div.point-spread",
				OddsDropdown: "select#odds",
			},
			BetSlip: struct {
				AddButton       string `yaml:"add_button"`
				RemoveButton    string `yaml:"remove_button"`
				StakeInput      string `yaml:"stake_input"`
				CalculateButton string `yaml:"calculate_button"`
				ClearButton     string `yaml:"clear_button"`
				PotentialPayout string `yaml:"potential_payout"`
				BetSlipItem     string `yaml:"bet_slip_item"`
			}{
				AddButton:       "button#addToBetSlip",
				RemoveButton:    "button#removeBetSlipItem",
				StakeInput:      "input#stake",
				CalculateButton: "button#calculate",
				ClearButton:     "button#clearBetSlip",
				PotentialPayout: "div#potentialPayout",
				BetSlipItem:     "div.bet-slip-item",
			},
			LiveBetting: struct {
				LiveBettingButton   string `yaml:"live_betting_button"`
				OddsChangeIndicator string `yaml:"odds_change_indicator"`
				LiveEventItem       string `yaml:"live_event_item"`
				LiveScore           string `yaml:"live_score"`
				LiveOddSelector     string `yaml:"live_odd_selector"`
				LiveEvent           string `yaml:"live_event"`
				InPlayBetButton     string `yaml:"in_play_bet_button"`
			}{
				LiveBettingButton:   "button#liveBetting",
				OddsChangeIndicator: "div.odds-change-indicator",
				LiveEventItem:       "div.live-event-item",
				LiveScore:           "div.live-score",
				LiveOddSelector:     "div.live-odd-selector",
				LiveEvent:           "div.live-event",
			},

			LineMovement: struct {
				LineChangeIndicator string `yaml:"line_change_indicator"`
				OddsHistory         string `yaml:"odds_history"`
				BettingLines        string `yaml:"betting_lines"`
			}{
				LineChangeIndicator: "div.line-change-indicator",
				OddsHistory:         "div.odds-history",
				BettingLines:        "div.betting-lines",
			},
			FilterOptions: struct {
				SportDropdown      string `yaml:"sport_dropdown"`
				MarketTypeDropdown string `yaml:"market_type_dropdown"`
				TimeFilter         string `yaml:"time_filter"`
				ResetFiltersButton string `yaml:"reset_filters_button"`
			}{
				SportDropdown:      "select#sportFilter",
				MarketTypeDropdown: "select#marketTypeFilter",
				TimeFilter:         "input#timeFilter",
				ResetFiltersButton: "button#resetFilters",
			},
			BetConfirmation: struct {
				ConfirmButton  string `yaml:"confirm_button"`
				ErrorMessage   string `yaml:"error_message"`
				SuccessMessage string `yaml:"success_message"`
				BetSummary     string `yaml:"bet_summary"`
			}{
				ConfirmButton:  "button#confirmBet",
				ErrorMessage:   "div#errorMessage",
				SuccessMessage: "div#successMessage",
				BetSummary:     "div#betSummary",
			},
			BetHistory: struct {
				HistoryPageLink string `yaml:"history_page_link"`
				BetRowSelector  string `yaml:"bet_row_selector"`
				EventColumn     string `yaml:"event_column"`
				StakeColumn     string `yaml:"stake_column"`
				OutcomeColumn   string `yaml:"outcome_column"`
				FilterByResult  string `yaml:"filter_by_result"`
				FilterByMarket  string `yaml:"filter_by_market"`
			}{
				HistoryPageLink: "a#betHistoryLink",
				BetRowSelector:  "div.bet-row",
				EventColumn:     "div.bet-row .event",
				StakeColumn:     "div.bet-row .stake",
				OutcomeColumn:   "div.bet-row .outcome",
				FilterByResult:  "select#filterByResult",
				FilterByMarket:  "select#filterByMarket",
			},
			Promotions: struct {
				PromotionBanner  string `yaml:"promotion_banner"`
				RedeemButton     string `yaml:"redeem_button"`
				PromoCodeInput   string `yaml:"promo_code_input"`
				ApplyPromoButton string `yaml:"apply_promo_button"`
			}{
				PromotionBanner:  "div#promotionBanner",
				RedeemButton:     "button#redeemPromo",
				PromoCodeInput:   "input#promoCode",
				ApplyPromoButton: "button#applyPromo",
			},
			CashOut: struct {
				CashOutButton           string `yaml:"cash_out_button"`
				OpenBet                 string `yaml:"open_bet"`
				CancellableBetIndicator string `yaml:"cancellable_bet_indicator"`
				CashoutOffer            string `yaml:"cashout_offer"`
				ConfirmCashoutButton    string `yaml:"confirm_cashout_button"`
			}{
				CashOutButton:           "button#cashOut",
				OpenBet:                 "div.open-bet",
				CancellableBetIndicator: "div.cancellable-bet",
				CashoutOffer:            "div.cashout-offer",
				ConfirmCashoutButton:    "button#confirmCashout",
			},
			NotificationCenter: struct {
				NotificationPopup   string `yaml:"notification_popup"`
				DismissButton       string `yaml:"dismiss_button"`
				NotificationMessage string `yaml:"notification_message"`
				NotificationType    string `yaml:"notification_type"`
			}{
				NotificationPopup:   "div#notificationPopup",
				DismissButton:       "button#dismissNotification",
				NotificationMessage: "div.notification-message",
				NotificationType:    "div.notification-type",
			},
		},
		BetButton:  "button#placeBet",
		BetHistory: "div#betHistory",
		Timeout: Timeout{
			BetOperation: 30000, // 30 seconds timeout for placing a bet
			PageLoad:     5000,  // 5 seconds timeout for page loading
			SelectorWait: 5000,  // 5 seconds timeout for waiting for selectors
		},
		Betting: Betting{
			Stake:   100,
			BetType: "match_result",
			Odds:    2.5,
			Team:    "Team A",
			EventID: "12345",
			Query:   "2023-10-15",
		},
		UserCredentials: UserCredentials{
			Username: fmt.Sprintf("user%d", index+1),
			Password: fmt.Sprintf("pass%d", index+1),
		},
	}

	// Apply any overrides from the provided override map
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
