package plex_css

import (
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/moznion/go-optional"
)

func ParseSimpleSelector(tokens *[]Token) (Selector, error) {
	selector := Selector{
		Classes: mapset.NewSet[string](),
	}

	len := len(*tokens)
	pos := 0

	if isWQStart(tokens, &pos, len) {
		namespace, tagname, err := ParseTypeSelector(tokens, &pos, len)
		if err != nil {
			return selector, err
		}
		selector.Namespace = namespace
		selector.TagName = tagname
		return selector, nil
	}

	id, class, pesudoClass, attr, err := ParseSubclassSelector(tokens, &pos, len)
	if err != nil {
		return selector, err
	}

	id.IfSome(func(v string) {
		selector.Id = v
	})
	class.IfSome(func(v string) {
		selector.Classes.Add(v)
	})
	pesudoClass.IfSome(func(v PesudoClass) {
		selector.PseudoClasses = append(selector.PseudoClasses, v)
	})
	attr.IfSome(func(v SelectorAttribute) {
		selector.Attributes[v.Value] = v
	})

	return selector, nil
}

/*
Grammer:

		<subclass-selector> = <id-selector> | <class-selector> |
	                      <attribute-selector> | <pseudo-class-selector>
*/
func ParseSubclassSelector(tokens *[]Token, pos *int, len int) (optional.Option[string], optional.Option[string], optional.Option[PesudoClass], optional.Option[SelectorAttribute], error) {

	switch (*tokens)[(*pos)].GetId() {
	case Token_Hash:
		id, err := ParseIdSelector(tokens, pos, len)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		return optional.Some(id), nil, nil, nil, nil
	case Token_Delim:
		v := (*tokens)[(*pos)].(*RuneToken)
		switch v.Value {
		case '.':
			class, err := ParseClassSelector(tokens, pos, len)
			if err != nil {
				return nil, nil, nil, nil, err
			}
			return nil, optional.Some(class), nil, nil, nil
		case ':':
			el, err := ParsePseudoClassSelector(tokens, pos, len)
			if err != nil {
				return nil, nil, nil, nil, err
			}

			return nil, nil, optional.Some(el), nil, nil
		}
	case TSimpleBlack:
		attr, err := ParseAttributeSelector(tokens, pos, len)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		return nil, nil, nil, optional.Some(attr), nil
	}

	return nil, nil, nil, nil, nil
}

func ParsePseudoElementSelector(tokens *[]Token, pos *int, len int) (PesudoElement, error) {
	return PesudoElement{}, fmt.Errorf("TODO: implement pesudo elements")
}

func ParsePseudoClassSelector(tokens *[]Token, pos *int, len int) (PesudoClass, error) {
	return PesudoClass{}, fmt.Errorf("TODO: implement pesudo classes")
}

/*
Grammer:

	<attribute-selector> = '[' <wq-name> ']' | '[' <wq-name> <attr-matcher> [ <string-token> | <ident-token> ] <attr-modifier>? ']'
*/
func ParseAttributeSelector(tokens *[]Token, pos *int, len int) (SelectorAttribute, error) {
	return SelectorAttribute{}, nil
}

/*
Grammer:

	<class-selector> = '.' <ident-token>
*/
func ParseClassSelector(tokens *[]Token, pos *int, len int) (string, error) {
	if (*pos)+1 > len {
		return "", fmt.Errorf("eof")
	}

	delim := (*tokens)[(*pos)]

	if delim.GetId() != Token_Delim {
		return "", fmt.Errorf("was expecting a delim token but found: %d", delim.GetId())
	}

	d := delim.(*RuneToken)

	if d.Value != '.' {
		return "", fmt.Errorf("was expecting a rune of '.' but got: %q", d.Value)
	}

	token := (*tokens)[(*pos)+1]

	if token.GetId() != Token_Ident {
		return "", fmt.Errorf("was expecting a ident token but found: %d", token.GetId())
	}

	v := token.(*StringToken)
	(*pos) += 2
	return v.Value, nil
}

/*
Grammer:

	<id-selector> = <hash-token>
*/
func ParseIdSelector(tokens *[]Token, pos *int, len int) (string, error) {
	if (*pos) > len {
		return "", fmt.Errorf("eof")
	}

	token := (*tokens)[*pos]

	if token.GetId() != Token_Hash {
		return "", fmt.Errorf("was expecting a hash token but found: %d", token.GetId())
	}

	v := token.(*FlagedStringToken)

	(*pos)++
	return v.Value, nil
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

func isWQStart(tokens *[]Token, pos *int, len int) bool {
	if (*pos) > len {
		return true
	}

	switch (*tokens)[*pos].GetId() {
	case Token_Ident:
		return true
	case Token_Delim:
		v := (*tokens)[*pos].(*RuneToken)
		return v.Value == '*' || v.Value == '|'
	default:
		return false
	}

}
