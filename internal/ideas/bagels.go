package ideas

import (
	"bufio"
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/crnvl96/plethora/internal/ui"
)

type bagels struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (b *bagels) SetStdin(r io.Reader)  { b.stdin = r }
func (b *bagels) SetStdout(w io.Writer) { b.stdout = w }
func (b *bagels) SetStderr(w io.Writer) { b.stderr = w }

func (b *bagels) Run() error {
	scanner := bufio.NewScanner(b.stdin)

	fmt.Fprintln(b.stdout, "Bagels!")
	fmt.Fprint(b.stdout, "Enter your name: ")

	if scanner.Scan() {
		name := scanner.Text()
		fmt.Fprintf(b.stdout, "Hello, %s!\n", name)
	}

	return nil
}

func onErr(err error) tea.Msg {
	return ui.ExecDoneMsg{Err: err}
}

func init() {
	Ideas["bagels"] = ui.Item{
		Header:   "Bagels",
		Body:     "Bagels is a deductive logic game. You must guess a secret three-digit number based on clues.",
		Command:  &bagels{},
		Callback: onErr,
	}
}
