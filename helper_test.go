package dj

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermute(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want [][]int
	}{
		{
			name: "single",
			in:   []int{1},
			want: [][]int{{1}},
		},

		{
			name: "double",
			in:   []int{1, 2},
			want: [][]int{{1, 2}, {2, 1}},
		},

		{
			name: "triple",
			in:   []int{1, 2, 3},
			want: [][]int{{1, 2, 3}, {2, 1, 3}, {3, 1, 2}, {1, 3, 2}, {2, 3, 1}, {3, 2, 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.ElementsMatch(t, tt.want, permute(tt.in, factorial(len(tt.in)), 0))
		})
	}
}
