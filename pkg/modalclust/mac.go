package modalclust

import (
	"encoding/json"

	"github.com/dgravesa/go-parallel/parallel"
)

// ModeDistThreshold is the allowable distance between two modes that they are still considered equivalent
var ModeDistThreshold float64 = 1e-01

// number of parallel goroutines to use in calculation; settable at runtime by SetNumGoroutines()
var numMACGoroutines int = parallel.WithCPUProportion(0.75).NumGoroutines()

// SetNumGoroutines sets the number of goroutines to use in MAC computation.
// TODO: make this variable per MAC call.
func SetNumGoroutines(numGR int) {
	numMACGoroutines = numGR
}

// MAC executes modal association clustering on a data slice
func MAC(data []DataPt, sigma float64) *MACResult {
	if data == nil {
		return nil
	}

	results := newMACResult()
	strategy := parallel.WithNumGoroutines(numMACGoroutines)
	resultsCh, done := results.newInsertChannel(numMACGoroutines)

	// execute MEM on each data point
	strategy.ForWithGrID(len(data), func(i, grID int) {
		mode := MEM(data, data[i], sigma)
		resultsCh <- macInsertPair{data[i], mode}
	})
	close(resultsCh)

	<-done
	return results
}

// MACResult is the result of a modal association clustering execution
type MACResult struct {
	clusters []Cluster
}

// Clusters returns the clusters of a modal association clustering result
func (r *MACResult) Clusters() []Cluster {
	return r.clusters
}

type resultJSON struct {
	Clusters    []clusterJSON `json:"clusters"`
	NumClusters int           `json:"count"`
}

// MarshalJSON marshals a cluster result to JSON bytes
func (r *MACResult) MarshalJSON() ([]byte, error) {
	rjson := resultJSON{}
	rjson.Clusters = []clusterJSON{}
	for _, cluster := range r.Clusters() {
		rjson.Clusters = append(rjson.Clusters, clusterJSON{
			Mode:       cluster.Mode(),
			Members:    cluster.Members(),
			NumMembers: len(cluster.Members()),
		})
	}
	rjson.NumClusters = len(rjson.Clusters)
	return json.Marshal(rjson)
}

func newMACResult() *MACResult {
	r := new(MACResult)
	r.clusters = []Cluster{}
	return r
}

type macInsertPair struct {
	datum DataPt
	mode  DataPt
}

func (r *MACResult) newInsertChannel(bufferCount int) (chan<- macInsertPair, <-chan bool) {
	resultsChannel := make(chan macInsertPair, bufferCount)
	doneChannel := make(chan bool)

	go func() {
		// insert processing loop
		for {
			pair, more := <-resultsChannel

			if more {
				r.insert(pair.datum, pair.mode)
			} else {
				doneChannel <- true
				return
			}
		}
	}()

	return resultsChannel, doneChannel
}

func (r *MACResult) insert(datum, mode DataPt) {
	// look for existing mode in cluster result
	for i, cluster := range r.clusters {
		if mode.dist(cluster.mode) < ModeDistThreshold {
			r.clusters[i].members = append(r.clusters[i].members, datum)
			return
		}
	}
	// create a new cluster with the given mode
	newCluster := Cluster{
		mode:    mode,
		members: []DataPt{datum},
	}
	r.clusters = append(r.clusters, newCluster)
}
