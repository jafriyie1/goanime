package animetries

// This code is the first attempts of
// a trie algorithm to aid in keyword
// search for anime shows

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"sort"

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

func PossibleShows(builtTrie *trie.Trie, searchedShow string) {
	possibleShows := builtTrie.PrefixSearch(searchedShow)
	for _, show := range possibleShows {
		fmt.Println(show)
	}
}

func ShowToEpisodeMap(a []AnimeInfo) []map[string]string {
	var animeEpisodeMap []map[string]string
	var mapsOfShows = make(map[string]string)
	for _, info := range a {
		mapsOfShows[info.Name] = info.Episodes
		animeEpisodeMap = append(animeEpisodeMap, mapsOfShows)
	}
	return animeEpisodeMap
}

func GetEpisodeFromMap(name string, m []map[string]string) string {
	var value string
	for _, maps := range m {
		if val, ok := maps[name]; ok {
			value = val
			break
		}
	}
	return value
}

func BinarySearchAnime(a []string, search string) (result int, count int) {
	mid := len(a) / 2
	switch {
	case len(a) == 0:
		result = -1
	case a[mid] > search:
		result, count = BinarySearchAnime(a[:mid], search)
	case a[mid] < search:
		result, count = BinarySearchAnime(a[mid+1:], search)
		result += mid + 1
	default:
		result = mid
	}
	count++
	return
}

//BuildAnimeTrie

func BuildAnimeTrie(f *File) ([]string, *trie.Trie, []map[string]string) {
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

	sort.Slice(anime, func(i, j int) bool {
		if anime[i].Name < anime[i].Name {
			return true
		}
		if anime[i].Name > anime[i].Name {
			return false
		}
		return anime[i].Name < anime[i].Name
	})
	animeEpisodeMap := ShowToEpisodeMap(anime)

	var justAnimeShows []string
	for i, a := range anime {
		animeTrie.Add(a.Name, i)
		justAnimeShows = append(justAnimeShows, a.Name)

	}

	sort.Strings(justAnimeShows)

	return justAnimeShows, animeTrie, animeEpisodeMap

}
