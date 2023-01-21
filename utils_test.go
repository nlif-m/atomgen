package main

import (
	"fmt"
	"testing"
)

func TestUnique(t *testing.T) {
	type testCase[T comparable] struct {
		input []T
		want  []T
	}
	tests := []testCase[string]{
		{[]string{"Hello", "Hello"}, []string{"Hello"}},
		{[]string{"Hello"}, []string{"Hello"}},
		{[]string{"Hello", "Hello", "Hi"}, []string{"Hello", "Hi"}},
		{[]string{""}, []string{""}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s,%s", test.input, test.want)
		t.Run(testname, func(t *testing.T) {
			answer := Unique(test.input)
			if !Compare(answer, test.want) {
				t.Errorf("got %v, want %v", answer, test.want)
			}

		})
	}

}
