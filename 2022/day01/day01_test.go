package main

import "testing"

func Test_maxCalBySingleElf(t *testing.T) {
	type args struct {
		input []string
	}
	tests := []struct {
		name             string
		args             args
		wantHighestTotal int
	}{
		{"sample", args{[]string{
			"1000", "2000", "3000", "",
			"4000", "",
			"5000", "6000", "",
			"7000", "8000", "9000", "",
			"10000", "",
		}}, 24000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHighestTotal := maxCalBySingleElf(tt.args.input); gotHighestTotal != tt.wantHighestTotal {
				t.Errorf("maxCalBySingleElf() = %v, want %v", gotHighestTotal, tt.wantHighestTotal)
			}
		})
	}
}
