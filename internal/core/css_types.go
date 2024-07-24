package plex

type CssLengthUnit = uint8

const (
	CssUnit_PX CssLengthUnit = 0
)

type Specificity struct {
	A int
	B int
	C int
}

// #region-start CssValue

type CssValue interface{}

type CssLengthValue struct {
	Value float32
	Unit  uint8
}

func (lv *CssLengthValue) ToPx() float32 {
	if lv.Unit == CssUnit_PX {
		return lv.Value
	}

	return 0.0
}

// #region-start CSSCore

type Declaration struct {
	name  string
	value CssValue
}

type Rule struct {
	origin      int
	selectors   []Selector
	declartions []Declaration
}

type Stylesheet struct {
	rules []Rule
}

// #region-start CssSelector

type Selector struct {
	TagName string
	Id      string
	Classes []string
}

// http://www.w3.org/TR/selectors/#specificity
func (s *Selector) Specificity() Specificity {
	a := 0
	c := 0

	if s.Id != "" {
		a++
	}

	b := len(s.Classes)

	if s.TagName != "" {
		c++
	}

	return Specificity{a, b, c}
}

func CreateNewSelector(tagName string, id string, classes []string) Selector {
	return Selector{
		TagName: tagName,
		Id:      id,
		Classes: classes,
	}
}
