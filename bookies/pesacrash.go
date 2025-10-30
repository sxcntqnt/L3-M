package bookies

import (
	"strings"
	"github.com/PuerkitoBio/goquery"
)

type PesaCrash struct {
	url string
}

func (b *PesaCrash) Name() string { return "pesacrash" }
func (b *PesaCrash) URL() string  { return b.url }
func (b *PesaCrash) SetURL(u string) { b.url = u }

func (b *PesaCrash) Verify(doc *goquery.Document) map[string]string {
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

