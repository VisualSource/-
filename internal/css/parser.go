package plex_css

import (
	"fmt"
	"strings"
)

type CssParser struct {
	pos   int
	len   int
	input []Token
}

func (p *CssParser) ParseStylesheet(value string, origin uint) (Stylesheet, error) {
	p.pos = 0
	tokenizer := Tokenizer{}

	tokens, err := tokenizer.Parse(value)
	if err != nil {
		return Stylesheet{}, err
	}
	p.len = len(tokens)
	p.input = tokens

	rules, atRules, err := p.ConsumeRulesList()
	if err != nil {
		return Stylesheet{}, err
	}

	return Stylesheet{
		Rules:    rules,
		AtRules:  atRules,
		TopLevel: true,
		Origin:   origin,
	}, nil
}
func (p *CssParser) ParseDeclarationsList(value string) ([]Declaration, error) {
	p.pos = 0
	tokenizer := Tokenizer{}

	tokens, err := tokenizer.Parse(value)
	if err != nil {
		return nil, err
	}
	p.len = len(tokens)
	p.input = tokens

	decs, _ := p.ConsumeDeclarationsList()

	return decs, nil
}
func (p *CssParser) ParseRule()                    {}
func (p *CssParser) ParseDeclaration()             {}
func (p *CssParser) ParseStyleBlockContent()       {}
func (p *CssParser) ParseComponentValue()          {}
func (p *CssParser) ParseComponentValues()         {}
func (p *CssParser) ParseComponentValuesByCommas() {}

func (p *CssParser) ConsumeRulesList() ([]Rule, []AtRule, error) {

	rules := []Rule{}
	atRules := []AtRule{}

	for {
		switch {
		case p.isCurrent(Token_Whitespace):
			p.pos++
		case p.isCurrent(Token_EOF):
			return rules, atRules, nil
		case p.isCurrent(Token_CDO) || p.isCurrent(Token_CDC):
			result, err := p.ConsumeQualifiedRule()
			if err == nil {
				rules = append(rules, result)
			}
		case p.isCurrent(Token_At_Keyword):
			result, err := p.ConsumeAtRule()
			if err == nil {
				atRules = append(atRules, result)
			}
		default:
			result, err := p.ConsumeQualifiedRule()
			if err == nil {
				rules = append(rules, result)
			}
		}
	}

}
func (p *CssParser) ConsumeAtRule() (AtRule, error) {

	rule := AtRule{}

	atToken := p.input[p.pos]

	if a, ok := atToken.(*StringToken); ok {
		rule.Name = string(a.Value)
	}
	p.pos++

	for {
		switch {
		// handle statement at rule
		case p.isCurrent(Token_Semicolon):
		case p.isCurrent(Token_EOF):
			return rule, fmt.Errorf("found EOF Token")
		// handle block at rule
		case p.isCurrent(Token_Clearly_Open):
			block, err := p.ConsumeSimpleBlock()
			if err != nil {
				return AtRule{}, err
			}
			if block.BlockType != Token_Clearly_Close {
				return AtRule{}, fmt.Errorf("expected '}'")
			}

			rule.Block = block
			return rule, nil
		// same as above
		case p.isCurrent(TSimpleBlack):
			if b, ok := p.input[p.pos].(*SimpleBlock); ok && b.BlockType == Token_Clearly_Close {
				rule.Block = *b
				return rule, nil
			}
		default:
			value, err := p.ConsumeComponentValue()
			if err != nil {
				return AtRule{}, err
			}
			if value.GetId() != Token_Whitespace {
				rule.Prelude = append(rule.Prelude, value)
			}
		}
	}

}
func (p *CssParser) ConsumeQualifiedRule() (Rule, error) {

	rule := Rule{}
	prelude := []Token{}

	for {
		switch {
		case p.isCurrent(Token_EOF):
			return Rule{}, fmt.Errorf("invalid rule")
		case p.isCurrent(Token_Clearly_Open):
			block, err := p.ConsumeSimpleBlock()
			if err != nil {
				return Rule{}, err
			}

			declarationParser := CssParser{}
			declarationParser.input = block.Tokens
			declarationParser.len = len(block.Tokens)
			declarations, _ := declarationParser.ConsumeDeclarationsList()
			rule.Block = declarations

			selector, err := ParseSimpleSelector(&prelude)

			if err != nil {
				return Rule{}, err
			}

			rule.Selector = append(rule.Selector, selector)

			return rule, nil
		case p.isCurrent(TSimpleBlack):
			if b, ok := p.input[p.pos].(*SimpleBlock); ok && b.BlockType == Token_Clearly_Close {
				declarationParser := CssParser{}
				declarationParser.input = (*b).Tokens
				declarationParser.len = len((*b).Tokens)
				declarations, _ := declarationParser.ConsumeDeclarationsList()
				rule.Block = declarations
				return rule, nil
			}
		default:
			result, err := p.ConsumeComponentValue()
			if err != nil {
				return Rule{}, err
			}
			prelude = append(prelude, result)
		}
	}
}
func (p *CssParser) ConsumeStyleBlockContents() {}
func (p *CssParser) ConsumeDeclarationsList() ([]Declaration, []AtRule) {

	declarations := []Declaration{}
	atRules := []AtRule{}

	for {
		switch {
		case p.isCurrent(Token_Whitespace) || p.isCurrent(Token_Semicolon):
			p.pos++
		case p.eof() || p.isCurrent(Token_EOF):
			return declarations, atRules
		case p.isCurrent(Token_At_Keyword):
			rule, err := p.ConsumeAtRule()
			if err == nil {
				atRules = append(atRules, rule)
			}
		case p.isCurrent(Token_Ident):
			/*
				FROM: https://www.w3.org/TR/css-syntax-3/#consume-a-list-of-declarations

				Initialize a temporary list initially filled with the current input token.
				As long as the next input token is anything other than a <semicolon-token> or <EOF-token>,
				consume a component value and append it to the temporary list.
				Consume a declaration from the temporary list.
				If anything was returned, append it to the list of declarations.

				ISSUE:
					the text says to put the ident token and the value from a ConsumeComponentValue call
					in a temp list and then call the consumeDeclaration from that temp list but
					doing so would result in invalid parsing. Example 'color: green;'
					'color' is the ident and the ':' would be the value from ConsumeComponentValue(https://www.w3.org/TR/css-syntax-3/#consume-a-component-value)
					and calling consumeDeclaration would return a declaration of color with no value
			*/
			if !p.isCurrent(Token_Semicolon) || !p.isCurrent(Token_EOF) {
				dec, err := p.ConsumeDeclaration()
				if err == nil {
					declarations = append(declarations, dec)
				}
			}
		default:
			p.pos++
			if !p.isCurrent(Token_Semicolon) || !p.isCurrent(Token_EOF) {
				// and throw away the returned value.
				p.ConsumeComponentValue()
			}
		}
	}

}
func (p *CssParser) ConsumeDeclaration() (Declaration, error) {

	ident := p.input[p.pos]
	var name string
	if i, ok := ident.(*StringToken); ok {
		name = string(i.Value)
	}

	decValue := []Token{}
	p.pos++

	for p.isCurrent(Token_Whitespace) {
		p.pos++
	}

	if !p.isCurrent(Token_Colon) {
		return Declaration{}, fmt.Errorf("expected to find token ':' but got token: %d", p.input[p.pos].GetId())
	}
	p.pos++ // eat ':'

	for p.isCurrent(Token_Whitespace) {
		p.pos++
	}

	// spec does says only to watch for EOF but should probably watch for ';' token?
	for !p.eof() && !p.isCurrent(Token_Semicolon) {
		value, err := p.ConsumeComponentValue()
		if err != nil {
			return Declaration{}, err
		}

		decValue = append(decValue, value)
	}

	for decValue[len(decValue)-1].GetId() == Token_Whitespace {
		decValue = decValue[:len(decValue)-1]
	}

	important := false
	if len(decValue) > 2 {
		markToken := decValue[len(decValue)-2]
		importantToken := decValue[len(decValue)-1]

		if isRune('!', &markToken) && isStringCaseInsensitive("important", &importantToken) {
			important = true
			decValue = decValue[:len(decValue)-2]
		}
	}

	for decValue[len(decValue)-1].GetId() == Token_Whitespace {
		decValue = decValue[:len(decValue)-1]
	}

	return Declaration{
		Value:     ParseCssValue(decValue),
		Name:      name,
		Important: important,
	}, nil
}

