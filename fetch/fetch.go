package fetch

import (
	"fmt"
	"sync"
        "time"

	"diago/config"
	"diago/report"

	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
)

// FetchPage fetches the page and returns a goquery document
// FetchPage fetches the page from the given URL and returns a goquery document.
func FetchPage(urlStr string) (*goquery.Document, error) {
    // Parse and clean URL
    parsedURL, err := url.Parse(urlStr)
    if err != nil {
        return nil, fmt.Errorf("failed to parse URL %q: %w", urlStr, err)
    }
    parsedURL.Fragment = ""

    // Custom HTTP client with timeout
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    // Create HTTP request with headers (some sites require a user-agent)
    req, err := http.NewRequest("GET", parsedURL.String(), nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request for %q: %w", parsedURL.String(), err)
    }
    req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; FetchBot/1.0; +https://yourdomain.com/bot)")

    // Execute the request
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch %q: %w", parsedURL.String(), err)
    }
    defer resp.Body.Close()

    // Check for HTTP errors
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return nil, fmt.Errorf("received HTTP %d for %q", resp.StatusCode, parsedURL.String())
    }

    // Parse the response body into a goquery document
    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to parse response from %q: %w", parsedURL.String(), err)
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
		"LiveEvent":     cfg.Selectors.LiveBetting.LiveEvent,
		"InPlayBetBtn":  cfg.Selectors.LiveBetting.InPlayBetButton,
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
