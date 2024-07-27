package plex_css_test

import (
	"testing"
	plex_css "visualsource/plex/internal/css"

	"github.com/gookit/goutil/dump"
)

// "Parse a list of declarations" is for the contents of a style attribute,
// which parses text into the contents of a single style rule.
func TestParseListOfDeclarations(t *testing.T) {

	parser := plex_css.CssParser{}

	dec, err := parser.ParseDeclarationsList("color:while; background-color: green !important ;")

	if err != nil {
		t.Fatalf("%s", err)
	}

	dump.P(dec)
}
