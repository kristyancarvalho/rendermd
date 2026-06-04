package ui

import (
	"reflect"
	"testing"

	"github.com/kristyancarvalho/rendermd/internal/layout"
)

func TestBuildSearchIndex(t *testing.T) {
	lines := []layout.Line{
		{Segments: []layout.Segment{
			{Text: "Hello", Style: layout.StyleNormal},
			{Text: " World", Style: layout.StyleStrong},
		}},
		{Segments: []layout.Segment{{Text: "Second", Style: layout.StyleNormal}}},
	}
	got := buildSearchIndex(lines)
	want := []string{"hello world", "second"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("search index: want %#v, got %#v", want, got)
	}
}

func TestFindHitsUsesSearchIndex(t *testing.T) {
	index := []string{"alpha beta", "gamma", "beta gamma"}
	got := findHits(index, "BETA")
	want := []int{0, 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("hits: want %#v, got %#v", want, got)
	}
}

func TestFindHitsEmptyQuery(t *testing.T) {
	if got := findHits([]string{"alpha"}, ""); got != nil {
		t.Errorf("empty query should produce nil hits, got %#v", got)
	}
}
