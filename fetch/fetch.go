package fetch

import (
	"fmt"
	"sync"

	"diago/config"
	"diago/report"

	"github.com/PuerkitoBio/goquery"
	"net/http"
)

// FetchPage fetches the page and returns a goquery document
func FetchPage(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", url, err)
	}
	return doc, nil
}

// VerifyBookieWithConfig checks all selectors from a config.Sportsbook
func VerifyBookieWithConfig(name, url string, cfg *config.Sportsbook) report.BookieReport {
	fmt.Printf("🔍 Checking %s at %s...\n", name, url)

	doc, err := FetchPage(cfg.BaseURL)
	if err != nil {
		return report.BookieReport{
			Name:    name,
			URL:     cfg.BaseURL,
			AllPass: false,
			Results: []report.SelectorResult{{Label: "Fetch error", Status: err.Error()}},
		}
	}

	results := []report.SelectorResult{}
	allPass := true

	checks := map[string]string{
		"UsernameInput": cfg.Selectors.Login.UsernameInput,
		"PasswordInput": cfg.Selectors.Login.PasswordInput,
		"LoginButton":   cfg.Selectors.Login.LoginButton,
		"SportDropdown": cfg.Selectors.EventSearch.SportDropdown,
		"DatePicker":    cfg.Selectors.EventSearch.DatePicker,
		"SearchButton":  cfg.Selectors.EventSearch.SearchButton,
		"Moneyline":     cfg.Selectors.OddsSelector.Moneyline,
		"Spread":        cfg.Selectors.OddsSelector.Spread,
		"Totals":        cfg.Selectors.OddsSelector.Totals,
		"BetButton":     cfg.BetButton,
		"BetHistory":    cfg.BetHistory,
		"LiveEvent":     cfg.LiveBetting.LiveEvent,
		"InPlayBetBtn":  cfg.LiveBetting.InPlayBetButton,
	}

	for label, selector := range checks {
		if selector == "" {
			continue
		}
		if doc.Find(selector).Length() > 0 {
			results = append(results, report.SelectorResult{Label: label, Status: "✅"})
		} else {
			results = append(results, report.SelectorResult{Label: label, Status: "❌"})
			allPass = false
		}
	}

	return report.BookieReport{
		Name:    name,
		URL:     cfg.BaseURL,
		AllPass: allPass,
		Results: results,
	}
}

// VerifyBookiesConcurrently fetches multiple bookies concurrently
func VerifyBookiesConcurrently(bookies []*config.Sportsbook) report.FullReport {
	var wg sync.WaitGroup
	resultsCh := make(chan report.BookieReport, len(bookies))

	for _, b := range bookies {
		wg.Add(1)
		go func(sb *config.Sportsbook) {
			defer wg.Done()
			r := VerifyBookieWithConfig(sb.Name, sb.BaseURL, sb)
			resultsCh <- r
		}(b)
	}

	wg.Wait()
	close(resultsCh)

	var summary []report.BookieReport
	var details []report.BookieReport
	for r := range resultsCh {
		details = append(details, r)
		summary = append(summary, r)
	}

	return report.FullReport{
		Summary: summary,
		Details: details,
	}
}

