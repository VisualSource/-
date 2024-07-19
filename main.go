package main

import (
	"fmt"
	plex_parser "visualsource/plex/internal/parser"
	plex_utils "visualsource/plex/internal/utils"
)

func main() {
	parser := plex_parser.CreateParser()

	dom, err := parser.Parse(`
		<html class="test">
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

	plex_utils.PrintDom(dom, 0)
}
