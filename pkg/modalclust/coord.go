package modalclust

import "math"

// Coord represents a float64-based N-dimensional data point
type Coord []float64

// Dist returns the distance between two coordinates
func (c Coord) Dist(other Coord) float64 {
	distSq := 0.0

	for i := 0; i < len(c); i++ {
		diff := c[i] - other[i]
		distSq += diff * diff
	}

	return math.Sqrt(distSq)
}

// Add returns the addition result of this and another coordinate
func (c Coord) Add(other Coord) Coord {
	result := make(Coord, len(c))

	for i := 0; i < len(c); i++ {
		result[i] = c[i] + other[i]
	}

	return result
}

// Scale returns the result of the coordinate scaled by a factor
func (c Coord) Scale(scalar float64) Coord {
	result := make(Coord, len(c))

	for i := 0; i < len(c); i++ {
		result[i] = scalar * c[i]
	}

	return result
}
