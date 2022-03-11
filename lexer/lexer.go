package lexer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"

	token "github.com/wreulicke/monkey/token"
)

const eof = -1

type Position struct {
	line   int
	column int
}

type Lexer struct {
	input    *bufio.Reader
	buffer   bytes.Buffer
	position *Position
	offset   int
	error    error
}

func New(input io.Reader) *Lexer {
	l := &Lexer{input: bufio.NewReader(input)}
	l.position = &Position{line: 1}
	return l
}

func (l *Lexer) Error(e string) {
	err := fmt.Errorf("%s in %d:%d", e, (*l).position.line, (*l).position.column)
	l.error = err
}

func (l *Lexer) TokenText() string {
	return l.buffer.String()
}

func (l *Lexer) Next() rune {
	r, w, err := l.input.ReadRune()
	if err == io.EOF {
		return eof
	}
	if r == '\n' {
		l.position = &Position{line: l.position.line + 1}
	}
	l.position.column += w
	l.offset += w
	l.buffer.WriteRune(r)
	return r
}

func (l *Lexer) Skip() rune {
	r, w, err := l.input.ReadRune()
	if err == io.EOF {
		return eof
	}
	if r == '\n' {
		l.position = &Position{line: l.position.line + 1}
	}
	l.position.column += w
	l.offset += w
	return r
}

func (l *Lexer) Peek() rune {
	lead, err := l.input.Peek(1)
	if err == io.EOF {
		return eof
	} else if err != nil {
		l.Error(err.Error())
		return 0
	}

	p, err := l.input.Peek(runeLen(lead[0]))

	if err == io.EOF {
		return eof
	} else if err != nil {
		l.Error("unexpected input error")
		return 0
	}

	ruNe, _ := utf8.DecodeRune(p)
	return ruNe
}

func (l *Lexer) readIdentifier() {
	next := l.Peek()
	for unicode.IsLetter(next) {
		l.Next()
		next = l.Peek()
	}
}

func (l *Lexer) readNumber(next rune) {
	if next == '0' && isDigit(l.Peek()) {
		l.Error("unexpected digit '0'")
		return
	} else if isDigit(next) {
		next := l.Peek()
		for {
			if !isDigit(next) {
				break
			}
			l.Next()
			next = l.Peek()
		}
		next = l.Peek()
		if next == '.' {
			l.Next()
			next = l.Peek()
			if !isDigit(next) {
				l.Error("unexpected token: expected digits")
				return
			}
			for {
				if !isDigit(next) {
					break
				}
				l.Next()
				next = l.Peek()
			}
		}
		next = l.Peek()
		if next == 'e' || next == 'E' {
			l.Next()
			next := l.Peek()
			if next == '+' || next == '-' {
				l.Next()
			}
			next = l.Peek()
			if !isDigit(next) {
				l.Error("digit expected for number exponent")
				return
			}
			l.Next()
			next = l.Peek()
			for {
				if !isDigit(next) {
					break
				}
				l.Next()
				next = l.Peek()
			}
		}
	} else {
		l.Error("error")
		return
	}
}

func (l *Lexer) readString(start rune) {
	for {
		next := l.Peek()
		if next == start {
			l.Skip()
			return
		}
		switch {
		case next == '\\':
			l.Skip()
			next := l.Peek()
			if next == start {
				l.Next()
			} else if next == 'b' {
				l.Skip()
				l.buffer.WriteRune('\b')
			} else if next == 'f' {
				l.Skip()
				l.buffer.WriteRune('\f')
			} else if next == 'n' {
				l.Skip()
				l.buffer.WriteRune('\n')
			} else if next == 'r' {
				l.Skip()
				l.buffer.WriteRune('\r')
			} else if next == 't' {
				l.Skip()
				l.buffer.WriteRune('\t')
			} else {
				l.Error("unsupported escape character")
				return
			}
		case unicode.IsControl(next):
			l.Error("cannot contain control characters in strings")
			return
		case next == eof:
			l.Error("unclosed string")
			return
		default:
			l.Next()
		}
	}
}

func (l *Lexer) skipWhitespace() {
	ruNe := l.Peek()
	for unicode.IsSpace(ruNe) {
		l.Next()
		ruNe = l.Peek()
	}
	l.buffer.Reset()
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	next := l.Peek()
	switch next {
	case '"':
		l.Skip()
		l.readString(next)
		return l.newToken(token.STRING)
	case '\'':
		l.Skip()
		l.readString(next)
		return l.newToken(token.STRING)
	}
	next = l.Next()
	switch next {
	case '=':
		if l.Peek() == '=' {
			l.Next()
			return l.newToken(token.EQ)
		}
		return l.newToken(token.ASSIGN)
	case '+':
		return l.newToken(token.PLUS)
	case '-':
		return l.newToken(token.MINUS)
	case '!':
		if l.Peek() == '=' {
			l.Next()
			return l.newToken(token.NOT_EQ)
		}
		return l.newToken(token.BANG)
	case '/':
		return l.newToken(token.SLASH)
	case '*':
		return l.newToken(token.ASTERISK)
	case '<':
		return l.newToken(token.LT)
	case '>':
		return l.newToken(token.GT)
	case ':':
		return l.newToken(token.COLON)
	case ';':
		return l.newToken(token.SEMICOLON)
	case '(':
		return l.newToken(token.LPAREN)
	case ')':
		return l.newToken(token.RPAREN)
	case ',':
		return l.newToken(token.COMMA)
	case '{':
		return l.newToken(token.LBRACE)
	case '}':
		return l.newToken(token.RBRACE)
	case '[':
		return l.newToken(token.LBRACKET)
	case ']':
		return l.newToken(token.RBRACKET)
	case '|':
		return l.newToken(token.PIPELINE)
	case eof:
		return l.newToken(token.EOF)
	default:
		if isLetter(next) {
			l.readIdentifier()
			return l.newToken(token.LookupIdent(l.TokenText()))
		} else if isDigit(next) {
			l.readNumber(next)
			return l.newToken(token.NUMBER)
		}
		return l.newToken(token.ILLEGAL)
	}
}

func (l *Lexer) newToken(tokenType token.TokenType) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: l.TokenText(),
	}
}

func runeLen(lead byte) int {
	if lead < 0xC0 {
		return 1
	} else if lead < 0xE0 {
		return 2
	} else if lead < 0xF0 {
		return 3
	}
	return 4
}
func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
