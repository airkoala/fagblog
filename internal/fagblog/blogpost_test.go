package fagblog

import (
	"testing"
)

func TestNormaliseHeadings(t *testing.T) {
	tests := []struct {
		name     string
		input    []Heading
		expected []Heading
	}{
		{
			name: "consecutive levels",
			input: []Heading{
				{Title: "H1", Level: 0},
				{Title: "H2", Level: 1},
				{Title: "H3", Level: 2},
			},
			expected: []Heading{
				{Title: "H1", Level: 0},
				{Title: "H2", Level: 1},
				{Title: "H3", Level: 2},
			},
		},
		{
			name: "non-consecutive levels",
			input: []Heading{
				{Title: "H1", Level: 1},
				{Title: "H2", Level: 3},
				{Title: "H3", Level: 5},
			},
			expected: []Heading{
				{Title: "H1", Level: 0},
				{Title: "H2", Level: 1},
				{Title: "H3", Level: 2},
			},
		},
		{
			name: "duplicate levels",
			input: []Heading{
				{Title: "H1", Level: 2},
				{Title: "H2", Level: 2},
				{Title: "H3", Level: 4},
			},
			expected: []Heading{
				{Title: "H1", Level: 0},
				{Title: "H2", Level: 0},
				{Title: "H3", Level: 1},
			},
		},
		{
			name: "single heading",
			input: []Heading{
				{Title: "H1", Level: 5},
			},
			expected: []Heading{
				{Title: "H1", Level: 0},
			},
		},
		{
			name:     "empty headings",
			input:    []Heading{},
			expected: []Heading{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying the test case
			headings := make([]Heading, len(tt.input))
			copy(headings, tt.input)
			
			normaliseHeadings(headings)
			
			if len(headings) != len(tt.expected) {
				t.Fatalf("Expected %d headings, got %d", len(tt.expected), len(headings))
			}
			
			for i := range headings {
				if headings[i].Level != tt.expected[i].Level {
					t.Errorf("Heading %d: expected level %d, got %d", i, tt.expected[i].Level, headings[i].Level)
				}
				if headings[i].Title != tt.expected[i].Title {
					t.Errorf("Heading %d: expected title %s, got %s", i, tt.expected[i].Title, headings[i].Title)
				}
			}
		})
	}
}
