package main

import (
	"fmt"
	plex "visualsource/plex/internal/core"
)

func main() {
	parser := plex.CreateHtmlParser()
	cssParser := plex.CssParser{}

	dom, err := parser.Parse(`
		<html class="test">
			<head>
				<style>
					h1, h2, h3 { margin: auto; color: #cc0000; }
					div.note { margin-bottom: 20px; padding: 10px; }
					#answer { display: none; }
				</style>
			</head>
			<body>
				<h1>Title</h1>
				<div id="main" class="test">
				    <!-- great text -->
					<p>Hello <em>world</em>!</p>
				</div>
			</body>
		</html>
	`)
	if err != nil {
		fmt.Printf("DOM parser error: %s\n", err)
		return
	}

	plex.PrintDom(dom, 0)

	if doc, ok := dom.(*plex.ElementNode); ok {
		styleTags := doc.QuerySelectorAll(plex.NewSelector("style", "", []string{}))

		stylesheets := []plex.Stylesheet{}
		for _, style := range styleTags {
			css, err := cssParser.Parse(style.GetTextContent(), 1)
			if err == nil {
				stylesheets = append(stylesheets, css)
			}
		}

		styletree := plex.StyleTree(&dom, stylesheets)

		fmt.Printf("%v+2", styletree)
	}

}
