package plex_test

import (
	"testing"
	plex_parser "visualsource/plex/internal/parser"
)

func TestParserExpect(t *testing.T) {
	parser := plex_parser.NewParser("+")

	err := parser.Expect([]rune{'+'})
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestNextChar(t *testing.T) {
	parser := plex_parser.NewParser("<div")

	char := parser.NextChar()
	if char != '<' {
		t.Fatalf("Failed to get next char")
	}
}

func TestParseName(t *testing.T) {
	parser := plex_parser.NewParser("id=\"name\"")

	result := parser.ParseName()
	t.Logf("%s", string(result))
	if string(result) != "id" {
		t.Fatalf("Failed to parse name")
	}
}

func TestParseAttr(t *testing.T) {
	parser := plex_parser.NewParser("id=\"name\"")

	key, value, err := parser.ParseAttr()

	if err != nil {
		t.Fatalf("Failed to parse attr %s:", err)
		return
	}

	if key != "id" && value != "name" {
		t.Fatalf("Invalid values. Key: %s Value: %s", key, value)
	}
}
