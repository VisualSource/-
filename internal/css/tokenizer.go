package plex_css

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Tokenizer struct {
	pos    int
	len    int
	data   []rune
	Tokens []Token
}

func CreateTestTokenizer(value string) Tokenizer {
	data := []rune(value)

	return Tokenizer{
		pos:  0,
		data: data,
		len:  len(data),
	}
}

func (t *Tokenizer) Parse(value string) ([]Token, error) {
	t.pos = 0
	t.data = []rune(value)
	t.len = len(t.data)
	t.Tokens = []Token{}

	for !t.eof() {
		err := t.ConsumeToken()
		if err != nil {
			return nil, err
		}
	}

	t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_EOF})

	return t.Tokens, nil
}

func (t *Tokenizer) ConsumeToken() error {
	char := t.data[t.pos]
	switch {
	case t.CheckNextTwo('/', '*'):
		t.ConsumeComment()
	case char == '"' || char == '\'':
		err := t.ConsumeString()
		if err != nil {
			return err
		}
	case unicode.IsSpace(char):
		t.ConsumeWhilespace()
		t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Whitespace})
	case char == '#':
		t.pos++

		if !t.eof() && isIdentCodePoint(t.data[t.pos]) || t.AreNextValidEscape(0) {
			flagType := "unrestricted"
			if t.DoNextStartIdentSequence() {
				flagType = "id"
			}

			ident := t.ConsumeIdent()

			t.Tokens = append(t.Tokens, &FlagedStringToken{
				Value: string(ident),
				Id:    Token_Hash,
				Flag:  flagType,
			})

			return nil
		}

		t.Tokens = append(t.Tokens, &RuneToken{Id: Token_Delim, Value: char})
	case char == '(':
		t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Pren_Open})
		t.pos++
	case char == ')':
		t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Pren_Close})
		t.pos++
	case char == '+':
		if t.DoNextStartNumber() {
			t.ConsumeNumeric()
			return nil
		}
		t.Tokens = append(t.Tokens, &RuneToken{Id: Token_Delim, Value: char})
		t.pos++
	case char == ',':
		t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Comma})
		t.pos++
	case char == '-':
		if t.DoNextStartNumber() {
			t.ConsumeNumeric()
			return nil
		}

		if t.pos+2 < t.len && t.data[t.pos+1] == '-' && t.data[t.pos+2] == '>' {
			t.pos += 3

			t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_CDC})

			return nil
		}

		if t.DoNextStartIdentSequence() {
			t.ConsumeIdentLike()
			return nil
		}

		t.Tokens = append(t.Tokens, &RuneToken{Id: Token_Delim, Value: char})
		t.pos++
	case char == '.':
		/*
			If the input stream starts with a number,
			reconsume the current input code point,
			consume a numeric token, and return it.
		*/
		if t.DoNextStartNumber() {
			t.ConsumeNumeric()
			return nil
		}

		t.Tokens = append(t.Tokens, &RuneToken{Id: Token_Delim, Value: char})
		t.pos++
	case char == ':':
		t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Colon})
		t.pos++
	case char == ';':
		t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Semicolon})
		t.pos++
	case char == '<':
		if t.IsNextRune('!', 1) && t.IsNextRune('-', 2) && t.IsNextRune('-', 3) {
			t.pos += 4
			t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_CDO})
			return nil
		}

		t.Tokens = append(t.Tokens, &RuneToken{Id: Token_Delim, Value: char})
		t.pos++
	case char == '@':
		/*
			If the next 3 input code points would start an ident sequence,
			consume an ident sequence, create an <at-keyword-token> with its value set to
			the returned value, and return it.
		*/
		t.pos++ //eat '@'
		if t.DoNextStartIdentSequence() {
			value := t.ConsumeIdent()

			t.Tokens = append(t.Tokens, &StringToken{
				Id:    Token_At_Keyword,
				Value: string(value),
			})

			return nil
		}

		t.Tokens = append(t.Tokens, &RuneToken{Id: Token_Delim, Value: '@'})
	case char == '[':
		t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Square_Bracket_Open})
		t.pos++
	case char == ']':
		t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Square_Bracket_Close})
		t.pos++
	case char == '\\':

		if t.AreNextValidEscape(0) {
			t.ConsumeIdentLike()
			return nil
		}

		/*
			If the input stream starts with a valid escape,
			reconsume the current input code point,
			consume an ident-like token, and return it.
		*/
		// Otherwise, this is a parse error.
		// Return a <delim-token> with its value set to the current input code point.
		t.Tokens = append(t.Tokens, &RuneToken{Id: Token_Delim, Value: char})
		t.pos++
		return fmt.Errorf("invalid escape")
	case char == '{':
		t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Clearly_Open})
		t.pos++
	case char == '}':
		t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Clearly_Close})
		t.pos++
	case unicode.IsDigit(char):
		t.ConsumeNumeric()
	case isIdentStartCodePoint(char):
		t.ConsumeIdentLike()
	default:
		t.Tokens = append(t.Tokens, &RuneToken{Id: Token_Delim, Value: char})
		t.pos++
	}

	return nil
}

