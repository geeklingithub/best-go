package list

import (
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestArrayListAppend(t *testing.T) {
	testCases := []struct {
		name      string
		list      *ArrayList[int]
		newVal    int
		wantSlice []int
	}{
		{
			name:      "append 234",
			list:      NewWithE[int]([]int{123}),
			newVal:    234,
			wantSlice: []int{123, 234},
		},
		{
			name:      "nil append 123",
			list:      New[int](0),
			newVal:    123,
			wantSlice: []int{123},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Append(tc.newVal)
			if err != nil {
				return
			}

			assert.Equal(t, tc.wantSlice, tc.list.elements)
		})
	}
}
