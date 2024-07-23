package plex

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/veandco/go-sdl2/sdl"
)

func isValidIdentifier(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_'
}

type CssParser struct {
	parser Parser
}

func (p *CssParser) Parse(sheet string, origin int) (Stylesheet, error) {
	p.parser.SetInput(sheet)
	p.parser.SetPos(0)

	rules := p.parseRules(origin)

	return Stylesheet{
		rules: rules,
	}, nil
}

func (p *CssParser) ParseSelector() Selector {
	selector := Selector{}

L:
	for !p.parser.EOF() {
		char := p.parser.NextChar()

		switch char {
		case '#':
			p.parser.ConsumeChar()
			selector.Id = p.parseIdentifier()
		case '.':
			p.parser.ConsumeChar()
			selector.Classes = append(selector.Classes, p.parseIdentifier())
		case '*':
			p.parser.ConsumeChar()
		default:
			if isValidIdentifier(char) {
				selector.TagName = p.parseIdentifier()
				continue
			}
			break L
		}

	}

	return selector
}

func (p *CssParser) parseIdentifier() string {
	result := p.parser.ConsumeWhile(isValidIdentifier)

	return string(result)
}

func (p *CssParser) parseRules(origin int) []Rule {
	rules := []Rule{}

	for {
		p.parser.ConsumeWhitespace()
		if p.parser.EOF() {
			break
		}

		rule, err := p.parseRule(origin)

		if err != nil {
			continue
		}

		rules = append(rules, rule)
	}

	return rules
}

func (p *CssParser) parseRule(origin int) (Rule, error) {
	selectors, err := p.parseSelectors()

	if err != nil {
		return Rule{}, err
	}

	declarations := p.parseDeclarations()

	return Rule{
		origin,
		selectors,
		declarations,
	}, nil
}

func (p *CssParser) parseSelectors() ([]Selector, error) {
	selectors := []Selector{}

L:
	for {
		selectors = append(selectors, p.ParseSelector())
		p.parser.ConsumeWhitespace()

		char := p.parser.NextChar()

		switch char {
		case ',':
			p.parser.ConsumeChar()
			p.parser.ConsumeWhitespace()
		case '{':
			break L
		default:
			return nil, fmt.Errorf("unepected charactor %s in selector list", string(char))
		}
	}

	sort.Slice(selectors, func(i, j int) bool {
		a := selectors[i].Specificity()
		b := selectors[j].Specificity()

		return a.A > b.A || a.B > b.B || a.C > b.C
	})

	return selectors, nil
}

func (p *CssParser) parseDeclarations() []Declaration {
	p.parser.ExpectRune('{')

	declarations := []Declaration{}

	for {
		p.parser.ConsumeWhitespace()

		if p.parser.NextChar() == '}' {
			p.parser.ConsumeChar()
			break
		}

		declaration, err := p.parseDeclaration()
		if err != nil {
			p.parser.ConsumeWhile(func(r rune) bool { return r != '}' })
			continue
		}

		declarations = append(declarations, declaration)
	}

	return declarations
}

func (p *CssParser) parseDeclaration() (Declaration, error) {
	name := p.parseIdentifier()
	p.parser.ConsumeWhitespace()
	err := p.parser.ExpectRune(':')
	if err != nil {
		return Declaration{}, err
	}
	p.parser.ConsumeWhitespace()

	value, err := p.parseValue()
	if err != nil {
		return Declaration{}, err
	}

	err = p.parser.ExpectRune(';')
	if err != nil {
		return Declaration{}, err
	}

	return Declaration{
		name,
		value,
	}, nil
}

func (p *CssParser) parseHexPair() (uint8, error) {
	s := p.parser.input[p.parser.pos : p.parser.pos+2]
	p.parser.pos += 2

	u8, err := strconv.ParseUint(string(s), 16, 8)
	if err != nil {
		return 0, err
	}
	return uint8(u8), nil
}

func (p *CssParser) parseValue() (CssValue, error) {

	char := p.parser.NextChar()

	if unicode.IsDigit(char) {

		float := p.parser.ConsumeWhile(func(r rune) bool {
			return unicode.IsDigit(r) || r == '.'
		})

		item, err := strconv.ParseFloat(string(float), 32)
		if err != nil {
			return nil, err
		}

		unitStr := strings.ToLower(p.parseIdentifier())

		var unit CssLengthUnit
		switch unitStr {
		case "px":
			unit = CssUnit_PX
		default:
			return nil, fmt.Errorf("unrecognized unit")
		}

		return CssLengthValue{
			Value: float32(item),
			Unit:  unit,
		}, nil

	} else if char == '#' {
		p.parser.ConsumeChar()

		r, err := p.parseHexPair()
		if err != nil {
			return nil, err
		}
		g, err := p.parseHexPair()
		if err != nil {
			return nil, err
		}
		b, err := p.parseHexPair()
		if err != nil {
			return nil, err
		}

		return sdl.Color{
			A: 255,
			R: r,
			G: g,
			B: b,
		}, nil
	}

	ident := p.parseIdentifier()

	return ident, nil
}
