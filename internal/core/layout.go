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

func (l *LayoutBox) layout(containing Dimensions) {
	switch l.boxType {
	case BoxType_Block:
		l.layoutBlock(containing)
	}
}

// Lay out a block-level element and its descendants.
func (l *LayoutBox) layoutBlock(containing Dimensions) {
	if l.node.IsNone() {
		return
	}
	l.calculateBlockWidth(containing)
	l.calculateBlockPosition(containing)
	l.layoutBlockChildren()
	l.calculateBlockHeight()

}

func (l *LayoutBox) calculateBlockHeight() {
	style := l.node.Unwrap()
	prop := style.GetProp("height")
	if prop.IsNone() {
		return
	}
	if i, ok := prop.Unwrap().(CssLengthValue); ok {
		l.dimensions.Content.H = i.ToPx()
	}
}

func (l *LayoutBox) layoutBlockChildren() {
	for i := 0; i < len(l.children); i++ {
		l.children[i].layout(l.dimensions)
		l.dimensions.Content.H += l.children[i].dimensions.MarginBox().H
	}
}

func (l *LayoutBox) calculateBlockPosition(containing Dimensions) {
	style := l.node.Unwrap()
	zero := optional.Some[CssValue](CssLengthValue{
		Value: 0,
		Unit:  CssUnit_PX,
	})

	l.dimensions.Margin.Top = AsCssLengthValue(style.Lookup("margin-top", "margin").Or(zero).Unwrap())
	l.dimensions.Margin.Bottom = AsCssLengthValue(style.Lookup("margin-bttom", "margin").Or(zero).Unwrap())

	l.dimensions.Border.Top = AsCssLengthValue(style.Lookup("border-top-width", "border-width").Or(zero).Unwrap())
	l.dimensions.Border.Bottom = AsCssLengthValue(style.Lookup("border-bottom-width", "border-with").Or(zero).Unwrap())

	l.dimensions.Padding.Top = AsCssLengthValue(style.Lookup("padding-top", "padding").Or(zero).Unwrap())
	l.dimensions.Padding.Bottom = AsCssLengthValue(style.Lookup("padding-bottom", "padding").Or(zero).Unwrap())

	l.dimensions.Content.X = containing.Content.X +
		l.dimensions.Margin.Left +
		l.dimensions.Border.Left +
		l.dimensions.Padding.Left

	l.dimensions.Content.Y = containing.Content.H +
		containing.Content.Y +
		l.dimensions.Margin.Top +
		l.dimensions.Border.Top +
		l.dimensions.Padding.Top
}

