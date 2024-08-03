package plex

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type FontCache map[string]*ttf.Font

func (c FontCache) RenderText(fontFamily string, fontSize int) {

}

// https://stackoverflow.com/questions/22886500/how-to-render-text-in-sdl2
func (c FontCache) RenderTextWraped(renderer *sdl.Renderer, target *sdl.FRect, text string, color sdl.Color, fontFamily string, fontSize int, containerWidth int) error {
	font, ok := c[fontFamily]

	if !ok {
		return fmt.Errorf("no font family")
	}

	surface, err := font.RenderUTF8Blended(text, color)
	if err != nil {
		return err
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}
	defer texture.Destroy()

	source := sdl.Rect{}
	if err = renderer.CopyF(texture, &source, target); err != nil {
		return err
	}

	return nil
}

func (c FontCache) LoadLocalFont(filePath string, familyName string, fontSize int) error {
	font, err := ttf.OpenFont(filePath, fontSize)
	if err != nil {
		return err
	}
	defer font.Close()

	c[familyName] = font

	return nil
}
