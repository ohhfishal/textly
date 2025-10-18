package compile_test

import (
	"fmt"
	"github.com/ohhfishal/textly/compile"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestValid(t *testing.T) {
	tests := []struct {
		Input  string
		Output string
	}{
		{
			Input:  "Hell[l]o",
			Output: "Hello",
		},
		{
			Input:  "\\# Test",
			Output: "# Test",
		},
		{
			Input:  "Hello \\\nWorld",
			Output: "Hello World",
		},
	}

	dir := t.TempDir()
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert := assert.New(t)

			f, err := os.CreateTemp(dir, fmt.Sprintf("test_%d.txt", i))
			defer f.Close()
			assert.Nil(err)

			_, err = f.WriteString(test.Input)
			assert.Nil(err)
			f.Close()

			f, err = os.Open(f.Name())
			assert.Nil(err)
			defer f.Close()

			defer os.Remove(f.Name())
			cmd := compile.Compile{
				Input:    f,
				Optimize: false,
				RunOptions: compile.RunOptions{
					Delay: 0,
					Beat:  0,
				},
			}

			var output MockTerminal
			err = cmd.Run(t.Context(), &output)
			assert.Nil(err)

			assert.Equal(
				test.Output,
				output.String(),
			)
		})
	}
}

type MockTerminal struct {
	buffer []rune
	cursor int
}

func (terminal *MockTerminal) Write(p []byte) (n int, err error) {
	for _, b := range p {
		r := rune(b)

		if r == '\b' {
			if terminal.cursor > 0 {
				terminal.cursor--
			}
		} else {
			for terminal.cursor >= len(terminal.buffer) {
				terminal.buffer = append(terminal.buffer, ' ')
			}
			terminal.buffer[terminal.cursor] = r
			terminal.cursor++
		}
	}

	return len(p), nil
}

func (terminal *MockTerminal) String() string {
	end := len(terminal.buffer)
	for end > 0 && terminal.buffer[end-1] == ' ' {
		end--
	}
	return string(terminal.buffer[:end])
}
