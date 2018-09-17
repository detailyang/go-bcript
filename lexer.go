package bscript

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrLexerInvalidHexString = errors.New("lexer: invalid hex string")
	ErrLexerReachEOF         = errors.New("lexer: reach end of source")
	ErrLexerUnknowCharacter  = errors.New("lexer: unknow character")
	ErrLexerNotFoundOPCode   = errors.New("lexer: not found opcode")
	ErrLexerNotHexString     = errors.New("lexer: not hex string")
	ErrLexerUnknowOPCode     = errors.New("lexer: unknow opcode")
	ErrLexerNumberOverFlow   = errors.New("lexer: overflow int64")
)

const (
	TOKEN_CODE = iota
	TOKEN_NUMBER
	TOKEN_HEXSTRING
	TOKEN_COMMENT
)

type Token struct {
	value interface{}
	raw   interface{}
	kind  int
}

func (o *Token) String() string {
	t := "OPCODE"

	switch o.kind {
	case TOKEN_CODE:
		t = "TOKEN OPCODE"
	case TOKEN_COMMENT:
		t = "TOKEN COMMENT"
	case TOKEN_HEXSTRING:
		t = "TOKEN HEXSTRING"
	case TOKEN_NUMBER:
		t = "TOKEN NUMBER"
	}
	return fmt.Sprintf("Token: type: %s value:%#v", t, o.value)
}

type Lexer struct {
	src    string
	pos    int
	line   int
	column int
}

func NewLexer(src string) *Lexer {
	return &Lexer{
		src:    src,
		pos:    0,
		line:   0,
		column: 0,
	}
}

func (l *Lexer) croak(err error) error {
	return fmt.Errorf("%s at %d:%d", err, l.line, l.column)
}

func (l *Lexer) lookahead(n int, consume bool) (byte, bool) {
	if l.pos+n > len(l.src)-1 {
		return 0, false
	}

	if consume {
		l.pos += n
		return l.src[l.pos], true
	}

	return l.src[l.pos+n], true
}

func (l *Lexer) consume(n int) (int, bool) {
	if l.pos+n > len(l.src)-1 {
		return 0, false
	}

	pos := l.pos
	l.pos += n
	return pos, true
}

func (l *Lexer) poke() (byte, error) {
	if l.pos > len(l.src)-1 {
		return 0, ErrLexerReachEOF
	}
	ch := l.src[l.pos]
	l.pos++

	return ch, nil
}

func (l *Lexer) scanOPCode() (*Token, error) {
	begin, _ := l.consume(3)

	for {
		ch, err := l.poke()
		if err != nil {
			l.pos++
			break
		}

		if !('A' <= ch && ch <= 'Z' || ('0' <= ch && ch <= '9')) {
			break
		}
	}

	code, err := NewOPCodeFromString(l.src[begin : l.pos-1])
	if err != nil {
		return nil, err
	}

	return &Token{
		value: code,
		raw:   code,
		kind:  TOKEN_CODE,
	}, nil
}

func (l *Lexer) scanComment() (*Token, error) {
	begin, _ := l.consume(1)

	for {
		ch, err := l.poke()
		if err != nil {
			l.pos++
			break
		}

		if ch == '\n' {
			break
		}
	}

	comment := l.src[begin : l.pos-1]

	return &Token{
		value: comment,
		raw:   comment,
		kind:  TOKEN_COMMENT,
	}, nil
}

func (l *Lexer) scanHexstring() (*Token, error) {
	begin, _ := l.consume(2)

	for {
		ch, err := l.poke()
		if err != nil {
			l.pos++
			break
		}

		if !(('a' <= ch && ch <= 'f') || ('0' <= ch && ch <= '9') || ('A' <= ch && ch <= 'F')) {
			break
		}
	}

	hexstring := l.src[begin : l.pos-1]
	if len(hexstring) < 2 {
		return nil, ErrLexerNotHexString
	}

	d, err := hex.DecodeString(hexstring[2:])
	if err != nil {
		return nil, ErrLexerInvalidHexString
	}

	return &Token{
		value: d,
		raw:   hexstring,
		kind:  TOKEN_HEXSTRING,
	}, nil
}

func (l *Lexer) scanNumber(neg bool) (*Token, error) {
	begin := l.pos
	if neg {
		l.pos++
	}

	for {
		ch, err := l.poke()
		if err != nil {
			l.pos++
			break
		}

		if !('0' <= ch && ch <= '9') {
			break
		}
	}

	numstr := l.src[begin : l.pos-1]
	num, err := strconv.ParseInt(numstr, 10, 64)
	if err != nil {
		return nil, err
	}

	return &Token{
		value: num,
		raw:   numstr,
		kind:  TOKEN_NUMBER,
	}, nil
}

type LexerPos struct {
	line   int
	column int
}

func (l *Lexer) Pos() LexerPos {
	return LexerPos{
		line:   l.line,
		column: l.column,
	}
}

func (l *Lexer) Scan() (*Token, error) {
	for {
		ch, ok := l.lookahead(0, false)
		if !ok {
			break
		}

		peek, _ := l.lookahead(1, false)
		peekpeek, _ := l.lookahead(2, false)

		switch ch {
		case '\f':
			fallthrough
		case '\r':
			fallthrough
		case ' ':
			fallthrough
		case '\t':
			l.column++
			_, ok := l.lookahead(1, true)
			if !ok {
				return nil, ErrLexerReachEOF
			}
		case '\n':
			l.line++
			l.column = 0
			_, ok := l.lookahead(1, true)
			if !ok {
				return nil, ErrLexerReachEOF
			}
		case 'O':
			if peek == 'P' && peekpeek == '_' {
				return l.scanOPCode()
			}
			return nil, ErrLexerUnknowCharacter
		case '#':
			return l.scanComment()
		case '-':
			if '0' <= peek && peek <= '9' {
				return l.scanNumber(true)
			}
			return nil, ErrLexerUnknowCharacter
		default:
			switch {
			case ch >= '0' && ch <= '9':
				if ch == '0' && peek == 'x' {
					return l.scanHexstring()
				}
				return l.scanNumber(false)
			default:
				return nil, ErrLexerUnknowCharacter
			}
		}
	}

	return nil, ErrLexerReachEOF
}
