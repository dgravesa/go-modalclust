package modalclust

import "math"

// StepDistThreshold is the threshold factor for continuing EM after one step
var StepDistThreshold float64 = 1e-01

// MAC executes modal association clustering on a data slice
func MAC(data []Coord, sigma float64) *Result {
	N := len(data)
	result := newResult()

	for i := 0; i < N; i++ {
		mode := EM(data, data[i], sigma)
		result.merge(data[i], mode)
	}

	return result
}

// EM executes expectation-maximization on a start coordinate and returns a local mode
func EM(data []Coord, start Coord, sigma float64) Coord {
	N := len(data)
	if N == 0 {
		return nil
	}

	current := start
	p := make([]float64, N)

	stepDist := math.MaxFloat64
	for stepDist > StepDistThreshold*sigma {
		// compute density impacts from each coordinate
		for i := 0; i < N; i++ {
			distOverSig := current.Dist(data[i]) / sigma
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
		next := data[0].Scale(p[0])
		for i := 1; i < N; i++ {
			nudge := data[i].Scale(p[i])
			next = next.Add(nudge)
		}

		// compute distance traveled in step
		stepDist = current.Dist(next)
		current = next
	}

	return current
}
