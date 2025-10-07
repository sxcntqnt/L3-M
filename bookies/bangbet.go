package bookies

import (
	"strings"
	"github.com/PuerkitoBio/goquery"
)

type BangBet struct {
	url string
}

func (b *BangBet) Name() string { return "bangbet" }
func (b *BangBet) URL() string  { return b.url }
func (b *BangBet) SetURL(u string) { b.url = u }

func (b *BangBet) Verify(doc *goquery.Document) map[string]string {
	results := map[string]string{}
	title := doc.Find("title").Text()
	if strings.Contains(strings.ToLower(title), strings.ToLower(b.Name())) {
		results["Title contains name"] = "✅"
	} else {
		results["Title contains name"] = "❌"
	}
	if doc.Find("body").HasClass("main-page") || doc.Find(".header").Length() > 0 {
		results["Structural element"] = "✅"
	} else {
		results["Structural element"] = "❌"
	}
	return results
}

