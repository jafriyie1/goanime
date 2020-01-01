package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/akamensky/argparse"
	"github.com/jafriyie1/goanime/animescrapper"
	"github.com/jafriyie1/goanime/animetries"
)

func runAnime(lowerEpisode, upperEpisode int, searchedShow string) {
	fmt.Println("Please wait....")
	wg := new(sync.WaitGroup)
	for i := lowerEpisode; i < upperEpisode+1; i++ {
		wg.Add(1)

		fmt.Println("Getting Episode", i)
		loopedEpisode := strconv.Itoa(i)
		upperLimit := strconv.Itoa(upperEpisode + 1)
		go animescrapper.OpenEpisodes(loopedEpisode, upperLimit, searchedShow, "", wg)

	}
	wg.Wait()
}

func main() {

	// Create new parser object
	parser := argparse.NewParser("print", "Prints provided string to stdout")
	// Create string flag
	show := parser.String("s", "show", &argparse.Options{Required: false, Help: "String to print"})
	start := parser.String("b", "begin", &argparse.Options{Required: false, Help: "String to print"})
	end := parser.String("e", "end", &argparse.Options{Required: false, Help: "String to print"})

	// Parse input
	err := parser.Parse(os.Args)
	//if err != nil {
	// In case of error print error and print usage
	// This can also be done by passing -h or --help flags
	//fmt.Print(parser.Usage(err))
	//}
	// Finally print the collected string
	//fmt.Println(*s)

	f, err := os.Open("../../Data/test/episodes.csv")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(bufio.NewReader(f))

	_, builtTrie, _ := animetries.BuildAnimeTrie(r)
	var option string

	if *show == "" || *start == "" || *end == "" {

		scanner := bufio.NewScanner(os.Stdin)
		searchedShow := animescrapper.GetShow(builtTrie)
		searchedShow = strings.TrimSpace(searchedShow)

		fmt.Println()
		season := ""
		fmt.Println("Here is a list of episodes for the given show and season (please wait):")

		animescrapper.GetEpisodeList(searchedShow, season)

		fmt.Println()
		fmt.Println("Scroll up to view episodes (and please ignore other messages).\nWould you like to watch one episode or mutliple (1 for episode, 2 for multiple)")

		if scanner.Scan() {
			option = scanner.Text()
		}
		option = strings.TrimSpace(option)

		var lowerLimitEpisode string
		var upperLimitEpisode string

		if option == "1" {
			lowerLimitEpisode, _ = animescrapper.GetRangeOfEpisodes(false)
			upperLimitEpisode = lowerLimitEpisode

		} else {
			lowerLimitEpisode, upperLimitEpisode = animescrapper.GetRangeOfEpisodes(true)
		}

		lowerEpisode, _ := strconv.Atoi(lowerLimitEpisode)
		upperEpisode, _ := strconv.Atoi(upperLimitEpisode)

		runAnime(lowerEpisode, upperEpisode, searchedShow)

	} else {
		var searchedShow string
		searchedShow = strings.Replace(*show, "'", "", -1)

		var lowerLimitEpisode string
		var upperLimitEpisode string
		lowerLimitEpisode = *start
		upperLimitEpisode = *end

		fmt.Println(searchedShow, lowerLimitEpisode, upperLimitEpisode)

		var lowerEpisode int
		var upperEpisode int

		lowerEpisode, _ = strconv.Atoi(lowerLimitEpisode)
		upperEpisode, _ = strconv.Atoi(upperLimitEpisode)

		runAnime(lowerEpisode, upperEpisode, searchedShow)
	}
}
