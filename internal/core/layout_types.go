package plex

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type BoxType uint8
type DisplayType uint8

// #region-start CONTS

const (
	BoxType_Block BoxType = iota
	BoxType_Inline
	BoxType_AnonymousBlock
)

const (
	DisplayType_Inline DisplayType = iota
	DisplayType_Block
	DisplayType_None
)

func (b DisplayType) ToBoxType() (BoxType, error) {
	switch b {
	case DisplayType_Block:
		return BoxType_Block, nil
	case DisplayType_Inline:
		return BoxType_Inline, nil
	default:
		return 255, fmt.Errorf("can not convert for display type of 'none'")
	}
}

type EdgeSizes struct {
	Left   float32
	Right  float32
	Top    float32
	Bottom float32
}

func CreateNewEdge(left, right, top, bottom float32) EdgeSizes {
	return EdgeSizes{
		Left:   left,
		Right:  right,
		Top:    top,
		Bottom: bottom,
	}
}

// #region-start Dimensions

type Dimensions struct {
	Content sdl.FRect
	Padding EdgeSizes
	Border  EdgeSizes
	Margin  EdgeSizes
}

// The area covered by the content area plus its padding.
func (d *Dimensions) PaddingBox() sdl.FRect {
	return expandedBy(d.Content, d.Padding)
}

// The area covered by the content area plus padding and borders.
func (d *Dimensions) BorderBox() sdl.FRect {
	return expandedBy(d.PaddingBox(), d.Border)
}

// The area covered by the content area plus padding, borders, and margin.
func (d *Dimensions) MarginBox() sdl.FRect {
	return expandedBy(d.BorderBox(), d.Margin)
}

func expandedBy(rect sdl.FRect, edge EdgeSizes) sdl.FRect {
	return sdl.FRect{
		X: rect.X - edge.Left,
		Y: rect.Y - edge.Top,
		W: rect.W + edge.Left + edge.Right,
		H: rect.H + edge.Top + edge.Bottom,
	}
}
