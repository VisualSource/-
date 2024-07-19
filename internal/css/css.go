package plex

import (
	"fmt"
	"unicode"
)

type CssValue interface{}

type Selector struct {
	tag_name string
	id       string
	class    []string
}
type Declaration struct {
	name  string
	value CssValue
}

type Rule struct {
	selectors   []Selector
	declartions []Declaration
}

type Stylesheet struct {
	rules []Rule
}

type CssParser struct {
	pos   int
	input []rune
}

func (p *CssParser) Parse(sheet string) (Stylesheet, error) {
	p.pos = 0
	p.input = []rune(sheet)

	rules, err := p.parseRules()
	if err != nil {
		return Stylesheet{}, err
	}

	return Stylesheet{
		rules,
	}, nil
}

func (p *CssParser) NextChar() rune {
	return p.input[p.pos]
}

func (p *CssParser) startsWith(s []rune) bool {
	isSame := true
	for i, v := range s {
		if p.input[p.pos+i] != v {
			isSame = false
			break
		}
	}

	return isSame
}

func (p *CssParser) Expect(s []rune) error {
	if p.startsWith(s) {
		p.pos += len(s)
		return nil
	}

	return fmt.Errorf("did not find Expected: '%s'", string(s))
}

func (p *CssParser) eof() bool {
	return p.pos >= (len(p.input) - 1)
}

func (p *CssParser) consumeChar() rune {
	r := p.NextChar()
	p.pos++
	return r
}

func (p *CssParser) consumeWhile(test func(rune) bool) []rune {
	result := []rune{}

	for !p.eof() && test(p.NextChar()) {
		result = append(result, p.consumeChar())
	}

	return result
}

func (p *CssParser) consumeWhitespace() {
	p.consumeWhile(unicode.IsSpace)
}

func (p *CssParser) parseRules() ([]Rule, error) {}

func (p *CssParser) parseRule() (Rule, error) {

}
