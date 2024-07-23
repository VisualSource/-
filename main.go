package main

import (
	"fmt"
	plex "visualsource/plex/internal/core"

	"github.com/gookit/goutil/dump"
)

func main() {
	dump.Config(func(opts *dump.Options) {
		opts.MaxDepth = 10
	})

	parser := plex.HtmlParser{}
	//cssParser := plex.CssParser{}

	dom, err := parser.Parse(`
		<html>
			<head>
				<style>
					html { display: block; }
					head { display: none; }
					body { display: block; margin: 8px; }
					div  { display: block; margin: 4px; background-color: #00ffff; height: 100px; }
				</style>
			</head>
			<body>
				<div></div>
			</body>
		</html>
	`)
	dump.P(dom)
	if err != nil {
		fmt.Printf("DOM parser error: %s\n", err)
		return
	}

	// load user agent css styles

	/*if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
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
		<html>
			<head>
				<style>
					html { display: block; }
					head { display: none; }
					body { display: block; margin: 8px; }
					div  { display: block; margin: 4px; background-color: #00ffff; height: 100px; }
				</style>
			</head>
			<body>
				<div></div>
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

	dim := plex.GetWindowDimentions(window)
	layout := plex.LayoutTree(styletree, dim)

	plex.Print(&layout, renderer, window)

	fmt.Printf("%#v\n", layout)

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
	}*/
}
