//go:build gopher

package main

import (
	"context"
	"os"
	"time"

	. "github.com/ohhfishal/gopher/runtime"
)

// Devel builds the gopher binary then runs it
func Devel(ctx context.Context, args RunArgs) error {
	return Run(ctx, NowAnd(OnFileChange(1*time.Second, ".go")),
		&Printer{},
		&GoBuild{
			Output: "target/dev",
		},
		&GoFormat{},
		&GoTest{},
		&GoVet{},
		&GoModTidy{},
		ExecCommand("echo", "---"),
		ExecCommand("echo", "DEVEL OK"),
	)
}

// cicd runs the entire ci/cd suite
func CICD(ctx context.Context, args RunArgs) error {
	return Run(ctx, Now(),
		&Printer{},
		&GoBuild{
			Output: "target/cicd",
		},
		&GoFormat{
			CheckOnly: true,
		},
		&GoTest{},
		&GoVet{},
		ExecCommand("echo", "CICD OK"),
	)
}

// Removes all local build artifacts.
func Clean(ctx context.Context, args RunArgs) error {
	return os.RemoveAll("target")
}
