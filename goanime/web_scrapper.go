package main 

import (
	"fmt"
	"runtime"
	"os/exec"
	"github.com/chromedp/chromedp/runner"
	"os"
	"bufio"
	"strings"
	"strconv"
	"github.com/chromedp/chromedp"
	"context"
	"regexp"
	"log"
	"reflect"

)



func main() {
	
	_, base_url, episode := getURL()
	fmt.Println(base_url)

	episode_search := base_url + episode + "?id=&s=rapidVideo"

	// chromedp 
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	var newerr error 
	var val string 

	// create headless chrome instance
	c, newerr := chromedp.New(ctxt, chromedp.WithRunnerOptions(
		runner.Flag("headless", true), 
		runner.Flag("disable-gpu", true)),
	)
	if newerr != nil { 
		log.Fatal(newerr)
	}

	// run task list 
	err := c.Run(ctxt, click(episode_search, &val))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	
	fmt.Println("out\n")
	fmt.Println(val)
	

	r, _  := regexp.Compile(`src="(.*?)"`)
	rapidVideoString := r.FindAllString(val, 1)
	urlRapidVideo := rapidVideoString[0][4:]
	
	fmt.Println(urlRapidVideo)
	fmt.Println(urlRapidVideo[:5])
	fmt.Println(reflect.TypeOf(urlRapidVideo))
	
	url := strings.Replace(urlRapidVideo, "\"", "", -1)
	openBrowser(url)

}

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


func getURL() (string, string, string) { 
	var line string 
	var episode string 

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please input anime name")
	
	if scanner.Scan() {
		line = scanner.Text()
	}
	fmt.Println("Please input episode number")
	if scanner.Scan() {
		episode = scanner.Text()
	}

	episodeToInt, err := strconv.ParseInt(episode, 10,8)
	if err != nil {
		log.Fatal(err)
	}
	if episodeToInt < 100{
		episode = "0"+episode
	}
	if episodeToInt < 10 {
		episode = "0"+episode
	}
	

	line = strings.Replace(line, " ", "-", -1)
	outline := "https://kissanime.ru/Anime/"+line+"/"
	episode = "Episode-"+ episode
	return line, outline, episode
}

func click(baseUrl string, val *string) chromedp.Tasks { 
	return chromedp.Tasks{
		chromedp.Navigate(baseUrl),
		chromedp.Click(`a.specialButton`, chromedp.NodeVisible),
		chromedp.OuterHTML(`iframe#my_video_1`, val, chromedp.NodeVisible),
	}
}