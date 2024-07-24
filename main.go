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

	err = plex.LoadLocalHtmlDocument("./test.html", renderer)

	if err != nil {
		fmt.Printf("Render Error: %s", err)
	}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
			case *sdl.KeyboardEvent:
				if t.Keysym.Sym == sdl.K_F5 && t.State == sdl.RELEASED {
					fmt.Println("Reloading html Document")
					err = plex.LoadLocalHtmlDocument("./test.html", renderer)
					if err != nil {
						fmt.Printf("Render Error: %s", err)
					}
				}
			}
		}
		sdl.Delay(33)
	}
}
