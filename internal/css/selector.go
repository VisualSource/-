package plex_css

import (
	"fmt"

	"github.com/moznion/go-optional"
)

func ParseForgivingSelectorList(value []Token) []Selector {

	valueLen := len(value)
	selectors := []Selector{}
	pos := 0

	for pos < valueLen {
		result, err := praseComplexSelector(&value, valueLen, &pos)
		if err != nil {
			fmt.Printf("Selector parse error: %s", err)
			continue
		}
		selectors = append(selectors, result)
	}

	return selectors
}

func praseComplexSelector(tokens *[]Token, len int, pos *int) (Selector, error) {

	selector := Selector{}

	parseCompundSelector(tokens, pos, len)

	for (*pos) < len {
		_, err := ParseCombinator(tokens, pos, len)
		if err != nil {
			return selector, err
		}
		parseCompundSelector(tokens, pos, len)
	}

	return selector, nil
}

/*
Grammer:

		<compound-selector> = [ <type-selector>? <subclass-selector>*
	                        [ <pseudo-element-selector> <pseudo-class-selector>* ]* ]!
*/
func parseCompundSelector(tokens *[]Token, pos *int, len int) {

	/*parseTypeSelector()
	for {
		parseTypeSelector()
	}

	for {
		parsePseudoElementSelector()
		for {
			parsePseudoClassSelector()
		}
	}*/
}

/*
Grammer:

	<combinator> = '>' | '+' | '~' | [ '|' '|' ]
*/
func ParseCombinator(tokens *[]Token, pos *int, len int) (uint, error) {
	if *pos > len {
		return 0, nil
	}
	if v, ok := (*tokens)[*pos].(*RuneToken); ok {
		switch v.Value {
		case '>':
			(*pos)++
			return 0, nil
		case '+':
			(*pos)++
			return 1, nil
		case '~':
			(*pos)++
			return 2, nil
		case '|':
			(*pos)++
			if (*pos) > len {
				return 0, nil
			}
			if b, ok := (*tokens)[(*pos)].(*RuneToken); ok && b.Value == '|' {
				(*pos)++
				return 3, nil
			}
			return 0, fmt.Errorf("was expecting token '|' but got %d", (*tokens)[(*pos)].GetId())
		}
	}

	return 0, nil
}

/*
Grammer:

	<type-selector> = <wq-name> | <ns-prefix>? '*'

This is required because the slector of *|* is not posiable with only <wq-name>
as it requires a prefix ident systax where '*' would be a delim
*/
func ParseTypeSelector(tokens *[]Token, pos *int, len int) (optional.Option[string], string, error) {

	namespace, tagname, err := ParseWqName(tokens, pos, len)
	if err == nil {
		return namespace, tagname, nil
	}

	nsPrefix, err := ParseNSPrefix(tokens, pos, len)

	if (*pos) >= len {
		return nil, "", fmt.Errorf("eof")
	}

	token := (*tokens)[(*pos)]

	if token.GetId() != Token_Delim {
		return nil, "", fmt.Errorf("expected a token of delim but got: %d", token.GetId())
	}

	v := token.(*RuneToken)
	if v.Value != '*' {
		return nil, "", fmt.Errorf("was expecting a rune of '*' but found %q", v.Value)
	}
	(*pos)++

	if err == nil {
		return optional.Some(nsPrefix), "*", nil
	}

	return nil, "*", nil
}

/*
Grammer:

	<wq-name> = <ns-prefix>? <ident-token>
*/
func ParseWqName(tokens *[]Token, pos *int, len int) (optional.Option[string], string, error) {

	nsPrefix, err := ParseNSPrefix(tokens, pos, len)

	if (*pos) >= len {
		return nil, "", fmt.Errorf("eof")
	}

	token := (*tokens)[(*pos)]

	if token.GetId() != Token_Ident {
		if err == nil {
			if nsPrefix != "" {
				(*pos) -= 2
			} else {
				(*pos) -= 1
			}
		}
		return nil, "", fmt.Errorf("expected ident token got: %d", token.GetId())
	}
	(*pos)++
	v := token.(*StringToken)

	if err == nil {
		return optional.Some(nsPrefix), v.Value, nil
	}

	return nil, v.Value, nil
}

/*
Grammer:

	<ns-prefix> = [ <ident-token> | '*' ]? '|'
*/
func ParseNSPrefix(tokens *[]Token, pos *int, len int) (string, error) {

	if (*pos) >= len {
		return "", fmt.Errorf("eof")
	}

	switch (*tokens)[*pos].GetId() {
	case Token_Ident:
		ident := (*tokens)[*pos].(*StringToken)
		if (*pos)+1 >= len {
			return "", fmt.Errorf("eof")
		}
		token := (*tokens)[(*pos)+1]

		if token.GetId() != Token_Delim {
			return "", fmt.Errorf("was expecting a delim token but got: %d", token.GetId())
		}
		pipe := token.(*RuneToken)
		if pipe.Value != '|' {
			return "", fmt.Errorf("was expecting a rune of '|' but got: %q", pipe.Value)
		}
		(*pos) += 2
		return ident.Value, nil
	case Token_Delim:
		pipe := (*tokens)[*pos].(*RuneToken)

		switch pipe.Value {
		case '*':
			if (*pos)+1 >= len {
				return "", fmt.Errorf("eof")
			}
			token := (*tokens)[*pos+1]

			if token.GetId() != Token_Delim {
				return "", fmt.Errorf("was expecting a delim token but got: %d", token.GetId())
			}
			pipe := token.(*RuneToken)
			if pipe.Value != '|' {
				return "", fmt.Errorf("was expecting a rune of '|' but got: %q", pipe.Value)
			}
			(*pos) += 2
			return "*", nil
		case '|':
			(*pos)++
			return "", nil
		}
	}

	return "", fmt.Errorf("invalid namepsace prefix")
}
