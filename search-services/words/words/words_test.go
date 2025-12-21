package words

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNorm(t *testing.T) {
	testCases := []struct {
		name     string
		given    string
		expected []string
	}{
		{name: "empty", given: "", expected: []string{}},
		{name: "single word", given: "casing", expected: []string{"case"}},
		{name: "sentence", given: "rock you like a huricane", expected: []string{"rock", "like", "hurican"}},
		{name: "only stop words", given: "I and you or me or them, who will?", expected: []string{}},
		{name: "weird", given: "Moscow!123'check-it'or   123, man,that,difficult:heck", expected: []string{"moscow", "check", "123", "man", "difficult", "heck"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.ElementsMatch(t, testCase.expected, Norm(testCase.given))
		})
	}
}