func (t *Tokenizer) ConsumeWhilespace() {
	for !t.eof() && unicode.IsSpace(t.data[t.pos]) {
		t.pos++
	}
}
func (t *Tokenizer) ConsumeComment() {
	t.pos += 2

	for !t.eof() && !t.CheckNextTwo('*', '/') {
		t.pos += 1
	}

	if t.CheckNextTwo('*', '/') {
		t.pos += 2
	}
}
func (t *Tokenizer) ConsumeNumeric() {
	value, dataType := t.ConsumeNumber()

	if t.DoNextStartIdentSequence() {
		ident := t.ConsumeIdent()

		t.Tokens = append(t.Tokens, &NumberToken{
			Id:       Token_Dimension,
			Value:    value,
			DataType: dataType,
			Unit:     string(ident),
		})
		return
	}

	if !t.eof() && t.data[t.pos] == '%' {
		t.pos++

		t.Tokens = append(t.Tokens, &NumberToken{
			Id:       Token_Percentage,
			DataType: dataType,
			Value:    value,
		})

		return
	}

	t.Tokens = append(t.Tokens, &NumberToken{
		Id:       Token_Number,
		Value:    value,
		DataType: dataType,
	})
}

// https://www.w3.org/TR/css-syntax-3/#consume-an-ident-like-token
func (t *Tokenizer) ConsumeIdentLike() {

	value := t.ConsumeIdent()

	if strings.ToLower(string(value)) == "url" && t.IsCurrent('(') {
		t.pos++

		// See https://github.com/w3c/csswg-drafts/issues/5416 for clearity on https://www.w3.org/TR/css-syntax-3/#consume-an-ident-like-token url function parse
		for {
			if unicode.IsSpace(t.data[t.pos]) && unicode.IsSpace(t.data[t.pos+1]) {
				t.pos++
				continue
			}

			if (t.data[t.pos] == '"' || t.data[t.pos] == '\'') || (t.data[t.pos+1] == '"' || t.data[t.pos+1] == '\'') {
				t.Tokens = append(t.Tokens, &StringToken{
					Id:    Token_Function,
					Value: string(value),
				})
				return
			}

			break
		}

		t.ConsumeUrl()
		return
	}

	if t.IsCurrent('(') {
		t.pos++
		t.Tokens = append(t.Tokens, &StringToken{
			Id:    Token_Function,
			Value: string(value),
		})
		return
	}

	t.Tokens = append(t.Tokens, &StringToken{
		Id:    Token_Ident,
		Value: string(value),
	})
}
func (t *Tokenizer) ConsumeString() error {
	delim := t.data[t.pos]
	t.pos++

	content := []rune{}

	for {
		switch {
		case t.eof():
			// This is a parse error. Return the <string-token>.
			t.Tokens = append(t.Tokens, &StringToken{Id: Token_String, Value: string(content)})
			return nil
		case t.data[t.pos] == delim:
			t.pos++
			// finish
			t.Tokens = append(t.Tokens, &StringToken{Id: Token_String, Value: string(content)})
			return nil
		case t.data[t.pos] == '\n':
			// This is a parse error. Reconsume the current input code point, create a <bad-string-token>, and return it.
			return fmt.Errorf("TODO implement newline in string. see: https://www.w3.org/TR/css-syntax-3/#consume-string-token")
		case t.data[t.pos] == '\\':
			return fmt.Errorf("TODO implement escape parseing. see: https://www.w3.org/TR/css-syntax-3/#consume-string-token")
		default:
			content = append(content, t.data[t.pos])
			t.pos++
		}
	}
}
func (t *Tokenizer) ConsumeUrl() error {
	t.ConsumeWhilespace()

	reper := []rune{}

L:
	for {
		switch {
		case t.IsCurrent(')'):
			t.pos++
			break L
		case unicode.IsSpace(t.data[t.pos]):
			t.ConsumeWhilespace()
		case t.IsCurrent('"') || t.IsCurrent('\'') || t.IsCurrent('(') || !unicode.IsPrint(t.data[t.pos]):
			for !t.eof() && !t.IsCurrent('}') {
				if t.AreNextValidEscape(0) {
					t.ConsumeEscaped()
				} else {
					t.pos++
				}
			}

			if t.IsCurrent('}') {
				t.pos++
			}

			t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Bad_Url})
			return nil
		case t.IsCurrent('\\'):
			t.pos++
			data, err := t.ConsumeEscaped()
			if err != nil {
				return err
			}
			reper = append(reper, data...)
		default:
			reper = append(reper, t.data[t.pos])
			t.pos++
		}
	}

	t.Tokens = append(t.Tokens, &StringToken{
		Id:    Token_Url,
		Value: string(reper),
	})

	return nil
}

