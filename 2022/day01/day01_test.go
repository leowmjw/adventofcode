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
			"10000",
		}}, 24000},
		{"weird format", args{[]string{
			"", "",
			"1000", "2000", "3000", "",
			"4000", "",
			"5000", "", "6000", "",
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

func Test_maxCalByTop3Elves(t *testing.T) {
	type args struct {
		input []string
	}
	tests := []struct {
		name             string
		args             args
		wantHighestTotal int
	}{
		//{"simple-single", args{[]string{"4000"}}, 4000},
		//{"simple-double", args{[]string{
		//	"4000", "",
		//	"6000", "",
		//}}, 10000},
		{"simple-triple", args{[]string{
			"4000", "",
			"6000", "",
			"10000", "",
		}}, 20000},
		{"sample", args{[]string{
			"1000", "2000", "3000", "",
			"4000", "",
			"5000", "6000", "",
			"7000", "8000", "9000", "",
			"10000",
		}}, 45000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHighestTotal := maxCalByTop3Elves(tt.args.input); gotHighestTotal != tt.wantHighestTotal {
				t.Errorf("maxCalByTop3Elves() = %v, want %v", gotHighestTotal, tt.wantHighestTotal)
			}
		})
	}
}
