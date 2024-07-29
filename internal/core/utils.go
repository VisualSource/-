package plex

import (
	plex_css "visualsource/plex/internal/css"

	"github.com/veandco/go-sdl2/sdl"
)

func MaxFloat32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func MinFloat32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func ClampFloat32(v, lo, hi float32) float32 {
	return MinFloat32(MaxFloat32(v, lo), hi)
}

func GetWindowDimentions(window *sdl.Window) Dimensions {
	w, _ := window.GetSize()

	return Dimensions{
		Content: sdl.FRect{
			X: 0,
			Y: 0,
			W: float32(w),
			H: 0,
		},
	}
}

func ParseStylesFromDocument(node Node) (StyledNode, sdl.Color) {
	cssParser := plex_css.CssParser{}

	var styletree StyledNode
	color := sdl.Color{A: 255, R: 255, G: 255, B: 255}

	if document, ok := node.(*ElementNode); ok {
		selector := plex_css.Selector{TagName: "style"}
		styleTags := document.QuerySelectorAll(&selector)
		stylesheets := []plex_css.Stylesheet{}

		for _, style := range styleTags {
			css, err := cssParser.ParseStylesheet(style.GetTextContent(), 1)
			if err == nil {
				stylesheets = append(stylesheets, css)
			}
		}

		styletree = StyleTree(node, stylesheets)

		bgColor := styletree.Lookup("background-color", "background")
		bgColor.IfSome(func(v CssValue) {
			if c, ok := v.(sdl.Color); ok {
				color = c
			}
		})
	}

	return styletree, color
}
