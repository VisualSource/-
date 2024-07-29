package plex

import (
	"sort"
	plex_css "visualsource/plex/internal/css"

	"github.com/moznion/go-optional"
)

type PropertyMap = map[string]CssValue

// #region-start StyleNode
type StyledNode struct {
	node     Node
	props    PropertyMap
	children []StyledNode
}

func (n *StyledNode) GetProp(key string) optional.Option[CssValue] {
	value := n.props[key]
	if value == nil {
		return nil
	}

	return optional.Some(value)
}

func (n *StyledNode) GetPropAsLength(key string) optional.Option[CssLengthValue] {
	value := n.props[key]
	if value == nil {
		return nil
	}
	if i, ok := value.(CssLengthValue); ok {
		return optional.Some(i)
	}

	return nil
}

func (n *StyledNode) GetDisplay() DisplayType {

	value := n.GetProp("display")

	if value.IsNone() {
		return DisplayType_Inline
	}

	display := value.Unwrap()

	if item, ok := display.(string); ok {
		switch item {
		case "block":
			return DisplayType_Block
		case "none":
			return DisplayType_None
		case "inline":
			fallthrough
		default:
			return DisplayType_Inline
		}
	}

	return DisplayType_Inline
}

func (n *StyledNode) LookupCssLength(props ...string) optional.Option[CssLengthValue] {

	for _, prop := range props {
		item := n.props[prop]

		if i, ok := item.(*CssLengthValue); ok {
			return optional.Some(*i)
		}
	}

	return nil
}

func (n *StyledNode) Lookup(props ...string) optional.Option[CssValue] {

	for _, prop := range props {
		item := n.props[prop]

		if item != nil {
			return optional.Some(item)
		}

	}

	return nil
}

func CreateStyleNode(node Node, props PropertyMap, children []StyledNode) StyledNode {
	return StyledNode{
		node:     node,
		props:    props,
		children: children,
	}
}

type MatchedRule struct {
	specificity plex_css.Specificity
	rule        plex_css.Rule
	Orgin       uint
}

// #region-start utility

func matchRule(el *ElementNode, rule plex_css.Rule) optional.Option[MatchedRule] {

	for _, selector := range rule.Selector {
		if el.Matches(&selector) {
			return optional.Some(MatchedRule{
				specificity: selector.GetSpecificity(),
				rule:        rule,
			})
		}
	}

	return nil
}

func matchRules(el *ElementNode, stylesheet *plex_css.Stylesheet) []MatchedRule {
	rules := []MatchedRule{}

	for _, rule := range stylesheet.Rules {
		result := matchRule(el, rule)
		if result.IsSome() {
			item := result.Unwrap()
			item.Orgin = stylesheet.Origin
			rules = append(rules, item)
		}
	}

	return rules
}

func specifiedValues(el *ElementNode, stylesheets []plex_css.Stylesheet) PropertyMap {
	values := PropertyMap{}

	rules := []MatchedRule{}
	for _, stylesheet := range stylesheets {
		rules = append(rules, matchRules(el, &stylesheet)...)
	}

	sort.Slice(rules, func(i, j int) bool {
		a := rules[i]
		b := rules[j]

		return a.Orgin < b.Orgin || a.specificity.A < b.specificity.A || a.specificity.B < b.specificity.B || a.specificity.C < b.specificity.C
	})

	for _, matched := range rules {
		for _, dec := range matched.rule.Block {
			values[dec.name] = dec.value
		}
	}

	return values
}

func StyleTree(root Node, stylesheet []plex_css.Stylesheet) StyledNode {

	var specified PropertyMap
	children := []StyledNode{}
	if node, ok := (root).(*ElementNode); ok {

		specified = specifiedValues(node, stylesheet)

		for _, child := range node.GetChildren() {
			children = append(children, StyleTree(child, stylesheet))
		}
	} else {
		specified = PropertyMap{}
	}

	return StyledNode{
		node:     root,
		props:    specified,
		children: children,
	}
}
