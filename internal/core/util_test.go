package core

import "testing"

func TestIndent(t *testing.T) {
	cases := []struct {
		name        string
		input       string
		n           int
		firstIndent int
		expected    string
	}{
		{
			name:        "single line, no indent",
			input:       "foo",
			n:           0,
			firstIndent: -1,
			expected:    "foo",
		},
		{
			name:        "single line, indent",
			input:       "foo",
			n:           2,
			firstIndent: -1,
			expected:    "  foo",
		},
		{
			name:        "multi-line, uniform indent",
			input:       "foo\nbar",
			n:           2,
			firstIndent: -1,
			expected:    "  foo\n  bar",
		},
		{
			name:        "multi-line, firstIndent",
			input:       "foo\nbar",
			n:           2,
			firstIndent: 4,
			expected:    "    foo\n  bar",
		},
		{
			name:        "multi-line, firstIndent zero",
			input:       "foo\nbar",
			n:           2,
			firstIndent: 0,
			expected:    "foo\n  bar",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := indent(tc.input, tc.n, tc.firstIndent)
			if got != tc.expected {
				t.Errorf("indent() = %q, want %q", got, tc.expected)
			}
		})
	}
}
