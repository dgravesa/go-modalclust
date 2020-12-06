package modalclust

import (
	"encoding/json"
	"sync"

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

	dataCh := newDataChannel(data)
	results := newMACResult()
	strategy := parallel.WithNumGoroutines(numMACGoroutines)

	// execute MEM on each data point
	strategy.For(numMACGoroutines, func(_ int) {
		for {
			datum, more := <-dataCh
			if !more {
				break
			}
			mode := MEM(data, datum, sigma)
			results.insert(datum, mode)
		}
	})

	return results
}

func newDataChannel(data []DataPt) <-chan DataPt {
	dataChannel := make(chan DataPt)

	go func() {
		for _, datum := range data {
			dataChannel <- datum
		}
		close(dataChannel)
	}()

	return dataChannel
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
