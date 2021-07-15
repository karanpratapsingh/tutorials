package length

import "testing"

var testCases = []struct {
	input    string
	expected int
}{
	{
		input:    "hey",
		expected: 3,
	},
	{
		input:    "hi",
		expected: 2,
	},
}

func Test_Length(t *testing.T) {
	for _, tc := range testCases {
		result := Length(tc.input)

		if result != tc.expected {
			t.Error("Expected", tc.expected, "Got", result)
		}
	}
}
