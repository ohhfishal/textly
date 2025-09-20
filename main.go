package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	// "time"

	"github.com/alecthomas/kong"
	"github.com/ohhfishal/textly/compile"
)

type Cmd struct {
	Compile compile.Compile `cmd:"" default:"withargs" help:""`
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	if err := Run(ctx, os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func Run(ctx context.Context, stdout io.Writer, args []string) error {
	var exit bool
	var cmd Cmd
	parser, err := kong.New(
		&cmd,
		kong.Exit(func(_ int) { exit = true }),
		kong.BindTo(ctx, new(context.Context)),
		kong.BindTo(stdout, new(io.Writer)),
	)
	if err != nil {
		return err
	}

	parser.Stdout = stdout
	parser.Stderr = stdout

	context, err := parser.Parse(
		os.Args[1:],
	)
	if err != nil || exit {
		return err
	}

	if err := context.Run(); err != nil {
		return err
	}
	return nil
}
