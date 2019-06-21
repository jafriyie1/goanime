package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	_ "go/goanime/cmd/API/statik"

	"github.com/gorilla/mux"
	"github.com/jafriyie1/goanime/animescrapper"
	"github.com/jafriyie1/goanime/animetries"
	"github.com/rakyll/statik/fs"
)

type Episodes struct {
	Episodes string `json:"episodes"`
}

type RapidURLInfo struct {
	URL     string `json:"url"`
	Episode string `json:"episode"`
	Show    string `json:"show"`
}

func runGoanime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := r.URL.Query()
	show := vars["name"][0]
	episode := vars["start"][0]
	wg := new(sync.WaitGroup)
	if vars["end"] != nil {
		episodeEnd := vars["end"][0]

		nEpisode, _ := strconv.Atoi(episode)
		nEpisodeEnd, _ := strconv.Atoi(episodeEnd)
		//diffSize := (nEpisodeEnd + 1) - nEpisode

		//urlSlice := make([]String, diffSize)
		var finalJson []RapidURLInfo

		for i := nEpisode; i < nEpisodeEnd+1; i++ {

			//c, ctxt := animescrapper.DoGoAnime()

			//loopedEpisode := nEpisode + i
			nLoopedEpisode := strconv.Itoa(i)
			fmt.Println("Getting Episode", nLoopedEpisode)
			fmt.Println(i)
			wg.Add(1)

			go func(loopedEpisode string, show string) {
				defer wg.Done()
				var finalInfo RapidURLInfo
				url := animescrapper.OpenEpisodesAPI(loopedEpisode, show)

				finalInfo.URL = url
				finalInfo.Episode = loopedEpisode
				finalInfo.Show = show
				//urlSlice[i] = url
				//fmt.Println(i)
				finalJson = append(finalJson, finalInfo)
			}(nLoopedEpisode, show)

			//time.Sleep(5*time.Second)
		}
		wg.Wait()
		json.NewEncoder(w).Encode(finalJson)

	} else {
		url := animescrapper.OpenEpisodesAPI(episode, show)

		var finalInfo RapidURLInfo

		finalInfo.URL = url
		finalInfo.Episode = episode
		finalInfo.Show = show

		json.NewEncoder(w).Encode(finalInfo)
	}

}

func getEpisodes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := r.URL.Query()
	show := vars["name"][0]
	episodes := animescrapper.GetEpisodeListAPI(show)
	var episodeSlice []Episodes

	for _, episode := range episodes {
		var episodeStruct Episodes
		episodeStruct.Episodes = episode
		episodeSlice = append(episodeSlice, episodeStruct)
	}

	json.NewEncoder(w).Encode(episodeSlice)
}

func getMatchedShows(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := r.URL.Query()
	queryForShow := vars["name"][0]
	fmt.Println(queryForShow)

	//queryForShow := query.Get("name")
	/*
		f, err := os.Open("episodes.csv")
		if err != nil {
			log.Fatal(err)
		}
	*/
	f, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	fileTrie := csv.NewReader(bufio.NewReader(f))

	_, builtTrie, _ := animetries.BuildAnimeTrie(fileTrie)

	shows := animetries.PossibleShowsAPI(builtTrie, queryForShow)
	//fmt.Println(shows)
	json.NewEncoder(w).Encode(shows)
}

func main() {
	router := mux.NewRouter()
	//posts = append(posts, Post{ID: "1", Title: "My first post", Body:      "This is the content of my first post"})
	router.HandleFunc("/search", getMatchedShows).Methods("GET")
	router.HandleFunc("/episodes", getEpisodes).Methods("GET")
	router.HandleFunc("/goanime", runGoanime).Methods("GET")

	http.ListenAndServe(":8000", router)
}
