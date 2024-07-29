package plex_css

type CssUnit = uint
type CssValueType = uint

const (
	TCssValue_COLOR     = 0
	TCssValue_KEYWORD   = 1
	TCssValue_DIMENTION = 2
	TCssValue_FUNCTION  = 3
)

const (
	CssUnit_EM   CssUnit = 0
	CssUnit_EX   CssUnit = 1
	CssUnit_CAP  CssUnit = 2
	CssUnit_CH   CssUnit = 3
	CssUnit_IC   CssUnit = 4
	CssUnit_REM  CssUnit = 5
	CssUnit_LH   CssUnit = 6
	CssUnit_RLH  CssUnit = 7
	CssUnit_VW   CssUnit = 8
	CssUnit_VH   CssUnit = 9
	CssUnit_VI   CssUnit = 10
	CssUnit_VB   CssUnit = 11
	CssUnit_VMIN CssUnit = 12
	CssUnit_VMAX CssUnit = 13
	CssUnit_PX   CssUnit = 14
)

type CssValue interface {
	GetType() CssValueType
}
type CssPropertyMap = map[string]CssValue

type CssColor struct {
	G, R, B, A uint
}

func (c *CssColor) GetType() CssValueType {
	return TCssValue_COLOR
}

type CssKeyword struct{}

func (c *CssKeyword) GetType() CssValueType {
	return TCssValue_KEYWORD
}

type CssDimention struct{}

func (c *CssDimention) GetType() CssValueType {
	return TCssValue_DIMENTION
}

type CssFunction struct{}

func (c *CssFunction) GetType() CssValueType {
	return TCssValue_FUNCTION
}
