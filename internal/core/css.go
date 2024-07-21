package plex

import (
	"fmt"
	"sort"
	"unicode"
)

type Specificity struct {
	a int
	b int
	c int
}

type CssValue interface{}

type Selector struct {
	tag_name string
	id       string
	class    []string
}

func NewSelector(tagName string, id string, classlist []string) Selector {
	return Selector{
		tag_name: tagName,
		id:       id,
		class:    classlist,
	}
}

// http://www.w3.org/TR/selectors/#specificity
func (s *Selector) Specificity() Specificity {
	a := 0
	c := 0

	if s.id != "" {
		a++
	}

	b := len(s.class)

	if s.tag_name != "" {
		c++
	}

	return Specificity{a, b, c}
}

type Declaration struct {
	name  string
	value CssValue
}

type Rule struct {
	origin      int
	selectors   []Selector
	declartions []Declaration
}

type Stylesheet struct {
	rules []Rule
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
	selector := Selector{tag_name: "", id: "", class: []string{}}

L:
	for !p.parser.EOF() {
		char := p.parser.NextChar()

		switch char {
		case '#':
			p.parser.ConsumeChar()
			selector.id = p.parseIdentifier()
		case '.':
			p.parser.ConsumeChar()
			selector.class = append(selector.class, p.parseIdentifier())
		case '*':
			p.parser.ConsumeChar()
		default:
			if isValidIdentifier(char) {
				selector.tag_name = p.parseIdentifier()
				continue
			}
			break L
		}

	}

	return selector
}

func isValidIdentifier(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_'
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

		return a.a > b.a || a.b > b.b || a.c > b.c
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

func (p *CssParser) parseValue() (CssValue, error) {

	result := p.parser.ConsumeWhile(func(r rune) bool { return r != ';' })

	return string(result), nil
}
