package plex

import (
	"github.com/moznion/go-optional"
	"github.com/veandco/go-sdl2/sdl"
)

type RenderSolidColor struct {
	Color sdl.Color
	Box   sdl.FRect
}

type RenderCommand interface {
}

func resolveColor(box *LayoutBox, key ...string) optional.Option[sdl.Color] {
	if box.boxType == BoxType_AnonymousBlock {
		return nil
	}

	if box.node.IsNone() {
		return nil
	}

	style := box.node.Unwrap()

	prop := style.Lookup(key...)

	if prop.IsNone() {
		return nil
	}

	switch v := prop.Unwrap().(type) {
	case string:
		return nil
	case sdl.Color:
		return optional.Some(v)
	default:
		return nil
	}
}

func renderBackground(list *[]RenderCommand, box *LayoutBox) {
	if box.node.IsNone() {
		return
	}
	bgColor := resolveColor(box, "background-color", "background")

	if bgColor.IsNone() {
		return
	}

	*list = append(*list, RenderSolidColor{
		Color: bgColor.Unwrap(),
		Box:   box.dimensions.BorderBox(),
	})
}

func renderBorder(list *[]RenderCommand, box *LayoutBox) {
	resolvedColor := resolveColor(box, "border-color")
	if resolvedColor.IsNone() {
		return
	}

	color := resolvedColor.Unwrap()
	borderBox := box.dimensions.BorderBox()
	// Left Border
	*list = append(*list, RenderSolidColor{
		Color: color,
		Box: sdl.FRect{
			X: borderBox.X,
			Y: borderBox.Y,
			W: box.dimensions.Border.Left,
			H: borderBox.H,
		},
	})

	// Right Border
	*list = append(*list, RenderSolidColor{
		Color: color,
		Box: sdl.FRect{
			X: borderBox.X + borderBox.W - box.dimensions.Border.Right,
			Y: borderBox.Y,
			W: box.dimensions.Border.Right,
			H: borderBox.H,
		},
	})

	// Top Border
	*list = append(*list, RenderSolidColor{
		Color: color,
		Box: sdl.FRect{
			X: borderBox.X,
			Y: borderBox.Y,
			W: borderBox.W,
			H: box.dimensions.Border.Top,
		},
	})

	// Bottom Border
	*list = append(*list, RenderSolidColor{
		Color: color,
		Box: sdl.FRect{
			X: borderBox.X,
			Y: borderBox.Y + borderBox.H - box.dimensions.Border.Bottom,
			W: borderBox.W,
			H: box.dimensions.Border.Bottom,
		},
	})

}

func renderLayout(list *[]RenderCommand, layout *LayoutBox) {
	renderBackground(list, layout)
	renderBorder(list, layout)

	// Render Text HERE

	for _, child := range layout.children {
		renderLayout(list, &child)
	}
}

func buildDisplayList(layout *LayoutBox) []RenderCommand {
	cmdList := []RenderCommand{}

	renderLayout(&cmdList, layout)

	return cmdList
}

func printItem(renderer *sdl.Renderer, width float32, height float32, cmd RenderCommand) {

	if v, ok := cmd.(RenderSolidColor); ok {
		renderer.SetDrawColor(v.Color.R, v.Color.G, v.Color.B, 255)
		renderer.FillRectF(&sdl.FRect{
			Y: v.Box.Y,
			X: v.Box.X,
			W: ClampFloat32(v.Box.W, 0.0, width),
			H: ClampFloat32(v.Box.H, 0.0, height),
		})
	}

}

func Print(layout *LayoutBox, renderer *sdl.Renderer, window *sdl.Window, windowBgColor sdl.Color) {
	displayList := buildDisplayList(layout)

	w, h := window.GetSize()

	fw := float32(w)
	fh := float32(h)

	renderer.SetDrawColor(windowBgColor.R, windowBgColor.G, windowBgColor.B, 255)
	renderer.FillRect(&sdl.Rect{
		H: h,
		W: w,
	})

	for _, child := range displayList {
		printItem(renderer, fw, fh, child)
	}

	renderer.Present()
}
