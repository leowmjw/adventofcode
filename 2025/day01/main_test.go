package main

import (
	"testing"
)

func Test_run(t *testing.T) {
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
			run(tt.args.input)
		})
	}
}
