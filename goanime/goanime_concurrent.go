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

	"github.com/derekparker/trie"
	"github.com/jafriyie1/animetries"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
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
	var line string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please input anime name.\nYou can also input a part of the show to search.")

	if scanner.Scan() {
		line = scanner.Text()
	}
	fmt.Println("Here are some shows that match your search:")
	fmt.Println()
	animetries.PossibleShows(b, line)
	fmt.Println()
	fmt.Println("Please copy or type the anime show and hit enter")
	if scanner.Scan() {
		line = scanner.Text()
	}
	return line

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

func getURL(line, episode, season string) (string, string, string) {
	//var line string
	//var episode string
	//var season string

	//scanner := bufio.NewScanner(os.Stdin)
	/*
		fmt.Println("Please input anime name")

		if scanner.Scan() {
			line = scanner.Text()
		}
	*/
	/*
		line = strings.TrimSpace(line)
		fmt.Println("Please input episode number")
		if scanner.Scan() {
			episode = scanner.Text()
		}
	*/

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
	if season != "" {
		season = "-" + season + "-Season"
	}
	line = strings.Replace(line, ":", "", -1)
	line = strings.Replace(line, ")", "", -1)
	line = strings.Replace(line, "(", "", -1)
	line = strings.Replace(line, " ", "-", -1)
	base_url := "https://kissanime.ru/Anime/" + line + season + "/"

	episode = "Episode-" + episode

	return line, base_url, episode
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

	return episodeStart, episodeEnd
}

func click(url string, val *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Click(`a.specialButton`, chromedp.NodeVisible),
		chromedp.OuterHTML(`iframe#my_video_1`, val, chromedp.NodeVisible),
	}
}

func doGoAnime(lowerLimitEpisode, upperLimitEpisode, searchedShow, season string, wg *sync.WaitGroup) {
	//defer wg.Done()
	_, baseURL, episode := getURL(searchedShow, lowerLimitEpisode, season)
	//fmt.Println(base_url)

	episodeSearch := baseURL + episode + "?id=&s=rapidVideo"

	// chromedp
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	var newerr error
	var val string

	fmt.Println(lowerLimitEpisode)
	fmt.Println("Please wait....")

	// create headless chrome instance
	c, newerr := chromedp.New(ctxt, chromedp.WithRunnerOptions(
		runner.Flag("headless", true),
		runner.Flag("disable-gpu", true)),
	)
	if newerr != nil {
		log.Fatal(newerr)
	}

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
	c, newerr := chromedp.New(ctxt, chromedp.WithRunnerOptions(
		runner.Flag("headless", true),
		runner.Flag("disable-gpu", true)),
	)
	if newerr != nil {
		log.Fatal(newerr)
	}

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

}

func main() {
	_, builtTrie, animeMap := animetries.BuildAnimeTrie()
	var option string

	scanner := bufio.NewScanner(os.Stdin)
	searchedShow := getShow(builtTrie)
	maxEpisode := animetries.GetEpisodeFromMap(searchedShow, animeMap)
	fmt.Println()
	fmt.Println(searchedShow + " has a maximum (or currently) " + maxEpisode + " episodes")
	fmt.Println("Would you like to watch one episode or mutliple (1 for episode, 2 for multiple)")

	if scanner.Scan() {
		option = scanner.Text()
	}
	option = strings.TrimSpace(option)
	var lowerLimitEpisode string
	//var upperLimitEpisode *string
	upperLimitEpisode := " "

	if option == "1" {
		lowerLimitEpisode = getOneEpisode()
	} else {
		fmt.Println("WARNING: You can only get a maximum of 2 episodes.\nOutside of that you will get wonky behavior.")

		lowerLimitEpisode, upperLimitEpisode = getRangeOfEpisodes()
	}

	season := getSeason()
	if upperLimitEpisode != " " {
		wg := new(sync.WaitGroup)
		fmt.Println("here")
		lowerEpisode, _ := strconv.Atoi(lowerLimitEpisode)
		//fmt.Println(lowerEpisode)

		upperEpisode, _ := strconv.Atoi(upperLimitEpisode)
		//fmt.Println(upperEpisode)
		for i := lowerEpisode; i < upperEpisode+1; i++ {
			wg.Add(1)
			fmt.Println(i)
			loopedEpisode := strconv.Itoa(i)
			//defer wg.Done()

			go doGoAnime(loopedEpisode, upperLimitEpisode, searchedShow, season, wg)
			//wg.Done()
			//wg.Wait()

			/*
				go func(i int) { // ensures each run gets distinct i
					fmt.Println("Sleeping", i, "seconds")
					time.Sleep(time.Duration(i) * time.Second)
					fmt.Println("Slept", i, "seconds")
					wg.Done()
				}(i)
			*/
		}
		wg.Wait()
	} else {
		doGoAnimeOneEpisode(lowerLimitEpisode, upperLimitEpisode, searchedShow, season)
	}

}
