package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetSets(t *testing.T) {
	type args struct {
		dictionary []string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "test_1",
			args: args{
				[]string{"пятак", "листок", "тяпка", "столик", "пятка", "слиток"},
			},
			want: map[string][]string{
				"листок": {"листок", "слиток", "столик"},
				"пятак":  {"пятак", "пятка", "тяпка"},
			},
		},
		{
			name: "test_2",
			args: args{
				[]string{"1122", "2121", "1221", "4334", "3344", "3434"},
			},
			want: map[string][]string{
				"1122": {"1122", "1221", "2121"},
				"4334": {"3344", "3434", "4334"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSets(tt.args.dictionary); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSets() = %v, want %v", got, tt.want)
			}
			fmt.Println("want ", tt.want)
		})
	}
}
