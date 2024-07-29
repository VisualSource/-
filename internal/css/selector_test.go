package plex_css_test

import (
	"testing"
	plex_css "visualsource/plex/internal/css"
)

// region-start ns-prefix
func TestParseNSPrefix_IDENT(t *testing.T) {
	tokens := []plex_css.Token{
		&plex_css.StringToken{
			Id:    plex_css.Token_Ident,
			Value: "prefix",
		},
		&plex_css.RuneToken{
			Id:    plex_css.Token_Delim,
			Value: '|',
		},
	}
	pos := 0
	len := len(tokens)

	result, err := plex_css.ParseNSPrefix(&tokens, &pos, len)

	if err != nil {
		t.Fatalf("parse error: %s", err)
	}

	if pos != len {
		t.Fatalf("Failed to inc pos")
	}

	if result != "prefix" {
		t.Fatalf("Invalid value")
	}
}

func TestParseNSPrefix_STAR(t *testing.T) {
	tokens := []plex_css.Token{
		&plex_css.RuneToken{
			Id:    plex_css.Token_Delim,
			Value: '*',
		},
		&plex_css.RuneToken{
			Id:    plex_css.Token_Delim,
			Value: '|',
		},
	}
	pos := 0
	len := len(tokens)

	result, err := plex_css.ParseNSPrefix(&tokens, &pos, len)

	if err != nil {
		t.Fatalf("parse error: %s", err)
	}
	if pos != len {
		t.Fatalf("Failed to inc pos")
	}

	if result != "*" {
		t.Fatalf("Invalid value")
	}
}

func TestParseNSPrefix_PIPEONLY(t *testing.T) {
	tokens := []plex_css.Token{
		&plex_css.RuneToken{
			Id:    plex_css.Token_Delim,
			Value: '|',
		},
	}
	pos := 0
	len := len(tokens)

	result, err := plex_css.ParseNSPrefix(&tokens, &pos, len)
	if err != nil {
		t.Fatalf("parse error: %s", err)
	}
	if pos != len {
		t.Fatalf("Failed to inc pos")
	}
	if result != "" {
		t.Fatalf("Invalid value")
	}
}

func TestParseNSPrefix_NONE(t *testing.T) {
	tokens := []plex_css.Token{
		&plex_css.StringToken{
			Id:    plex_css.Token_Ident,
			Value: "prefix",
		},
	}
	pos := 0
	len := len(tokens)

	_, err := plex_css.ParseNSPrefix(&tokens, &pos, len)
	if err == nil {
		t.Fatalf("Expected error")
	}
	if pos != 0 {
		t.Fatalf("Failed to inc pos")
	}
}

// region-start wq-name

func TestParseWqName_ALL(t *testing.T) {
	tokens := []plex_css.Token{
		&plex_css.StringToken{
			Id:    plex_css.Token_Ident,
			Value: "prefix",
		},
		&plex_css.RuneToken{
			Id:    plex_css.Token_Delim,
			Value: '|',
		},
		&plex_css.StringToken{
			Id:    plex_css.Token_Ident,
			Value: "suffix",
		},
	}
	pos := 0
	len := len(tokens)

	namespace, tagname, err := plex_css.ParseWqName(&tokens, &pos, len)

	if err != nil {
		t.Fatalf("parse error: %s", err)
	}

	if pos != 3 {
		t.Fatalf("Failed to inc pos variable")
	}
	if tagname != "suffix" {
		t.Fatalf("Invalid value")
	}

	if namespace.IsNone() {
		t.Fatalf("Expected prefix")
	}

	if namespace.Unwrap() != "prefix" {
		t.Fatalf("Invalid prefix value")
	}
}

func TestParseWqName_NONAMESPACE(t *testing.T) {
	tokens := []plex_css.Token{
		&plex_css.StringToken{
			Id:    plex_css.Token_Ident,
			Value: "suffix",
		},
	}
	pos := 0
	len := len(tokens)

	namespace, tagname, err := plex_css.ParseWqName(&tokens, &pos, len)

	if err != nil {
		t.Fatalf("parse error: %s", err)
	}

	if pos != 1 {
		t.Fatalf("Failed to inc pos variable")
	}
	if tagname != "suffix" {
		t.Fatalf("Invalid value")
	}

	if namespace.IsSome() {
		t.Fatalf("Expected prefix")
	}
}

// region-start type-selector

func TestParseParseTypeSelector_WQNAME(t *testing.T) {
	tokens := []plex_css.Token{
		&plex_css.StringToken{
			Id:    plex_css.Token_Ident,
			Value: "prefix",
		},
		&plex_css.RuneToken{
			Id:    plex_css.Token_Delim,
			Value: '|',
		},
		&plex_css.StringToken{
			Id:    plex_css.Token_Ident,
			Value: "suffix",
		},
	}
	pos := 0
	len := len(tokens)

	namespace, tagname, err := plex_css.ParseTypeSelector(&tokens, &pos, len)

	if err != nil {
		t.Fatalf("parse error: %s", err)
	}

	if pos != 3 {
		t.Fatalf("Failed to inc pos variable")
	}
	if tagname != "suffix" {
		t.Fatalf("Invalid value")
	}

	if namespace.IsNone() {
		t.Fatalf("Expected prefix")
	}

	if namespace.Unwrap() != "prefix" {
		t.Fatalf("Invalid prefix value")
	}
}

func TestParseParseTypeSelector_SUFFIX_STAR(t *testing.T) {
	tokens := []plex_css.Token{
		&plex_css.StringToken{
			Id:    plex_css.Token_Ident,
			Value: "prefix",
		},
		&plex_css.RuneToken{
			Id:    plex_css.Token_Delim,
			Value: '|',
		},
		&plex_css.RuneToken{
			Id:    plex_css.Token_Delim,
			Value: '*',
		},
	}
	pos := 0
	len := len(tokens)

	namespace, tagname, err := plex_css.ParseTypeSelector(&tokens, &pos, len)

	if err != nil {
		t.Fatalf("parse error: %s", err)
	}

	if pos != 3 {
		t.Fatalf("Failed to inc pos variable")
	}
	if tagname != "*" {
		t.Fatalf("Invalid value")
	}

	if namespace.IsNone() {
		t.Fatalf("Expected prefix")
	}

	if namespace.Unwrap() != "prefix" {
		t.Fatalf("Invalid prefix value")
	}
}

// region-start combinator

func TestParseCombinator(t *testing.T) {
	tokens := []plex_css.Token{
		&plex_css.RuneToken{
			Id:    plex_css.Token_Delim,
			Value: '>',
		},
	}
	pos := 0
	len := len(tokens)

	token, err := plex_css.ParseCombinator(&tokens, &pos, len)
	if err != nil {
		t.Fatalf("parse error: %s", err)
	}
	if pos != 1 {
		t.Fatalf("Failed to inc pos")
	}

	if token != 0 {
		t.Fatalf("Invalid output")
	}
}
