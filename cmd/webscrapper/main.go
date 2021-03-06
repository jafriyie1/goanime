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

	filePath := fmt.Sprintf("../../Data/test/episodes.csv")
	file, fileErr := os.Create(filePath)
	for j := 0; j <= 11; j++ {
		//ctxt, cancel := context.WithTimeout(context.Background(), 3*time.Hour)
		ctxt, cancel := chromedp.NewContext(context.Background())
		ctxt, cancel = context.WithTimeout(ctxt, 180*time.Second)
		defer cancel()

		var incr int
		incr = (j * 19) + 1
		end := incr + 19

		if fileErr != nil {
			log.Fatal(fileErr)
		}

		for i := incr; i <= end; i++ {

			pageNumber = i
			strNumber = strconv.Itoa(pageNumber)
			url = baseURL + strNumber
			fmt.Println(url)

			/*
				if newerr != nil {
					log.Fatal(newerr)
				}
			*/
			err := chromedp.Run(ctxt, animescrapper.ClickForEpisodeList(url, &val))
			if err != nil {
				csvwriter := csv.NewWriter(file)

				for _, val := range episodeSlice {
					if err := csvwriter.Write(val); err != nil {
						log.Fatal("Error occured during writing data to csv")
					}
				}
				csvwriter.Flush()
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

	}

	csvwriter := csv.NewWriter(file)

	for _, val := range episodeSlice {
		if err := csvwriter.Write(val); err != nil {
			log.Fatal("Error occured during writing data to csv")
		}
	}
	csvwriter.Flush()

}
