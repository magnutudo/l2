package main

import (
	"reflect"
	"testing"
)

func TestGetFields(t *testing.T) {
	type args struct {
		arr         [][]string
		leftBorder  int
		rightBorder int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "getFields_test_1",
			args: args{
				arr:         [][]string{{"asd", "fgh", "jkl"}, {"qwe", "rty", "uiop"}, {"zxc", "vbn", "m"}},
				leftBorder:  1,
				rightBorder: 2,
			},
			want: [][]string{{"fgh", "jkl"}, {"rty", "uiop"}, {"vbn", "m"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFields(tt.args.arr, tt.args.leftBorder, tt.args.rightBorder); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
