package main

// This code is the first attempts of
// a trie algorithm to aid in keyword
// search for anime shows

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/derekparker/trie"
)

type AnimeInfo struct {
	AnimeID  string
	Name     string
	Genre    string
	Type     string
	Episodes string
	Rating   string
	Members  string
}

func main() {
	f, _ := os.Open("../Data/anime.csv")

	//create new trie
	animeTrie := trie.New()
	var anime []AnimeInfo

	r := csv.NewReader(bufio.NewReader(f))

	for {
		line, error := r.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		anime = append(anime, AnimeInfo{
			line[0], line[1], line[2], line[3],
			line[4], line[5], line[6],
		})
	}
	//fmt.Println(anime)

	for i, a := range anime {
		animeTrie.Add(a.Name, i)
		//fmt.Println(a.Name)
	}
	//fmt.Println(animeTrie.FuzzySearch("One Piece"))
	fmt.Println(animeTrie.PrefixSearch("Boku no Hero")[1])
	/*df := dataframe.ReadCSV(f)
	namesOfShows := df.Select([]string{"name"})
	fmt.Println(namesOfShows[1])

	for i, shows := range namesOfShows {
		animeTrie.Add(shows, i)
	}

	animeTrie.PrefixSearch("boku")
	*/
}
