package main

import (
	"testing"
)

func Test_part1(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
	}{
		{"happy", args{"test.txt"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			part1(tt.args.input)
		})
	}
}

func Test_part2(t *testing.T) {
	type args struct {
		input string
	}
	// 50 + R50 (count+1) + R100 (count+1) + R100 (count+1)  + R9 =
	tests := []struct {
		name string
		args args
	}{
		{"happy", args{"part2.txt"}},
		{"happy2", args{"part2a.txt"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			part2(tt.args.input)
		})
	}
}
