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
func (c FontCache) RenderTextWraped(renderer *sdl.Renderer, text string, color sdl.Color, fontFamily string, fontSize int, containerWidth int) error {
	font, ok := c[fontFamily]

	if !ok {
		return fmt.Errorf("No font family")
	}

	surface, err := font.RenderUTF8BlendedWrapped(text, color, containerWidth)
	if err != nil {
		return err
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}

	out := sdl.FRect{}
	if err = renderer.CopyF(texture, nil, &out); err != nil {
		return err
	}

	return nil
}

func (c FontCache) LoadLocalFont(filePath string, fontSize int) {
	/*font, err := ttf.OpenFont(filePath, fontSize)
	if err != nil {
		return nil, err
	}
	defer font.Close()

	f[filePath] = font

	surface, err := font.RenderUTF8BlendedWrapped("Test", sdl.Color{R: 255, G: 255, B: 255, A: 255}, 800)
	if err != nil {
		return nil, err
	}
	defer surface.Free()

	return font, nil*/
}
