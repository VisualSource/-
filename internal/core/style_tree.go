package plex

import (
	"sort"
	plex_css "visualsource/plex/internal/css"

	"github.com/moznion/go-optional"
)

// #region-start StyleNode
type StyledNode struct {
	node     Node
	props    plex_css.CssPropertyMap
	children []StyledNode
}

func (n *StyledNode) GetDisplay() DisplayType {

	value := n.props.GetProp("display")

	if value.IsNone() {
		return DisplayType_Inline
	}

	display := value.Unwrap()

	if item, ok := display.GetValue().(*plex_css.CssKeyword); ok {
		switch item.Value {
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

func CreateStyleNode(node Node, props plex_css.CssPropertyMap, children []StyledNode) StyledNode {
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

func specifiedValues(el *ElementNode, stylesheets []plex_css.Stylesheet) plex_css.CssPropertyMap {
	values := plex_css.CssPropertyMap{}

	rules := []MatchedRule{}
	for _, stylesheet := range stylesheets {
		rules = append(rules, matchRules(el, &stylesheet)...)
	}

	sort.Slice(rules, func(i, j int) bool {
		a := rules[i]
		b := rules[j]
		return a.Orgin < b.Orgin || a.specificity.Less(&b.specificity)
	})

	for _, matched := range rules {
		for _, dec := range matched.rule.Block {
			values[dec.Name] = dec
		}
	}

	return values
}

func StyleTree(root Node, stylesheet []plex_css.Stylesheet) StyledNode {

	var specified plex_css.CssPropertyMap
	children := []StyledNode{}
	if node, ok := (root).(*ElementNode); ok {

		specified = specifiedValues(node, stylesheet)

		for _, child := range node.GetChildren() {
			children = append(children, StyleTree(child, stylesheet))
		}
	} else {
		specified = plex_css.CssPropertyMap{}
	}

	return StyledNode{
		node:     root,
		props:    specified,
		children: children,
	}
}
