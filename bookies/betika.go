package bookies

import (
	"strings"
	"github.com/PuerkitoBio/goquery"
)

type Betika struct {
	url string
}

func (b *Betika) Name() string { return "betika" }
func (b *Betika) URL() string  { return b.url }
func (b *Betika) SetURL(u string) { b.url = u }

func (b *Betika) Verify(doc *goquery.Document) map[string]string {
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

