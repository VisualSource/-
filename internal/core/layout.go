package plex

// https://github.com/moznion/go-optional
import (
	plex_css "visualsource/plex/internal/css"

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
	prop := style.props.GetProp("height")
	if prop.IsNone() {
		return
	}

	c := prop.Unwrap()

	if i, ok := c.GetValue().(*plex_css.CssDimention); ok {
		l.dimensions.Content.H = i.AsPx()
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
	zero := plex_css.CreateCssDimention(0, plex_css.CssUnit_PX)

	style.props.Lookup("margin-top", "margin")
	l.dimensions.Margin.Top = style.props.ResolveLookupAsDimention("margin-top", "margin").Or(zero).UnwrapAsPtr().AsPx()
	l.dimensions.Margin.Bottom = style.props.ResolveLookupAsDimention("margin-bottom", "margin").Or(zero).UnwrapAsPtr().AsPx()

	l.dimensions.Border.Top = style.props.ResolveLookupAsDimention("border-top-width", "border-width").Or(zero).UnwrapAsPtr().AsPx()
	l.dimensions.Border.Bottom = style.props.ResolveLookupAsDimention("border-bottom-width", "border-with").Or(zero).UnwrapAsPtr().AsPx()

	l.dimensions.Padding.Top = style.props.ResolveLookupAsDimention("padding-top", "padding").Or(zero).UnwrapAsPtr().AsPx()
	l.dimensions.Padding.Bottom = style.props.ResolveLookupAsDimention("padding-bottom", "padding").Or(zero).UnwrapAsPtr().AsPx()

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

	var width plex_css.CssValue
	widthDec := style.props.GetProp("width") //.Or(optional.Some[CssValue]("auto")).Unwrap()
	if widthDec.IsSome() {
		v := widthDec.Unwrap()
		width = v.GetValueAt(0).Or(optional.Some[plex_css.CssValue](&plex_css.CssKeyword{Value: "auto"})).Unwrap()
	} else {
		width = &plex_css.CssKeyword{Value: "auto"}
	}

	zero := optional.Some[plex_css.CssValue](&plex_css.CssDimention{Value: 0, Unit: plex_css.CssUnit_PX})

	marginLeft := style.props.ResolveLookupToCssValue("margin-left", "margin").Or(zero).Unwrap()
	marginRight := style.props.ResolveLookupToCssValue("margin-right", "margin").Or(zero).Unwrap()
	borderLeft := style.props.ResolveLookupToCssValue("border-left-width", "border-width").Unwrap()
	borderRight := style.props.ResolveLookupToCssValue("border-right-width", "border-width").Unwrap()
	paddingLeft := style.props.ResolveLookupToCssValue("padding-left", "padding").Unwrap()
	paddingRight := style.props.ResolveLookupToCssValue("padding-right", "padding").Unwrap()

	items := []*plex_css.CssValue{&marginLeft, &marginRight, &borderLeft, &borderRight, &paddingLeft, &paddingRight, &width}
	var total float32 = 0.0
	for _, i := range items {
		if plex_css.IsCssValue(*i, plex_css.TCssValue_DIMENTION) {
			v := (*i).(*plex_css.CssDimention)
			total += v.AsPx()
		}
	}

	if !plex_css.IsCssKeyword(width, "auto") && total > containing.Content.W {
		if plex_css.IsCssKeyword(marginLeft, "auto") {
			marginLeft = &plex_css.CssDimention{
				Value: 0.0,
				Unit:  plex_css.CssUnit_PX,
			}
		}
		if plex_css.IsCssKeyword(marginRight, "auto") {
			marginRight = &plex_css.CssDimention{
				Value: 0.0,
				Unit:  plex_css.CssUnit_PX,
			}
		}
	}

	underflow := containing.Content.W - total

	widthIsAuto := plex_css.IsCssKeyword(width, "auto")
	marginLeftIsAuto := plex_css.IsCssKeyword(marginLeft, "auto")
	marginRightIsAuto := plex_css.IsCssKeyword(marginRight, "auto")

	if !widthIsAuto && !marginLeftIsAuto && !marginRightIsAuto {
		marginRight = &plex_css.CssDimention{
			Value: 0.0,
			Unit:  plex_css.CssUnit_PX,
		}
	} else if !widthIsAuto && !marginLeftIsAuto && marginRightIsAuto {
		marginRight = &plex_css.CssDimention{
			Value: underflow,
			Unit:  plex_css.CssUnit_PX,
		}
	} else if !widthIsAuto && marginLeftIsAuto && !marginRightIsAuto {
		marginLeft = &plex_css.CssDimention{
			Value: underflow,
			Unit:  plex_css.CssUnit_PX,
		}
	} else if !widthIsAuto && marginLeftIsAuto && marginRightIsAuto {
		marginLeft = &plex_css.CssDimention{
			Value: underflow / 2.0,
			Unit:  plex_css.CssUnit_PX,
		}
		marginRight = &plex_css.CssDimention{
			Value: underflow / 2.0,
			Unit:  plex_css.CssUnit_PX,
		}
	} else if widthIsAuto {
		if marginLeftIsAuto {
			marginLeft = &plex_css.CssDimention{
				Value: 0.0,
				Unit:  plex_css.CssUnit_PX,
			}
		}
		if marginRightIsAuto {
			marginRight = &plex_css.CssDimention{
				Value: 0.0,
				Unit:  plex_css.CssUnit_PX,
			}
		}

		if underflow >= 0.0 {
			width = &plex_css.CssDimention{
				Value: underflow,
				Unit:  plex_css.CssUnit_PX,
			}
		} else {
			width = &plex_css.CssDimention{
				Value: underflow,
				Unit:  plex_css.CssUnit_PX,
			}
			a := marginRight.(*plex_css.CssDimention)

			// right margin has been set to a css length value
			marginRight = &plex_css.CssDimention{
				Value: a.AsPx() + underflow,
				Unit:  plex_css.CssUnit_PX,
			}
		}
	}

	l.dimensions.Content.W = plex_css.ResolveDimentionToPXFloat(width)
	l.dimensions.Padding.Left = plex_css.ResolveDimentionToPXFloat(paddingLeft)
	l.dimensions.Padding.Right = plex_css.ResolveDimentionToPXFloat(paddingRight)
	l.dimensions.Border.Left = plex_css.ResolveDimentionToPXFloat(borderLeft)
	l.dimensions.Border.Right = plex_css.ResolveDimentionToPXFloat(borderRight)
	l.dimensions.Margin.Left = plex_css.ResolveDimentionToPXFloat(marginLeft)
	l.dimensions.Margin.Right = plex_css.ResolveDimentionToPXFloat(marginRight)
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
