package plex

import (
	"os"
	plex_css "visualsource/plex/internal/css"

	"github.com/gookit/goutil/dump"
	"github.com/veandco/go-sdl2/sdl"
)

func LoadLocalStylesheet(filepath string) (plex_css.Stylesheet, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return plex_css.Stylesheet{}, err
	}

	parser := plex_css.CssParser{}

	result, err := parser.ParseStylesheet(string(content), 1)

	if err != nil {
		return plex_css.Stylesheet{}, err
	}

	return result, nil
}

func LoadLocalHtmlDocument(filepath string, renderer *sdl.Renderer, stylesheets []plex_css.Stylesheet) error {
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
	style, bgColor := ParseStylesFromDocument(dom, stylesheets)
	layout := LayoutTree(style, dim)

	dump.P(layout)

	Print(&layout, renderer, window, bgColor)

	return nil
}
