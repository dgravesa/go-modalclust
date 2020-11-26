package modalclust

// Cluster represents one cluster of a modal association clustering result
type Cluster struct {
	mode    DataPoint
	members []DataPoint
}

// Mode returns the local maximum of a cluster
func (c *Cluster) Mode() DataPoint {
	return c.mode
}

// Members returns the cluster membership array
func (c *Cluster) Members() []DataPoint {
	return c.members
}

type clusterJSON struct {
	Mode       DataPoint   `json:"mode"`
	Members    []DataPoint `json:"members"`
	NumMembers int         `json:"size"`
}
