package main

import "testing"

func TestGetRidOfWrongCharacters(t *testing.T) {
	tables := []struct {
		testCase string
		answer   string
	}{
		{"Hello world  peace", "Hello%20world%20%20peace"},
		{"HelloWorld", "HelloWorld"},
		{"Hello World", "Hello%20World"},
		{"Hello&World", "Hello%26World"},
	}
	for _, table := range tables {
		rightName := GetRidOfWrongCharacters(table.testCase)
		if rightName != table.answer {
			t.Errorf("RightName was incorrect, got: %s, want: %s.", rightName, table.answer)
		}
	}
}
