package plex

// https://github.com/moznion/go-optional
import (
	"github.com/moznion/go-optional"
)

// #region-start Types

// #region-start LayoutBox
type LayoutBox struct {
	dimensions Dimensions
	boxType    BoxType
	node       optional.Option[StyledNode]
	children   []LayoutBox
}

func CreateLayoutBox(dimensions Dimensions, boxType BoxType, node optional.Option[StyledNode], children []LayoutBox) LayoutBox {
	return LayoutBox{
		dimensions: dimensions,
		boxType:    boxType,
		node:       node,
		children:   children,
	}
}

func (l *LayoutBox) GetInlineContainer() *LayoutBox {
	switch l.boxType {
	case BoxType_Inline:
		fallthrough
	case BoxType_AnonymousBlock:
		return l
	case BoxType_Block:

		// create anony block if not exists
		// else create return block

		if len(l.children) < 1 || l.children[len(l.children)-1].boxType != BoxType_AnonymousBlock {
			l.children = append(l.children, LayoutBox{
				boxType: BoxType_AnonymousBlock,
			})
		}

		last := l.children[len(l.children)-1]

		return &last
	}

	panic("Should not be here")
}

func (l *LayoutBox) GetLayout(containing Dimensions) {
	switch l.boxType {
	case BoxType_Block:
		l.layoutBlock(containing)
	}
}

func (l *LayoutBox) layoutBlock(containing Dimensions) {
	l.CalculateBlockWidth(containing)
	l.CalculateBlockPosition(containing)
	l.LayoutBlockChildren()
	l.CalculateBlockHeight()
}

