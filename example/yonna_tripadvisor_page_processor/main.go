package main

import (
	"encoding/json"
	"fmt"
	"go_spider/core/common/page"
	"go_spider/core/pipeline"
	"go_spider/core/spider"

	"gopkg.in/mgo.v2/bson"

	"github.com/PuerkitoBio/goquery"
)

var (
	startURI string
	nextURI  string
)

type TripAdvisorProcessor struct {
}

/*CModel is a commentary model structure */
type CModel struct {
	Author     string   `json:"author"`
	Title      string   `json:"title"`
	Commentary string   `json:"commentary"`
	Date       string   `"json:"date"`
	Images     []string `json:"images , omitempty`
}

/*NewTripAdvisorProcessor return new instance of TA processor  */
func NewTripAdvisorProcessor() *TripAdvisorProcessor {
	return &TripAdvisorProcessor{}
}

/*Process is life cycle spider
Parse html dom here and record the parse result that we want to Page.
Package goquery (http://godoc.org/github.com/PuerkitoBio/goquery) is used to parse html.*/
func (t *TripAdvisorProcessor) Process(p *page.Page) {

	if p.IsSucc() {
		/* Parse commentaries informations */
		query := p.GetHtmlParser()
		review := query.Find(".review-container")
		/*html, _ := review.Html()
		fmt.Println(html)
		*/page := query.Find(".ui_pagination")

		review.Each(func(i int, s *goquery.Selection) {
			c := restaurantParser(s)
			d, _ := json.Marshal(c)
			p.AddField(bson.NewObjectId().String(), string(d))
		})
		cPage := nextURI
		getNextPage(page)

		if cPage != nextURI && nextURI != "" {
			p.AddTargetRequest("https://www.tripadvisor.fr"+nextURI, "html")

		}

	}

}

/*Finish is last task of spider function */
func (t *TripAdvisorProcessor) Finish() {
	fmt.Printf(`
    (
	)
   (
/\  .-"""-.  /\
//\\/  ,,,  \//\\
|/\| ,;;;;;, |/\|
//\\\;-"""-;///\\
//  \/   .   \/  \\
(| ,-_| \ | / |_-, |)
//'__\.-.-./__'\\
// /.-(() ())-.\ \\
(\ |)   '---'   (| /)
' (|           |) '
jgs     \)           (/
	SPIDERING ENDED
				`)
}

func restaurantParser(s *goquery.Selection) CModel {

	/* Parse website informations */
	date, _ := s.Find(".ratingDate").Attr("title")
	auth := s.Find(".info_text").Text()
	com := s.Find(".partial_entry").Text()
	title := s.Find(".noQuotes").Text()
	image := s.Find(".centeredImg")

	nCom := CModel{
		Author:     auth,
		Date:       date,
		Title:      title,
		Commentary: com,
	}

	/* Get commentaries images urls */
	image.Each(func(i int, s *goquery.Selection) {
		imgURI, _ := s.Attr("data-lazyurl")
		nCom.Images = append(nCom.Images, imgURI)
	})

	return nCom
}

func hotelParser(s *goquery.Selection) CModel {

	/* Parse website informations */
	date, _ := s.Find(".ratingDate").Attr("title")
	auth := s.Find(".info_text").Text()
	com := s.Find(".partial_entry").Text()
	title := s.Find(".noQuotes").Text()
	image := s.Find(".centeredImg")

	nCom := CModel{
		Author:     auth,
		Date:       date,
		Title:      title,
		Commentary: com,
	}

	/* Get commentaries images urls */
	image.Each(func(i int, s *goquery.Selection) {
		imgURI, _ := s.Attr("data-lazyurl")
		nCom.Images = append(nCom.Images, imgURI)
	})

	return nCom
}

func getNextPage(p *goquery.Selection) {

	p.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if s.Text() == "Suivant" {
			u, _ := s.Attr("href")
			if u != "" {
				nextURI = u
				return false

			}
			if u != startURI {
				nextURI = u
				return false
			}
		}
		return true
	})
}

func main() {
	startURI = "https://www.tripadvisor.fr/Hotel_Review-g3367394-d1235967-Reviews-B_B_Hotel_le_Mans_Nord_1-Saint_Saturnin_Le_Mans_Sarthe_Pays_de_la_Loire.html"

	spider.NewSpider(NewTripAdvisorProcessor(), "spyders").
		AddUrl(startURI, "html").

		//AddPipeline(pipeline.NewPipelineConsole()).            // Print result on screen
		AddPipeline(pipeline.NewPipelineFile("./commentary")). // Print result in file
		OpenFileLog("./tmp").                                  // Error info or other useful info in spider will be logged in file of defalt path like "WD/log/log.2014-9-1".

		SetThreadnum(100). // Crawl request by three Coroutines
		Run()
}
