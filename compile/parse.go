package compile

import (
	"context"
	"errors"
	"fmt"
	"io"
)

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

func Parse(ctx context.Context, tokens <-chan Token) (program *Program, err error) {
	var reader = TokenReader{
		Channel: tokens,
	}
	instructions, err := parse(ctx, &reader)
	if err != nil {
		return nil, err
	}
	return &Program{
		Instructions: instructions,
	}, nil
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
	// TODO: move into main parse function
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
			return nil, fmt.Errorf("invalid bracket section: %w", err)
		}
		return bracketInstructions, nil
	case TokenCommandStart:
		command, err := parseCommand(ctx, reader)
		if err != nil {
			return nil, fmt.Errorf("invalid command: %w", err)
		}
		return command, nil
	case TokenEOF:
		return []Instruction{}, io.EOF
	default:
		return nil, fmt.Errorf("switch: unknown token: %s", token.String())
	}
}

func parseCommand(ctx context.Context, reader *TokenReader) ([]Instruction, error) {
	// TODO: Pop all the tokens until the bracket end, then parse off those tokens!
	//       To enable a lot more types of commands
	var instructions []Instruction
	for {
		next := reader.Pop()
		switch next.Type {
		case TokenCommandClose:
			// Happy case
			return instructions, nil
		case TokenCharacter:
			// TODO: Support whitespace
			// {.. clear}
			if next.Value == "." {
				instructions = append(instructions, Instruction{
					Opcode: OpSleep,
					Arg:    1,
				})
			} else if next.Value == "c" {
				chars := []string{"l", "e", "a", "r"}
				for _, char := range chars {
					next := reader.Pop()
					if next.Type != TokenCharacter {
						// TODO: Make this error message good
						return nil, fmt.Errorf("invalid keyword")
					} else if next.Value != char {
						return nil, fmt.Errorf("invalid keyword: expected '%s'", char)
					}
				}
				instructions = append(instructions, Instruction{
					Opcode: OpClear,
				})
			}
		default:
			return nil, fmt.Errorf("expected '}' or character got: %s", next.String())
		}
	}
}

func parseBracket(ctx context.Context, reader *TokenReader) ([]Instruction, error) {
	var instructions []Instruction
	var chars int
	for {
		next := reader.Pop()
		switch next.Type {
		case TokenBracketClose:
			// Happy case
			// TODO: Append the deleting command
			return append(instructions, Instruction{
				Opcode: OpDelete,
				Arg:    chars,
			}), nil
		case TokenNewline:
			fallthrough
		case TokenEOF:
			return nil, fmt.Errorf(`expected: "]" got: "%s"`, next)
		case TokenCommandStart:
			command, err := parseCommand(ctx, reader)
			if err != nil {
				return nil, fmt.Errorf("invalid command: %w", err)
			}
			instructions = append(instructions, command...)
		case TokenBracketStart:
			bracketInstructions, err := parseBracket(ctx, reader)
			if err != nil {
				return nil, fmt.Errorf("invalid bracket section: %w", err)
			}
			instructions = append(instructions, bracketInstructions...)
		case TokenCharacter:
			instructions = append(instructions, Instruction{
				Opcode: OpPrint,
				Arg:    next.Value,
			})
			chars++
		default:
			return nil, fmt.Errorf("bracket: unknown token: %s", next.String())
		}
	}
}
