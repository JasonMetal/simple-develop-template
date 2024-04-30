package slices

import (
	"math"
	"testing"
)

var tests = []struct {
	s    []int
	v    int
	want int
}{
	{
		nil,
		0,
		-1,
	},
	{
		[]int{},
		0,
		-1,
	},
	{
		[]int{1, 2, 3},
		2,
		1,
	},
	{
		[]int{1, 2, 2, 3},
		2,
		1,
	},
	{
		[]int{1, 2, 3, 2},
		2,
		1,
	},
}

var equalIntTests = []struct {
	s1, s2 []int
	want   bool
}{
	{
		[]int{1},
		nil,
		false,
	},
	{
		[]int{},
		nil,
		true,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3},
		true,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3, 4},
		false,
	},
}

var equalFloatTests = []struct {
	s1, s2       []float64
	wantEqual    bool
	wantEqualNaN bool
}{
	{
		[]float64{1, 2},
		[]float64{1, 2},
		true,
		true,
	},
	{
		[]float64{1, 2, math.NaN()},
		[]float64{1, 2, math.NaN()},
		false,
		true,
	},
}

var compareIntTests = []struct {
	s1, s2 []int
	want   int
}{
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3, 4},
		-1,
	},
	{
		[]int{1, 2, 3, 4},
		[]int{1, 2, 3},
		+1,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 4, 3},
		-1,
	},
	{
		[]int{1, 4, 3},
		[]int{1, 2, 3},
		+1,
	},
}

var compareFloatTests = []struct {
	s1, s2 []float64
	want   int
}{
	{
		[]float64{1, 2, math.NaN()},
		[]float64{1, 2, math.NaN()},
		0,
	},
	{
		[]float64{1, math.NaN(), 3},
		[]float64{1, math.NaN(), 4},
		-1,
	},
	{
		[]float64{1, math.NaN(), 3},
		[]float64{1, 2, 4},
		-1,
	},
	{
		[]float64{1, math.NaN(), 3},
		[]float64{1, 2, math.NaN()},
		0,
	},
	{
		[]float64{1, math.NaN(), 3, 4},
		[]float64{1, 2, math.NaN()},
		+1,
	},
}

func TestCompare(t *testing.T) {
	intWant := func(want bool) string {
		if want {
			return "0"
		}
		return "!= 0"
	}
	for _, test := range equalIntTests {
		if got := Compare(test.s1, test.s2); (got == 0) != test.want {
			t.Errorf("Compare(%v, %v) = %d, want %s", test.s1, test.s2, got, intWant(test.want))
		}
	}
	for _, test := range equalFloatTests {
		if got := Compare(test.s1, test.s2); (got == 0) != test.wantEqualNaN {
			t.Errorf("Compare(%v, %v) = %d, want %s", test.s1, test.s2, got, intWant(test.wantEqualNaN))
		}
	}

	for _, test := range compareIntTests {
		if got := Compare(test.s1, test.s2); got != test.want {
			t.Errorf("Compare(%v, %v) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}

	for _, test := range compareFloatTests {
		if got := Compare(test.s1, test.s2); got != test.want {
			t.Errorf("Compare(%v, %v) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}
}

func TestContains(t *testing.T) {
	for _, tt := range tests {
		if got := Contains(tt.s, tt.v); got != (tt.want != -1) {
			t.Errorf("Contains() = %v, want %v", got, tt.want)
		}
	}
}

func TestIndex(t *testing.T) {
	for _, test := range tests {
		if got := Index(test.s, test.v); got != test.want {
			t.Errorf("Index(%v, %v) = %d, want %d", test.s, test.v, got, test.want)
		}
	}
}
