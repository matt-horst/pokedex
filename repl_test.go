package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "	this 	is a 	test",
			expected: []string{"this", "is", "a", "test"},
		},
		{
			input: "\n\nthis is\tanother\t\n",
			expected: []string{"this", "is", "another"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("actaul: %v !=  expected: %v", actual, c.expected)
			t.FailNow()
		}

		for i := range actual {
			actualWord := actual[i]
			expectedWord := c.expected[i]

			if actualWord != expectedWord {
				t.Errorf("actaul: %v !=  expected: %v", actual, c.expected)
				t.FailNow()
			}
		}
	}
}
