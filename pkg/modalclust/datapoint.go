package modalclust

import "math"

// DataPt represents a float64-based N-dimensional data point.
type DataPt []float64

// dist returns the distance between two coordinates
func (dp DataPt) dist(other DataPt) float64 {
	distSq := 0.0

	for i := 0; i < len(dp); i++ {
		diff := dp[i] - other[i]
		distSq += diff * diff
	}

	return math.Sqrt(distSq)
}

// storeAdd computes the result of dp + other and stores the outcome in result.
// This function assumes that result has already been allocated and has same dimension as dp.
func (dp DataPt) storeAdd(other DataPt, result *DataPt) {
	for i := 0; i < len(dp); i++ {
		(*result)[i] = dp[i] + other[i]
	}
}

// storeScale computes the result of scalar * dp and stores the outcome in result.
// This function assumes that result has already been allocated and has same dimension as dp.
func (dp DataPt) storeScale(scalar float64, result *DataPt) {
	for i := 0; i < len(dp); i++ {
		(*result)[i] = scalar * dp[i]
	}
}

// copyTo copies dp into other.
// This function assumes that other has already been allocated and has same dimension as dp.
func (dp DataPt) copyTo(other *DataPt) {
	for i := 0; i < len(dp); i++ {
		(*other)[i] = dp[i]
	}
}
