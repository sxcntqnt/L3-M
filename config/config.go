package config

import "gopkg.in/yaml.v3"


// Sportsbook represents a single bookie's configuration
type Sportsbook struct {
	Name            string          `yaml:"name"`
	BaseURL         string          `yaml:"base_url"`
	BrowserPath     string          `yaml:"browser_path"`
	Username        string          `yaml:"username"`
	Password        string          `yaml:"password"`
	Region          string          `yaml:"region"`
	Selectors       Selectors       `yaml:"selectors"`
	BetButton       string          `yaml:"bet_button"`
	BetHistory      string          `yaml:"bet_history"`
	Timeout         Timeout         `yaml:"timeout"`
	Betting         Betting         `yaml:"betting"`
	UserCredentials UserCredentials `yaml:"user_credentials"`
}

// Selectors holds CSS selectors for login, event search, and odds
type Selectors struct {
	Login struct {
		UsernameInput string `yaml:"username_input"`
		PasswordInput string `yaml:"password_input"`
		LoginButton   string `yaml:"login_button"`
	} `yaml:"login"`

	Session struct {
		LogoutButton    string `yaml:"logout_button"`
		SessionUserInfo string `yaml:"session_user_info"`
	} `yaml:"session"`

	EventSearch struct {
		SportDropdown string `yaml:"sport_dropdown"`
		DatePicker    string `yaml:"date_picker"`
		SearchButton  string `yaml:"search_button"`
		EventResults  string `yaml:"event_results"`
		EventItem     string `yaml:"event_item"`
		EventTitle    string `yaml:"event_title"`
		EventTeam     string `yaml:"event_team"`
	} `yaml:"event_search"`

	OddsSelector struct {
		Moneyline  string `yaml:"moneyline"`
		Spread    string `yaml:"spread"`
		Totals     string `yaml:"totals"`
		OddsDropdown string `yaml:"odds_dropdown"`
	} `yaml:"odds_selector"`

	BetSlip struct {
		AddButton       string `yaml:"add_button"`
		RemoveButton    string `yaml:"remove_button"`
		StakeInput      string `yaml:"stake_input"`
		CalculateButton string `yaml:"calculate_button"`
		ClearButton     string `yaml:"clear_button"`
		PotentialPayout string `yaml:"potential_payout"`
		BetSlipItem     string `yaml:"bet_slip_item"`
	} `yaml:"bet_slip"`

	LiveBetting struct {
		LiveBettingButton   string `yaml:"live_betting_button"`
		OddsChangeIndicator string `yaml:"odds_change_indicator"`
		LiveEventItem       string `yaml:"live_event_item"`
		LiveScore           string `yaml:"live_score"`
		LiveOddSelector     string `yaml:"live_odd_selector"`
		LiveEvent           string `yaml:"live_event"`
		InPlayBetButton     string `yaml:"in_play_bet_button"`
	} `yaml:"live_betting"`

	LineMovement struct {
		LineChangeIndicator string `yaml:"line_change_indicator"`
		OddsHistory         string `yaml:"odds_history"`
		BettingLines        string `yaml:"betting_lines"`
	} `yaml:"line_movement"`

	FilterOptions struct {
		SportDropdown      string `yaml:"sport_dropdown"`
		MarketTypeDropdown string `yaml:"market_type_dropdown"`
		TimeFilter         string `yaml:"time_filter"`
		ResetFiltersButton string `yaml:"reset_filters_button"`
	} `yaml:"filter_options"`

	BetConfirmation struct {
		ConfirmButton  string `yaml:"confirm_button"`
		ErrorMessage   string `yaml:"error_message"`
		SuccessMessage string `yaml:"success_message"`
		BetSummary     string `yaml:"bet_summary"`
	} `yaml:"bet_confirmation"`

	BetHistory struct {
		HistoryPageLink string `yaml:"history_page_link"`
		BetRowSelector  string `yaml:"bet_row_selector"`
		EventColumn     string `yaml:"event_column"`
		StakeColumn     string `yaml:"stake_column"`
		OutcomeColumn   string `yaml:"outcome_column"`
		FilterByResult  string `yaml:"filter_by_result"`
		FilterByMarket  string `yaml:"filter_by_market"`
	} `yaml:"bet_history"`

	Promotions struct {
		PromotionBanner  string `yaml:"promotion_banner"`
		RedeemButton     string `yaml:"redeem_button"`
		PromoCodeInput   string `yaml:"promo_code_input"`
		ApplyPromoButton string `yaml:"apply_promo_button"`
	} `yaml:"promotions"`

	CashOut struct {
		CashOutButton           string `yaml:"cash_out_button"`
		OpenBet                 string `yaml:"open_bet"`
		CancellableBetIndicator string `yaml:"cancellable_bet_indicator"`
		CashoutOffer            string `yaml:"cashout_offer"`
		ConfirmCashoutButton    string `yaml:"confirm_cashout_button"`
	} `yaml:"cash_out"`

	NotificationCenter struct {
		NotificationPopup   string `yaml:"notification_popup"`
		DismissButton       string `yaml:"dismiss_button"`
		NotificationMessage string `yaml:"notification_message"`
		NotificationType    string `yaml:"notification_type"`
	} `yaml:"notification_center"`
}

// Timeout holds timeout settings for various operations
type Timeout struct {
	BetOperation int `yaml:"bet_operation"`
	PageLoad     int `yaml:"page_load"`
	SelectorWait int `yaml:"selector_wait"`
}

// Betting holds betting related information
type Betting struct {
	Stake   int     `yaml:"stake"`
	BetType string  `yaml:"bet_type"`
	Odds    float64 `yaml:"odds"`
	Team    string  `yaml:"team"`
	EventID string  `yaml:"event_id"`
	Query   string  `yaml:"query"`
}

// UserCredentials stores the login details
type UserCredentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// OverrideMap holds per-bookie overrides loaded from YAML
type OverrideMap map[string]map[string]interface{}

// ApplyOverrides merges overrides into a Sportsbook struct
func (sb *Sportsbook) ApplyOverrides(overrides map[string]interface{}) {
	data, _ := yaml.Marshal(overrides)
	_ = yaml.Unmarshal(data, sb)
}

