package plex

import "github.com/veandco/go-sdl2/sdl"

type BoxType = uint8
type DisplayType = uint8

// #region-start CONTS

const (
	BoxType_Block          BoxType = 0
	BoxType_Inline         BoxType = 1
	BoxType_AnonymousBlock BoxType = 2
)

const (
	DisplayType_Inline DisplayType = 0
	DisplayType_Block  DisplayType = 1
	DisplayType_None   DisplayType = 2
)

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

func DisplayTypeToBoxType(display DisplayType) BoxType {
	switch display {
	case DisplayType_Block:
		return BoxType_Block
	case DisplayType_Inline:
		return BoxType_Inline
	default:
		panic("root node has display none")
	}
}
