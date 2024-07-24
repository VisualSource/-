package main

import (
	"fmt"
	plex "visualsource/plex/internal/core"

	"github.com/gookit/goutil/dump"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	dump.Config(func(opts *dump.Options) {
		opts.MaxDepth = 10
	})

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Plex", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	parser := plex.HtmlParser{}

	dom, err := parser.Parse(`
		<html>
			<head>
				<style>
					html { display: block; }
					head { display: none; }
					body { display: block; margin: 8px; }
					div  { display: block; background-color: #00ffff; height: 100px; }
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

	dim := plex.GetWindowDimentions(window)
	style := plex.ParseStylesFromDocument(dom)
	layout := plex.LayoutTree(style, dim)
	dump.P(layout)
	plex.Print(&layout, renderer, window)

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
