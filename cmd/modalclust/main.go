package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"strconv"
	"strings"
	"time"

	"github.com/dgravesa/go-modalclust/pkg/modalclust"
)

// command line arguments
var inputName string
var sigma float64
var numGoroutines int
var printRuntime bool
var traceName string
var cpuProfName string

var macStart time.Time
var macRuntime time.Duration
var traceFile *os.File
var cpuProfFile *os.File

func main() {
	flag.StringVar(&inputName, "input", "", "name of the input file")
	flag.Float64Var(&sigma, "sigma", 0.3, "sigma value to use for clustering")
	flag.IntVar(&numGoroutines, "numgr", 0, "specify number of goroutines to use in MAC computation")
	flag.BoolVar(&printRuntime, "runtime", false, "print time to generate cluster result")
	flag.StringVar(&traceName, "trace", "", "output trace of MAC to file")
	flag.StringVar(&cpuProfName, "cpuprofile", "", "output CPU profile of MAC to file")
	flag.Parse()

	if inputName == "" {
		log.Fatalln("no input name provided")
	}

	// read data from file
	data := parseFileData(inputName)

	// set number of goroutines if specified
	if numGoroutines != 0 {
		modalclust.SetNumGoroutines(numGoroutines)
	}

	// execute clustering
	preMAC()
	result := modalclust.MAC(data, sigma)
	postMAC()

	// output results to json
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(resultJSON))

	if printRuntime {
		fmt.Println("execution time:", macRuntime)
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

func preMAC() {
	if traceName != "" {
		var err error
		traceFile, err = os.Create(traceName)

		if err != nil {
			log.Fatalln(err)
		}

		err = trace.Start(traceFile)

		if err != nil {
			log.Fatalln(err)
		}
	}

	if cpuProfName != "" {
		var err error
		cpuProfFile, err = os.Create(cpuProfName)

		if err != nil {
			log.Fatalln(err)
		}

		err = pprof.StartCPUProfile(cpuProfFile)

		if err != nil {
			log.Fatalln(err)
		}
	}

	macStart = time.Now()
}

func postMAC() {
	macStop := time.Now()
	macRuntime = macStop.Sub(macStart)

	if traceName != "" {
		trace.Stop()
		err := traceFile.Close()

		if err != nil {
			log.Fatalln(err)
		}
	}

	if cpuProfName != "" {
		pprof.StopCPUProfile()
		err := cpuProfFile.Close()

		if err != nil {
			log.Fatalln(err)
		}
	}
}
