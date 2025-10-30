package bookies

import (
	"strings"
	"github.com/PuerkitoBio/goquery"
)

type Mcheza struct {
	url string
}

func (b *Mcheza) Name() string { return "mcheza" }
func (b *Mcheza) URL() string  { return b.url }
func (b *Mcheza) SetURL(u string) { b.url = u }

func (b *Mcheza) Verify(doc *goquery.Document) map[string]string {
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

