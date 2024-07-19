package plex

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
