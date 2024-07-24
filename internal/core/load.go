package plex

import (
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

func LoadLocalHtmlDocument(filepath string, renderer *sdl.Renderer) error {
	window, err := renderer.GetWindow()
	if err != nil {
		return err
	}

	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	doc := string(content)
	parser := HtmlParser{}

	dom, err := parser.Parse(doc)
	if err != nil {
		return err
	}

	dim := GetWindowDimentions(window)
	style, bgColor := ParseStylesFromDocument(dom)
	layout := LayoutTree(style, dim)

	Print(&layout, renderer, window, bgColor)

	return nil
}
