package dj_test

import (
	"fmt"
	"testing"

	"github.com/go-dj/dj"
)

func Benchmark_Power(b *testing.B) {
	for _, slice := range dj.MapN(10, func(i int) []int { return dj.RangeN(i + 1) }) {
		b.Run(fmt.Sprintf("size=%d", len(slice)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				dj.Power(slice)
			}
		})
	}
}
