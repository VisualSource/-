package plex

import (
	"fmt"
	"strings"
)

func PrintDom(n Node, indent int) {

	if p, ok := n.(*ElementNode); ok {
		fmt.Printf("└-%s %s TAG:%s\n", strings.Repeat("-", indent), p.GetType(), p.GetTagName())
	} else if p, ok := n.(*TextNode); ok {
		fmt.Printf("└-%s %s TextContent: \"%s\"\n", strings.Repeat("-", indent), p.GetType(), p.GetTextContent())
	} else if p, ok := n.(*CommentNode); ok {
		fmt.Printf("└-%s %s TextContent: \"%s\"\n", strings.Repeat("-", indent), p.GetType(), p.GetTextContent())
	}

	for _, v := range n.GetChildren() {
		PrintDom(v, indent+5)
	}
}
