package fetch

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"time"

	"diago/config"
	"diago/report"

	"github.com/PuerkitoBio/goquery"
)

// FetchPage fetches a URL and returns a parsed goquery document.
func FetchPage(urlStr string) (*goquery.Document, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL %q: %w", urlStr, err)
	}
	parsedURL.Fragment = ""

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %q: %w", parsedURL.String(), err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; FetchBot/2.0; +https://yourdomain.com/bot)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %q: %w", parsedURL.String(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("received HTTP %d for %q", resp.StatusCode, parsedURL.String())
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response from %q: %w", parsedURL.String(), err)
	}

	return doc, nil
}

// VerifyBookieWithConfig checks all selectors dynamically from config.Sportsbook
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

	// Traverse all selector fields dynamically via reflection
	traverseSelectors(reflect.ValueOf(cfg.Selectors), "", doc, &results, &allPass)

	// Also check top-level selectors if present
	topLevelChecks := map[string]string{
		"BetButton":  cfg.BetButton,
		"BetHistory": cfg.BetHistory,
	}
	for label, selector := range topLevelChecks {
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

// traverseSelectors recursively inspects nested structs and verifies each selector
func traverseSelectors(v reflect.Value, prefix string, doc *goquery.Document, results *[]report.SelectorResult, allPass *bool) {
	v = reflect.Indirect(v)
	if !v.IsValid() {
		return
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		label := fieldType.Name

		// Build hierarchical labels like "Login.UsernameInput"
		fullLabel := label
		if prefix != "" {
			fullLabel = prefix + "." + label
		}

		switch field.Kind() {
		case reflect.String:
			selector := field.String()
			if selector == "" {
				continue
			}
			if doc.Find(selector).Length() > 0 {
				*results = append(*results, report.SelectorResult{Label: fullLabel, Status: "✅"})
			} else {
				*results = append(*results, report.SelectorResult{Label: fullLabel, Status: "❌"})
				*allPass = false
			}
		case reflect.Struct:
			traverseSelectors(field, fullLabel, doc, results, allPass)
		}
	}
}

// VerifyBookiesConcurrently fetches multiple bookies concurrently.
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
	for r := range resultsCh {
		summary = append(summary, r)
	}

	return report.FullReport{
		Summary: summary,
		Details: summary,
	}
}
