package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	plex "visualsource/plex/internal/core"
	plex_css "visualsource/plex/internal/css"

	"github.com/gookit/goutil/dump"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	WindowTitle  = "Plex"
	WindowWidth  = 800
	WindowHeight = 600
	FrameRate    = 60
)

var runningMutex sync.Mutex

func parseArgs() string {
	var htmlDocumentPath string = ""
	var debug bool

	flag.StringVar(&htmlDocumentPath, "o", "./test.html", "Specify docuemnt to open")
	flag.BoolVar(&debug, "debug", true, "Specify debug flag")

	flag.Parse()

	if debug {
		dump.Config(func(opts *dump.Options) {
			opts.MaxDepth = 10
		})
	}

	return htmlDocumentPath
}

func run() int {
	htmlFile := parseArgs()

	//var fontCache = plex.FontCache{}
	var window *sdl.Window
	var renderer *sdl.Renderer
	var err error

	stylesheet, err := plex.LoadLocalStylesheet("./resources/useragent.css")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load stylesheet %s\n", err)
		return 1
	}

	sdl.Do(func() {
		err = ttf.Init()
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init ttf %s\n", err)
		return 1
	}
	defer ttf.Quit()

	/*sdl.Do(func() {
		err = sdl.Init(sdl.INIT_EVERYTHING)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init sdl %s\n", err)
		return 1
	}
	defer sdl.Quit()*/

	sdl.Do(func() {
		window, err = sdl.CreateWindow(
			WindowTitle,
			sdl.WINDOWPOS_UNDEFINED,
			sdl.WINDOWPOS_UNDEFINED,
			WindowWidth,
			WindowHeight,
			sdl.WINDOW_SHOWN,
		)
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}

	defer func() {
		sdl.Do(func() { window.Destroy() })
	}()

	sdl.Do(func() {
		renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	sdl.Do(func() {
		renderer.Clear()
	})
	defer func() {
		sdl.Do(func() {
			renderer.Destroy()
		})
	}()

	if htmlFile != "" {
		sdl.Do(func() {
			/*font, err := ttf.OpenFont("./resources/Ubuntu-Medium.ttf", 32)
			if err == nil {
				defer font.Close()

				surface, err := font.RenderUTF8Blended("Test", sdl.Color{R: 255, G: 255, B: 255})
				if err == nil {
					defer surface.Free()
					texture, err := renderer.CreateTextureFromSurface(surface)
					if err == nil {
						defer texture.Destroy()
						err = renderer.Copy(texture, nil, &sdl.Rect{X: 0, Y: 0, W: surface.W, H: surface.H})
						if err != nil {
							fmt.Fprintf(os.Stderr, "Failed Copy %s\n", err)
						}
						renderer.Present()
					} else {
						fmt.Fprintf(os.Stderr, "Failed CreateTextureFromSurface %s\n", err)
					}
				} else {
					fmt.Fprintf(os.Stderr, "Failed RenderUTF8Blende %s\n", err)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Failed OpenFont %s\n", err)
			}*/
			err = plex.LoadLocalHtmlDocument(htmlFile, renderer, []plex_css.Stylesheet{stylesheet})
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to font %s\n", err)
			return 1
		}
	}

	running := true
	for running {
		sdl.Do(func() {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch t := event.(type) {
				case *sdl.QuitEvent:
					runningMutex.Lock()
					running = false
					runningMutex.Unlock()
				case *sdl.KeyboardEvent:
					if t.Keysym.Sym == sdl.K_F5 && t.State == sdl.RELEASED {
						fmt.Println("Reloading html document")
						err = plex.LoadLocalHtmlDocument("./test.html", renderer, []plex_css.Stylesheet{stylesheet})
						if err != nil {
							fmt.Printf("Render Error: %s", err)
						}
					}
				}
			}
		})

		sdl.Do(func() {
			sdl.Delay(1000 / FrameRate)
		})
	}

	return 0
}

func main() {
	var exitcode int

	sdl.Main(func() {
		exitcode = run()
		ttf.Quit()
		sdl.Quit()
	})

	os.Exit(exitcode)
}
