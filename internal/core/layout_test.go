package plex_test

import (
	"testing"
	plex "visualsource/plex/internal/core"

	"github.com/veandco/go-sdl2/sdl"
)

func TestDimensionsPadding(t *testing.T) {
	d := plex.Dimensions{
		Content: sdl.FRect{
			X: 1.0,
			Y: 1.0,
			W: 1.0,
			H: 1.0,
		},
		Padding: plex.EdgeSizes{
			Left:   1.0,
			Right:  1.0,
			Top:    1.0,
			Bottom: 1.0,
		},
		Margin: plex.EdgeSizes{},
		Border: plex.EdgeSizes{},
	}

	r := d.PaddingBox()

	if r.H != 3.0 || r.W != 3.0 || r.X != 0.0 || r.Y != 0.0 {
		t.Fatalf("invalid values: %v", r)
	}
}

func TestDimensionsBorderBox(t *testing.T) {
	d := plex.Dimensions{
		Content: sdl.FRect{
			X: 1.0,
			Y: 1.0,
			W: 1.0,
			H: 1.0,
		},
		Padding: plex.EdgeSizes{
			Left:   1.0,
			Right:  1.0,
			Top:    1.0,
			Bottom: 1.0,
		},
		Border: plex.EdgeSizes{
			Left:   1.0,
			Right:  1.0,
			Top:    1.0,
			Bottom: 1.0,
		},
		Margin: plex.EdgeSizes{},
	}

	r := d.BorderBox()

	t.Logf("%v", r)

}

func TestDimensionsMarginBox(t *testing.T) {
	d := plex.Dimensions{
		Content: sdl.FRect{
			X: 1.0,
			Y: 1.0,
			W: 1.0,
			H: 1.0,
		},
		Padding: plex.EdgeSizes{
			Left:   1.0,
			Right:  1.0,
			Top:    1.0,
			Bottom: 1.0,
		},
		Border: plex.EdgeSizes{
			Left:   1.0,
			Right:  1.0,
			Top:    1.0,
			Bottom: 1.0,
		},
		Margin: plex.EdgeSizes{
			Left:   1.0,
			Right:  1.0,
			Top:    1.0,
			Bottom: 1.0,
		},
	}

	r := d.MarginBox()

	t.Logf("%v", r)
}

func TestBuildLayout(t *testing.T) {

	rootNode := plex.CreateElementNode("html", plex.AttributeMap{"id": "root"}, []plex.Node{})
	styleNode := plex.CreateStyleNode(&rootNode, plex.PropertyMap{
		"width":   "auto",
		"height":  "auto",
		"display": "block",
	}, []plex.StyledNode{})

	layout := plex.BuildLayoutTree(styleNode)

	t.Logf("%#v", layout)
}

/*func TestCalculateBlockWidth(t *testing.T){

	rootNode := plex.CreateElementNode("html", plex.AttributeMap{"id": "root"}, []plex.Node{})
	styleNode := plex.CreateStyleNode(&rootNode, plex.PropertyMap{
		"width":   "auto",
		"height":  "auto",
		"display": "block",
	}, []plex.StyledNode{})

	layout := plex.CreateLayoutBox(plex.Dimensions{},plex.BoxType_Block,optional.Some(styleNode),[]plex.LayoutBox{})

	//layout.CalculateBlockWidth()
}*/
