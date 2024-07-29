package plex_css

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/moznion/go-optional"
)

type Stylesheet struct {
	Rules    []Rule
	AtRules  []AtRule
	TopLevel bool
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
	Prelude []Token
	Block   []Declaration
}

type SelectorAttribute struct {
	Operation uint8
	Value     string
}
type Selector struct {
	TagName        string
	Namespace      optional.Option[string]
	Attributes     map[string]SelectorAttribute
	PseudoClasses  []string
	PseudoElements []string
	Classes        mapset.Set[string]
}
