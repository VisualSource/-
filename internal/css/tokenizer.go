package plex_css

import (
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
		t.ConsumeString()
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
				Value: ident,
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
		/*
			If the input stream starts with a number, reconsume the current input code point,
			consume a numeric token, and return it.
		*/

		t.Tokens = append(t.Tokens, &RuneToken{Id: Token_Delim, Value: char})
		t.pos++
	case char == ',':
		t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Comma})
		t.pos++
	case char == '-':
		/*
			If the input stream starts with a number, reconsume the current input code point,
				consume a numeric token, and return it.

			Otherwise, if the next 2 input code points are
				U+002D HYPHEN-MINUS U+003E GREATER-THAN SIGN (->),
				consume them and return a <CDC-token>.

			Otherwise, if the input stream starts with an ident sequence,
				reconsume the current input code point,
				consume an ident-like token,
				and return it.
		*/
		t.Tokens = append(t.Tokens, &RuneToken{Id: Token_Delim, Value: char})
		t.pos++
	case char == '.':
		/*
			If the input stream starts with a number,
			reconsume the current input code point,
			consume a numeric token, and return it.
		*/

		t.Tokens = append(t.Tokens, &RuneToken{Id: Token_Delim, Value: char})
		t.pos++
	case char == ':':
		t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Colon})
		t.pos++
	case char == ';':
		t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Semicolon})
		t.pos++
	case char == '<':

		if t.data[t.pos+1] == '!' && t.data[t.pos+2] == '-' && t.data[t.pos+3] == '-' {

			/*
				If the next 3 input code points are U+0021 EXCLAMATION MARK U+002D HYPHEN-MINUS U+002D HYPHEN-MINUS (!--), consume them and return a <CDO-token>.
			*/

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
		// create at keyword token else
		t.Tokens = append(t.Tokens, &RuneToken{Id: Token_Delim, Value: char})
		t.pos++
	case char == '[':
		t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Square_Bracket_Open})
		t.pos++
	case char == ']':
		t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Square_Bracket_Close})
		t.pos++
	case char == '\\':

		/*
			If the input stream starts with a valid escape,
			reconsume the current input code point,
			consume an ident-like token, and return it.
		*/
		// Otherwise, this is a parse error.
		// Return a <delim-token> with its value set to the current input code point.
		t.Tokens = append(t.Tokens, &RuneToken{Id: Token_Delim, Value: char})
		t.pos++
	case char == '{':
		t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Clearly_Open})
		t.pos++
	case char == '}':
		t.Tokens = append(t.Tokens, &EmptyToken{Id: Token_Clearly_Close})
		t.pos++
	case unicode.IsDigit(char):
		t.ConsumeNumeric()
	case unicode.IsLetter(char) || char == '_':
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
func (t *Tokenizer) ConsumeNumeric()   {}
func (t *Tokenizer) ConsumeIdentLike() {}
func (t *Tokenizer) ConsumeString() {
	delim := t.data[t.pos]
	t.pos++

	content := []rune{}

	for {
		switch {
		case t.eof():
			// This is a parse error. Return the <string-token>.
			t.Tokens = append(t.Tokens, &StringToken{Id: Token_String, Value: content})
			return
		case t.data[t.pos] == delim:
			t.pos++
			// finish
			t.Tokens = append(t.Tokens, &StringToken{Id: Token_String, Value: content})
			return
		case t.data[t.pos] == '\n':
			// This is a parse error. Reconsume the current input code point, create a <bad-string-token>, and return it.
			panic("TODO implement newline in string. see: https://www.w3.org/TR/css-syntax-3/#consume-string-token")
		case t.data[t.pos] == '\\':
			panic("TODO implement escape parseing. see: https://www.w3.org/TR/css-syntax-3/#consume-string-token")
		default:
			content = append(content, t.data[t.pos])
			t.pos++
		}
	}

}
func (t *Tokenizer) ConsumeUrl()    {}
func (t *Tokenizer) ConsumeBadUrl() {}
func (t *Tokenizer) ConsumeEscaped() []rune {

	target := []rune{}

	if t.eof() {
		target = append(target, '�')
	} else if unicode.IsDigit(t.data[t.pos]) || unicode.IsLetter(t.data[t.pos]) {
		// prase hex digit
	} else {
		target = append(target, t.data[t.pos])
	}

	return target
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
func (t *Tokenizer) ConsumeNumber() {}

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
	if t.pos+2 >= t.len-1 {
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

func isIdentStartCodePoint(value rune) bool {
	return unicode.IsLetter(value) || value == '_'
}

func isIdentCodePoint(value rune) bool {
	return isIdentStartCodePoint(value) || unicode.IsDigit(value) || value == '-'
}

/*
func eof(len int, pos int) bool {
	return pos >= len-1
}

func startsWith(s *[]rune, data *[]rune, dataLen int, startPos int) bool {
	isSame := true
	for i := 0; i < len(*s); i++ {
		d := startPos + 1
		if d > dataLen || (*data)[startPos+i] != (*s)[i] {
			isSame = false
			break
		}
	}
	return isSame
}

func RunnerTokenizer(input string) ([]Token, error) {
	var commentStart = []rune{'/', '*'}
	var commentEnd = []rune{'*', '/'}

	data := []rune(input)
	Tokens := []Token{}
	dataLen := len(data)
	pos := 0
	for pos < dataLen {
		if eof(dataLen, pos) {
			Tokens = append(Tokens, &EmptyToken{Id: Token_EOF})
			break
		}
		char := data[pos]

		switch {
		// https://www.w3.org/TR/css-syntax-3/#consume-comments
		case startsWith(&commentStart, &data, dataLen, pos):
			pos += 2 // eat '/*'

			for !startsWith(&commentEnd, &data, dataLen, pos) {
				pos++
				if eof(dataLen, pos) {
					return nil, fmt.Errorf("expected '' but found EOF token")
				}
			}

			pos += 2 // eat ''

		// https://www.w3.org/TR/css-syntax-3/#whitespace-diagram
		// https://www.w3.org/TR/css-syntax-3/#whitespace
		case unicode.IsSpace(char):
			pos++ // eat char
			for unicode.IsSpace(data[pos]) {
				pos++
				if eof(dataLen, pos) {
					break
				}
			}
			//fmt.Printf("%U\n", data[pos])
			Tokens = append(Tokens, &EmptyToken{Id: Token_Whitespace})
		// https://www.w3.org/TR/css-syntax-3/#escape-diagram
		// https://www.w3.org/TR/css-syntax-3/#string-token-diagram
		// https://www.w3.org/TR/css-syntax-3/#consume-string-token
		case char == '"' || char == '\'':
			pos++
			i := char

			codePoints := []rune{}

			for {
				if eof(dataLen, pos) {
					pos++
					// create string token
					Tokens = append(Tokens, &StringToken{
						Id:    Token_String,
						Value: codePoints,
					})
					break
				}
				if data[pos] == i {
					pos++ // eat '\'' or '\"'
					Tokens = append(Tokens, &StringToken{
						Id:    Token_String,
						Value: codePoints,
					})
					break
				}

				if data[pos] == '\n' {
					pos++
					Tokens = append(Tokens, &EmptyToken{Id: Token_Bad_String})
					break // parser error create bad string token
				}
				if data[pos] == '\\' {
					pos++
					// https://www.w3.org/TR/css-syntax-3/#escape-diagram
					if eof(dataLen, pos) {
						continue
					}
					if data[pos] == '\n' {
						pos++
						continue
					}

					/*
						Otherwise, (the stream starts with a valid escape)
						consume an escaped code point and append the returned
						code point to the <string-token>’s value.

						https://www.w3.org/TR/css-syntax-3/#consume-escaped-code-point

					return nil, fmt.Errorf("escaped code point is not fully implemented")
				}

				codePoints = append(codePoints, data[pos])
				pos++
			}
		// https://www.w3.org/TR/css-syntax-3/#consume-token
		// #
		case char == '#':
			pos++ //eat '#'

			Tokens = append(Tokens, &RuneToken{Id: Token_Delim, Value: '#'})
		case char == '(':
			Tokens = append(Tokens, &EmptyToken{Id: Token_Pren_Open})
			pos++
		case char == ')':
			Tokens = append(Tokens, &EmptyToken{Id: Token_Pren_Close})
			pos++
		case char == '+':
			//
		case char == ',':
			Tokens = append(Tokens, &EmptyToken{Id: Token_Comma})
			pos++
		case char == '-':
			//
		case char == '.':
			//
		case char == ':':
			Tokens = append(Tokens, &EmptyToken{Id: Token_Colon})
			pos++
		case char == ';':
			Tokens = append(Tokens, &EmptyToken{Id: Token_Semicolon})
			pos++
		case char == '<':
			//
		case char == '@':
			//
		case char == '[':
			Tokens = append(Tokens, &EmptyToken{Id: Token_Square_Bracket_Open})
			pos++
		case char == '\\':
			//
		case char == ']':
			Tokens = append(Tokens, &EmptyToken{Id: Token_Square_Bracket_Close})
			pos++
		case char == '{':
			Tokens = append(Tokens, &EmptyToken{Id: Token_Clearly_Open})
			pos++
		case char == '}':
			Tokens = append(Tokens, &EmptyToken{Id: Token_Clearly_Close})
			pos++
		// https://www.w3.org/TR/css-syntax-3/#digit
		case unicode.IsDigit(char):
			// https://www.w3.org/TR/css-syntax-3/#consume-a-numeric-token
		// https://www.w3.org/TR/css-syntax-3/#ident-start-code-point
		case unicode.IsLetter(char) || char == '_':
			// https://www.w3.org/TR/css-syntax-3/#consume-an-ident-like-token
		default:
			Tokens = append(Tokens, &RuneToken{
				Id:    Token_Delim,
				Value: char,
			})
			pos++
		}
	}

	return Tokens, nil
}

// https://www.w3.org/TR/css-syntax-3/#escape-diagram
// https://www.w3.org/TR/css-syntax-3/#string-token-diagram
// https://www.w3.org/TR/css-syntax-3/#consume-string-token
func consumeString(data *[]rune, dataLen int, pos int) {
	chars := []rune{}
	delim := (*data)[pos]
	offset := pos + 1

	for {
		if eof(dataLen, pos) {
			offset++
			// create string token
			/*Tokens = append(Tokens, &StringToken{
				Id:    Token_String,
				Value: codePoints,
			})
			break
		}
		if (*data)[offset] == delim {
			offset++ // eat '\'' or '\"'
			/*Tokens = append(Tokens, &StringToken{
				Id:    Token_String,
				Value: codePoints,
			})
			break
		}
		if (*data)[offset] == '\n' {
			offset++
			//Tokens = append(Tokens, &EmptyToken{Id: Token_Bad_String})
			break // parser error create bad string token
		}
		if (*data)[offset] == '\\' {
			offset++
			// https://www.w3.org/TR/css-syntax-3/#escape-diagram
			if eof(dataLen, offset) {
				continue
			}
			if (*data)[offset] == '\n' {
				offset++
				continue
			}


				Otherwise, (the stream starts with a valid escape)
				consume an escaped code point and append the returned
				code point to the <string-token>’s value.

				https://www.w3.org/TR/css-syntax-3/#consume-escaped-code-point

			//return nil, fmt.Errorf("escaped code point is not fully implemented")
		}

		//codePoints = append(codePoints, data[pos])
		offset++

	}

}*/
