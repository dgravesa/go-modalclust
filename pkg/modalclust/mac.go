package modalclust

import (
	"encoding/json"
	"sync"

	"github.com/dgravesa/go-parallel/parallel"
)

// ModeDistThreshold is the allowable distance between two modes that they are still considered equivalent
var ModeDistThreshold float64 = 1e-01

// number of parallel goroutines to use in calculation; settable at runtime by SetNumGoroutines()
var numMACGoroutines int = parallel.DefaultNumGoroutines()

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

	input := newInputRetriever(data)
	results := newMACResult()
	strategy := parallel.WithNumGoroutines(numMACGoroutines)

	// execute MEM on each data point
	strategy.For(numMACGoroutines, func(_ int) {
		for {
			datum, more := input.fetch()
			if !more {
				break
			}
			mode := MEM(data, datum, sigma)
			results.insert(datum, mode)
		}
	})

	return results
}

// MACResult is the result of a modal association clustering execution
type MACResult struct {
	clusters      []Cluster
	clustersMutex *sync.RWMutex
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

type inputRetriever struct {
	data         []DataPt
	dataLen      int
	currentIndex int
	fetchMutex   *sync.Mutex
}

func newInputRetriever(data []DataPt) inputRetriever {
	return inputRetriever{
		data:         data,
		dataLen:      len(data),
		currentIndex: 0,
		fetchMutex:   new(sync.Mutex),
	}
}

func (ir *inputRetriever) fetch() (DataPt, bool) {
	ir.fetchMutex.Lock()
	defer ir.fetchMutex.Unlock()

	if ir.currentIndex >= ir.dataLen {
		return nil, false
	}

	datum := ir.data[ir.currentIndex]
	ir.currentIndex++
	return datum, true
}

func newMACResult() *MACResult {
	r := new(MACResult)
	r.clusters = []Cluster{}
	r.clustersMutex = new(sync.RWMutex)
	return r
}

type macInsertPair struct {
	datum DataPt
	mode  DataPt
}

// insert is a thread-safe call to merge a datum-mode pair into a clustering result
func (r *MACResult) insert(datum, mode DataPt) {
	// look for existing mode in cluster result
	r.clustersMutex.RLock()
	for i := 0; i < len(r.clusters); i++ {
		cluster := &r.clusters[i]

		if mode.dist(cluster.mode) < ModeDistThreshold {
			// insert into existing cluster
			cluster.insert(datum)
			r.clustersMutex.RUnlock()
			return
		}
	}
	r.clustersMutex.RUnlock()

	// create a new cluster with the given mode
	cluster := makeCluster(mode, datum)
	r.clustersMutex.Lock()
	r.clusters = append(r.clusters, cluster)
	r.clustersMutex.Unlock()
}
