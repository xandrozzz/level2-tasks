package main

import "testing"

type unpackTest struct {
	arg, expected string
	error         bool
}

var unpackTests = []unpackTest{
	{"a4bc2d5e", "aaaabccddddde", false},
	{"abcd", "abcd", false},
	{"45", "", true},
	{"", "", false},
}

func TestUnpack(t *testing.T) {
	for _, test := range unpackTests {
		if output, err := Unpack(test.arg); output != test.expected || (err != nil) != test.error {
			t.Errorf("Output %v, %v was not equal to expected %v, %v", output, err != nil, test.expected, test.error)
		}
	}
}
