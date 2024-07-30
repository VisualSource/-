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

func ParseStylesFromDocument(node Node) (StyledNode, plex_css.CssColor) {
	cssParser := plex_css.CssParser{}

	var styletree StyledNode
	color := plex_css.CSS_COLOR_KEYWORDS["white"]

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

		bgColor := styletree.props.Lookup("background-color")
		bgColor.IfSome(func(v plex_css.Declaration) {
			if c, ok := v.GetValue().(*plex_css.CssColor); ok {
				color = *c
			} else if c, ok := v.GetValue().(*plex_css.CssKeyword); ok {
				c.ResolveColor().IfSome(func(v plex_css.CssColor) {
					color = v
				})
			}
		})
	}

	return styletree, color
}
