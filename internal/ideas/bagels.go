package ideas

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/crnvl96/plethora/internal/ui"
)

const (
	numDigits  = 3
	maxGuesses = 10
)

type bagelsGenerator interface {
	getSecretNum() string
	getClues(guess, secret string) string
}

type defaultBagelsGenerator struct{}

func (d *defaultBagelsGenerator) getSecretNum() string {
	numbers := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})
	return strings.Join(numbers[:numDigits], "")
}

func (d *defaultBagelsGenerator) getClues(guess, secret string) string {
	if guess == secret {
		return "You got it!"
	}

	clues := []string{}
	tempSecret := []rune(secret)

	for i := 0; i < len(guess); i++ {
		if guess[i] == secret[i] {
			clues = append(clues, "Fermi")
			tempSecret[i] = ' '
		}
	}

	for i := 0; i < len(guess); i++ {
		if guess[i] != secret[i] {
			for j, s := range tempSecret {
				if s == rune(guess[i]) {
					clues = append(clues, "Pico")
					tempSecret[j] = ' '
					break
				}
			}
		}
	}

	if len(clues) == 0 {
		return "Bagels"
	} else {
		sort.Strings(clues)
		return strings.Join(clues, " ")
	}
}

func newDefaultBagelsGenerator() *defaultBagelsGenerator {
	return &defaultBagelsGenerator{}
}

type bagelsExecCommand struct {
	generator bagelsGenerator
	stdin     io.Reader
	stdout    io.Writer
	stderr    io.Writer
}

func (b *bagelsExecCommand) SetStderr(w io.Writer) { b.stderr = w }
func (b *bagelsExecCommand) SetStdout(w io.Writer) { b.stdout = w }
func (b *bagelsExecCommand) SetStdin(r io.Reader)  { b.stdin = r }
func (b *bagelsExecCommand) Run() error {
	scanner := bufio.NewScanner(b.stdin)
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		fmt.Fprintln(b.stderr, "\nInterrupted. Thanks for playing!")
		os.Exit(0)
	}()

	fmt.Fprintf(b.stdout, `Bagels, a deductive login game.
	By Al Sweigart al@inventwithpython.com

	I am thinking of a %d-digit number with no repeated digits.
	Try to guess what it is. Here are some clues:

	When I say:     That means:
	Pico			One digit is correct but in the wrong position
	Fermi			One digit is correct and in the right position
	Bagels			No digit is correct

	For example, if the secret number was 248 and your guess was 843, the clues would be Fermi Pico.
	`, numDigits)

	for {
		secretNum := b.generator.getSecretNum()
		fmt.Fprintln(b.stdout, "I have thought up a number")
		fmt.Fprintf(b.stdout, "You have %d guesses to get it\n", maxGuesses)

		var numGuesses int
		for numGuesses = 1; numGuesses <= maxGuesses; numGuesses++ {
			guess := ""
			for {
				fmt.Fprintf(b.stdout, "Guess #%d: ", numGuesses)
				if !scanner.Scan() {
					return errors.New("failed to read input")
				}

				guess = strings.TrimSpace(scanner.Text())
				if len(guess) == numDigits {
					if _, err := strconv.Atoi(guess); err == nil {
						break
					}
				}
			}

			clues := b.generator.getClues(guess, secretNum)
			fmt.Fprintln(b.stdout, clues)
			if guess == secretNum {
				break
			}
		}

		if numGuesses > maxGuesses {
			fmt.Fprintln(b.stdout, "You ran out of guesses")
			fmt.Fprintf(b.stdout, "The answer was %s\n", secretNum)
		}

		fmt.Fprintln(b.stdout, "Want to play again? (y/n)")
		if !scanner.Scan() {
			break
		}

		answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if !strings.HasPrefix(answer, "y") {
			break
		}
	}

	fmt.Fprintln(b.stdout, "Thanks for playing")
	return nil
}

func newBagelsExecCommand() *bagelsExecCommand {
	g := newDefaultBagelsGenerator()
	return &bagelsExecCommand{
		generator: g,
		stdin:     os.Stdin,
		stdout:    os.Stdout,
		stderr:    os.Stderr,
	}
}

type bagelsExecCallback struct{}

func (b *bagelsExecCallback) OnErr(err error) tea.Msg { return ui.DoneMsg{Err: err} }

func newBagelsExecCallback() *bagelsExecCallback {
	return &bagelsExecCallback{}
}

func init() {
	bcmd := newBagelsExecCommand()
	bcb := newBagelsExecCallback().OnErr

	Ideas["bagels"] = ui.Item{
		Component: ui.Component{
			Title:       "Bagels",
			Description: "Bagels is a deductive logic game. You must guess a secret three-digit number based on clues.",
		},
		Command:  bcmd,
		Callback: bcb,
	}
}
