package plex

import (
	"fmt"
	"unicode"
	plex_dom "visualsource/plex/internal/dom"
)

// https://limpet.net/mbrubeck/2014/08/13/toy-layout-engine-3-css.html
var OPEN_BRACKET = [1]rune{'<'}
var CLOSED_BRACKET = [1]rune{'>'}
var EQUAL = [1]rune{'='}
var ELEMENT_CLOSED_BRACKET = [2]rune{'<', '/'}
var COMMENT_START = [4]rune{'<', '!', '-', '-'}
var COMMENT_END = [3]rune{'-', '-', '>'}

type Parser struct {
	pos   int
	input []rune
}

func CreateParser() Parser {
	return Parser{pos: 0, input: []rune{}}
}

func NewParser(input string) Parser {
	return Parser{
		pos:   0,
		input: []rune(input),
	}
}

func (p *Parser) Parse(document string) (plex_dom.Node, error) {
	p.input = []rune(document)
	p.pos = 0

	nodes, err := p.parseNodes()
	if err != nil {
		return nil, err
	}

	if len(nodes) == 1 {
		return nodes[0], nil
	}

	root := plex_dom.CreateElementNode("html", plex_dom.AttributeMap{}, nodes)

	return &root, nil
}

// Read the current character without consuming it.
func (p *Parser) NextChar() rune {
	return p.input[p.pos]
}

func (p *Parser) startsWith(s []rune) bool {
	isSame := true
	for i, v := range s {
		if p.input[p.pos+i] != v {
			isSame = false
			break
		}
	}

	return isSame
}

func (p *Parser) Expect(s []rune) error {
	if p.startsWith(s) {
		p.pos += len(s)
		return nil
	}

	return fmt.Errorf("did not find Expected: '%s'", string(s))
}

func (p *Parser) eof() bool {
	return p.pos >= (len(p.input) - 1)
}

func (p *Parser) consumeChar() rune {
	r := p.NextChar()
	p.pos++
	return r
}

func (p *Parser) consumeWhile(test func(rune) bool) []rune {
	result := []rune{}

	for !p.eof() && test(p.NextChar()) {
		result = append(result, p.consumeChar())
	}

	return result
}

func (p *Parser) consumeWhitespace() {
	p.consumeWhile(unicode.IsSpace)
}

func (p *Parser) ParseName() []rune {
	return p.consumeWhile(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	})
}

// Parse a single node.
func (p *Parser) parseNode() (plex_dom.Node, error) {
	if p.startsWith(COMMENT_START[:]) {
		return p.parseComment()
	}
	if p.startsWith(OPEN_BRACKET[:]) {
		return p.parseElement()
	}
	return p.parseText()
}

func (p *Parser) parseComment() (plex_dom.Node, error) {

	err := p.Expect(COMMENT_START[:])
	if err != nil {
		return nil, err
	}

	content := []rune{}
	for !p.eof() && !p.startsWith(COMMENT_END[:]) {
		content = append(content, p.consumeChar())
	}

	err = p.Expect(COMMENT_END[:])
	if err != nil {
		return nil, err
	}

	node := plex_dom.CreateCommentNode(string(content))

	return &node, nil
}

func (p *Parser) parseElement() (plex_dom.Node, error) {

	err := p.Expect(OPEN_BRACKET[:])

	if err != nil {
		return nil, err
	}

	tagName := p.ParseName()
	attrs, err := p.ParseAttributes()

	if err != nil {
		return nil, err
	}

	err = p.Expect(CLOSED_BRACKET[:])
	if err != nil {
		return nil, err
	}

	children, err := p.parseNodes()

	if err != nil {
		return nil, err
	}

	err = p.Expect(ELEMENT_CLOSED_BRACKET[:])
	if err != nil {
		return nil, err
	}

	err = p.Expect(tagName)
	if err != nil {
		return nil, err
	}

	err = p.Expect(CLOSED_BRACKET[:])
	if err != nil {
		return nil, err
	}

	result := plex_dom.CreateElementNode(string(tagName), attrs, children)

	return &result, nil
}

func (p *Parser) parseText() (plex_dom.Node, error) {
	result := p.consumeWhile(func(r rune) bool { return r != '<' })

	textNode := plex_dom.CreateTextNode(string(result))

	return &textNode, nil
}

func (p *Parser) ParseAttr() (string, string, error) {
	name := p.ParseName()

	err := p.Expect(EQUAL[:])
	if err != nil {
		return "", "", err
	}

	value, err := p.ParseAttrValue()

	if err != nil {
		return "", "", err
	}

	return string(name), value, nil
}

func (p *Parser) ParseAttrValue() (string, error) {
	openQuote := p.consumeChar()
	if openQuote != '"' && openQuote != '\'' {
		return "", fmt.Errorf("was expecting a '\"' or ''' but found '%s'", string(openQuote))
	}

	value := p.consumeWhile(func(r rune) bool { return r != openQuote })

	closeQuote := p.consumeChar()

	if closeQuote != '"' && closeQuote != '\'' {
		return "", fmt.Errorf("was expecting a '\"' or ''' but found '%s'", string(closeQuote))
	}

	return string(value), nil
}

func (p *Parser) ParseAttributes() (plex_dom.AttributeMap, error) {
	attributes := plex_dom.AttributeMap{}

	for {
		p.consumeWhitespace()

		if p.NextChar() == '>' {
			break
		}

		name, value, err := p.ParseAttr()

		if err != nil {
			return nil, err
		}

		attributes[name] = value
	}

	return attributes, nil
}

func (p *Parser) parseNodes() ([]plex_dom.Node, error) {
	nodes := []plex_dom.Node{}

	for {
		p.consumeWhitespace()

		if p.eof() || p.startsWith(ELEMENT_CLOSED_BRACKET[:]) {
			break
		}

		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)

	}

	return nodes, nil
}
