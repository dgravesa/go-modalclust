package modalclust

import (
	"math"
)

// StepDistThreshold is the threshold factor for continuing EM after one step
var StepDistThreshold float64 = 1e-05

// MEM executes expectation-maximization on a start coordinate and returns a local mode
func MEM(data []DataPt, start DataPt, sigma float64) DataPt {
	if data == nil || start == nil {
		return nil
	}

	N := len(data)
	p := make([]float64, N)

	dim := len(start)
	current := make(DataPt, dim)
	next := make(DataPt, dim)
	nudge := make(DataPt, dim)

	start.copyTo(&current)

	stepDist := math.MaxFloat64
	for stepDist > StepDistThreshold*sigma {
		// compute density impacts from each coordinate
		for i := 0; i < N; i++ {
			distOverSig := current.dist(data[i]) / sigma
			p[i] = math.Exp(-0.5 * distOverSig * distOverSig)
		}

		psum := 0.0
		for i := 0; i < N; i++ {
			psum += p[i]
		}

		// normalize p
		for i := 0; i < N; i++ {
			p[i] /= psum
		}

		// compute next position
		data[0].storeScale(p[0], &next)
		for i := 1; i < N; i++ {
			data[i].storeScale(p[i], &nudge)
			next.storeAdd(nudge, &next)
		}

		// compute distance traveled in step
		stepDist = current.dist(next)
		next.copyTo(&current)
	}

	return current
}
