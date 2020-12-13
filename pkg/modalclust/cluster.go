package modalclust

import (
	"sync"
)

// Cluster represents one cluster of a modal association clustering result
type Cluster struct {
	mode        DataPt
	members     []DataPt
	insertMutex *sync.Mutex
}

// Mode returns the local maximum of a cluster
func (c *Cluster) Mode() DataPt {
	return c.mode
}

// Members returns the cluster membership array
func (c *Cluster) Members() []DataPt {
	return c.members
}

type clusterJSON struct {
	Mode       DataPt   `json:"mode"`
	Members    []DataPt `json:"members"`
	NumMembers int      `json:"size"`
}

func newCluster(mode DataPt, members ...DataPt) *Cluster {
	return &Cluster{
		mode:        mode,
		members:     members,
		insertMutex: new(sync.Mutex),
	}
}

func (c *Cluster) insert(datum DataPt) {
	c.insertMutex.Lock()
	c.members = append(c.members, datum)
	c.insertMutex.Unlock()
}
