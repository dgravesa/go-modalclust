package modalclust

import (
	"encoding/json"
)

// ModeDistThreshold is the allowable distance between two modes that they are still considered equivalent
var ModeDistThreshold float64 = 1e-02

// Result is the result of a modal association clustering execution
type Result struct {
	clusters []Cluster
}

// Clusters returns the clusters of a modal association clustering result
func (r *Result) Clusters() []Cluster {
	return r.clusters
}

type resultJSON struct {
	Clusters    []clusterJSON `json:"clusters"`
	NumClusters int           `json:"count"`
}

// MarshalJSON marshals a cluster result to JSON bytes
func (r *Result) MarshalJSON() ([]byte, error) {
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

func newResult() *Result {
	r := new(Result)
	r.clusters = []Cluster{}
	return r
}

func (r *Result) merge(datum, mode DataPt) {
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
