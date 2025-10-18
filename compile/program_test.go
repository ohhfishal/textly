package compile_test

import (
	"bytes"
	"testing"
	"time"

	. "github.com/ohhfishal/textly/compile"
	"github.com/stretchr/testify/assert"
)

func TestProgramRun(t *testing.T) {
	tests := []struct {
		name     string
		program  Program
		options  RunOptions
		expected string
	}{
		{
			name: "simple print",
			program: Program{
				Instructions: []Instruction{
					{Opcode: OpPrint, Arg: "hello"},
				},
			},
			options:  RunOptions{Delay: 0},
			expected: "hello",
		},
		{
			name: "multiple prints",
			program: Program{
				Instructions: []Instruction{
					{Opcode: OpPrint, Arg: "hello"},
					{Opcode: OpPrint, Arg: " world"},
				},
			},
			options:  RunOptions{Delay: 0},
			expected: "hello world",
		},
		{
			name: "print with delete",
			program: Program{
				Instructions: []Instruction{
					{Opcode: OpPrint, Arg: "hello"},
					{Opcode: OpDelete, Arg: 2},
				},
			},
			options:  RunOptions{Delay: 0},
			expected: "hello\b \b\b \b",
		},
		{
			name: "list mode converts spaces to newlines",
			program: Program{
				Instructions: []Instruction{
					{Opcode: OpPrint, Arg: "hello world test"},
				},
			},
			options:  RunOptions{List: true, Delay: 0},
			expected: "hello\nworld\ntest",
		},
		{
			name: "zero case",
			program: Program{
				Instructions: []Instruction{},
			},
			options:  RunOptions{Delay: 0},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tt.program.Run(&buf, tt.options)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}

func TestOptimize(t *testing.T) {
	tests := []struct {
		name     string
		program  Program
		options  OptimizeOptions
		expected []Instruction
	}{
		{
			name: "combine consecutive prints",
			program: Program{
				Instructions: []Instruction{
					{Opcode: OpPrint, Arg: "a"},
					{Opcode: OpPrint, Arg: "b"},
					{Opcode: OpPrint, Arg: "c"},
				},
			},
			options: OptimizeOptions{},
			expected: []Instruction{
				{Opcode: OpPrint, Arg: "abc"},
			},
		},
		{
			name: "flatten print followed by delete when enabled",
			program: Program{
				Instructions: []Instruction{
					{Opcode: OpPrint, Arg: "hello"},
					{Opcode: OpDelete, Arg: 2},
				},
			},
			options: OptimizeOptions{Flatten: true},
			expected: []Instruction{
				{Opcode: OpPrint, Arg: "hel"},
			},
		},
		{
			name: "flatten allows more prints to be combined",
			program: Program{
				Instructions: []Instruction{
					{Opcode: OpPrint, Arg: "testing"},
					{Opcode: OpDelete, Arg: 3},
					{Opcode: OpPrint, Arg: "ed"},
				},
			},
			options: OptimizeOptions{Flatten: true},
			expected: []Instruction{
				{Opcode: OpPrint, Arg: "tested"},
			},
		},
		{
			name: "disabled flatten",
			program: Program{
				Instructions: []Instruction{
					{Opcode: OpPrint, Arg: "hello"},
					{Opcode: OpDelete, Arg: 2},
				},
			},
			options: OptimizeOptions{Flatten: false},
			expected: []Instruction{
				{Opcode: OpPrint, Arg: "hello"},
				{Opcode: OpDelete, Arg: 2},
			},
		},
		{
			name: "zero case",
			program: Program{
				Instructions: []Instruction{},
			},
			options:  OptimizeOptions{},
			expected: []Instruction{},
		},
		{
			name: "complex case",
			program: Program{
				Instructions: []Instruction{
					{Opcode: OpPrint, Arg: "hel"},
					{Opcode: OpPrint, Arg: "lo"},
					{Opcode: OpDelete, Arg: 1},
					{Opcode: OpPrint, Arg: "p"},
					{Opcode: OpDelete, Arg: 1},
					{Opcode: OpPrint, Arg: "o"},
				},
			},
			options: OptimizeOptions{Flatten: true},
			expected: []Instruction{
				{Opcode: OpPrint, Arg: "hello"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.program.Optimize(tt.options)
			assert.Equal(t, len(tt.expected), len(tt.program.Instructions), "wrong number of instructions")
			for i, instr := range tt.program.Instructions {
				assert.Equal(t, tt.expected[i].Opcode, instr.Opcode, "instruction %d opcode", i)
				assert.Equal(t, tt.expected[i].Arg, instr.Arg, "instruction %d arg", i)
			}
		})
	}
}

func TestProgramRunTiming(t *testing.T) {
	tests := []struct {
		name        string
		program     Program
		options     RunOptions
		minDuration time.Duration
	}{
		{
			name: "delay between characters",
			program: Program{
				Instructions: []Instruction{
					{Opcode: OpPrint, Arg: "abc"},
				},
			},
			options:     RunOptions{Delay: 10 * time.Millisecond},
			minDuration: 30 * time.Millisecond,
		},
		{
			name: "beat for sleep",
			program: Program{
				Instructions: []Instruction{
					{Opcode: OpSleep, Arg: 2},
				},
			},
			options:     RunOptions{Beat: 10 * time.Millisecond},
			minDuration: 20 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			start := time.Now()
			tt.program.Run(&buf, tt.options)
			duration := time.Since(start)

			assert.GreaterOrEqual(t, duration, tt.minDuration)
		})
	}
}