func (p *CssParser) ConsumeComponentValue() (Token, error) {

	if p.isCurrent(Token_Clearly_Open) || p.isCurrent(Token_Square_Bracket_Open) || p.isCurrent(Token_Pren_Open) {
		result, err := p.ConsumeSimpleBlock()

		return &result, err
	}

	if p.isCurrent(Token_Function) {
		return p.ConsumeFunction()
	}

	token := p.input[p.pos]
	p.pos++

	return token, nil
}

// Note: This algorithm assumes that the current input token has already been checked to be an <{-token>, <[-token>, or <(-token>.
func (p *CssParser) ConsumeSimpleBlock() (SimpleBlock, error) {
	var bracketEnd TokenType = Token_Clearly_Close
	switch p.input[p.pos].GetId() {
	case Token_Square_Bracket_Open:
		bracketEnd = Token_Square_Bracket_Close
	case Token_Pren_Open:
		bracketEnd = Token_Pren_Close
	}
	p.pos++

	tokens := []Token{}

	for {
		switch {
		case p.isCurrent(bracketEnd):
			p.pos++ // eat end bracket
			return SimpleBlock{Tokens: tokens, BlockType: bracketEnd}, nil
		case p.eof():
			return SimpleBlock{Tokens: tokens, BlockType: bracketEnd}, fmt.Errorf("found EOF")
		default:
			result, err := p.ConsumeComponentValue()

			if err != nil {
				return SimpleBlock{}, err
			}

			tokens = append(tokens, result)
		}
	}
}
func (p *CssParser) ConsumeFunction() (Token, error) {

	var name string
	token := p.input[p.pos]
	if v, ok := token.(*StringToken); ok {
		name = v.Value
	}

	p.pos++

	args := []Token{}

	for !p.eof() && !p.isCurrent(Token_Pren_Close) {

		result, err := p.ConsumeComponentValue()

		if err != nil {
			return nil, err
		}
		if result.GetId() != Token_Whitespace {
			args = append(args, result)
		}
	}
	p.pos++

	if p.eof() {
		return &FunctionBlock{Args: args, Name: name}, fmt.Errorf("found EOF")
	}

	return &FunctionBlock{Args: args, Name: name}, nil
}

func (p *CssParser) isCurrent(t TokenType) bool {
	if p.eof() {
		return false
	}
	return p.input[p.pos].GetId() == t
}

func (p *CssParser) eof() bool {
	return p.pos >= p.len
}

func isStringCaseInsensitive(value string, t *Token) bool {
	if (*t).GetId() != Token_Ident {
		return false
	}

	if v, ok := (*t).(*StringToken); ok && strings.ToLower(string(v.Value)) == value {
		return true
	}

	return false
}

func isRune(v rune, t *Token) bool {
	id := (*t).GetId()
	if id != Token_Delim {
		return false
	}

	if mark, ok := (*t).(*RuneToken); ok && mark.Value == v {
		return true
	}

	return false
}
