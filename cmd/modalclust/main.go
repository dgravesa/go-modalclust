package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgravesa/go-modalclust/pkg/modalclust"
)

func main() {
	var inputName string
	var sigma float64
	var printRuntime bool

	flag.StringVar(&inputName, "InputName", "", "name of the input file")
	flag.Float64Var(&sigma, "Sigma", 0.3, "sigma value to use for clustering")
	flag.BoolVar(&printRuntime, "PrintRuntime", false, "print time to generate cluster result")
	flag.Parse()

	if inputName == "" {
		log.Fatalln("no input name provided")
	}

	// read data from file
	data := parseFileData(inputName)

	// execute clustering
	t1 := time.Now()
	result := modalclust.MAC(data, sigma)
	t2 := time.Now()

	// output results to json
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(resultJSON))

	if printRuntime {
		fmt.Println("execution time:", t2.Sub(t1))
	}
}

func parseFileData(fname string) []modalclust.DataPt {
	data := []modalclust.DataPt{}

	f, err := os.Open(fname)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	N := 0
	lnum := 1

	s := bufio.NewScanner(f)
	for s.Scan() {
		fields := strings.Split(s.Text(), ",")
		if N == 0 {
			// initialize N
			N = len(fields)
		} else if N != len(fields) {
			// mismatched number of dimensions
			log.Fatalf("mismatched dimension on line %d: expected %d-D, but found %d-D\n",
				lnum, N, len(fields))
		}

		c := make([]float64, len(fields))
		for i := 0; i < len(fields); i++ {
			val, err := strconv.ParseFloat(fields[i], 64)
			if err != nil {
				log.Fatalf("parse error on line %d: %s\n", lnum, err)
			}
			c[i] = val
		}

		data = append(data, c)

		lnum++
	}

	return data
}