func (l *LayoutBox) CalculateBlockWidth(containing Dimensions) {
	if l.node.IsNone() {
		return
	}

	style := l.node.Unwrap()

	width := style.GetProp("width").Or(optional.Some[CssValue]("auto")).Unwrap()

	var zero optional.Option[CssValue] = optional.Some[CssValue](CssLengthValue{
		Value: 0.0,
		Unit:  CssUnit_PX,
	})

	margin_left := style.Lookup("margin-left", "margin").Or(zero).Unwrap()
	margin_right := style.Lookup("margin-right", "margin").Or(zero).Unwrap()

	border_left := style.Lookup("border-left-width", "border-width").Or(zero).Unwrap()
	border_right := style.Lookup("border-right-width", "border-width").Or(zero).Unwrap()

	padding_left := style.Lookup("padding-left", "padding").Or(zero).Unwrap()
	padding_right := style.Lookup("padding-right", "padding").Or(zero).Unwrap()

	totals := [7]CssValue{margin_left, margin_right, border_left, border_right, padding_left, padding_right, width}
	var total float32 = 0.0

	for _, i := range totals {
		if length, ok := i.(CssLengthValue); ok {
			total += length.ToPx()
		}
	}

	var isWidthAuto bool = false
	if item, ok := width.(string); ok && item == "auto" {
		isWidthAuto = true
	}

	if !isWidthAuto && total > containing.Content.W {
		if ml, ok := margin_left.(string); ok && ml == "auto" {
			margin_left = CssLengthValue{
				Value: 0.0,
				Unit:  CssUnit_PX,
			}
		}
		if mr, ok := margin_right.(string); ok && mr == "auto" {
			margin_right = CssLengthValue{
				Value: 0.0,
				Unit:  CssUnit_PX,
			}
		}
	}

	underflow := containing.Content.W - total

	var isMarginLeftAuto bool = false
	if ml, ok := margin_left.(string); ok && ml == "auto" {
		isMarginLeftAuto = true
	}
	var isMarginRightAuto bool = false
	if mr, ok := margin_right.(string); ok && mr == "auto" {
		isMarginRightAuto = true
	}

	if !isWidthAuto && !isMarginLeftAuto && !isMarginRightAuto {
		if mr, ok := margin_right.(*CssLengthValue); ok {
			mr.Value += underflow
		} else {
			margin_right = CssLengthValue{
				Value: underflow,
				Unit:  CssUnit_PX,
			}
		}
	} else if !isWidthAuto && !isMarginLeftAuto && isMarginRightAuto {
		margin_right = CssLengthValue{
			Value: underflow,
			Unit:  CssUnit_PX,
		}
	} else if !isWidthAuto && isMarginLeftAuto && !isMarginRightAuto {
		margin_left = CssLengthValue{
			Value: underflow,
			Unit:  CssUnit_PX,
		}
	} else if !isWidthAuto && isMarginLeftAuto && isMarginRightAuto {
		margin_left = CssLengthValue{
			Value: underflow / 2.0,
			Unit:  CssUnit_PX,
		}
		margin_right = CssLengthValue{
			Value: underflow / 2.0,
			Unit:  CssUnit_PX,
		}
	} else if isWidthAuto {
		if isMarginLeftAuto {
			margin_left = CssLengthValue{
				Value: 0.0,
				Unit:  CssUnit_PX,
			}
		}
		if isMarginRightAuto {
			margin_right = CssLengthValue{
				Value: 0.0,
				Unit:  CssUnit_PX,
			}
		}

		if underflow >= 0.0 {
			width = CssLengthValue{
				Value: underflow,
				Unit:  CssUnit_PX,
			}
		} else {
			width = CssLengthValue{
				Value: 0.0,
				Unit:  CssUnit_PX,
			}
			if mr, ok := margin_right.(*CssLengthValue); ok {
				mr.Value += underflow
			} else {
				margin_right = CssLengthValue{
					Value: underflow,
					Unit:  CssUnit_PX,
				}
			}
		}
	}

	if w, ok := width.(CssLengthValue); ok {
		l.dimensions.Content.W = w.ToPx()
	}

	if pl, ok := padding_left.(CssLengthValue); ok {
		l.dimensions.Padding.Left = pl.ToPx()
	}
	if pr, ok := padding_right.(CssLengthValue); ok {
		l.dimensions.Padding.Right = pr.ToPx()
	}

	if pl, ok := border_left.(CssLengthValue); ok {
		l.dimensions.Border.Left = pl.ToPx()
	}
	if pr, ok := border_right.(CssLengthValue); ok {
		l.dimensions.Border.Right = pr.ToPx()
	}

	if pl, ok := margin_left.(CssLengthValue); ok {
		l.dimensions.Margin.Left = pl.ToPx()
	}
	if pr, ok := margin_right.(CssLengthValue); ok {
		l.dimensions.Margin.Right = pr.ToPx()
	}

}
func (l *LayoutBox) CalculateBlockPosition(containing Dimensions) {
	if l.node.IsNone() {
		return
	}

	style := l.node.Unwrap()
	zero := optional.Some(CssLengthValue{
		Value: 0.0,
		Unit:  CssUnit_PX,
	})

	l.dimensions.Margin.Top = style.LookupCssLength("margin-top", "margin").Or(zero).Unwrap().Value
	l.dimensions.Margin.Bottom = style.LookupCssLength("margin-top", "margin").Or(zero).Unwrap().Value

	l.dimensions.Border.Top = style.LookupCssLength("border-top-width", "border-width").Or(zero).Unwrap().Value
	l.dimensions.Border.Bottom = style.LookupCssLength("border-top-width", "border-width").Or(zero).Unwrap().Value

	l.dimensions.Padding.Top = style.LookupCssLength("padding-top", "padding").Or(zero).Unwrap().Value
	l.dimensions.Padding.Bottom = style.LookupCssLength("padding-top", "padding").Or(zero).Unwrap().Value

	l.dimensions.Content.X = containing.Content.X + l.dimensions.Margin.Left + l.dimensions.Border.Left + l.dimensions.Padding.Left
	l.dimensions.Content.Y = containing.Content.H + containing.Content.Y + l.dimensions.Margin.Top + l.dimensions.Border.Top + l.dimensions.Padding.Top
}
func (l *LayoutBox) LayoutBlockChildren() {
	for _, child := range l.children {
		child.layoutBlock(l.dimensions)

		l.dimensions.Content.H += child.dimensions.MarginBox().H
	}
}
func (l *LayoutBox) CalculateBlockHeight() {
	if l.node.IsNone() {
		return
	}
	style := l.node.Unwrap()
	height := style.GetPropAsLength("height")

	if height.IsNone() {
		return
	}

	l.dimensions.Content.H = height.Unwrap().Value
}

// #region start functions

func getContainerDisplay(display DisplayType) BoxType {
	switch display {
	case DisplayType_Block:
		return BoxType_Block
	case DisplayType_Inline:
		return BoxType_Inline
	default:
		panic("root node has display: none")
	}
}

func BuildLayoutTree(node StyledNode) LayoutBox {

	root := LayoutBox{
		boxType: getContainerDisplay(node.GetDisplay()),
		node:    optional.Some(node),
	}

	for _, child := range node.children {
		switch child.GetDisplay() {
		case DisplayType_Block:
			root.children = append(root.children, BuildLayoutTree(child))
		case DisplayType_Inline:
			inline := root.GetInlineContainer()
			inline.children = append(inline.children, BuildLayoutTree(child))
		}
	}

	return root
}

func LayoutTree(node StyledNode, dim Dimensions) LayoutBox {
	root := BuildLayoutTree(node)

	root.layoutBlock(dim)

	return root
}
