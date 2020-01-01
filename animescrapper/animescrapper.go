// This program allows one to watch anime
// without having to go through the hoops
// of using popular anime sites

package animescrapper

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"

	"github.com/derekparker/trie"
	"github.com/jafriyie1/goanime/animetries"
)

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func GetShow(b *trie.Trie) string {
	var show string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please input anime name.\nYou can also input a part of the show to search.")

	if scanner.Scan() {
		show = scanner.Text()
	}

	show = strings.TrimSpace(show)

	upperCase := byte(unicode.ToUpper(rune(show[0])))
	upperLine := string(upperCase) + show[1:]

	fmt.Println("Here are some shows that match your search:")
	fmt.Println()
	animetries.PossibleShows(b, upperLine)
	fmt.Println()
	fmt.Println("Please copy or type the anime show and hit enter")
	if scanner.Scan() {
		show = scanner.Text()
	}
	return show

}

func GetSeason() string {
	var season string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Which season would you like to watch (1st, 2nd, 3rd, 4th, etc.)?\n(Hit enter if there is no season) ")
	if scanner.Scan() {
		season = scanner.Text()
	}
	return season
}

func GetURL(show, episode, season string) (string, string, string) {

	season = strings.TrimSpace(season)

	episodeToInt, err := strconv.ParseInt(episode, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	if episodeToInt < 100 {
		episode = "0" + episode
	}
	if episodeToInt < 10 {
		episode = "0" + episode
	}
	if season == "1st" {
		season = ""
	}
	if season != "" {
		season = "-" + season + "-Season"
	}
	show = strings.Replace(show, ":", "", -1)
	show = strings.Replace(show, ")", "", -1)
	show = strings.Replace(show, "(", "", -1)
	show = strings.Replace(show, " ", "-", -1)
	base_url := "https://kissanime.ru/Anime/" + show + season + "/"

	episode = "Episode-" + episode

	return show, base_url, episode
}

func GetOneEpisode() string {
	var episode string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Please input episode number")
	if scanner.Scan() {
		episode = scanner.Text()
	}

	episode = strings.TrimSpace(episode)
	return episode

}

func GetRangeOfEpisodes(option bool) (string, string) {
	var episodeStart string
	var episodeEnd string

	scanner := bufio.NewScanner(os.Stdin)

	if option == false {
		fmt.Println("Which episode would you like to watch?\nNote: do not place any zeros in front of the episode (i.e. type 25 for Episode 025)")
		fmt.Println("Please input episode number")

		if scanner.Scan() {
			episodeStart = scanner.Text()
		}
		episodeStart = strings.TrimSpace(episodeStart)
		episodeEnd = episodeStart

	} else {
		fmt.Println("What range of episodes would you like to watch?")
		fmt.Println("Please input starting episode number")
		if scanner.Scan() {
			episodeStart = scanner.Text()
		}
		fmt.Println("Please input ending episode number")
		if scanner.Scan() {
			episodeEnd = scanner.Text()
		}

		episodeStart = strings.TrimSpace(episodeStart)
		episodeEnd = strings.TrimSpace(episodeEnd)
	}
	//fmt.Println("Start episode", episodeStart)

	return episodeStart, episodeEnd
}

func Click(url string, val *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		//chromedp.WaitVisible(`#footer`),
		chromedp.Click(`a.specialButton`, chromedp.NodeVisible),
		chromedp.Click(`a.specialButton`, chromedp.NodeVisible),
		//chromedp.WaitVisible(`div#divMyVideo`),
		chromedp.OuterHTML(`iframe#my_video_1`, val, chromedp.NodeVisible),
	}
}

func ClickForEpisodeList(url string, val *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#head`),
		chromedp.OuterHTML(`table.listing`, val, chromedp.NodeVisible),
	}
}

