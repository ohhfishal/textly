package compile

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

// Prevent infinte loops. If you ever get the panic contact devs
const _MAX_OPTIMIZATIONS = 10

type Opcode string

const (
	OpPrint     = "print"     // print(content str)
	OpDelete    = "delete"    // delete(count int) // Number of characters to backspace
	OpSleep     = "sleep"     // sleep(seconds int)
	OpClear     = "clear"     // clear()
	OpPushColor = "pushColor" // pushColor(color string)
	OpPopColor  = "popColor"  // popColor()
)

const (
	ClearANSI = "\033[H\033[2J"
)

const (
	Red   = "\033[31m"
	Reset = "\033[0m"
)

type Instruction struct {
	Opcode Opcode `json:"opcode"`
	Arg    any    `json:"arg"`
}

func (instruction Instruction) String() string {
	arg := instruction.Arg
	if arg == nil {
		arg = ""
	}

	return strings.ReplaceAll(
		fmt.Sprintf(
			"%s(%v)",
			instruction.Opcode,
			arg,
		),
		"\n",
		"\\n",
	)
}

type Program struct {
	Instructions []Instruction
}

type RunOptions struct {
	List  bool          `short:"l" help:"Output a single word per line."`
	Delay time.Duration `default:"0.1s"`
	Beat  time.Duration `default:"1s"`
}

func (program Program) Run(stdout io.Writer, options RunOptions) error {
	colors := []string{Reset}
	for _, instruction := range program.Instructions {
		switch instruction.Opcode {
		case OpPrint:
			for _, char := range instruction.Arg.(string) {
				if options.List && char == ' ' {
					char = '\n'
				}
				fmt.Fprintf(stdout, "%c", char)
				time.Sleep(options.Delay)
			}
		case OpDelete:
			for range instruction.Arg.(int) {
				fmt.Fprint(stdout, "\b \b")
				time.Sleep(options.Delay)
			}
		case OpSleep:
			for range instruction.Arg.(int) {
				time.Sleep(options.Beat)
			}
		case OpClear:
			fmt.Fprintf(stdout, ClearANSI)
		case OpPushColor:
			color := instruction.Arg.(string)
			fmt.Fprint(stdout, color)
			colors = append(colors, color)
		case OpPopColor:
			if len(colors) == 0 {
				return errors.New("call to popColor when the stack is empty")
			}
			reset := colors[len(colors)-2]
			colors = colors[:len(colors)-1]
			fmt.Fprint(stdout, reset)

		default:
			return fmt.Errorf("unknown op: %s", instruction.Opcode)
		}

	}
	return nil
}

type OptimizeOptions struct {
	Flatten bool `help:"Premptively delete before printing to stdout."`
}

func (program *Program) Optimize(options OptimizeOptions) {
	instructions, cont := optimize(program.Instructions, options)
	i := 0
	for cont {
		if i >= _MAX_OPTIMIZATIONS {
			panic(fmt.Errorf("stuck in an infinite loop optimizing: %v", instructions))
		}
		instructions, cont = optimize(instructions, options)
		i++
	}
	program.Instructions = instructions
}

func optimize(original []Instruction, opts OptimizeOptions) ([]Instruction, bool) {
	var instructions []Instruction
	if original == nil || len(original) == 0 {
		return []Instruction{}, false
	}

	cur := &original[0]
	for _, next := range original[1:] {
		switch {
		case cur.Opcode == OpSleep && cur.Opcode == next.Opcode:
			cur.Arg = cur.Arg.(int) + next.Arg.(int)
		case cur.Opcode == OpPrint && cur.Opcode == next.Opcode:
			cur.Arg = cur.Arg.(string) + next.Arg.(string)
		case opts.Flatten && cur.Opcode == OpPrint && next.Opcode == OpDelete:
			arg := cur.Arg.(string)
			cur.Arg = arg[:len(arg)-next.Arg.(int)]
		default:
			instructions = append(instructions, *cur)
			cur = &next
		}
	}

	if cur != nil {
		instructions = append(instructions, *cur)
	}
	return instructions, len(instructions) != len(original)

}
