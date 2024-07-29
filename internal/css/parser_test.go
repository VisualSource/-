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

func TestParseStylesheet(t *testing.T) {
	dump.Config(func(opts *dump.Options) {
		opts.MaxDepth = 10
	})

	parser := plex_css.CssParser{}

	stylesheet, err := parser.ParseStylesheet(`
	body {
		background-color: lightblue;
	}

	h1 {
		color: white;
		text-align: center;
	}

	p {
		font-family: verdana;
		font-size: 20px;
	}

	selector::pseudo-element {
		property: value;
	}

	p {
		background-color: red !important;
	}

	a[target="_blank"] {
		background-color: yellow;
	}

	#div1 {
		position: absolute;
		left: 50px;
		width: calc(100% - 100px);
		border: 1px solid black;
		background-color: yellow;
		padding: 5px;
	}


	@media screen and (min-width: 480px) {
		body {
			background-color: lightgreen;
		}
	}
	`)
	if err != nil {
		t.Fatalf("%s", err)
	}

	dump.P(stylesheet)
}
