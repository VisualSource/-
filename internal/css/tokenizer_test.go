package plex_css_test

import (
	"testing"
	plex_css "visualsource/plex/internal/css"
)

func TestConsumeComment(t *testing.T) {
	tokenizer := plex_css.CreateTestTokenizer("/*Im a commet*/")
	tokenizer.ConsumeComment()

	if len(tokenizer.Tokens) != 0 {
		t.Fatalf("Ther was more tokens then expected")
	}
}

func TestConsumeWhileSpace(t *testing.T) {
	tokenizer := plex_css.CreateTestTokenizer("\n\r\n\r\f \t")
	tokenizer.ConsumeWhilespace()
}

func TestConsumeTokenComment(t *testing.T) {
	parser := plex_css.Tokenizer{}

	result, err := parser.Parse("/*Im a commet*/")
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected only 1 token got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[0].GetId())
	}
}

func TestConsumeTokenWhilespace(t *testing.T) {
	parser := plex_css.Tokenizer{}

	result, err := parser.Parse("\n\r\n\r\f \t")
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_Whitespace {
		t.Fatalf("Did not find Whilespace Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[0].GetId())
	}
}

func TestConsumeTokenStringDouble(t *testing.T) {
	parser := plex_css.Tokenizer{}

	result, err := parser.Parse("\"I Am a string\"")
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_String {
		t.Fatalf("Did not find Whilespace Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[0].GetId())
	}

	if e, ok := result[0].(*plex_css.StringToken); ok {

		if !equal(e.Value, []rune{'I', ' ', 'A', 'm', ' ', 'a', ' ', 's', 't', 'r', 'i', 'n', 'g'}) {
			t.Fatalf("Value does not match input")
		}

	} else {
		t.Fatalf("Token is not a String token get: %v", result[0])
	}
}
func TestConsumeTokenStringSingle(t *testing.T) {
	parser := plex_css.Tokenizer{}

	content := "'I Am a string'"

	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_String {
		t.Fatalf("Did not find Whilespace Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[0].GetId())
	}

	if e, ok := result[0].(*plex_css.StringToken); ok {

		if !equal(e.Value, []rune{'I', ' ', 'A', 'm', ' ', 'a', ' ', 's', 't', 'r', 'i', 'n', 'g'}) {
			t.Fatalf("Value does not match input")
		}

	} else {
		t.Fatalf("Token is not a String token get: %v", result[0])
	}
}

func TestConsumeTokenHashSingle(t *testing.T) {
	parser := plex_css.Tokenizer{}

	content := "#"

	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_Delim {
		t.Fatalf("Did not find Hash Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[1].GetId())
	}

	if e, ok := result[0].(*plex_css.RuneToken); ok {
		if e.Value != '#' {
			t.Fatalf("Value does not match input")
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", result[0])
	}
}

func TestConsumeTokenHashLong(t *testing.T) {
	parser := plex_css.Tokenizer{}

	content := "#ImAId"

	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_Hash {
		t.Fatalf("Did not find Hash Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[1].GetId())
	}

	if e, ok := result[0].(*plex_css.FlagedStringToken); ok {
		if !equal(e.Value, []rune{'I', 'm', 'A', 'I', 'd'}) {
			t.Fatalf("Value does not match input")
		}
		if e.Flag != "id" {
			t.Fatalf("Expected flag to be of 'id' not %s", e.Flag)
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", result[0])
	}
}

func TestConsumeTokenHashShort(t *testing.T) {
	parser := plex_css.Tokenizer{}

	content := "#ID"

	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_Hash {
		t.Fatalf("Did not find Hash Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[1].GetId())
	}

	if e, ok := result[0].(*plex_css.FlagedStringToken); ok {
		if !equal(e.Value, []rune{'I', 'D'}) {
			t.Fatalf("Value does not match input")
		}
		if e.Flag != "unrestricted" {
			t.Fatalf("Expected flag to be of 'unrestricted' not %s", e.Flag)
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", result[0])
	}
}

func TestConsumeSingleRune(t *testing.T) {
	parser := plex_css.Tokenizer{}

	items := []string{"[", "]", "(", ")", "{", "}", ":", ";", ","}
	test := []plex_css.TokenType{
		plex_css.Token_Square_Bracket_Open,
		plex_css.Token_Square_Bracket_Close,
		plex_css.Token_Pren_Open,
		plex_css.Token_Pren_Close,
		plex_css.Token_Clearly_Open,
		plex_css.Token_Clearly_Close,
		plex_css.Token_Colon,
		plex_css.Token_Semicolon,
		plex_css.Token_Comma,
	}

	for i, b := range items {
		results, err := parser.Parse(b)
		if err != nil {
			t.Fatalf("There was an error: %s", err)
		}
		if len(results) != 2 {
			t.Fatalf("Expected only 2 tokens got %d", len(results))
		}

		if results[0].GetId() != test[i] {
			t.Fatalf("Did not find expected %d. Found %d", test[i], results[0].GetId())
		}
		if results[1].GetId() != plex_css.Token_EOF {
			t.Fatalf("Did not find EOF Token Found %d", results[1].GetId())
		}
	}
}

func equal(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
