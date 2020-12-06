package modalclust_test

import (
	"testing"

	"github.com/dgravesa/go-modalclust/pkg/modalclust"
)

func BenchmarkMEM(b *testing.B) {
	start := sampleData100[50]

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = modalclust.MEM(sampleData100, start, 0.5)
	}
	b.StopTimer()
}
