package plex

import (
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
	h, w := window.GetSize()

	return Dimensions{
		Content: sdl.FRect{
			X: 0,
			Y: 0,
			W: float32(w),
			H: float32(h),
		},
	}
}
