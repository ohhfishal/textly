package compile

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

type Compile struct {
	Input           *os.File        `arg:""`
	Lex             bool            `short:"L" help:"Only run the lexer and print all tokens to standard out."`
	Dump            bool            `short:"D" help:"Print all instructions to standard out then return."`
	Optimize        bool            `negatable:"" default:"true" help:"Enable optimizations (default: enabled)"`
	OptimizeOptions OptimizeOptions `embed:""`
	RunOptions      RunOptions      `embed:""`
}

func (cmd *Compile) Run(ctx context.Context, stdout io.Writer) error {
	var wg sync.WaitGroup

	errs := make(chan error)
	tokens := make(chan Token, 10)
	waitChan := make(chan any, 1)

	// Lexer
	wg.Go(func() {
		defer cmd.Input.Close() //nolint:errcheck
		defer close(tokens)
		reader := bufio.NewReader(cmd.Input)
		if err := Lex(ctx, reader, tokens); err != nil {
			errs <- err
		}
	})

	// Parser
	wg.Go(func() {
		if cmd.Lex {
			for token := range tokens {
				fmt.Println(stdout, token)
			}
			return
		}
		program, err := Parse(ctx, tokens)
		if err != nil {
			errs <- err
			return
		}

		// TODO: Optimize while reading from the channel?
		if cmd.Optimize {
			program.Optimize(cmd.OptimizeOptions)
		}

		if cmd.Dump {
			for i, instruction := range program.Instructions {
				if _, err := fmt.Fprintf(stdout, "%3d: %s\n", i, instruction); err != nil {
					errs <- err
				}
			}
			return
		}

		if err := program.Run(stdout, cmd.RunOptions); err != nil {
			errs <- err
			return
		}
	})

	go func() {
		wg.Wait()
		waitChan <- nil
		close(waitChan)
	}()

	select {
	case <-ctx.Done():
		return errors.New("context closed")
	case err := <-errs:
		return err
	case <-waitChan:
	}

	return nil
}
