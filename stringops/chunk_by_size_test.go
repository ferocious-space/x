package stringops

import (
	"reflect"
	"testing"
)

func TestChunkSize(t *testing.T) {
	type args struct {
		items []string
		size  int
	}
	tests := []struct {
		name       string
		args       args
		wantChunks [][]string
	}{
		{
			name: "ChunkSize",
			args: args{
				items: []string{"1", "2", "3", "4", "5", "6", "7", "8"},
				size:  3,
			},
			wantChunks: [][]string{{"1", "2", "3"}, {"4", "5", "6"}, {"7", "8"}},
		},
		{
			name: "ChunkSizelong",
			args: args{
				items: []string{"1", "22", "333", "4444", "55555", "666666", "7777777", "88888888"},
				size:  10,
			},
			wantChunks: [][]string{{"1", "22", "333", "4444"}, {"55555"}, {"666666"}, {"7777777"}, {"88888888"}},
		},
		{
			name: "ChunkSizeShort",
			args: args{
				items: []string{},
			},
			wantChunks: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotChunks := ChunkSize(tt.args.items, tt.args.size); !reflect.DeepEqual(gotChunks, tt.wantChunks) {
				t.Errorf("ChunkSize() = %v, want %v", gotChunks, tt.wantChunks)
			}
		})
	}
}

func Test_listLen(t *testing.T) {
	type args struct {
		items []string
	}
	tests := []struct {
		name    string
		args    args
		wantSum int
	}{
		{
			name: "listLen",
			args: args{
				items: []string{"1", "2", "3", "4", "5", "6", "7", "8"},
			},
			wantSum: 8,
		},
		{
			name: "listLenlong",
			args: args{
				items: []string{"1", "22", "333", "4444", "55555", "666666", "7777777", "88888888"},
			},
			wantSum: 36,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSum := listLen(tt.args.items); gotSum != tt.wantSum {
				t.Errorf("listLen() = %v, want %v", gotSum, tt.wantSum)
			}
		})
	}
}
