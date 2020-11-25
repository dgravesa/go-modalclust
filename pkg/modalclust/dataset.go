package modalclust

import "math"

// StepDistThreshold is the threshold factor for continuing EM after one step
var StepDistThreshold float64 = 1e-01

// Dataset contains a dataset for executing MAC and MEM
type Dataset struct {
	data []DataPt
	N    int
}

// MakeDataset creates a dataset instance from a data point array
func MakeDataset(data []DataPt) *Dataset {
	if data == nil {
		return &Dataset{nil, 0}
	}
	return &Dataset{data, len(data)}
}

// MAC executes modal association clustering on a data slice
func (ds *Dataset) MAC(sigma float64) *Result {
	if ds.N == 0 {
		return nil
	}

	result := newResult()

	for i := 0; i < ds.N; i++ {
		mode := ds.MEM(ds.data[i], sigma)
		result.merge(ds.data[i], mode)
	}

	return result
}

// MEM executes expectation-maximization on a start coordinate and returns a local mode
func (ds *Dataset) MEM(start DataPt, sigma float64) DataPt {
	if ds.N == 0 {
		return nil
	}

	current := start
	p := make([]float64, ds.N)

	stepDist := math.MaxFloat64
	for stepDist > StepDistThreshold*sigma {
		// compute density impacts from each coordinate
		for i := 0; i < ds.N; i++ {
			distOverSig := current.Dist(ds.data[i]) / sigma
			p[i] = math.Exp(-0.5 * distOverSig * distOverSig)
		}

		psum := 0.0
		for i := 0; i < ds.N; i++ {
			psum += p[i]
		}

		// normalize p
		for i := 0; i < ds.N; i++ {
			p[i] /= psum
		}

		// compute next position
		next := ds.data[0].Scale(p[0])
		for i := 1; i < ds.N; i++ {
			nudge := ds.data[i].Scale(p[i])
			next = next.Add(nudge)
		}

		// compute distance traveled in step
		stepDist = current.Dist(next)
		current = next
	}

	return current
}