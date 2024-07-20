package plex

import (
	"fmt"
	"unicode"
)

type Parser struct {
	pos   int
	input []rune
}

func (p *Parser) SetPos(pos int) {
	p.pos = pos
}

func (p *Parser) SetInput(value string) {
	p.input = []rune(value)
}

// Read the current character without consuming it.
func (p *Parser) NextChar() rune {
	return p.input[p.pos]
}

func (p *Parser) StartsWith(s []rune) bool {
	isSame := true
	for i, v := range s {
		if p.input[p.pos+i] != v {
			isSame = false
			break
		}
	}

	return isSame
}

func (p *Parser) ExpectRune(c rune) error {
	if p.input[p.pos] == c {
		p.pos += 1
		return nil
	}

	return fmt.Errorf("was expecting rune %s but got '%s'", string(c), string(p.input[p.pos]))
}

func (p *Parser) Expect(s []rune) error {
	if p.StartsWith(s) {
		p.pos += len(s)
		return nil
	}

	return fmt.Errorf("did not find Expected: '%s'", string(s))
}

func (p *Parser) EOF() bool {
	return p.pos >= (len(p.input) - 1)
}

func (p *Parser) ConsumeChar() rune {
	r := p.NextChar()
	p.pos++
	return r
}

func (p *Parser) ConsumeWhile(test func(rune) bool) []rune {
	result := []rune{}

	for !p.EOF() && test(p.NextChar()) {
		result = append(result, p.ConsumeChar())
	}

	return result
}

func (p *Parser) ConsumeWhitespace() {
	p.ConsumeWhile(unicode.IsSpace)
}
