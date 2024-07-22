package plex

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

func RenderLayout(surface *sdl.Renderer, layout LayoutBox) {
	if layout.node.IsSome() {
		style := layout.node.Unwrap()
		if layout.boxType != BoxType_AnonymousBlock {
			if background := style.Lookup("background-color", "background"); background.IsSome() {
				value := background.Unwrap()
				if color, ok := value.(sdl.Color); ok {
					borderBox := layout.dimensions.BorderBox()
					fmt.Printf("Border Box: %v\n", borderBox)
					surface.SetDrawColor(color.R, color.G, color.B, color.A)
					surface.FillRectF(&sdl.FRect{X: 0.0, Y: 0.0, W: 100.0, H: 100.0})
				}
			}
		}

		if borderColor := style.Lookup("border-color"); borderColor.IsSome() {
			color := borderColor.Unwrap()
			borderBox := layout.dimensions.BorderBox()
			if color, ok := color.(sdl.Color); ok {
				surface.SetDrawColor(color.R, color.B, color.B, color.A)
				// Left border
				err := surface.DrawRectF(&sdl.FRect{
					X: borderBox.X,
					Y: borderBox.Y,
					W: layout.dimensions.Border.Left,
					H: borderBox.H,
				})
				if err != nil {
					fmt.Printf("%s", err)
				}

				// Right border
				err = surface.DrawRectF(&sdl.FRect{
					X: borderBox.X + borderBox.W - layout.dimensions.Border.Right,
					Y: borderBox.Y,
					W: layout.dimensions.Border.Right,
					H: borderBox.H,
				})
				if err != nil {
					fmt.Printf("%s", err)
				}
				// Top border
				err = surface.DrawRectF(&sdl.FRect{
					X: borderBox.X,
					Y: borderBox.Y,
					W: borderBox.W,
					H: layout.dimensions.Border.Top,
				})
				if err != nil {
					fmt.Printf("%s", err)
				}
				// Bottom border
				err = surface.DrawRectF(&sdl.FRect{
					X: borderBox.X,
					Y: borderBox.Y + borderBox.H - layout.dimensions.Border.Bottom,
					W: borderBox.W,
					H: layout.dimensions.Border.Bottom,
				})

				if err != nil {
					fmt.Printf("%s", err)
				}
			}
		}
	}

	for _, child := range layout.children {
		RenderLayout(surface, child)
	}
}
