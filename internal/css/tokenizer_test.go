package plex_css_test

import (
	"testing"
	plex_css "visualsource/plex/internal/css"
)

func TestTokenizerWhilespace(t *testing.T) {
	result, err := plex_css.Tokenizer("\n\r\n\r\f \t")

	if err != nil {
		t.Fatalf("Got an error: %s", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected slice to have len of 1 but got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_Whitespace {
		t.Fatalf("Expected type to be of Token Whitespace but get %d", result[0].GetId())
	}
}

func TestConsumeComments(t *testing.T) {
	result, err := plex_css.Tokenizer("/* I Am a comment \n in css*/")

	if err != nil {
		t.Fatalf("Got an error: %s", err)
	}

	if len(result) != 0 {
		t.Fatalf("Expected slice to have len of 0 but got %d", len(result))
	}
}

func TestConsumeCommentsAndWhilteSpace(t *testing.T) {
	result, err := plex_css.Tokenizer("  /* I Am a comment \n in css*/\n\t")

	if err != nil {
		t.Fatalf("Got an error: %s", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected slice to have len of 2 but got %d", len(result))
	}
}

func TestConsumeStringDouble(t *testing.T) {
	result, err := plex_css.Tokenizer("\"Im a string\"")

	if err != nil {
		t.Fatalf("Got an error: %s", err)
	}
	if len(result) != 1 {
		t.Fatalf("Expected slice to have len of 0 but got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_String {
		t.Fatalf("Expected type to be of Token Whitespace but get %d", result[0].GetId())
	}

	if i, ok := (result[0]).(*plex_css.StringToken); ok {
		target := []rune{'I', 'm', ' ', 'a', ' ', 's', 't', 'r', 'i', 'n', 'g'}
		for i, a := range i.Value {
			if target[i] != a {
				t.Fatalf("Found '%s' expected '%s'", string(a), string(target[i]))
			}
		}
	} else {
		t.Fatalf("Expected a string Token got %v", i)
	}

}

func TestConsumeStringSingle(t *testing.T) {
	result, err := plex_css.Tokenizer("'Im a string'")

	if err != nil {
		t.Fatalf("Got an error: %s", err)
	}
	if len(result) != 1 {
		t.Fatalf("Expected slice to have len of 0 but got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_String {
		t.Fatalf("Expected type to be of Token Whitespace but get %d", result[0].GetId())
	}

	if i, ok := (result[0]).(*plex_css.StringToken); ok {
		target := []rune{'I', 'm', ' ', 'a', ' ', 's', 't', 'r', 'i', 'n', 'g'}
		for i, a := range i.Value {
			if target[i] != a {
				t.Fatalf("Found '%s' expected '%s'", string(a), string(target[i]))
			}
		}
	} else {
		t.Fatalf("Expected a string Token got %v", i)
	}

}
