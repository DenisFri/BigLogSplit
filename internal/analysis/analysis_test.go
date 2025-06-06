package analysis

import (
	"regexp"
	"testing"
)

func TestFilterLine(t *testing.T) {
	reFoo := regexp.MustCompile("foo")
	tests := []struct {
		name     string
		line     string
		mode     string
		patterns []*regexp.Regexp
		want     bool
	}{
		{"none mode", "anything", "none", []*regexp.Regexp{reFoo}, true},
		{"include match", "foo bar", "include", []*regexp.Regexp{reFoo}, true},
		{"include no match", "bar", "include", []*regexp.Regexp{reFoo}, false},
		{"exclude match", "foo bar", "exclude", []*regexp.Regexp{reFoo}, false},
		{"exclude no match", "bar", "exclude", []*regexp.Regexp{reFoo}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterLine(tt.line, tt.mode, tt.patterns)
			if got != tt.want {
				t.Errorf("FilterLine(%q, %q) = %v, want %v", tt.line, tt.mode, got, tt.want)
			}
		})
	}
}
