package modalclust

// Cluster represents one cluster of a modal association clustering result
type Cluster struct {
	mode    DataPt
	members []DataPt
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
