package compile

import (
	"context"
	"errors"
	"fmt"
	"io"
)

type Opcode string

const (
	OpPrint  = "print"
	OpDelete = "delete"
)

type Instruction struct {
	Opcode Opcode `json:"opcode"`
	Arg    any    `json:"arg"`
}

type Program struct {
	Instructions []Instruction
}

type TokenReader struct {
	Channel <-chan Token
	cache   *Token
}

func (reader *TokenReader) Peek() Token {
	if reader.cache != nil {
		return *reader.cache
	}
	next := <-reader.Channel
	reader.cache = &next
	return next
}

func (reader *TokenReader) Pop() Token {
	next := reader.Peek()
	reader.cache = nil
	return next
}

func Parse(ctx context.Context, tokens <-chan Token) (instructions []Instruction, err error) {
	var reader = TokenReader{
		Channel: tokens,
	}
	return parse(ctx, &reader)
}

func parse(ctx context.Context, reader *TokenReader) ([]Instruction, error) {
	var instructions []Instruction
	for {
		newInstructions, err := parseSwitch(ctx, reader, reader.Pop())
		if errors.Is(err, io.EOF) {
			return instructions, nil
		} else if err != nil {
			return nil, err
		}
		instructions = append(instructions, newInstructions...)
	}
}

func parseSwitch(ctx context.Context, reader *TokenReader, token Token) ([]Instruction, error) {
	switch token.Type {
	case TokenNewline:
		fallthrough
	case TokenCharacter:
		return []Instruction{Instruction{
			Opcode: OpPrint,
			Arg:    token.Value,
		}}, nil
	case TokenBracketStart:
		bracketInstructions, err := parseBracket(ctx, reader)
		if err != nil {
			return nil, fmt.Errorf("invalid bracket setction: %w", err)
		}
		return bracketInstructions, nil
	case TokenEOF:
		return []Instruction{}, io.EOF
	default:
		return nil, fmt.Errorf("unknown token: %s", token.String())
	}
}
func parseBracket(ctx context.Context, reader *TokenReader) ([]Instruction, error) {
	var instructions []Instruction
	var chars int
	for {
		next := reader.Pop()
		switch {
		case next.Type == TokenBracketClose:
			// Happy case
			// TODO: Append the deleting command
			return append(instructions, Instruction{
				Opcode: OpDelete,
				Arg:    chars,
			}), nil
		case next.Type == TokenNewline:
			fallthrough
		case next.Type == TokenEOF:
			return nil, fmt.Errorf(`expected: "]" got: "%s"`, next)
		case next.Type == TokenBracketStart:
			return nil, errors.New("not implemented: nested bracked")
		case next.Type == TokenCharacter:
			instructions = append(instructions, Instruction{
				Opcode: OpPrint,
				Arg:    next.Value,
			})
			chars++
		default:
			return nil, fmt.Errorf("unknown token: %s", next.String())
		}
	}
}
