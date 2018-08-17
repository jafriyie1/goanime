package main

import (
	"bufio"
	"encoding/csv"
	"os"

	"github.com/kniren/gota/dataframe"
)

func main() {
	f, _ := os.Open("../Data/anime.csv")

	r := csv.NewReader(bufio.NewReader(f))
	df := dataframe.ReadCSV(f)
}
