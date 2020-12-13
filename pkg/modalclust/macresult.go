package modalclust

import (
	"encoding/json"
	"sync"
)

// MACResult is the result of a modal association clustering execution
type MACResult struct {
	clusters      []*Cluster
	clustersMutex *sync.RWMutex
}

// Clusters returns the clusters of a modal association clustering result
func (r *MACResult) Clusters() []*Cluster {
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
	r.clusters = []*Cluster{}
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
	cluster := r.findCluster(mode)
	r.clustersMutex.RUnlock()

	// insert into existing cluster
	if cluster != nil {
		cluster.insert(datum)
		return
	}

	r.clustersMutex.Lock()
	if cluster = r.findCluster(mode); cluster != nil {
		// insert into recently created cluster
		cluster.insert(datum)
	} else {
		// create a new cluster with the given mode
		cluster = newCluster(mode, datum)
		r.clusters = append(r.clusters, cluster)
	}
	r.clustersMutex.Unlock()
}

func (r *MACResult) findCluster(mode DataPt) *Cluster {
	for i := 0; i < len(r.clusters); i++ {
		cluster := r.clusters[i]

		if mode.dist(cluster.mode) < ModeDistThreshold {
			return cluster
		}
	}
	return nil
}