func (t *Tokenizer) ConsumeEscaped() ([]rune, error) {

	target := []rune{}

	if t.eof() {
		target = append(target, '�')
	} else if unicode.IsDigit(t.data[t.pos]) || unicode.IsLetter(t.data[t.pos]) {
		// prase hex digit
		return target, fmt.Errorf("TODO: handle hex escape code")
	} else {
		target = append(target, t.data[t.pos])
	}

	return target, nil
}
func (t *Tokenizer) ConsumeIdent() []rune {

	result := []rune{}
	// TODO: add escape parsing. See https://www.w3.org/TR/css-syntax-3/#consume-an-ident-sequence
	for !t.eof() && isIdentCodePoint(t.data[t.pos]) {
		result = append(result, t.data[t.pos])
		t.pos++
	}

	return result
}
func (t *Tokenizer) ConsumeNumber() (float32, NumberType) {
	var valueType NumberType = NumberType_Integer

	reper := []rune{}

	if t.data[t.pos] == '+' || t.data[t.pos] == '-' {
		reper = append(reper, t.data[t.pos])
		t.pos++
	}

	for !t.eof() && unicode.IsDigit(t.data[t.pos]) {
		reper = append(reper, t.data[t.pos])
		t.pos++
	}

	if !t.eof() && t.data[t.pos] == '.' && unicode.IsDigit(t.data[t.pos+1]) {
		reper = append(reper, t.data[t.pos])
		t.pos++
		reper = append(reper, t.data[t.pos])
		t.pos++
		valueType = NumberType_Number
		for !t.eof() && unicode.IsDigit(t.data[t.pos]) {
			reper = append(reper, t.data[t.pos])
			t.pos++
		}
	}

	if !t.eof() && (t.data[t.pos] == 'e' || t.data[t.pos] == 'E') && (unicode.IsDigit(t.data[t.pos+1]) ||
		(t.data[t.pos+1] == '+' || t.data[t.pos+1] == '-' && unicode.IsDigit(t.data[t.pos+2]))) {

		reper = append(reper, t.data[t.pos])
		t.pos++

		if t.data[t.pos] == '+' || t.data[t.pos] == '-' {
			reper = append(reper, t.data[t.pos])
			t.pos++
		}

		valueType = NumberType_Number

		for !t.eof() && unicode.IsDigit(t.data[t.pos]) {
			reper = append(reper, t.data[t.pos])
			t.pos++
		}
	}

	value, err := strconv.ParseFloat(string(reper), 32)

	if err != nil {
		return 0, valueType
	}

	return float32(value), valueType
}

func (t *Tokenizer) CheckNextTwo(a rune, b rune) bool {
	if t.pos+1 >= t.len {
		return false
	}

	return t.data[t.pos] == a && t.data[t.pos+1] == b
}
func (t *Tokenizer) CheckNextThree(test func(rune) bool) bool {
	if t.pos+3 > t.len-1 {
		return false
	}

	for i := 0; i < 3; i++ {
		if !test(t.data[t.pos+i]) {
			return false
		}
	}

	return true
}
func (t *Tokenizer) IsCurrent(v rune) bool {
	if t.pos >= t.len {
		return false
	}

	return t.data[t.pos] == v
}

func (t *Tokenizer) IsNextRune(v rune, offset int) bool {
	if t.pos+offset >= t.len {
		return false
	}

	return t.data[t.pos+offset] == v
}

func (t *Tokenizer) eof() bool {
	return t.pos >= t.len
}

// https://www.w3.org/TR/css-syntax-3/#starts-with-a-valid-escape
func (t *Tokenizer) AreNextValidEscape(offset int) bool {
	if (offset+t.pos)+1 >= t.len-1 {
		return false
	}

	if t.data[(offset+t.pos)] != '\\' {
		return false
	}

	if t.data[(offset+t.pos)+1] == '\n' {
		return false
	}

	return true
}

// https://www.w3.org/TR/css-syntax-3/#would-start-an-identifier
func (t *Tokenizer) DoNextStartIdentSequence() bool {
	if t.pos+2 >= t.len {
		return false
	}
	char := t.data[t.pos]
	switch {
	case char == '-':
		if isIdentStartCodePoint(t.data[t.pos+1]) || t.data[t.pos+1] == '-' {
			return true
		}

		if t.AreNextValidEscape(1) {
			return true
		}

		return false
	case isIdentStartCodePoint(char):
		return true
	case char == '\\':
		return t.AreNextValidEscape(0)
	default:
		return false
	}
}

func (t *Tokenizer) DoNextStartNumber() bool {
	if t.pos+3 > t.len {
		return false
	}

	if t.data[t.pos] == '+' || t.data[t.pos] == '-' {
		if unicode.IsDigit(t.data[t.pos+1]) {
			return true
		}

		if t.data[t.pos+1] == '.' && unicode.IsDigit(t.data[t.pos+2]) {
			return true
		}

		return false
	}

	if t.data[t.pos] == '.' {
		return unicode.IsDigit(t.data[t.pos+1])
	}

	if unicode.IsDigit(t.data[t.pos]) {
		return true
	}

	return false
}

func isIdentStartCodePoint(value rune) bool {
	return unicode.IsLetter(value) || value == '_'
}

func isIdentCodePoint(value rune) bool {
	return isIdentStartCodePoint(value) || unicode.IsDigit(value) || value == '-'
}
