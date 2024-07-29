package plex_css

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/moznion/go-optional"
)

type Stylesheet struct {
	Rules    []Rule
	AtRules  []AtRule
	TopLevel bool
	Origin   uint
}

type Declaration struct {
	Name      string
	Value     []Token
	Important bool
}

type AtRule struct {
	Name    string
	Prelude []Token
	Block   SimpleBlock
}

type Rule struct {
	Selector []Selector
	Block    []Declaration
}

type SelectorAttribute struct {
	Operation uint8
	Value     string
	Modifier  rune
}

type PesudoClass struct{}
type PesudoElement struct{}

type Selector struct {
	TagName        string
	Id             string
	Namespace      optional.Option[string]
	Attributes     map[string]SelectorAttribute
	PseudoClasses  []PesudoClass
	PseudoElements []PesudoElement
	Classes        mapset.Set[string]
}

func (s *Selector) GetSpecificity() Specificity {
	spec := Specificity{}

	// count the number of ID selectors in the selector (= A)
	if s.Id != "" {
		spec.A++
	}

	// count the number of class selectors, attributes selectors, and pseudo-classes in the selector (= B)

	if s.Classes != nil {
		for range s.Classes.Iter() {
			spec.B++
		}
	}

	spec.B += uint(len(s.Attributes))
	// TODO: count pseudo-classes

	// count the number of type selectors and pseudo-elements in the selector (= C)

	if s.TagName != "" && s.TagName != "*" {
		spec.C++
	}

	// TODO: count pseudo-elemets

	return spec
}

type Specificity struct {
	A uint
	B uint
	C uint
}

func (s *Specificity) Compare(p *Specificity) bool {
	return s.A > p.A || s.B > p.B || s.C > p.C
}
