package modalclust

import (
	"github.com/dgravesa/go-parallel/parallel"
)

// ModeDistThreshold is the allowable distance between two modes that they are still considered equivalent
var ModeDistThreshold float64 = 1e-01

// number of parallel goroutines to use in calculation; settable at runtime by SetNumGoroutines()
var numMACGoroutines int = parallel.DefaultNumGoroutines()

// SetNumGoroutines sets the number of goroutines to use in MAC computation.
// TODO: make this variable per MAC call.
func SetNumGoroutines(numGR int) {
	numMACGoroutines = numGR
}

// MAC executes modal association clustering on a data slice
func MAC(data []DataPt, sigma float64) *MACResult {
	if data == nil {
		return nil
	}

	results := newMACResult()
	executor := parallel.WithNumGoroutines(numMACGoroutines).
		WithStrategy(parallel.StrategyAtomicCounter)

	// execute MEM on each data point
	executor.For(len(data), func(i, _ int) {
		mode := MEM(data, data[i], sigma)
		results.insert(data[i], mode)
	})

	return results
}
