package modalclust

// DataPoint is an interface that is compatible with MAC and MEM computations
type DataPoint interface {
	Dist(other DataPoint) float64
	Add(other DataPoint) DataPoint
	Scale(scalar float64) DataPoint
}
