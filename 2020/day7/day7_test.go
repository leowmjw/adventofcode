package day7

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestPart1(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Part1(tt.args.filename); got != tt.want {
				t.Errorf("Part1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitRules(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"sample", args{"testdata/sample.txt"}, []string{
			"light red;1:bright white,2:muted yellow",
			"dark orange;3:bright white,4:muted yellow",
			"bright white;1:shiny gold",
			"muted yellow;2:shiny gold,9:faded blue",
			"shiny gold;1:dark olive,2:vibrant plum",
			"dark olive;3:faded blue,4:dotted black",
			"vibrant plum;5:faded blue,6:dotted black",
			"faded blue;0:none",
			"dotted black;0:none",
		}},
		{"transitive-muted-crimson", args{"testdata/transitive.txt"}, []string{
			"wavy purple;2:shiny gold,4:mirrored maroon",
			"pale magenta;2:muted orange,4:muted crimson,4:striped turquoise",
			"muted fuchsia;4:plaid indigo,2:shiny gold",
			"bob blue;2:vibrant fuchsia",
			"shiny red;2:clear tomato,4:muted crimson,2:plaid gray,1:bright gold",
			"mirrored tan;4:pale orange,3:plaid olive,5:muted crimson,5:posh salmon",
			"shiny gold;4:posh coral,2:clear violet",
			"clear green;3:muted crimson",
			"mirrored plum;1:bright chartreuse,4:mirrored purple,1:dim turquoise,4:shiny gold",
			"muted crimson;3:mirrored coral,4:light silver,2:shiny gold",
			"vibrant fuchsia;3:dim cyan,2:muted crimson",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitRules(tt.args.filename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitRules() = %v, want %v", got, tt.want)
				fmt.Println("GOT:")
				spew.Dump(got)
				fmt.Println("WANT:")
				spew.Dump(tt.want)
			}
		})
	}
}

func Test_countBagContainsShinyGold(t *testing.T) {
	type args struct {
		allRules []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"sample", args{[]string{
			"light red;1:bright white,2:muted yellow",
			"dark orange;3:bright white,4:muted yellow",
			"bright white;1:shiny gold",
			"muted yellow;2:shiny gold,9:faded blue",
			"shiny gold;1:dark olive,2:vibrant plum",
			"dark olive;3:faded blue,4:dotted black",
			"vibrant plum;5:faded blue,6:dotted black",
			"faded blue;0:none",
			"dotted black;0:none",
		}}, 4},
		{"transitive-muted-crimson", args{[]string{
			"wavy purple;2:shiny gold,4:mirrored maroon",
			"pale magenta;2:muted orange,4:muted crimson,4:striped turquoise",
			"muted fuchsia;4:plaid indigo,2:shiny gold",
			"bob blue;2:vibrant fuchsia",
			"shiny red;2:clear tomato,4:muted crimson,2:plaid gray,1:bright gold",
			"mirrored tan;4:pale orange,3:plaid olive,5:muted crimson,5:posh salmon",
			"shiny gold;4:posh coral,2:clear violet",
			"clear green;3:muted crimson",
			"mirrored plum;1:bright chartreuse,4:mirrored purple,1:dim turquoise,4:shiny gold",
			"muted crimson;3:mirrored coral,4:light silver,2:shiny gold",
			"vibrant fuchsia;3:dim cyan,2:muted crimson",
		}}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countBagContainsShinyGold(tt.args.allRules); got != tt.want {
				t.Errorf("countBagContainsShinyGold() = %v, want %v", got, tt.want)
			}
		})
	}
}
