package day4

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_splitPassports(t *testing.T) {
	sample, err := os.Open("testdata/sample.txt")
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer sample.Close()

	type args struct {
		in io.ReadCloser
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"sample", args{sample}, []string{
			"ecl:gry pid:860033327 eyr:2020 hcl:#fffffd byr:1937 iyr:2017 cid:147 hgt:183cm",
			"iyr:2013 ecl:amb cid:350 eyr:2023 pid:028048884 hcl:#cfa07d byr:1929",
			"hcl:#ae17e1 iyr:2013 eyr:2024 ecl:brn pid:760753108 byr:1931 hgt:179cm",
			"hcl:#cfa07d eyr:2025 pid:166559648 iyr:2011 ecl:brn hgt:59in",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitPassports(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitPassports() = %v, want %v", got, tt.want)
				fmt.Println("GOT:")
				spew.Dump(got)
				fmt.Println("WANT:")
				spew.Dump(tt.want)
			}
		})
	}
}

func Test_countValidPassport(t *testing.T) {
	type args struct {
		content []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"sample", args{[]string{
			"ecl:gry pid:860033327 eyr:2020 hcl:#fffffd byr:1937 iyr:2017 cid:147 hgt:183cm",
			"iyr:2013 ecl:amb cid:350 eyr:2023 pid:028048884 hcl:#cfa07d byr:1929",
			"hcl:#ae17e1 iyr:2013 eyr:2024 ecl:brn pid:760753108 byr:1931 hgt:179cm",
			"hcl:#cfa07d eyr:2025 pid:166559648 iyr:2011 ecl:brn hgt:59in",
		}}, 2},
		{"sample-type", args{[]string{
			"cid:94 byr:1934 hgt:59in eyr:2022 hcl:#623a2f pid:573884719 iyr:2016 ecl:oth",  // OK
			"pid:206185815 ecl:grn hcl:#cfa07d eyr:2027 iyr:2018 byr:1989 hgt:176cm",        // OK
			"hgt:175cm byr:1999 pid:409477026 hcl:#cfa07d ecl:amb eyr:2021 iyr:2017 cid:75", // OK
			"ecl:#564a01 hgt:136 iyr:1984 pid:#646419 eyr:2032 hcl:hzl",                     // Less field, not OK
			"hcl:#ae17e1 iyr:2013 eyr:2024 ecl:brn pid:760753108 byr:1931 hgt:179cm",        // OK
		}}, 4},
		{"false-negative", args{[]string{
			"iyr:1928 cid:150 pid:476113241 eyr:2039 hcl:a5ac0f ecl:#25f8d2 byr:2027 hgt:190", // ecl no #; Should be NOT OK
			"hcl:#007d7c pid:195125455 cid:213 hgt:154cm eyr:2021 ecl:grn byr:1981",           // Should be Rejected, has extra cid
			"ecl:oth hgt:185cm pid:958027833 hcl:#b6652a iyr:2028 cid:171 eyr:1994",           // Should be Rejected, has extra cid
			"iyr:2006 hgt:103 ecl:#2d77e5 cid:214 byr:2018 hcl:6c53a4 pid:190cm eyr:1940",     // ecl no #; hcl needs #; Should be NOT OK
			"ecl:gmt hgt:75cm byr:2007 eyr:2037 iyr:2028 hcl:4591f6 cid:118",                  // Should be Rejected, has extra cid
			"ecl:hzl eyr:2027 iyr:2019 pid:125201586 byr:1947 cid:74 hcl:#341e13",             // Should be Rejected, has extra cid
		}}, 0},
		{"false-positive", args{[]string{
			"pid:302395756 ecl:grn hcl:z byr:2005 hgt:111 eyr:2031 cid:147",         // missing field
			"hgt:172cm byr:1923 pid:741415636 ecl:grn eyr:2022 iyr:2013",            // missing field
			"pid:457776708 byr:1992 hcl:#b6652a hgt:157cm eyr:2024 iyr:2011",        // missing field
			"hcl:#ceb3a1 iyr:2013 pid:592603167 cid:95 ecl:blu eyr:2022",            // missing field
			"hgt:150cm ecl:grn hcl:8f3824 pid:735766540 eyr:2029 byr:2000 iyr:2015", // Don't Accept hcl without: #
			"hcl:z ecl:hzl byr:2003 hgt:118 eyr:2008 iyr:2022 pid:157cm",            // Don't Accept hcl: z
		}}, 0},
		{"sample-invalid", args{[]string{
			"eyr:1972 cid:100 hcl:#18171d ecl:amb hgt:170 pid:186cm iyr:2018 byr:1926",
			"iyr:2019 hcl:#602927 eyr:1967 hgt:170cm ecl:grn pid:012533040 byr:1946",
			"hcl:dab227 iyr:2012 ecl:brn hgt:182cm pid:021572410 eyr:2020 byr:1992 cid:277",
			"hgt:59cm ecl:zzz eyr:2038 hcl:74454a iyr:2023 pid:3556412378 byr:2007",
			"hgt:74in ecl:grn pid:1655089174 iyr:2012 eyr:2030 byr:1980 hcl:#623a2f",
		}}, 0},
		{"sample-valid", args{[]string{
			"pid:087499704 hgt:74in ecl:grn iyr:2012 eyr:2030 byr:1980 hcl:#623a2f",
			"eyr:2029 ecl:blu cid:129 byr:1989 iyr:2014 pid:896056539 hcl:#a97842 hgt:165cm",
			"hcl:#888785 hgt:164cm byr:2001 iyr:2015 cid:88 pid:545766238 ecl:hzl eyr:2022",
			"iyr:2010 hgt:158cm hcl:#b6652a ecl:blu byr:1944 eyr:2021 pid:093154719",
		}}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countValidPassport(tt.args.content); got != tt.want {
				t.Errorf("countValidPassport() = %v, want %v", got, tt.want)
				fmt.Println("GOT:")
				spew.Dump(got)
				fmt.Println("WANT:")
				spew.Dump(tt.want)

			}
		})
	}
}
