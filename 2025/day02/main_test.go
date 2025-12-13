package main

import "testing"

func Test_splitCount(t *testing.T) {
	type args struct {
		line    string
		invalid []string
	}
	tests := []struct {
		name string
		args args
	}{
		{"happy 31", args{"11-22,95-115,998-1012", []string{"11", "22"}}},
		{"happy #2", args{"1188511880-1188511890,222220-222224", []string{""}}},
		{"happy #3", args{"1698522-1698528,446443-446449,38593856-38593862", []string{""}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			splitCount(tt.args.line)
		})
	}
}
