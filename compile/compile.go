package compile

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Compile struct {
	Input *os.File      `arg:""`
	Delay time.Duration `default:"0.1s"`
}

func (cmd *Compile) Run(ctx context.Context, stdout io.Writer) error {
	var wg sync.WaitGroup

	errs := make(chan error)
	tokens := make(chan Token, 10)
	waitChan := make(chan any, 1)

	// Lexer
	wg.Go(func() {
		defer cmd.Input.Close()
		defer close(tokens)
		reader := bufio.NewReader(cmd.Input)
		if err := Lex(ctx, reader, tokens); err != nil {
			errs <- err
		}
	})

	// Parser
	wg.Go(func() {
		instructions, err := Parse(ctx, tokens)
		if err != nil {
			errs <- err
			return
		}
		instructions = Optimize(instructions)

		// TODO: Wrap in function and enable via option
		for _, instruction := range instructions {
			switch instruction.Opcode {
			case OpPrint:
				for _, char := range instruction.Arg.(string) {
					fmt.Fprintf(stdout, "%c", char)
					time.Sleep(cmd.Delay)
				}
			case OpDelete:
				for range instruction.Arg.(int) {
					fmt.Fprint(stdout, "\b \b")
					time.Sleep(cmd.Delay)
				}
			}
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
