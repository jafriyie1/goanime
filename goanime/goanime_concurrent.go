// This program allows one to watch anime
// without having to go through the hoops
// of using popular anime sites
package main

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
	"github.com/chromedp/chromedp/runner"
	"github.com/derekparker/trie"
	"github.com/jafriyie1/animetries"
)

func openBrowser(url string) {
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

func getShow(b *trie.Trie) string {
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

func getSeason() string {
	var season string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Which season would you like to watch (1st, 2nd, 3rd, 4th, etc.)?\n(Hit enter if there is no season) ")
	if scanner.Scan() {
		season = scanner.Text()
	}
	return season
}

func getURL(show, episode, season string) (string, string, string) {

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

func getOneEpisode() string {
	var episode string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Please input episode number")
	if scanner.Scan() {
		episode = scanner.Text()
	}

	episode = strings.TrimSpace(episode)
	return episode

}

func getRangeOfEpisodes() (string, string) {
	var episodeStart string
	var episodeEnd string
	scanner := bufio.NewScanner(os.Stdin)

	if episodeStart == episodeEnd {
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

	return episodeStart, episodeEnd
}

func click(url string, val *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		//chromedp.WaitVisible(`#footer`),
		chromedp.Click(`a.specialButton`, chromedp.NodeVisible),
		chromedp.Click(`a.specialButton`, chromedp.NodeVisible),
		chromedp.WaitVisible(`div#divMyVideo`),
		chromedp.OuterHTML(`iframe#my_video_1`, val, chromedp.NodeVisible),
	}
}

func clickForEpisodeList(url string, val *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#head`),
		chromedp.OuterHTML(`table.listing`, val, chromedp.NodeVisible),
		//chromedp.Click(`a.specialButton`, chromedp.NodeVisible),
		//chromedp.Click(`a.specialButton`, chromedp.NodeVisible),
		//chromedp.WaitVisible(`div#divMyVideo`),
		//chromedp.OuterHTML(`iframe#my_video_1`, val, chromedp.NodeVisible),
	}
}

func nothing() {
	return
}

func doGoAnime() (*chromedp.CDP, context.Context) {

	// chromedp
	ctxt, _ := context.WithCancel(context.Background())
	//defer cancel()

	// create headless chrome instance
	/*
		c, _ := chromedp.New(ctxt, chromedp.WithRunnerOptions(
			runner.Flag("headless", true),
			runner.Flag("disable-gpu", true),
			runner.Flag("no-sandbox", true)),
		)
	*/
	c, _ := chromedp.New(ctxt, chromedp.WithRunnerOptions(
		//runner.Headless(pathBrowser, 9222),
		//runner.Flag("headless", true),
		//runner.Flag("disable-gpu", true),
		//runner.Flag("no-sandbox", true)))
		//runner.Headless(path, 9222),
		runner.Flag("headless", true),
		runner.Flag("disable-gpu", true),
		runner.Flag("no-first-run", true),
		runner.Flag("no-default-browser-check", true),
		//runner.Port(9222),
	))
	/*
		if newerr != nil {
			log.Fatal(newerr)
		}
	*/

	return c, ctxt

}

func concurrentEpisodes(lowerLimitEpisode, upperLimitEpisode, searchedShow, season string, wg *sync.WaitGroup, c *chromedp.CDP, ctxt context.Context) {
	//defer wg.Done()
	//var newerr error
	var val string

	_, baseURL, episode := getURL(searchedShow, lowerLimitEpisode, season)
	//fmt.Println(base_url)

	episodeSearch := baseURL + episode + "?id=&s=rapidVideo"
	fmt.Println(episodeSearch)
	// run task list

	err := c.Run(ctxt, click(episodeSearch, &val))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	r, _ := regexp.Compile(`src="(.*?)"`)
	rapidVideoString := r.FindAllString(val, 1)
	urlRapidVideo := rapidVideoString[0][4:]

	url := strings.Replace(urlRapidVideo, "\"", "", -1)
	openBrowser(url)
	wg.Done()

}

func getEpisodeList(searchedShow, season string, c *chromedp.CDP, ctxt context.Context) {
	//defer wg.Done()
	//var newerr error
	var val string
	_, baseURL, _ := getURL(searchedShow, "1", season)
	//fmt.Println(base_url)

	// run task list
	var episodeSlice []string
	c.Run(ctxt, clickForEpisodeList(baseURL, &val))

	stringVals := strings.NewReader(val)
	doc, _ := goquery.NewDocumentFromReader(stringVals)

	doc.Find("tbody").Each(func(i int, s *goquery.Selection) {
		episodeVals := s.Find("a").Text()
		trimEpisodeVals := strings.Trim(episodeVals, "\n\t\r")
		//fmt.Print(trimEpisodeVals)
		episodeSlice = append(episodeSlice, trimEpisodeVals)
	})

	for _, c := range episodeSlice {
		fmt.Println(c)
	}

	c.Shutdown(ctxt)

}

func doGoAnimeOneEpisode(lowerLimitEpisode, upperLimitEpisode, searchedShow, season string) {
	_, baseURL, episode := getURL(searchedShow, lowerLimitEpisode, season)
	//fmt.Println(base_url)

	episodeSearch := baseURL + episode + "?id=&s=rapidVideo"

	// chromedp
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	var newerr error
	var val string

	fmt.Println("Please wait....")

	// create headless chrome instance
	//chromedp.Withlog(log.Printf)
	log.SetFlags(0)

	//path := `/Applications/Google Chrome.app`
	/*
		if runtime.GOOS != "windows" {
			path = "/usr/bin/google-chrome"
		}
	*/

	c, newerr := chromedp.New(ctxt, chromedp.WithRunnerOptions(
		//runner.Headless(pathBrowser, 9222),
		//runner.Flag("headless", true),
		//runner.Flag("disable-gpu", true),
		//runner.Flag("no-sandbox", true)))
		//runner.Headless(path, 9222),
		runner.Flag("headless", true),
		runner.Flag("disable-gpu", true),
		runner.Flag("no-first-run", true),
		runner.Flag("no-default-browser-check", true),
	))

	if newerr != nil {
		log.Fatal(newerr)
	}

	// run task list
	err := c.Run(ctxt, click(episodeSearch, &val))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish

	r, _ := regexp.Compile(`src="(.*?)"`)
	rapidVideoString := r.FindAllString(val, 1)
	urlRapidVideo := rapidVideoString[0][4:]

	url := strings.Replace(urlRapidVideo, "\"", "", -1)
	openBrowser(url)

}

func main() {
	_, builtTrie, _ := animetries.BuildAnimeTrie()
	var option string

	scanner := bufio.NewScanner(os.Stdin)
	searchedShow := getShow(builtTrie)
	searchedShow = strings.TrimSpace(searchedShow)
	//maxEpisode := animetries.GetEpisodeFromMap(searchedShow, animeMap)
	fmt.Println()
	season := getSeason()
	fmt.Println("Here is a list of episodes for the given show and season (please wait):")

	cOne, ctxtOne := doGoAnime()
	getEpisodeList(searchedShow, season, cOne, ctxtOne)
	cOne.Shutdown(ctxtOne)
	/*
		if cErrOne != nil {
			log.Fatal(cErrOne)
		}
	*/
	fmt.Println()
	time.Sleep(5)
	//fmt.Println(searchedShow + " has a maximum (or currently) " + maxEpisode + " episodes")
	fmt.Println("Scroll up to view episodes (and please ignore other messages).\nWould you like to watch one episode or mutliple (1 for episode, 2 for multiple)")

	if scanner.Scan() {
		option = scanner.Text()
	}
	option = strings.TrimSpace(option)
	var lowerLimitEpisode string
	//var upperLimitEpisode *string
	upperLimitEpisode := " "

	if option == "1" {
		//lowerLimitEpisode = getOneEpisode()
		lowerLimitEpisode, lowerLimitEpisode = getRangeOfEpisodes()
	} else {
		fmt.Println("WARNING: You can only get a maximum of 2 episodes.\nOutside of that you will get wonky behavior.")

		lowerLimitEpisode, upperLimitEpisode = getRangeOfEpisodes()
	}
	upperLimitEpisode = lowerLimitEpisode

	wg := new(sync.WaitGroup)
	if upperLimitEpisode != " " {

		fmt.Println("here")
		lowerEpisode, _ := strconv.Atoi(lowerLimitEpisode)
		//fmt.Println(lowerEpisode)

		upperEpisode, _ := strconv.Atoi(upperLimitEpisode)
		c, ctxt := doGoAnime()
		//fmt.Println(upperEpisode)
		for i := lowerEpisode; i < upperEpisode+1; i++ {
			wg.Add(1)
			fmt.Println(i)
			loopedEpisode := strconv.Itoa(i)
			//defer wg.Done()

			go concurrentEpisodes(loopedEpisode, upperLimitEpisode, searchedShow, season, wg, c, ctxt)

		}
		wg.Wait()
		cErr := c.Shutdown(ctxt)
		if cErr != nil {
			log.Fatal(cErr)
		}
	} else {
		//wg.Add(1)
		doGoAnimeOneEpisode(lowerLimitEpisode, upperLimitEpisode, searchedShow, season)
		//wg.Wait()
	}

}
