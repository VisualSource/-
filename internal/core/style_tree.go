package plex

import (
	"sort"

	"github.com/moznion/go-optional"
)

type PropertyMap = map[string]CssValue
type StyledNode struct {
	node     *Node
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

func (n *StyledNode) GetDisplay() DisplayType {

	display := n.props["display"]

	if item, ok := display.(*string); ok {
		switch *item {
		case "block":
			return DisplayBlock
		case "inline-block":
			return DisplayInlineBlock
		case "flex":
			return DisplayFlex
		case "none":
			return DisplayNone
		case "inline":
			fallthrough
		default:
			return DisplayInline
		}
	}

	return DisplayNone
}

type MatchedRule struct {
	specificity Specificity
	rule        Rule
}

func matchRule(el *ElementNode, rule Rule) optional.Option[MatchedRule] {

	for _, selector := range rule.selectors {
		if el.Matches(selector) {
			return optional.Some(MatchedRule{
				specificity: selector.Specificity(),
				rule:        rule,
			})
		}
	}

	return nil
}

func matchRules(el *ElementNode, stylesheet *Stylesheet) []MatchedRule {
	rules := []MatchedRule{}

	for _, rule := range stylesheet.rules {
		result := matchRule(el, rule)
		if result.IsSome() {
			rules = append(rules, result.Unwrap())
		}
	}

	return rules
}

func specifiedValues(el *ElementNode, stylesheets []Stylesheet) PropertyMap {
	values := map[string]CssValue{}

	rules := []MatchedRule{}
	for _, stylesheet := range stylesheets {
		rules = append(rules, matchRules(el, &stylesheet)...)
	}

	sort.Slice(rules, func(i, j int) bool {
		a := rules[i]
		b := rules[j]

		return a.rule.origin < b.rule.origin || a.specificity.a < b.specificity.a || a.specificity.b < b.specificity.b || a.specificity.c < b.specificity.c
	})

	for _, matched := range rules {
		for _, dec := range matched.rule.declartions {
			values[dec.name] = dec.value
		}
	}

	return values
}

func StyleTree(root *Node, stylesheet []Stylesheet) StyledNode {

	var specified PropertyMap
	children := []StyledNode{}
	if node, ok := (*root).(*ElementNode); ok {

		specified = specifiedValues(node, stylesheet)

		for _, child := range node.GetChildren() {
			children = append(children, StyleTree(&child, stylesheet))
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
