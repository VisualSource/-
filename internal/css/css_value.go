package plex_css

import (
	"github.com/gookit/goutil/mathutil"
	"github.com/moznion/go-optional"
)

type CssUnit = uint
type CssValueType = uint

const (
	TCssValue_COLOR CssValueType = iota
	TCssValue_KEYWORD
	TCssValue_DIMENTION
	TCssValue_FUNCTION
	TCssValue_EXPRESSION
)

const (
	CssUnit_NO_UNIT CssUnit = iota
	CssUnit_PRESENT
	CssUnit_EM
	CssUnit_EX
	CssUnit_CAP
	CssUnit_CH
	CssUnit_IC
	CssUnit_REM
	CssUnit_LH
	CssUnit_RLH
	CssUnit_VW
	CssUnit_VH
	CssUnit_VI
	CssUnit_VB
	CssUnit_VMIN
	CssUnit_VMAX
	CssUnit_PX
)

func strToUnit(value string) CssUnit {
	switch value {
	case "em":
		return CssUnit_EM
	case "ex":
		return CssUnit_EX
	case "cap":
		return CssUnit_CAP
	case "ch":
		return CssUnit_CH
	case "ic":
		return CssUnit_IC
	case "rem":
		return CssUnit_REM
	case "lh":
		return CssUnit_LH
	case "rlh":
		return CssUnit_RLH
	case "vw":
		return CssUnit_VW
	case "vh":
		return CssUnit_VH
	case "vi":
		return CssUnit_VI
	case "px":
		return CssUnit_PX
	default:
		return CssUnit_NO_UNIT
	}
}

// https://stackoverflow.com/questions/54197913/parse-hex-string-to-image-color
func parseHexValue(s string) CssValue {
	color := CssColor{}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		return 0
	}

	switch len(s) {
	case 8: // Full RGBA
		color.R = int(hexToByte(s[0])<<4 + hexToByte(s[1]))
		color.G = int(hexToByte(s[2])<<4 + hexToByte(s[3]))
		color.B = int(hexToByte(s[4])<<4 + hexToByte(s[5]))
		color.A = int(hexToByte(s[6])<<4 + hexToByte(s[7]))
		return &color
	case 6: // Full RGB
		color.R = int(hexToByte(s[0])<<4 + hexToByte(s[1]))
		color.G = int(hexToByte(s[2])<<4 + hexToByte(s[3]))
		color.B = int(hexToByte(s[4])<<4 + hexToByte(s[5]))
		return &color
	default:
		return nil
	}
}

func parseValue(tokens *[]Token, pos *int) CssValue {
	token := (*tokens)[*pos]

	switch token.GetId() {
	case Token_Whitespace:
		(*pos)++
		return nil
	case Token_Hash:
		v := token.(*FlagedStringToken)
		(*pos)++
		return parseHexValue(v.Value)
	case Token_Ident:
		v := token.(*StringToken)
		(*pos)++
		return &CssKeyword{Value: v.Value}
	case Token_Dimension:
		v := token.(*NumberToken)
		(*pos)++
		return &CssDimention{
			Value: v.Value,
			Unit:  strToUnit(v.Unit),
		}
	case Token_Number:
		v := token.(*NumberToken)
		(*pos)++
		return &CssDimention{
			Value: v.Value,
			Unit:  CssUnit_NO_UNIT,
		}
	case Token_Percentage:
		v := token.(*NumberToken)
		(*pos)++
		return &CssDimention{
			Value: v.Value,
			Unit:  CssUnit_PRESENT,
		}
	case TFunction:
		f := token.(*FunctionBlock)
		(*pos)++

		args := ParseCssValue(f.Args)

		return &CssFunction{
			Name: f.Name,
			Args: args,
		}
	default:
		return nil
	}
}

func ParseCssValue(tokens []Token) []CssValue {
	values := []CssValue{}
	len := len(tokens)
	pos := 0

	for {
		if pos >= len {
			return values
		}
		switch tokens[pos].GetId() {
		case Token_Whitespace, Token_Ident, TFunction, Token_Hash:
			value := parseValue(&tokens, &pos)
			if value != nil {
				values = append(values, value)
			}
		case Token_Dimension, Token_Number, Token_Percentage:
			value := parseValue(&tokens, &pos)
			if value == nil {
				continue
			}

			for pos < len && tokens[pos].GetId() == Token_Whitespace {
				pos++
			}

			if pos < len && tokens[pos].GetId() == Token_Delim {
				v := tokens[pos].(*RuneToken)
				if v.Value == '+' || v.Value == '-' || v.Value == '*' || v.Value == '/' {
					pos++
					for pos < len && tokens[pos].GetId() == Token_Whitespace {
						pos++
					}
					right := parseValue(&tokens, &pos)
					if right == nil {
						continue
					}

					values = append(values, &CssExpression{
						Left:  value,
						Right: right,
						Op:    uint8(v.Value),
					})
					continue
				}
			}
			values = append(values, value)
		}
	}
}

