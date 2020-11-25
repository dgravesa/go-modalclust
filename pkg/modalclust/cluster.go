package modalclust

// Cluster represents one cluster of a modal association clustering result
type Cluster struct {
	mode    Coord
	members []Coord
}

// Mode returns the local maximum of a cluster
func (c *Cluster) Mode() Coord {
	return c.mode
}

// Members returns the cluster membership array
func (c *Cluster) Members() []Coord {
	return c.members
}

type clusterJSON struct {
	Mode       Coord   `json:"mode"`
	Members    []Coord `json:"members"`
	NumMembers int     `json:"size"`
}
