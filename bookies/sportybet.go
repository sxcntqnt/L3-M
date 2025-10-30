package bookies

import (
	"strings"
        "diago/utils"
	"github.com/PuerkitoBio/goquery"
)

type SportyBet struct {
	url string
}

func (b *SportyBet) Name() string { return "sportybet" }
func (b *SportyBet) URL() string  { return b.url }
func (b *SportyBet) SetURL(u string) { b.url = u }

func (b *SportyBet) Verify(doc *goquery.Document) map[string]string {
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

func init() {
	utils.Register(&OneXBet{})
	utils.Register(&Betika{})
	utils.Register(&BetWinner{})
	utils.Register(&PesaCrash{})
	utils.Register(&InstaBets{})
	utils.Register(&Mcheza{})
	utils.Register(&PepetaBet{})
	utils.Register(&StarBet{})
	utils.Register(&TwentyTwoBet{})
	utils.Register(&BetKing{})
	utils.Register(&BolYeSports{})
	utils.Register(&Forzza{})
	utils.Register(&JamboBet{})
	utils.Register(&MelBet{})
	utils.Register(&PesaLand{})
	utils.Register(&Shabiki{})
	utils.Register(&EightEightyEightStarz{})
	utils.Register(&BetKwiff{})
	utils.Register(&BongoBongo{})
	utils.Register(&GameGuys{})
	utils.Register(&JantaBets{})
	utils.Register(&MojaBet{})
	utils.Register(&Pinnacle{})
	utils.Register(&StrikeBet{})
	utils.Register(&BangBet{})
	utils.Register(&BetLion{})
	utils.Register(&CaptainsBet{})
	utils.Register(&GeniusBet{})
	utils.Register(&KenyaCharity{})
	utils.Register(&MossBets{})
	utils.Register(&Pitch90Bet{})
	utils.Register(&SokaBet{})
	utils.Register(&MegaPari{})
	utils.Register(&BetaFriq{})
	utils.Register(&BetNare{})
	utils.Register(&ChezaCash{})
	utils.Register(&HelaBet{})
	utils.Register(&KiliBet{})
	utils.Register(&MozzartBet{})
	utils.Register(&PlayBet{})
	utils.Register(&SolBet{})
	utils.Register(&UltraBet{})
	utils.Register(&BetBureau{})
	utils.Register(&BetPawa{})
	utils.Register(&HollywoodBets{})
	utils.Register(&Kwachua{})
	utils.Register(&OdDiBet{})
	utils.Register(&PlayMaster{})
	utils.Register(&Sportika{})
	utils.Register(&WorldSportBetting{})
	utils.Register(&BetFlame{})
	utils.Register(&BetSafe{})
	utils.Register(&DafaBet{})
	utils.Register(&IBet{})
	utils.Register(&KwikBet{})
	utils.Register(&PalmsBet{})
	utils.Register(&SportPesa{})
	utils.Register(&BetGr8{})
	utils.Register(&BetWay{})
	utils.Register(&DimbaKenya{})
	utils.Register(&InBetKenya{})
	utils.Register(&LigiBet{})
	utils.Register(&PariMatch{})
	utils.Register(&SaharaGames{})
	utils.Register(&SportyBet{})
}