func OpenEpisodesAPI(lowerLimitEpisode string, searchedShow string) string {

	//ctxt, cancel := DoGoAnime()

	ctxt, cancel := chromedp.NewContext(context.Background())
	ctxt, cancel = context.WithTimeout(ctxt, 180*time.Second)
	defer cancel()

	defer cancel()
	var val string

	_, baseURL, episode := GetURL(searchedShow, lowerLimitEpisode, "")

	episodeSearch := baseURL + episode + "?id=&s=rapidVideo"

	// run task list
	log.SetFlags(0)

	err := chromedp.Run(ctxt, Click(episodeSearch, &val))
	if err != nil {
		fmt.Println("Couldn't get the episode. Please re run the program.")
		//newerr := chromedp.Shutdown(ctxt)
		//if newerr != nil {
		//	log.Fatal(newerr)
		//}
		log.Fatal(err)
		return "Couldn't get the episode. Please re run the program."
		//os.Exit(1)
	}

	r, _ := regexp.Compile(`src="(.*?)"`)
	rapidVideoString := r.FindAllString(val, 1)
	urlRapidVideo := rapidVideoString[0][4:]

	url := strings.Replace(urlRapidVideo, "\"", "", -1)
	//OpenBrowser(url)
	return url

}

func OpenEpisodes(lowerLimitEpisode, upperLimitEpisode, searchedShow, season string, wg *sync.WaitGroup) {
	defer wg.Done()
	var val string

	ctxt, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctxt, cancel = context.WithTimeout(ctxt, 180*time.Second)
	defer cancel()

	_, baseURL, episode := GetURL(searchedShow, lowerLimitEpisode, season)

	episodeSearch := baseURL + episode + "?id=&s=rapidVideo"
	//fmt.Println(episodeSearch)

	// run task list
	log.SetFlags(0)

	err := chromedp.Run(ctxt, Click(episodeSearch, &val))
	if err != nil {
		fmt.Println("Couldn't get the episode. Please re run the program.")

		log.Fatal(err)
		os.Exit(1)
	}

	r, _ := regexp.Compile(`src="(.*?)"`)
	rapidVideoString := r.FindAllString(val, 1)
	urlRapidVideo := rapidVideoString[0][4:]

	url := strings.Replace(urlRapidVideo, "\"", "", -1)
	OpenBrowser(url)

}

func GetEpisodeList(searchedShow, season string) {
	log.SetFlags(0)
	var val string
	_, baseURL, _ := GetURL(searchedShow, "1", season)

	ctxt, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var episodeSlice []string
	var dateSlice []string

	err := chromedp.Run(ctxt, ClickForEpisodeList(baseURL, &val))
	if err != nil {
		fmt.Println("Couldn't get the episode list. Please re run the program.")
		log.Fatal(err)
	}

	stringVals := strings.NewReader(val)
	doc, _ := goquery.NewDocumentFromReader(stringVals)
	var temp []string
	doc.Find("tbody").Each(func(i int, s *goquery.Selection) {
		episodeVals := s.Find("a").Text()
		dateVals := s.Find("td").Text()

		trimEpisodeVals := strings.Trim(episodeVals, "\n\t\r")
		trimDateVals := strings.Trim(dateVals, "\n\t\r")
		temp = strings.Split(trimDateVals, "\n")
		episodeSlice = append(episodeSlice, trimEpisodeVals)
		dateSlice = append(dateSlice, trimDateVals)
	})

	for _, episodes := range temp {
		fmt.Println(episodes)
	}

}

func GetEpisodeListAPI(searchedShow string) []string {
	log.SetFlags(0)
	var val string
	_, baseURL, _ := GetURL(searchedShow, "1", "")

	//ctxt, cancel := context.WithTimeout(context.Background(), 180*time.Second)

	ctxt, cancel := chromedp.NewContext(context.Background())
	//ctxt, cancel = context.WithTimeout(ctxt, 180*time.Second)
	defer cancel()

	var episodeSlice []string
	var dateSlice []string

	err := chromedp.Run(ctxt, ClickForEpisodeList(baseURL, &val))
	if err != nil {
		fmt.Println("Couldn't get the episode list. Please re run the program.")
		log.Fatal(err)
	}

	stringVals := strings.NewReader(val)
	doc, _ := goquery.NewDocumentFromReader(stringVals)
	var temp []string
	doc.Find("tbody").Each(func(i int, s *goquery.Selection) {
		episodeVals := s.Find("a").Text()
		dateVals := s.Find("td").Text()

		trimEpisodeVals := strings.Trim(episodeVals, "\n\t\r")
		trimDateVals := strings.Trim(dateVals, "\n\t\r")
		temp = strings.Split(trimDateVals, "\n")
		episodeSlice = append(episodeSlice, trimEpisodeVals)
		dateSlice = append(dateSlice, trimDateVals)
	})

	return temp
}
