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

// region start ConsumeString
func TestConsumeToken_StringDouble(t *testing.T) {
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
func TestConsumeToken_StringSingle(t *testing.T) {
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

// region start ConsumeToken'#'

func TestConsumeTokenHash_Single(t *testing.T) {
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

func TestConsumeToken_HashLong(t *testing.T) {
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

func TestConsumeToken_HashShort(t *testing.T) {
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

// region start ConsumeSingleRune

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

// region start ConsumeNumber

func TestConsumeNumber_Int(t *testing.T) {
	tokenizer := plex_css.CreateTestTokenizer("566")

	value, dataType := tokenizer.ConsumeNumber()

	if dataType != "integer" {
		t.Fatalf("Expected dataType to be of integer not %s", dataType)
	}

	if value != 566 {
		t.Fatalf("Expected value to be 566 not %f", value)
	}
}

func TestConsumeNumber_Float(t *testing.T) {
	tokenizer := plex_css.CreateTestTokenizer("566.3")

	value, dataType := tokenizer.ConsumeNumber()

	if dataType != "number" {
		t.Fatalf("Expected dataType to be of number not %s", dataType)
	}

	if value != 566.3 {
		t.Fatalf("Expected value to be 566 not %f", value)
	}
}

func TestConsumeNumber_Exponent(t *testing.T) {
	tokenizer := plex_css.CreateTestTokenizer("566e4")

	value, dataType := tokenizer.ConsumeNumber()

	if dataType != "number" {
		t.Fatalf("Expected dataType to be of number not %s", dataType)
	}

	if value != 566e4 {
		t.Fatalf("Expected value to be 566e4 not %f", value)
	}
}

func TestConsumeNumber_Signed(t *testing.T) {
	tokenizer := plex_css.CreateTestTokenizer("-566")

	value, dataType := tokenizer.ConsumeNumber()

	if dataType != "integer" {
		t.Fatalf("Expected dataType to be of integer not %s", dataType)
	}

	if value != -566 {
		t.Fatalf("Expected value to be -566 not %f", value)
	}
}

// region start ConsumeToken'+'

func TestConsumeToken_Plus(t *testing.T) {
	parser := plex_css.Tokenizer{}

	content := "+"
	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}
	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_Delim {
		t.Fatalf("Did not find Delim Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[1].GetId())
	}

	if e, ok := result[0].(*plex_css.RuneToken); ok {
		if e.Value != '+' {
			t.Fatalf("Value does not match input")
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", result[0])
	}
}

func TestConsumeToken_PlusNumber(t *testing.T) {
	parser := plex_css.Tokenizer{}

	content := "+4445.0"
	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}
	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_Number {
		t.Fatalf("Did not find Delim Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[1].GetId())
	}

	if e, ok := result[0].(*plex_css.NumberToken); ok {
		if e.Value != 4445.0 {
			t.Fatalf("Value does not match input")
		}
		if e.DataType != "number" {
			t.Fatalf("Value does not match input")
		}
		if !equal(e.Unit, []rune{}) {
			t.Fatalf("Value does not match input")
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", result[0])
	}
}

// region start CosumeToken'-'

func TestConsumeToken_Minus(t *testing.T) {
	parser := plex_css.Tokenizer{}

	content := "-4445.0"
	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}
	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_Number {
		t.Fatalf("Did not find Delim Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[1].GetId())
	}

	if e, ok := result[0].(*plex_css.NumberToken); ok {
		if e.Value != -4445.0 {
			t.Fatalf("Value does not match input")
		}
		if e.DataType != "number" {
			t.Fatalf("Value does not match input")
		}
		if !equal(e.Unit, []rune{}) {
			t.Fatalf("Value does not match input")
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", result[0])
	}
}

func TestConsumeToken_MinusCDC(t *testing.T) {
	parser := plex_css.Tokenizer{}

	content := "-->"
	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}
	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_CDC {
		t.Fatalf("Did not find Delim Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[1].GetId())
	}
}

func TestConsumeToken_MinusIdent(t *testing.T) {
	parser := plex_css.Tokenizer{}

	content := "-dent"
	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}
	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_Ident {
		t.Fatalf("Did not find Delim Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[1].GetId())
	}

	if e, ok := result[0].(*plex_css.StringToken); ok {
		if string(e.Value) != content {
			t.Fatalf("Value does not match input")
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", result[0])
	}
}

func TestConsumeToken_MinusSelf(t *testing.T) {
	parser := plex_css.Tokenizer{}

	content := "-"
	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}
	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_Delim {
		t.Fatalf("Did not find Delim Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[1].GetId())
	}

	if e, ok := result[0].(*plex_css.RuneToken); ok {
		if e.Value != '-' {
			t.Fatalf("Value does not match input")
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", result[0])
	}
}

// region start ConsumeIdentLike

func TestConsumeIdentLike_Ident(t *testing.T) {
	tokenizer := plex_css.CreateTestTokenizer("ident")

	tokenizer.ConsumeIdentLike()

	if len(tokenizer.Tokens) != 1 {
		t.Fatalf("Expected only 1 tokens got %d", len(tokenizer.Tokens))
	}

	if tokenizer.Tokens[0].GetId() != plex_css.Token_Ident {
		t.Fatalf("Did not find Ident Token Found %d", tokenizer.Tokens[0].GetId())
	}

	if e, ok := tokenizer.Tokens[0].(*plex_css.StringToken); ok {
		if string(e.Value) != "ident" {
			t.Fatalf("Value does not match input")
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", tokenizer.Tokens[0])
	}
}

func TestConsumeIdentLike_Function(t *testing.T) {
	tokenizer := plex_css.CreateTestTokenizer("var()")

	tokenizer.ConsumeIdentLike()

	if len(tokenizer.Tokens) != 1 {
		t.Fatalf("Expected only 1 tokens got %d", len(tokenizer.Tokens))
	}

	if tokenizer.Tokens[0].GetId() != plex_css.Token_Function {
		t.Fatalf("Did not find Ident Token Found %d", tokenizer.Tokens[0].GetId())
	}

	if e, ok := tokenizer.Tokens[0].(*plex_css.StringToken); ok {
		if string(e.Value) != "var" {
			t.Fatalf("Value does not match input")
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", tokenizer.Tokens[0])
	}
}

func TestConsumeIdentLike_URLFunction(t *testing.T) {
	tokenizer := plex_css.CreateTestTokenizer("url('test')")

	tokenizer.ConsumeIdentLike()

	if len(tokenizer.Tokens) != 1 {
		t.Fatalf("Expected only 1 tokens got %d", len(tokenizer.Tokens))
	}

	if tokenizer.Tokens[0].GetId() != plex_css.Token_Function {
		t.Fatalf("Did not find Ident Token Found %d", tokenizer.Tokens[0].GetId())
	}

	if e, ok := tokenizer.Tokens[0].(*plex_css.StringToken); ok {
		if string(e.Value) != "url" {
			t.Fatalf("Value does not match input")
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", tokenizer.Tokens[0])
	}
}

func TestConsumeIdentLike_URL(t *testing.T) {
	tokenizer := plex_css.CreateTestTokenizer("url(test.url)")

	tokenizer.ConsumeIdentLike()

	if len(tokenizer.Tokens) != 1 {
		t.Fatalf("Expected only 1 tokens got %d", len(tokenizer.Tokens))
	}

	if tokenizer.Tokens[0].GetId() != plex_css.Token_Url {
		t.Fatalf("Did not find Ident Token Found %d", tokenizer.Tokens[0].GetId())
	}

	if e, ok := tokenizer.Tokens[0].(*plex_css.StringToken); ok {
		if string(e.Value) != "test.url" {
			t.Fatalf("Value does not match input")
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", tokenizer.Tokens[0])
	}
}

// region end ConsumeIdentLike

// region start ConsumeURL

func TestConsumeUrl(t *testing.T) {
	tokenizer := plex_css.CreateTestTokenizer("  test.url    )")

	tokenizer.ConsumeUrl()

	if len(tokenizer.Tokens) != 1 {
		t.Fatalf("Expected only 1 tokens got %d", len(tokenizer.Tokens))
	}

	if tokenizer.Tokens[0].GetId() != plex_css.Token_Url {
		t.Fatalf("Did not find Ident Token Found %d", tokenizer.Tokens[0].GetId())
	}

	if e, ok := tokenizer.Tokens[0].(*plex_css.StringToken); ok {
		if string(e.Value) != "test.url" {
			t.Fatalf("Value does not match input")
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", tokenizer.Tokens[0])
	}
}

func TestConsumeBadUrl(t *testing.T) {
	tokenizer := plex_css.CreateTestTokenizer("  test.url\"  }  )")

	tokenizer.ConsumeUrl()

	if len(tokenizer.Tokens) != 1 {
		t.Fatalf("Expected only 1 tokens got %d", len(tokenizer.Tokens))
	}

	if tokenizer.Tokens[0].GetId() != plex_css.Token_Bad_Url {
		t.Fatalf("Did not find Bad URL Token Found %d", tokenizer.Tokens[0].GetId())
	}
}

// regions start ConsumeToken'.'

func TestConsumeToken_StopSelf(t *testing.T) {
	parser := plex_css.Tokenizer{}

	content := "."
	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}
	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_Delim {
		t.Fatalf("Did not find Delim Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[1].GetId())
	}

	if e, ok := result[0].(*plex_css.RuneToken); ok {
		if e.Value != '.' {
			t.Fatalf("Value does not match input")
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", result[0])
	}
}

func TestConsumeToken_StopNumber(t *testing.T) {
	parser := plex_css.Tokenizer{}

	content := ".45"
	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}
	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_Number {
		t.Fatalf("Did not find Delim Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[1].GetId())
	}

	if e, ok := result[0].(*plex_css.NumberToken); ok {
		if e.Value != 0.45 {
			t.Fatalf("Value does not match input")
		}
		if e.DataType != "number" {
			t.Fatalf("Value does not match input")
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", result[0])
	}
}

// regions start ConsumeToken'<'
func TestConsumeToken_LessThanSelf(t *testing.T) {
	parser := plex_css.Tokenizer{}

	content := "<"
	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}
	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_Delim {
		t.Fatalf("Did not find Delim Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[1].GetId())
	}

	if e, ok := result[0].(*plex_css.RuneToken); ok {
		if e.Value != '<' {
			t.Fatalf("Value does not match input")
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", result[0])
	}
}

func TestConsumeToken_LessThanCDO(t *testing.T) {
	parser := plex_css.Tokenizer{}

	content := "<!--"
	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}
	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_CDO {
		t.Fatalf("Did not find Delim Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[1].GetId())
	}
}

// #region-start ConsumeToken'@'

func TestConsumeToken_At(t *testing.T) {
	parser := plex_css.Tokenizer{}

	content := "@"
	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}
	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_Delim {
		t.Fatalf("Did not find Delim Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[1].GetId())
	}

	if e, ok := result[0].(*plex_css.RuneToken); ok {
		if e.Value != '@' {
			t.Fatalf("Value does not match input")
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", result[0])
	}
}

func TestConsumeToken_AtRule(t *testing.T) {
	parser := plex_css.Tokenizer{}

	content := "@charset"
	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}
	if len(result) != 2 {
		t.Fatalf("Expected only 2 tokens got %d", len(result))
	}

	if result[0].GetId() != plex_css.Token_At_Keyword {
		t.Fatalf("Did not find Delim Token Found %d", result[0].GetId())
	}
	if result[1].GetId() != plex_css.Token_EOF {
		t.Fatalf("Did not find EOF Token Found %d", result[1].GetId())
	}

	if e, ok := result[0].(*plex_css.StringToken); ok {
		if string(e.Value) != "charset" {
			t.Fatalf("Value does not match input")
		}
	} else {
		t.Fatalf("Token is not a Rune token got: %v", result[0])
	}
}

func TestParse(t *testing.T) {
	parser := plex_css.Tokenizer{}

	result, err := parser.Parse(`
		body {
			background-color: #F3F4F7;
			display: flex;
			flex-direction: row;
			justify-content: center;
			align-items: center;
			height: 100vh;
			overflow: hidden;
			gap: 20px;
		}
		span:nth-child(1):hover ~ button::before { 
			width: 100%;
		}
		@keyframes scale {
			0% {
				transform: translateY(var(--origin, 0%));
			}
			100% {
				transform: translateY(var(--destination, -50%));
			}
		}
	`)
	if err != nil {
		t.Fatalf("There was an error: %s", err)
	}

	if len(result) != 127 {
		t.Fatalf("Missing token expected 127 tokens got: %d", len(result))
	}
}

// #region start UTILS

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
