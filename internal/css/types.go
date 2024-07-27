package plex_css

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
	Block   SimpleBlock
}
