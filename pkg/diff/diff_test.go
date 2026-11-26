package diff

import (
	"strings"
	"testing"
)

func TestUnified(t *testing.T) {
	tests := []struct {
		name     string
		from     string
		to       string
		opts     Options
		expected []string
	}{
		{
			name: "identical strings produce empty diff",
			from: "apiVersion: v1\nkind: Pod\n",
			to:   "apiVersion: v1\nkind: Pod\n",
			opts: Options{Color: false, Context: 3},
		},
		{
			name: "detects deleted lines",
			from: "apiVersion: v1\nkind: Pod\nmetadata:\n  creationTimestamp: \"2026\"\n  name: mypod\n",
			to:   "apiVersion: v1\nkind: Pod\nmetadata:\n  name: mypod\n",
			opts: Options{Color: false, Context: 1},
			expected: []string{
				"--- raw",
				"+++ neat",
				"-  creationTimestamp: \"2026\"",
				"   name: mypod",
			},
		},
		{
			name: "colored output includes ANSI codes",
			from: "line1\nline2\n",
			to:   "line1\nline3\n",
			opts: Options{Color: true, Context: 1},
			expected: []string{
				"\033[31m-line2\033[0m",
				"\033[32m+line3\033[0m",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unified("raw", "neat", tt.from, tt.to, tt.opts)
			if len(tt.expected) == 0 {
				if got != "" {
					t.Errorf("expected empty diff, got:\n%s", got)
				}
				return
			}
			for _, exp := range tt.expected {
				if !strings.Contains(got, exp) {
					t.Errorf("expected output to contain %q, but got:\n%s", exp, got)
				}
			}
		})
	}
}
