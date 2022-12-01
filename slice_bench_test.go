package xn_test

import (
	"fmt"
	"testing"

	"github.com/jameshoulahan/xn"
)

func Benchmark_Power(b *testing.B) {
	for _, slice := range xn.MapN(10, func(i int) []int { return xn.RangeN(i + 1) }) {
		b.Run(fmt.Sprintf("size=%d", len(slice)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				xn.Power(slice)
			}
		})
	}
}
