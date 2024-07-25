package plex_css

type TokenType = uint8

const (
	Token_Ident                TokenType = 0
	Token_Function             TokenType = 1
	Token_At_Ketword           TokenType = 2
	Token_Hash                 TokenType = 3
	Token_String               TokenType = 4
	Token_Bad_String           TokenType = 5
	Token_Url                  TokenType = 6
	Token_Bad_Url              TokenType = 7
	Token_Delim                TokenType = 8
	Token_Percentage           TokenType = 9
	Token_Dimension            TokenType = 10
	Token_Whitespace           TokenType = 11
	Token_CDO                  TokenType = 12
	Token_CDC                  TokenType = 13
	Token_Colon                TokenType = 14
	Token_Semicolon            TokenType = 15
	Token_Comma                TokenType = 16
	Token_Square_Bracket_Open  TokenType = 17
	Token_Square_Bracket_Close TokenType = 18
	Token_Pren_Open            TokenType = 19
	Token_Pren_Close           TokenType = 20
	Token_Clearly_Open         TokenType = 21
	Token_Clearly_Close        TokenType = 22
)

type Token interface {
	GetId() TokenType
}

type StringToken struct {
	Id    TokenType
	Value []rune
}

func (t *StringToken) GetId() TokenType {
	return t.Id
}

type EmptyToken struct {
	Id TokenType
}

func (t *EmptyToken) GetId() TokenType {
	return t.Id
}

type RuneToken struct {
	Id    TokenType
	Value rune
}

func (t *RuneToken) GetId() TokenType {
	return t.Id
}

type FlagedStringToken struct {
	Id    TokenType
	Value []rune
	Flag  []rune
}

func (f *FlagedStringToken) GetId() TokenType {
	return f.Id
}
