package main

import "testing"

func TestGetRidOfWrongCharacters(t *testing.T) {
	tables := []struct {
		testCase string
		answer   string
	}{
		{"Hello world  peace", "Hello_world__peace"},
		{"HelloWorld", "HelloWorld"},
		{"Hello&World", "Hello_and_World"},
	}
	for _, table := range tables {
		rightName := GetRidOfWrongCharacters(table.testCase)
		if rightName != table.answer {
			t.Errorf("RightName was incorrect, got: %s, want: %s.", rightName, table.answer)
		}
	}
}
