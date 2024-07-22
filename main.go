package main

import (
	"fmt"
	plex "visualsource/plex/internal/core"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	// load user agent css styles

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	parser := plex.CreateHtmlParser()
	cssParser := plex.CssParser{}

	dom, err := parser.Parse(`
		<html class="test">
			<head>
				<style>
					h1, h2, h3 { margin: auto; color: #cc0000; }
					div.note { margin-bottom: 20px; padding: 10px; }
					head { display: none; }
					html { background-color: #ffffff; height: 100px; }
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

	var styletree plex.StyledNode
	if doc, ok := dom.(*plex.ElementNode); ok {
		styleTags := doc.QuerySelectorAll(plex.NewSelector("style", "", []string{}))

		stylesheets := []plex.Stylesheet{}
		for _, style := range styleTags {
			css, err := cssParser.Parse(style.GetTextContent(), 1)
			if err == nil {
				stylesheets = append(stylesheets, css)
			}
		}

		styletree = plex.StyleTree(&dom, stylesheets)
	}

	layout := plex.BuildLayoutTree(styletree)

	plex.RenderLayout(renderer, layout)

	fmt.Printf("%v\n", layout)

	renderer.Present()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
			}
		}
		sdl.Delay(33)
	}
}
