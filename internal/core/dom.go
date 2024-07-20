package plex

import (
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

type NodeType string

// https://www.sohamkamani.com/golang/enums/
const (
	Text    NodeType = "Text"
	Element NodeType = "Element"
	Comment NodeType = "Comment"
)

type Node interface {
	GetType() NodeType
	GetChildren() []Node
}

// #region-start TextNode

type TextNode struct {
	children []Node
	content  string
}

func (n *TextNode) GetType() NodeType {
	return Text
}

func (n *TextNode) GetChildren() []Node {
	return n.children
}

func (n *TextNode) GetTextContent() string {
	return n.content
}

func (n *TextNode) SetTextContent(value string) {
	n.content = value
}

func IsTextNode(n Node) bool {
	return n.GetType() == Text
}

func CreateTextNode(textContent string) TextNode {
	return TextNode{
		children: []Node{},
		content:  textContent,
	}
}

// #region-start ElementNode

// https://go.dev/blog/maps
type AttributeMap = map[string]string

type ElementNode struct {
	children []Node
	attr     AttributeMap
	tag_name string
}

func (n *ElementNode) QuerySelector(selector Selector) (*ElementNode, error) {

	for _, child := range n.children {
		if child.GetType() == Text || child.GetType() == Comment {
			continue
		}

		if p, ok := child.(*ElementNode); ok {
			if p.Matches(selector) {
				return p, nil
			}

			result, err := p.QuerySelector(selector)
			if err != nil {
				return nil, err
			}
			if result != nil {
				return result, nil
			}
		}
	}

	return nil, nil
}

func (n *ElementNode) QuerySelectorAll(selector Selector) []*ElementNode {
	result := []*ElementNode{}

	for _, child := range n.GetChildren() {
		if child.GetType() == Text || child.GetType() == Comment {
			continue
		}

		if p, ok := child.(*ElementNode); ok {

			if p.Matches(selector) {
				result = append(result, p)
			}

			children := p.QuerySelectorAll(selector)

			result = append(result, children...)
		}
	}

	return result
}

func (n *ElementNode) GetTextContent() string {

	textNodes := []*TextNode{}

	for _, v := range n.GetChildren() {
		if v.GetType() != Text {
			continue
		}
		if node, ok := v.(*TextNode); ok {
			textNodes = append(textNodes, node)
		}
	}

	var output string
	for _, node := range textNodes {
		output += node.GetTextContent()
	}

	return output
}

func (n *ElementNode) GetId() string {
	return n.GetAttribute("id")
}

func (n *ElementNode) GetClassList() mapset.Set[string] {
	class := n.GetAttribute("class")
	items := strings.Split(class, " ")

	return mapset.NewSet(items...)
}

func (n *ElementNode) Matches(selector Selector) bool {

	if selector.id != "" && selector.id == n.GetAttribute("id") {
		return true
	}

	if selector.tag_name != "" && selector.tag_name == n.tag_name {
		return true
	}

	classes := n.GetClassList()

	return classes.ContainsAny(selector.class...)
}

func (n *ElementNode) GetTagName() string {
	return n.tag_name
}

func (n *ElementNode) GetType() NodeType {
	return Element
}

func (n *ElementNode) GetChildren() []Node {
	return n.children
}

func (n *ElementNode) SetAttribute(key string, value string) {
	n.attr[key] = value
}

func (n *ElementNode) GetAttribute(key string) string {
	return n.attr[key]
}

func (n *ElementNode) RemoveAttribute(key string) {
	delete(n.attr, key)
}

func (n *ElementNode) HasAttribute(key string) bool {
	_, ok := n.attr[key]

	return ok
}

func IsElementNode(n Node) bool {
	return n.GetType() == Element
}

func CreateElementNode(tag_name string, attrs AttributeMap, children []Node) ElementNode {
	return ElementNode{
		tag_name: tag_name,
		attr:     attrs,
		children: children,
	}
}

// #region-start CommentNode

type CommentNode struct {
	children []Node
	content  string
}

func (n *CommentNode) GetType() NodeType {
	return Comment
}

func (n *CommentNode) GetChildren() []Node {
	return n.children
}

func (n *CommentNode) GetTextContent() string {
	return n.content
}

func (n *CommentNode) SetTextContent(value string) {
	n.content = value
}

func IsCommentNode(n Node) bool {
	return n.GetType() == Comment
}

func CreateCommentNode(content string) CommentNode {
	return CommentNode{
		content:  content,
		children: []Node{},
	}
}
