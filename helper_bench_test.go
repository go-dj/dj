package xn

import (
	"fmt"
	"testing"
)

func Benchmark_permute(b *testing.B) {
	for _, slice := range MapN(10, func(n int) []int { return RangeN(n + 1) }) {
		b.Run(fmt.Sprintf("size=%d", len(slice)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				permute(slice, factorial(len(slice)), 0)
			}
		})
	}
}
