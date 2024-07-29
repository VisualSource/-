package plex

import (
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

type NodeType uint8

// https://www.sohamkamani.com/golang/enums/
const (
	NodeType_Text    NodeType = 0
	NodeType_Element NodeType = 1
	NodeType_Comment NodeType = 2
)

type Node interface {
	GetType() NodeType
	GetChildren() []Node
}

// #region-start TextNode

type TextNode struct {
	content  string
	children []Node
}

func (n *TextNode) GetType() NodeType {
	return NodeType_Text
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

func CreateTextNode(textContent string) TextNode {
	return TextNode{
		content: textContent,
	}
}

// #region-start ElementNode

// https://go.dev/blog/maps
type AttributeMap = map[string]string

type ElementNode struct {
	tagName  string
	attr     AttributeMap
	children []Node
}

func (n *ElementNode) QuerySelector(selector *Selector) (*ElementNode, error) {

	for _, child := range n.children {
		if child.GetType() == NodeType_Text || child.GetType() == NodeType_Comment {
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

func (n *ElementNode) QuerySelectorAll(selector *Selector) []*ElementNode {
	result := []*ElementNode{}

	for _, child := range n.GetChildren() {
		if child.GetType() == NodeType_Text || child.GetType() == NodeType_Comment {
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
		if v.GetType() != NodeType_Text {
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

// https://www.w3.org/TR/selectors-4/#match-against-element
func (n *ElementNode) Matches(selector *Selector) bool {

	if selector.Id != "" && selector.Id == n.GetId() {
		return true
	}

	if selector.TagName != "" && selector.TagName == n.tagName {
		return true
	}

	classes := n.GetClassList()

	return classes.ContainsAny(selector.Classes...)
}

func (n *ElementNode) GetTagName() string {
	return n.tagName
}

func (n *ElementNode) GetType() NodeType {
	return NodeType_Element
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

func CreateElementNode(tagName string, attrs AttributeMap, children []Node) ElementNode {
	return ElementNode{
		tagName:  tagName,
		attr:     attrs,
		children: children,
	}
}

// #region-start CommentNode

type CommentNode struct {
	content  string
	children []Node
}

func (n *CommentNode) GetType() NodeType {
	return NodeType_Comment
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

func CreateCommentNode(content string) CommentNode {
	return CommentNode{
		content: content,
	}
}
