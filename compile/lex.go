package compile

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
)

type TokenType uint8

const (
	TokenCharacter TokenType = iota + 1
	TokenBracketStart
	TokenBracketClose
	TokenCommandStart
	TokenCommandClose
	TokenNewline
	TokenDecorator
	TokenEOF
)

type Token struct {
	Type   TokenType
	Line   int
	Column int
	Value  string
}

func (tokenType TokenType) String() string {
	switch tokenType {
	case TokenBracketStart:
		return "[      "
	case TokenBracketClose:
		return "]      "
	case TokenCommandStart:
		return "{      "
	case TokenCommandClose:
		return "}      "
	case TokenEOF:
		return "EOF    "
	case TokenNewline:
		return "NEWLINE"
	case TokenCharacter:
		return "CHAR   "
	case TokenDecorator:
		return "@      "
	default:
		return "UNKNOWN"
	}
}

func (token Token) String() string {
	return fmt.Sprintf(
		`%s: "%s" (Line=%d, Col=%d)`,
		token.Type.String(),
		token.Value,
		token.Line,
		token.Column,
	)
}

type RuneReader struct {
	reader *bufio.Reader
	cache  *rune
}

func (reader *RuneReader) Pop() (rune, error) {
	if reader.cache != nil {
		char := *reader.cache
		reader.cache = nil
		return char, nil
	}

	char, _, err := reader.reader.ReadRune()
	if err != nil {
		return 0, err
	}
	return char, nil
}

func (reader *RuneReader) Peek() (rune, error) {
	if reader.cache != nil {
		return *reader.cache, nil
	}
	char, _, err := reader.reader.ReadRune()
	if err != nil {
		return 0, err
	}
	reader.cache = &char
	return char, nil
}

func Lex(ctx context.Context, file *bufio.Reader, tokens chan<- Token) error {
	// TODO: Respect the context
	// TODO: Extract to make easier to test? Or mock channel
	reader := RuneReader{reader: file}

	var line, column int
	var escaped bool
	var comment bool
	for {
		char, err := reader.Pop()
		if errors.Is(err, io.EOF) {
			tokens <- Token{
				Type:   TokenEOF,
				Line:   line,
				Column: column,
			}
			return nil
		} else if err != nil {
			return fmt.Errorf("reading character: %w", err)
		}

		switch {
		case comment && char != '\n':
		case comment && char == '\n':
			fallthrough
		case char == '\n':
			if !escaped {
				tokens <- Token{
					Type:   TokenNewline,
					Value:  "\n",
					Line:   line,
					Column: column,
				}
			}
			escaped = false
			comment = false
			line++
			column = -1
		case !escaped && char == '@':
			tokens <- Token{
				Type:   TokenDecorator,
				Value:  "@",
				Line:   line,
				Column: column,
			}
		case !escaped && char == '#':
			comment = true
		case !escaped && char == '\\':
			escaped = true
		case !escaped && char == '{':
			tokens <- Token{
				Type:   TokenCommandStart,
				Value:  "{",
				Line:   line,
				Column: column,
			}
		case !escaped && char == '}':
			tokens <- Token{
				Type:   TokenCommandClose,
				Value:  "}",
				Line:   line,
				Column: column,
			}
		case !escaped && char == '[':
			tokens <- Token{
				Type:   TokenBracketStart,
				Value:  "[",
				Line:   line,
				Column: column,
			}
		case !escaped && char == ']':
			tokens <- Token{
				Type:   TokenBracketClose,
				Value:  "]",
				Line:   line,
				Column: column,
			}
		default:
			tokens <- Token{
				Type:   TokenCharacter,
				Value:  fmt.Sprintf("%c", char),
				Line:   line,
				Column: column,
			}
			escaped = false
		}
		column++
	}
}