type CssValue interface {
	GetType() CssValueType
}

func IsCssValue(v CssValue, t CssValueType) bool {
	if v == nil {
		return false
	}
	return v.GetType() == t
}

func IsCssKeyword(v CssValue, keyword string) bool {
	if !IsCssValue(v, TCssValue_KEYWORD) {
		return true
	}
	if i, ok := v.(*CssKeyword); ok {
		return i.Value == keyword
	}
	return false
}

type CssPropertyMap map[string]Declaration

func (c CssPropertyMap) GetProp(key string) optional.Option[Declaration] {
	value, ok := c[key]
	if !ok {
		return nil
	}
	return optional.Some(value)
}

func (c CssPropertyMap) Lookup(keys ...string) optional.Option[Declaration] {
	for _, key := range keys {
		value := c.GetProp(key)
		if value != nil {
			return value
		}
	}
	return nil
}

func (c CssPropertyMap) ResolveLookupToCssValue(keys ...string) optional.Option[CssValue] {
	v := c.Lookup(keys...)

	if v.IsNone() {
		return nil
	}
	r := v.Unwrap()
	return optional.Some(r.GetValue())
}

func (c CssPropertyMap) ResolveLookupAsDimention(keys ...string) optional.Option[CssDimention] {
	result := c.Lookup(keys...)
	if result.IsNone() {
		return nil
	}

	item := result.Unwrap()

	if a, ok := item.GetValue().(*CssDimention); ok {
		return optional.Some(*a)
	}
	return nil
}

type CssExpression struct {
	Left  CssValue
	Op    uint8
	Right CssValue
}

func (c *CssExpression) GetType() CssValueType {
	return TCssValue_EXPRESSION
}

type CssColor struct {
	G, R, B, A int
}

/*
Source: https://www.w3.org/TR/css-color-4/#rgb-to-hsl
*/
func (c *CssColor) RgbToHsl() (int, int, int) {
	valueMax := max(c.R, c.G, c.B)
	valueMin := min(c.R, c.G, c.B)

	var hue int
	sat := 0
	light := (valueMin + valueMax) / 2

	d := valueMax - valueMin

	if d != 0 {
		if light == 0 || light == 1 {
			sat = 0
		} else {
			sat = (valueMax - light) / min(light, 1-light)
		}
		switch valueMax {
		case c.R:
			var m int
			if c.G < c.B {
				m = 6
			} else {
				m = 0
			}

			hue = (c.G-c.B)/d + m
		case c.G:
			hue = (c.B-c.R)/d + 2
		case c.B:
			hue = (c.R-c.G)/d + 4
		}

		hue = hue * 60
	}

	if sat < 0 {
		hue += 180
		sat = mathutil.Abs(sat)
	}

	if hue >= 360 {
		hue -= 360
	}

	return hue, sat * 100, light * 100
}

func (c *CssColor) GetType() CssValueType {
	return TCssValue_COLOR
}

type CssKeyword struct {
	Value string
}

func (c *CssKeyword) ResolveColor() optional.Option[CssColor] {
	v, ok := CSS_COLOR_KEYWORDS[c.Value]

	if !ok {
		return nil
	}

	return optional.Some(v)
}

func (c *CssKeyword) GetType() CssValueType {
	return TCssValue_KEYWORD
}

func ResolveCssValueToColor(c CssValue) optional.Option[CssColor] {
	if i, ok := c.(*CssKeyword); ok {
		color := i.ResolveColor()
		if color != nil {
			return color
		}
	}
	if d, ok := c.(*CssColor); ok {
		return optional.Some(*d)
	}

	return nil
}

type CssDimention struct {
	Value float32
	Unit  CssUnit
}

func ResolveDimentionToPXFloat(c CssValue) float32 {
	if v, ok := c.(*CssDimention); ok {
		return v.AsPx()
	}

	return 0.0
}

func (c *CssDimention) AsPx() float32 {
	switch c.Unit {
	case CssUnit_PX:
		return c.Value
	default:
		return 0.0
	}
}

func (c *CssDimention) GetType() CssValueType {
	return TCssValue_DIMENTION
}

func CreateCssDimention(value float32, t CssUnit) optional.Option[CssDimention] {
	return optional.Some(CssDimention{Value: value, Unit: t})
}

type CssFunction struct {
	Name string
	Args []CssValue
}

func (c *CssFunction) GetType() CssValueType {
	return TCssValue_FUNCTION
}
