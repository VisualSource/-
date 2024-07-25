package plex_css

import (
	"fmt"
	"unicode"
)

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

func Tokenizer(input string) ([]Token, error) {
	var commentStart = []rune{'/', '*'}
	var commentEnd = []rune{'*', '/'}

	data := []rune(input)
	tokens := []Token{}
	dataLen := len(data)
	pos := 0
	for pos < dataLen {
		if eof(dataLen, pos) {
			tokens = append(tokens, &EmptyToken{Id: Token_EOF})
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
					return nil, fmt.Errorf("expected '*/' but found EOF token")
				}
			}

			pos += 2 // eat '*/'

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
			tokens = append(tokens, &EmptyToken{Id: Token_Whitespace})
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
					tokens = append(tokens, &StringToken{
						Id:    Token_String,
						Value: codePoints,
					})
					break
				}
				if data[pos] == i {
					pos++ // eat '\'' or '\"'
					tokens = append(tokens, &StringToken{
						Id:    Token_String,
						Value: codePoints,
					})
					break
				}

				if data[pos] == '\n' {
					pos++
					tokens = append(tokens, &EmptyToken{Id: Token_Bad_String})
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
					*/
					return nil, fmt.Errorf("escaped code point is not fully implemented")
				}

				codePoints = append(codePoints, data[pos])
				pos++
			}
		// https://www.w3.org/TR/css-syntax-3/#consume-token
		// #
		case char == '#':
			pos++ //eat '#'

			tokens = append(tokens, &RuneToken{Id: Token_Delim, Value: '#'})
		case char == '(':
			tokens = append(tokens, &EmptyToken{Id: Token_Pren_Open})
			pos++
		case char == ')':
			tokens = append(tokens, &EmptyToken{Id: Token_Pren_Close})
			pos++
		case char == '+':
			//
		case char == ',':
			tokens = append(tokens, &EmptyToken{Id: Token_Comma})
			pos++
		case char == '-':
			//
		case char == '.':
			//
		case char == ':':
			tokens = append(tokens, &EmptyToken{Id: Token_Colon})
			pos++
		case char == ';':
			tokens = append(tokens, &EmptyToken{Id: Token_Semicolon})
			pos++
		case char == '<':
			//
		case char == '@':
			//
		case char == '[':
			tokens = append(tokens, &EmptyToken{Id: Token_Square_Bracket_Open})
			pos++
		case char == '\\':
			//
		case char == ']':
			tokens = append(tokens, &EmptyToken{Id: Token_Square_Bracket_Close})
			pos++
		case char == '{':
			tokens = append(tokens, &EmptyToken{Id: Token_Clearly_Open})
			pos++
		case char == '}':
			tokens = append(tokens, &EmptyToken{Id: Token_Clearly_Close})
			pos++
		// https://www.w3.org/TR/css-syntax-3/#digit
		case unicode.IsDigit(char):
			// https://www.w3.org/TR/css-syntax-3/#consume-a-numeric-token
		// https://www.w3.org/TR/css-syntax-3/#ident-start-code-point
		case unicode.IsLetter(char) || char == '_':
			// https://www.w3.org/TR/css-syntax-3/#consume-an-ident-like-token
		default:
			tokens = append(tokens, &RuneToken{
				Id:    Token_Delim,
				Value: char,
			})
			pos++
		}
	}

	return tokens, nil
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
			/*tokens = append(tokens, &StringToken{
				Id:    Token_String,
				Value: codePoints,
			})*/
			break
		}
		if (*data)[offset] == delim {
			offset++ // eat '\'' or '\"'
			/*tokens = append(tokens, &StringToken{
				Id:    Token_String,
				Value: codePoints,
			})*/
			break
		}
		if (*data)[offset] == '\n' {
			offset++
			//tokens = append(tokens, &EmptyToken{Id: Token_Bad_String})
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

			/*
				Otherwise, (the stream starts with a valid escape)
				consume an escaped code point and append the returned
				code point to the <string-token>’s value.

				https://www.w3.org/TR/css-syntax-3/#consume-escaped-code-point
			*/
			//return nil, fmt.Errorf("escaped code point is not fully implemented")
		}

		//codePoints = append(codePoints, data[pos])
		offset++

	}

}
