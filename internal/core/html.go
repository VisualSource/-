package plex

import (
	"fmt"
	"unicode"
)

// https://limpet.net/mbrubeck/2014/08/13/toy-layout-engine-3-css.html
var OPEN_BRACKET = [1]rune{'<'}
var CLOSED_BRACKET = [1]rune{'>'}
var EQUAL = [1]rune{'='}
var ELEMENT_CLOSED_BRACKET = [2]rune{'<', '/'}
var COMMENT_START = [4]rune{'<', '!', '-', '-'}
var COMMENT_END = [3]rune{'-', '-', '>'}

type HtmlParser struct {
	parser Parser
}

func CreateHtmlParser() HtmlParser {
	return HtmlParser{
		parser: Parser{},
	}
}

func (p *HtmlParser) Parse(document string) (Node, error) {
	p.parser.SetInput(document)
	p.parser.SetPos(0)

	nodes, err := p.parseNodes()
	if err != nil {
		return nil, err
	}

	if len(nodes) == 1 {
		return nodes[0], nil
	}

	root := CreateElementNode("html", AttributeMap{}, nodes)

	return &root, nil
}

func (hp *HtmlParser) ParseName() []rune {
	return hp.parser.ConsumeWhile(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	})
}

// Parse a single node.
func (hp *HtmlParser) parseNode() (Node, error) {
	if hp.parser.StartsWith(COMMENT_START[:]) {
		return hp.parseComment()
	}
	if hp.parser.StartsWith(OPEN_BRACKET[:]) {
		return hp.parseElement()
	}
	return hp.parseText()
}

func (hp *HtmlParser) parseComment() (Node, error) {

	err := hp.parser.Expect(COMMENT_START[:])
	if err != nil {
		return nil, err
	}

	content := []rune{}
	for !hp.parser.EOF() && !hp.parser.StartsWith(COMMENT_END[:]) {
		content = append(content, hp.parser.ConsumeChar())
	}

	err = hp.parser.Expect(COMMENT_END[:])
	if err != nil {
		return nil, err
	}

	node := CreateCommentNode(string(content))

	return &node, nil
}

func (hp *HtmlParser) parseElement() (Node, error) {

	err := hp.parser.Expect(OPEN_BRACKET[:])

	if err != nil {
		return nil, err
	}

	tagName := hp.ParseName()
	attrs, err := hp.ParseAttributes()

	if err != nil {
		return nil, err
	}

	err = hp.parser.Expect(CLOSED_BRACKET[:])
	if err != nil {
		return nil, err
	}

	children, err := hp.parseNodes()

	if err != nil {
		return nil, err
	}

	err = hp.parser.Expect(ELEMENT_CLOSED_BRACKET[:])
	if err != nil {
		return nil, err
	}

	err = hp.parser.Expect(tagName)
	if err != nil {
		return nil, err
	}

	err = hp.parser.Expect(CLOSED_BRACKET[:])
	if err != nil {
		return nil, err
	}

	result := CreateElementNode(string(tagName), attrs, children)

	return &result, nil
}

func (hp *HtmlParser) parseText() (Node, error) {
	result := hp.parser.ConsumeWhile(func(r rune) bool { return r != '<' })

	textNode := CreateTextNode(string(result))

	return &textNode, nil
}

func (hp *HtmlParser) ParseAttr() (string, string, error) {
	name := hp.ParseName()

	err := hp.parser.Expect(EQUAL[:])
	if err != nil {
		return "", "", err
	}

	value, err := hp.ParseAttrValue()

	if err != nil {
		return "", "", err
	}

	return string(name), value, nil
}

func (hp *HtmlParser) ParseAttrValue() (string, error) {
	openQuote := hp.parser.ConsumeChar()
	if openQuote != '"' && openQuote != '\'' {
		return "", fmt.Errorf("was expecting a '\"' or ''' but found '%s'", string(openQuote))
	}

	value := hp.parser.ConsumeWhile(func(r rune) bool { return r != openQuote })

	closeQuote := hp.parser.ConsumeChar()

	if closeQuote != '"' && closeQuote != '\'' {
		return "", fmt.Errorf("was expecting a '\"' or ''' but found '%s'", string(closeQuote))
	}

	return string(value), nil
}

func (hp *HtmlParser) ParseAttributes() (AttributeMap, error) {
	attributes := AttributeMap{}

	for {
		hp.parser.ConsumeWhitespace()

		if hp.parser.NextChar() == '>' {
			break
		}

		name, value, err := hp.ParseAttr()

		if err != nil {
			return nil, err
		}

		attributes[name] = value
	}

	return attributes, nil
}

func (hp *HtmlParser) parseNodes() ([]Node, error) {
	nodes := []Node{}

	for {
		hp.parser.ConsumeWhitespace()

		if hp.parser.EOF() || hp.parser.StartsWith(ELEMENT_CLOSED_BRACKET[:]) {
			break
		}

		node, err := hp.parseNode()
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)

	}

	return nodes, nil
}
