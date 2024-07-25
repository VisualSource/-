package plex_css

type CssParser struct {
	pos    int
	input  []rune
	tokens []Token
}

func (c *CssParser) tokenizer() {

	switch c.input[c.pos] {
	case '"':
		fallthrough
	case '\'':
		// handle string input
		char := c.input[c.pos]
		c.pos++

		data := []rune{}
		id := Token_String
		for c.input[c.pos] != char {
			nextChar := c.input[c.pos]
			if nextChar == '\n' {
				id = Token_Bad_String
				break
			}
			data = append(data, nextChar)
			c.pos++
		}

		c.pos++

		c.tokens = append(c.tokens, &StringToken{
			Value: data,
			Id:    id,
		})

	case '{':
		c.pos++
		c.tokens = append(c.tokens, &EmptyToken{Id: Token_Clearly_Open})
	case '}':
		c.pos++
		c.tokens = append(c.tokens, &EmptyToken{Id: Token_Clearly_Close})
	case '(':
		c.pos++
		c.tokens = append(c.tokens, &EmptyToken{Id: Token_Pren_Open})
	case ')':
		c.pos++
		c.tokens = append(c.tokens, &EmptyToken{Id: Token_Pren_Close})
	case '[':
		c.pos++
		c.tokens = append(c.tokens, &EmptyToken{Id: Token_Square_Bracket_Open})
	case ']':
		c.pos++
		c.tokens = append(c.tokens, &EmptyToken{Id: Token_Square_Bracket_Close})
	case '+':
	case ',':
		c.pos++
		c.tokens = append(c.tokens, &EmptyToken{Id: Token_Comma})
	case ':':
		c.pos++
		c.tokens = append(c.tokens, &EmptyToken{Id: Token_Colon})
	case ';':
		c.pos++
		c.tokens = append(c.tokens, &EmptyToken{Id: Token_Semicolon})
	}

}
