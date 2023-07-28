package collections

import (
	"reflect"
	"sort"
	"testing"
)

func TestChunkBy(t *testing.T) {
	type args[T Slice[N], N any] struct {
		items T
		chunk int
	}
	type testCase[T Slice[N], N any] struct {
		name       string
		args       args[T, N]
		wantChunks []T
	}
	var tests []testCase[[]int, int]
	tests = append(tests, testCase[[]int, int]{
		name: "ChunkBy",
		args: args[[]int, int]{
			items: []int{1, 2, 3, 4, 5, 6, 7, 8},
			chunk: 3,
		},
		wantChunks: [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8}},
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotChunks := ChunkBy(tt.args.items, tt.args.chunk); !reflect.DeepEqual(gotChunks, tt.wantChunks) {
				t.Errorf("ChunkBy() = %v, want %v", gotChunks, tt.wantChunks)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args[T comparable] struct {
		s []T
		e T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "Contains",
			args: args[int]{s: []int{1, 2, 3, 4, 5}, e: 3},
			want: true,
		},
		{
			name: "NotContains",
			args: args[int]{s: []int{1, 2, 3, 4, 5}, e: 6},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapKeys(t *testing.T) {
	type args[K comparable, V any] struct {
		m map[K]V
	}
	type testCase[K comparable, V any] struct {
		name    string
		args    args[K, V]
		wantOut []K
	}
	tests := []testCase[int, int]{
		{
			name:    "MapKeys",
			args:    args[int, int]{m: map[int]int{1: 1, 2: 2, 3: 3}},
			wantOut: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := MapKeys(tt.args.m)
			sort.Ints(gotOut)
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("MapKeys() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestMapValues(t *testing.T) {
	type args[K comparable, V any] struct {
		m map[K]V
	}
	type testCase[K comparable, V any] struct {
		name    string
		args    args[K, V]
		wantOut []V
	}
	tests := []testCase[int, int]{
		{
			name:    "MapValues",
			args:    args[int, int]{m: map[int]int{1: 1, 2: 2, 3: 3}},
			wantOut: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := MapValues(tt.args.m)
			sort.Ints(gotOut)
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("MapValues() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestMax(t *testing.T) {
	type args[T Ordered] struct {
		a []T
	}
	type testCase[T Ordered] struct {
		name  string
		args  args[T]
		wantM T
	}
	tests := []testCase[int]{
		{
			name:  "Max",
			args:  args[int]{a: []int{1, 2, 3, 4, 5}},
			wantM: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotM := Max(tt.args.a...); !reflect.DeepEqual(gotM, tt.wantM) {
				t.Errorf("Max() = %v, want %v", gotM, tt.wantM)
			}
		})
	}
}

func TestMin(t *testing.T) {
	type args[T Ordered] struct {
		a []T
	}
	type testCase[T Ordered] struct {
		name  string
		args  args[T]
		wantM T
	}
	tests := []testCase[int]{
		{
			name:  "Min",
			args:  args[int]{a: []int{3, 2, 1, 4, 5}},
			wantM: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotM := Min(tt.args.a...); !reflect.DeepEqual(gotM, tt.wantM) {
				t.Errorf("Min() = %v, want %v", gotM, tt.wantM)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	type args[T any, N any] struct {
		s         []T
		reduceFN  func(N, T) N
		initValue N
	}
	type testCase[T any, N any] struct {
		name string
		args args[T, N]
		want N
	}
	tests := []testCase[int, int]{
		{
			name: "Reduce",
			args: args[int, int]{s: []int{1, 2, 3, 4, 5},
				reduceFN: func(acc, e int) int {
					return acc + e
				},
				initValue: 0,
			},
			want: 15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reduce(tt.args.s, tt.args.reduceFN, tt.args.initValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceApply(t *testing.T) {
	type args[T any] struct {
		s       []T
		applyFN func(T) T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "SliceApply",
			args: args[int]{s: []int{1, 2, 3, 4, 5},
				applyFN: func(e int) int {
					return e * e
				},
			},
			want: []int{1, 4, 9, 16, 25},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceApply(tt.args.s, tt.args.applyFN); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceApply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceFilter(t *testing.T) {
	type args[T any] struct {
		s        []T
		filterFN func(T) bool
	}
	type testCase[T any] struct {
		name       string
		args       args[T]
		wantOutput []T
	}
	tests := []testCase[int]{
		{
			name: "SliceFilter",
			args: args[int]{s: []int{1, 2, 3, 4, 5},
				filterFN: func(e int) bool {
					return e%2 == 0
				},
			},
			wantOutput: []int{2, 4},
		},
		{
			name: "SliceNilFilter",
			args: args[int]{s: []int{1, 2, 3, 4, 5},
				filterFN: nil,
			},
			wantOutput: []int{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOutput := SliceFilter(tt.args.s, tt.args.filterFN); !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("SliceFilter() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestUniqueNonEmptyElementsOf(t *testing.T) {
	type args[T comparable] struct {
		s []T
	}
	type testCase[T comparable] struct {
		name    string
		args    args[T]
		wantOut []T
	}
	tests := []testCase[int]{
		{
			name:    "UniqueNonEmptyElementsOf",
			args:    args[int]{s: []int{1, 2, 3, 4, 5, 5, 4, 3, 2, 1}},
			wantOut: []int{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := UniqueNonEmptyElementsOf(tt.args.s); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("UniqueNonEmptyElementsOf() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
