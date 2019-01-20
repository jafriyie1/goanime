package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
	"github.com/jafriyie1/goanime/animescrapper"
)

func main() {
	log.SetFlags(0)
	var val string
	var episodeSlice [][]string

	baseURL := "https://kissanime.ru/AnimeList/MostPopular?page="
	episodeSlice = append(episodeSlice, []string{"anime_id", "name", "genre", "type", "episodes", "ratings", "members"})

	var pageNumber int
	var strNumber string
	var url string

	ctxt, cancel := context.WithTimeout(context.Background(), 3*time.Hour)
	defer cancel()

	c, newerr := chromedp.New(ctxt, chromedp.WithRunnerOptions(

		runner.Flag("headless", true),
		runner.Flag("disable-gpu", true),
		runner.Flag("no-first-run", true),
		runner.Flag("no-sandbox", true),
		//runner.Flag("no-default-browser-check", true),
	))

	//c, newerr := chromedp.New(ctxt)
	var incr int
	incr = (9 * 19) + 1
	end := incr + 19

	file, fileErr := os.Create("../../Data/merge/test7.csv")

	if fileErr != nil {
		log.Fatal(fileErr)
	}

	for i := incr; i <= end; i++ {

		pageNumber = i
		strNumber = strconv.Itoa(pageNumber)
		url = baseURL + strNumber
		fmt.Println(url)

		if newerr != nil {
			log.Fatal(newerr)
		}
		err := c.Run(ctxt, animescrapper.ClickForEpisodeList(url, &val))
		if err != nil {
			log.Fatal(err)
		}

		stringVals := strings.NewReader(val)
		doc, _ := goquery.NewDocumentFromReader(stringVals)

		doc.Find("tbody tr").Each(func(_ int, tr *goquery.Selection) {

			// for each <tr> found, find the <td>s inside
			// ix is the index
			tr.Find("td a").Each(func(ix int, td *goquery.Selection) {

				if ix == 0 {
					val, _ := td.Attr("href")
					t := strings.Replace(val, "/Anime/", "", -1)
					finalString := strings.Replace(t, "-", " ", -1)
					episodeSlice = append(episodeSlice, []string{"", finalString, "", "", "", "", ""})
					log.Printf("index: %d content: '%s'", ix, finalString)
				}
			})
		})

		time.Sleep(time.Second * 5)
	}

	cErr := c.Shutdown(ctxt)
	if cErr != nil {
		log.Fatal("cErr")
	}

	csvwriter := csv.NewWriter(file)

	for _, val := range episodeSlice {
		if err := csvwriter.Write(val); err != nil {
			log.Fatal("Error occured during writing data to csv")
		}
	}
	csvwriter.Flush()

}
