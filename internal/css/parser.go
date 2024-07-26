package plex_css

import "fmt"

type CssParser struct {
	pos        int
	len        int
	input      []Token
	isTopLevel bool
}

func (p *CssParser) Parse(value string) (Stylesheet, error) {
	p.pos = 0
	tokenizer := Tokenizer{}

	result, err := tokenizer.Parse(value)
	p.len = len(result)
	p.input = result
	if err != nil {
		return Stylesheet{}, err
	}

	// parse stylesheet

	return Stylesheet{}, nil
}

func (p *CssParser) IfCurrentIs(t TokenType) bool {
	if p.eof() {
		return false
	}
	return p.input[p.pos].GetId() == t
}

func (p *CssParser) ParseRulesList() ([]string, error) {

	rules := []string{}

L:
	for !p.eof() {
		switch {
		case p.IfCurrentIs(Token_Whitespace):
			p.pos++
		case p.IfCurrentIs(Token_EOF):
			p.pos++
			break L
		case p.IfCurrentIs(Token_CDO) || p.IfCurrentIs(Token_CDC):
			p.pos++
			if p.isTopLevel {
				result, err := p.ConsumeQualifiedRule()
				if err != nil {
					return nil, err
				}
				rules = append(rules, result)
			}
		case p.IfCurrentIs(Token_At_Keyword):
			rules = append(rules, p.ConsumeAtRule())
		default:
			result, err := p.ConsumeQualifiedRule()
			if err != nil {
				return nil, err
			}
			rules = append(rules, result)
		}
	}

	return rules, nil
}

func (p *CssParser) ConsumeQualifiedRule() (string, error) {

	prelude := []Token{}

	for !p.eof() {
		switch {
		case p.IfCurrentIs(Token_EOF):
			return "", fmt.Errorf("Failed to parse qualified rule")
		case p.IfCurrentIs(Token_Clearly_Open):
			p.pos++
			p.ConsumeSimpleBlock()

			return "", nil
		default:
			prelude = append(prelude, p.input[p.pos])
			p.pos++
		}

	}

	return "", nil
}

func (p *CssParser) ConsumeAtRule() string {
	return ""
}

func (p *CssParser) ConsumeSimpleBlock() {}

func (p *CssParser) ConsumeFunction() error {

	for !p.eof() {
		switch {
		case p.IfCurrentIs(Token_Pren_Close):
			return nil
		case p.IfCurrentIs(Token_EOF):
			return fmt.Errorf("Failed to parse function")
		default:

		}
	}

}

func (p *CssParser) eof() bool {
	return p.pos >= p.len
}
