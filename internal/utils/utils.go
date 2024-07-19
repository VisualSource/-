package plex

import (
	"fmt"
	"strings"
	plex_dom "visualsource/plex/internal/dom"
)

func PrintDom(n plex_dom.Node, indent int) {

	if p, ok := n.(*plex_dom.ElementNode); ok {
		fmt.Printf("└-%s %s TAG:%s\n", strings.Repeat("-", indent), p.GetType(), p.GetTagName())
	} else if p, ok := n.(*plex_dom.TextNode); ok {
		fmt.Printf("└-%s %s TextContent: \"%s\"\n", strings.Repeat("-", indent), p.GetType(), p.GetTextContent())
	} else if p, ok := n.(*plex_dom.CommentNode); ok {
		fmt.Printf("└-%s %s TextContent: \"%s\"\n", strings.Repeat("-", indent), p.GetType(), p.GetTextContent())
	}

	for _, v := range n.GetChildren() {
		PrintDom(v, indent+5)
	}
}
