package plex_css_test

import (
	"testing"
	plex_css "visualsource/plex/internal/css"

	"github.com/gookit/goutil/dump"
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

func TestParseCssValue_BLOCK(t *testing.T) {
	tokens := []plex_css.Token{
		&plex_css.NumberToken{
			Id:       plex_css.Token_Dimension,
			Value:    1,
			Unit:     "px",
			DataType: plex_css.NumberType_Number,
		},
		&plex_css.EmptyToken{Id: plex_css.Token_Whitespace},
		&plex_css.StringToken{
			Id:    plex_css.Token_Ident,
			Value: "solid",
		},
		&plex_css.EmptyToken{Id: plex_css.Token_Whitespace},
		&plex_css.StringToken{
			Id:    plex_css.Token_Ident,
			Value: "black",
		},
	}
	value := plex_css.ParseCssValue(tokens)

	dump.P(value)
}

func TestParseCssValue_FUNC(t *testing.T) {
	tokens := []plex_css.Token{
		&plex_css.FunctionBlock{
			Name: "calc",
			Args: []plex_css.Token{
				&plex_css.NumberToken{
					Id:       plex_css.Token_Dimension,
					Value:    1,
					Unit:     "px",
					DataType: plex_css.NumberType_Number,
				},
				&plex_css.EmptyToken{Id: plex_css.Token_Whitespace},
				&plex_css.RuneToken{Id: plex_css.Token_Delim, Value: '+'},
				&plex_css.EmptyToken{Id: plex_css.Token_Whitespace},
				&plex_css.NumberToken{
					Id:       plex_css.Token_Dimension,
					Value:    1,
					Unit:     "px",
					DataType: plex_css.NumberType_Number,
				},
			},
		},
	}
	value := plex_css.ParseCssValue(tokens)

	dump.P(value)
}
