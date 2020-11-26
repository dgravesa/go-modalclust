package modalclust

import (
	"encoding/json"

	"github.com/dgravesa/go-parallel/parallel"
)

// ModeDistThreshold is the allowable distance between two modes that they are still considered equivalent
var ModeDistThreshold float64 = 1e-01

// MAC executes modal association clustering on a data slice
func MAC(data []DataPt, sigma float64) *MACResult {
	if data == nil {
		return nil
	}

	N := len(data)

	// initialize per-thread results
	strategy := parallel.WithCPUProportion(0.7)
	numGoroutines := strategy.NumGoroutines()
	results := []*MACResult{}
	for i := 0; i < numGoroutines; i++ {
		results = append(results, newMACResult())
	}

	// execute MEM on each data point
	strategy.ForWithGrID(N, func(i, grID int) {
		mode := MEM(data, data[i], sigma)
		results[grID].insert(data[i], mode)
	})

	for i := 1; i < numGoroutines; i++ {
		results[0].merge(results[i])
	}

	return results[0]
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

func (r *MACResult) insert(datum, mode DataPt) {
	// look for existing mode in cluster result
	for i, cluster := range r.clusters {
		if mode.Dist(cluster.mode) < ModeDistThreshold {
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

func (r *MACResult) merge(other *MACResult) {
	for _, oc := range other.clusters {
		// look for existing mode in my cluster result
		found := false
		for ri, rc := range r.clusters {
			if rc.mode.Dist(oc.mode) < ModeDistThreshold {
				r.clusters[ri].members = append(r.clusters[ri].members, oc.members...)
				found = true
				break
			}
		}
		// create a new cluster in my cluster result
		if !found {
			r.clusters = append(r.clusters, oc)
		}
	}
}
