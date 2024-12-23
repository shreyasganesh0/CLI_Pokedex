package main

import "testing"

func TestCleanInput(t *testing.T){
    cases := []struct {
	input    string
	expected []string
    }{
	{
		input:    "  hello  world  ",
		expected: []string{"hello", "world"},
	},
    }

for _, c := range cases {
	actual := cleanInput(c.input)
    actual_len := len(actual)
    if actual_len != len(c.expected){
        t.Errorf("Sizes dont match between %v and %v", actual, c.expected)
        }
	for i := range actual {
		word := actual[i]
		expectedWord := c.expected[i]
        if expectedWord != word{
            t.Errorf("Word %v doesnt match %v", word, expectedWord)
        }
	}
}
}
