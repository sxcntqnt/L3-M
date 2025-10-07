package bookies

import (
	"strings"
	"github.com/PuerkitoBio/goquery"
)

type OneXBet struct {
	url string
}

func (b *OneXBet) Name() string { return "1xbet" }
func (b *OneXBet) URL() string  { return b.url }
func (b *OneXBet) SetURL(u string) { b.url = u }

func (b *OneXBet) Verify(doc *goquery.Document) map[string]string {
	results := map[string]string{}
	title := doc.Find("title").Text()
	if strings.Contains(strings.ToLower(title), strings.ToLower(b.Name())) {
		results["Title contains name"] = "✅"
	} else {
		results["Title contains name"] = "❌"
	}
	if doc.Find("body").HasClass("main-page") || doc.Find(".header").Length() > 0 { // Advanced: check for common structural elements
		results["Structural element"] = "✅"
	} else {
		results["Structural element"] = "❌"
	}
	return results
}

