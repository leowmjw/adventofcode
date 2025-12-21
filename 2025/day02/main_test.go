package main

import "testing"

func Test_isRepeatingPattern(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"11 - 1 twice", "11", true},
		{"22 - 2 twice", "22", true},
		{"99 - 9 twice", "99", true},
		{"111 - 1 three times", "111", true},
		{"999 - 9 three times", "999", true},
		{"1010 - 10 twice", "1010", true},
		{"6464 - 64 twice", "6464", true},
		{"123123 - 123 twice", "123123", true},
		{"12341234 - 1234 twice", "12341234", true},
		{"123123123 - 123 three times", "123123123", true},
		{"1212121212 - 12 five times", "1212121212", true},
		{"1111111 - 1 seven times", "1111111", true},
		{"565656 - 56 three times", "565656", true},
		{"222222 - 2 six times or 22 three times", "222222", true},
		{"446446 - 446 twice", "446446", true},
		{"1188511885 - 11885 twice", "1188511885", true},
		{"824824824 - 824 three times", "824824824", true},
		{"2121212121 - 21 five times", "2121212121", true},
		{"101 - not repeating", "101", false},
		{"123 - not repeating", "123", false},
		{"1234 - not repeating", "1234", false},
		{"12345 - not repeating", "12345", false},
		{"565653 - not repeating", "565653", false},
		{"1698522 - not repeating", "1698522", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRepeatingPattern(tt.s); got != tt.want {
				t.Errorf("isRepeatingPattern(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func Test_findRepeatedAtLeastTwice(t *testing.T) {
	tests := []struct {
		name   string
		start  int
		finish int
		want   int
	}{
		{"11-22: has 11 and 22", 11, 22, 11 + 22},
		{"95-115: has 99 and 111", 95, 115, 99 + 111},
		{"998-1012: has 999 and 1010", 998, 1012, 999 + 1010},
		{"1188511880-1188511890: has 1188511885", 1188511880, 1188511890, 1188511885},
		{"222220-222224: has 222222", 222220, 222224, 222222},
		{"1698522-1698528: no invalid IDs", 1698522, 1698528, 0},
		{"446443-446449: has 446446", 446443, 446449, 446446},
		{"38593856-38593862: has 38593859", 38593856, 38593862, 38593859},
		{"565653-565659: has 565656", 565653, 565659, 565656},
		{"824824821-824824827: has 824824824", 824824821, 824824827, 824824824},
		{"2121212118-2121212124: has 2121212121", 2121212118, 2121212124, 2121212121},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findRepeatedAtLeastTwice(tt.start, tt.finish); got != tt.want {
				t.Errorf("findRepeatedAtLeastTwice(%d, %d) = %v, want %v", tt.start, tt.finish, got, tt.want)
			}
		})
	}
}

func Test_part2_example_total(t *testing.T) {
	// From PART2.md: Adding up all the invalid IDs in this example produces 4174379265.
	ranges := [][2]int{
		{11, 22},
		{95, 115},
		{998, 1012},
		{1188511880, 1188511890},
		{222220, 222224},
		{1698522, 1698528},
		{446443, 446449},
		{38593856, 38593862},
		{565653, 565659},
		{824824821, 824824827},
		{2121212118, 2121212124},
	}

	total := 0
	for _, r := range ranges {
		total += findRepeatedAtLeastTwice(r[0], r[1])
	}

	want := 4174379265
	if total != want {
		t.Errorf("Total invalid IDs = %v, want %v", total, want)
	}
}
