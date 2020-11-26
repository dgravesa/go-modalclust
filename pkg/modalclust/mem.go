package modalclust

import "math"

// StepDistThreshold is the threshold factor for continuing EM after one step
var StepDistThreshold float64 = 1e-05

// MEM executes expectation-maximization on a start coordinate and returns a local mode
func MEM(data []DataPoint, start DataPoint, sigma float64) DataPoint {
	if data == nil || start == nil {
		return nil
	}

	N := len(data)

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
