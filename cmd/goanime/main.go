package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jafriyie1/goanime/animescrapper"
	"github.com/jafriyie1/goanime/animetries"
)

func main() {
	_, builtTrie, _ := animetries.BuildAnimeTrie()
	var option string

	scanner := bufio.NewScanner(os.Stdin)
	searchedShow := animescrapper.getShow(builtTrie)
	searchedShow = strings.TrimSpace(searchedShow)

	fmt.Println()
	season := animescrapper.getSeason()
	fmt.Println("Here is a list of episodes for the given show and season (please wait):")

	animescrapper.getEpisodeList(searchedShow, season)

	fmt.Println()
	time.Sleep(5)
	fmt.Println("Scroll up to view episodes (and please ignore other messages).\nWould you like to watch one episode or mutliple (1 for episode, 2 for multiple)")

	if scanner.Scan() {
		option = scanner.Text()
	}
	option = strings.TrimSpace(option)
	var lowerLimitEpisode string

	upperLimitEpisode := " "

	if option == "1" {
		lowerLimitEpisode, lowerLimitEpisode = animescrapper.getRangeOfEpisodes()
	} else {
		fmt.Println("WARNING: You can only get a maximum of 2 episodes.\nOutside of that you will get wonky behavior.")

		lowerLimitEpisode, upperLimitEpisode = animescrapper.getRangeOfEpisodes()
	}
	upperLimitEpisode = lowerLimitEpisode

	wg := new(sync.WaitGroup)

	lowerEpisode, _ := strconv.Atoi(lowerLimitEpisode)
	upperEpisode, _ := strconv.Atoi(upperLimitEpisode)

	c, ctxt := animescrapper.doGoAnime()
	fmt.Println("Please wait....")
	for i := lowerEpisode; i < upperEpisode+1; i++ {
		wg.Add(1)
		loopedEpisode := strconv.Itoa(i)

		go animescrapper.concurrentEpisodes(loopedEpisode, upperLimitEpisode, searchedShow, season, wg, c, ctxt)

	}
	wg.Wait()
	cErr := c.Shutdown(ctxt)

	if cErr != nil {
		log.Fatal(cErr)
	}

}
