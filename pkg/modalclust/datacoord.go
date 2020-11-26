package modalclust

import "math"

// DataCoord is a float64-based N-dimensional DataPoint implementation
type DataCoord []float64

// Dist returns the distance between two coordinates
func (dp DataCoord) Dist(other DataPoint) float64 {
	otherDC := other.(*DataCoord)

	distSq := 0.0

	for i := 0; i < len(dp); i++ {
		diff := dp[i] - (*otherDC)[i]
		distSq += diff * diff
	}

	return math.Sqrt(distSq)
}

// Add returns the addition result of this and another coordinate
func (dp DataCoord) Add(other DataPoint) DataPoint {
	otherDC := other.(*DataCoord)

	result := make(DataCoord, len(dp))

	for i := 0; i < len(dp); i++ {
		result[i] = dp[i] + (*otherDC)[i]
	}

	return &result
}

// Scale returns the result of the coordinate scaled by a factor
func (dp DataCoord) Scale(scalar float64) DataPoint {
	result := make(DataCoord, len(dp))

	for i := 0; i < len(dp); i++ {
		result[i] = scalar * dp[i]
	}

	return &result
}
