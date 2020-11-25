package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/dgravesa/go-modalclust/pkg/modalclust"
)

func main() {
	var inputName string
	var sigma float64

	flag.StringVar(&inputName, "InputName", "", "name of the input file")
	flag.Float64Var(&sigma, "Sigma", 0.3, "sigma value to use for clustering")
	flag.Parse()

	data := parseFileData(inputName)

	result := modalclust.MAC(data, sigma)

	// output results to json
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(resultJSON))
}

func parseFileData(fname string) []modalclust.Coord {
	// TODO: actually parse file data
	return []modalclust.Coord{
		{0.02, 0.18},
		{-0.10, 0.03},
		{0.05, -0.08},
		{6.67, 3.21},
		{6.71, 3.24},
		{0.13, 0.11},
		{6.62, 3.19},
	}
}
