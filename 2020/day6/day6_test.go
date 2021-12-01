package day6

import (
	"reflect"
	"testing"
)

func Test_splitCustomDeclarationForms(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"sample", args{"testdata/sample.txt"}, []string{
			"abc",
			"abc",
			"abac",
			"aaaa",
			"b",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitCustomDeclarationForms(tt.args.filename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitCustomDeclarationForms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPart1(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"sample", args{"testdata/sample.txt"}, 11},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Part1(tt.args.filename); got != tt.want {
				t.Errorf("Part1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_countUniqueYesPerForm(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"sample #1", args{"abc"}, 3},
		{"sample #2", args{"abac"}, 3},
		{"sample #3", args{"aaaa"}, 1},
		{"sample #3", args{"b"}, 1},
		{"sample #4", args{"rurrrbwr"}, 4},
		{"sample #5", args{"jjjjcj"}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countUniqueYesPerForm(tt.args.content); got != tt.want {
				t.Errorf("countUniqueYesPerForm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitCustomDeclarationFormsByIndividual(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"sample", args{"testdata/sample.txt"}, []string{
			"abc",
			"a b c",
			"ab ac",
			"a a a a",
			"b",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitCustomDeclarationFormsByIndividual(tt.args.filename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitCustomDeclarationFormsByIndividual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPart2(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"sample", args{"testdata/sample.txt"}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Part2(tt.args.filename); got != tt.want {
				t.Errorf("Part2() = %v, want %v", got, tt.want)
			}
		})
	}
}
