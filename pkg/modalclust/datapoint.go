package modalclust

import "math"

// DataPt represents a float64-based N-dimensional data point
type DataPt []float64

// Dist returns the distance between two coordinates
func (dp DataPt) Dist(other DataPt) float64 {
	distSq := 0.0

	for i := 0; i < len(dp); i++ {
		diff := dp[i] - other[i]
		distSq += diff * diff
	}

	return math.Sqrt(distSq)
}

// Add returns the addition result of this and another coordinate
func (dp DataPt) Add(other DataPt) DataPt {
	result := make(DataPt, len(dp))

	for i := 0; i < len(dp); i++ {
		result[i] = dp[i] + other[i]
	}

	return result
}

// Scale returns the result of the coordinate scaled by a factor
func (dp DataPt) Scale(scalar float64) DataPt {
	result := make(DataPt, len(dp))

	for i := 0; i < len(dp); i++ {
		result[i] = scalar * dp[i]
	}

	return result
}
