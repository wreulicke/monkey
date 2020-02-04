package token

type Token struct {
	Type    TokenType
	Literal string
}

var typeNames = []string{
	"ILLEGAL",
	"EOF",

	"IDENT",
	"NUMBER",
	"STRING",

	"ASSIGN",
	"PLUS",
	"MINUS",
	"BANG",
	"ASTERISK",
	"SLASH",
	"PIPELINE",

	"EQ",
	"NOT_EQ",

	"LT",
	"GT",

	"COMMA",
	"SEMICOLON",
	"COLON",

	"LBRACKET",
	"RBRACKET",

	"LPAREN",
	"RPAREN",

	"LBRACE",
	"RBRACE",

	"FUNCTION",
	"LET",
	"RETURN",
	"TRUE",
	"FALSE",
	"IF",
	"ELSE",
}

type TokenType int

func (t TokenType) String() string {
	return typeNames[t]
}

const (
	ILLEGAL TokenType = iota
	EOF

	IDENT
	NUMBER
	STRING

	ASSIGN
	PLUS
	MINUS
	BANG
	ASTERISK
	SLASH
	PIPELINE

	EQ
	NOT_EQ

	LT
	GT

	COMMA
	SEMICOLON
	COLON

	LBRACKET
	RBRACKET

	LPAREN
	RPAREN

	LBRACE
	RBRACE

	FUNCTION
	LET
	RETURN
	TRUE
	FALSE
	IF
	ELSE
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
