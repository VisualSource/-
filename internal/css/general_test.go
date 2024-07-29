package plex_css_test

import (
	"testing"
	plex_css "visualsource/plex/internal/css"
)

func TestSpecificity_STAR(t *testing.T) {

	selector := plex_css.Selector{
		TagName: "*",
	}

	spec := selector.GetSpecificity()

	if spec.A != 0 || spec.B != 0 || spec.C != 0 {
		t.Fatalf("invalid")
	}
}

func TestSpecificity_ID(t *testing.T) {

	selector := plex_css.Selector{
		TagName: "div",
		Id:      "x34y",
	}

	spec := selector.GetSpecificity()

	if spec.A != 1 || spec.B != 0 || spec.C != 1 {
		t.Fatalf("invalid")
	}
}