func (l *LayoutBox) calculateBlockWidth(containing Dimensions) {
	style := l.node.Unwrap()

	width := style.GetProp("width").Or(optional.Some[CssValue]("auto")).Unwrap()

	zero := optional.Some[CssValue](CssLengthValue{
		Value: 0,
		Unit:  CssUnit_PX,
	})

	marginLeft := style.Lookup("margin-left", "margin").Or(zero).Unwrap()
	marginRight := style.Lookup("margin-right", "margin").Or(zero).Unwrap()
	borderLeft := style.Lookup("border-left-width", "border-width").Or(zero).Unwrap()
	borderRight := style.Lookup("border-right-width", "border-width").Or(zero).Unwrap()
	paddingLeft := style.Lookup("padding-left", "padding").Or(zero).Unwrap()
	paddingRight := style.Lookup("padding-right", "padding").Or(zero).Unwrap()

	items := [7]CssValue{marginLeft, marginRight, borderLeft, borderRight, paddingLeft, paddingRight, width}
	var total float32 = 0.0

	for i := 0; i < len(items); i++ {
		if value, ok := items[i].(CssLengthValue); ok {
			total += value.ToPx()
		}
	}

	if !IsCssKeyword(width, "auto") && total > containing.Content.W {
		if IsCssKeyword(marginLeft, "auto") {
			marginLeft = CssLengthValue{
				Value: 0.0,
				Unit:  CssUnit_PX,
			}
		}
		if IsCssKeyword(marginRight, "auto") {
			marginRight = CssLengthValue{
				Value: 0.0,
				Unit:  CssUnit_PX,
			}
		}
	}

	underflow := containing.Content.W - total

	widthIsAuto := IsCssKeyword(width, "auto")
	marginLeftIsAuto := IsCssKeyword(marginLeft, "auto")
	marginRightIsAuto := IsCssKeyword(marginRight, "auto")

	if !widthIsAuto && !marginLeftIsAuto && !marginRightIsAuto {
		marginRight = CssLengthValue{
			Value: 0.0,
			Unit:  CssUnit_PX,
		}
	} else if !widthIsAuto && !marginLeftIsAuto && marginRightIsAuto {
		marginRight = CssLengthValue{
			Value: underflow,
			Unit:  CssUnit_PX,
		}
	} else if !widthIsAuto && marginLeftIsAuto && !marginRightIsAuto {
		marginLeft = CssLengthValue{
			Value: underflow,
			Unit:  CssUnit_PX,
		}
	} else if !widthIsAuto && marginLeftIsAuto && marginRightIsAuto {
		marginLeft = CssLengthValue{
			Value: underflow / 2.0,
			Unit:  CssUnit_PX,
		}
		marginRight = CssLengthValue{
			Value: underflow / 2.0,
			Unit:  CssUnit_PX,
		}
	} else if widthIsAuto {
		if marginLeftIsAuto {
			marginLeft = CssLengthValue{
				Value: 0.0,
				Unit:  CssUnit_PX,
			}
		}
		if marginRightIsAuto {
			marginRight = CssLengthValue{
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
				Value: underflow,
				Unit:  CssUnit_PX,
			}

			// right margin has been set to a css length value
			marginRight = CssLengthValue{
				Value: marginRight.(CssLengthValue).Value + underflow,
				Unit:  CssUnit_PX,
			}
		}
	}

	l.dimensions.Content.W = AsCssLengthValue(width)
	l.dimensions.Padding.Left = AsCssLengthValue(paddingLeft)
	l.dimensions.Padding.Right = AsCssLengthValue(paddingRight)
	l.dimensions.Border.Left = AsCssLengthValue(borderLeft)
	l.dimensions.Border.Right = AsCssLengthValue(borderRight)
	l.dimensions.Margin.Left = AsCssLengthValue(marginLeft)
	l.dimensions.Margin.Right = AsCssLengthValue(marginRight)
}

func (l *LayoutBox) getInlineContainer() *LayoutBox {

	switch l.boxType {
	case BoxType_Block:
		// if where are no children
		if len(l.children) <= 0 {
			l.children = append(l.children, LayoutBox{
				boxType: BoxType_AnonymousBlock,
			})
			// if there is a child but is not a AnonymousBlock
		} else if l.children[len(l.children)-1].boxType != BoxType_AnonymousBlock {
			l.children = append(l.children, LayoutBox{
				boxType: BoxType_AnonymousBlock,
			})
		}

		last := l.children[len(l.children)-1]

		return &last
	default:
		return l
	}
}

func createNewLayoutBox(boxType BoxType, dim Dimensions, node optional.Option[StyledNode]) LayoutBox {
	return LayoutBox{
		boxType:    boxType,
		dimensions: dim,
		node:       node,
		children:   []LayoutBox{},
	}
}

func LayoutTree(node StyledNode, containing Dimensions) LayoutBox {
	root := buildLayoutTree(node)
	root.layout(containing)
	return root
}

func buildLayoutTree(node StyledNode) LayoutBox {

	boxType := DisplayTypeToBoxType(node.GetDisplay())

	root := createNewLayoutBox(boxType, Dimensions{}, optional.Some(node))

	for _, child := range node.children {
		switch child.GetDisplay() {
		case DisplayType_Block:
			item := buildLayoutTree(child)
			root.children = append(root.children, item)
		case DisplayType_Inline:
			inline := root.getInlineContainer()
			item := buildLayoutTree(child)
			inline.children = append(inline.children, item)
		}
	}

	return root
}
