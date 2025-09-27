package ideas

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSecretNum(t *testing.T) {
	g := newDefaultBagelsGenerator()
	secret := g.getSecretNum()
	numDigitsSecret := len(secret)

	err := "Secret number should have the correct number of digits"
	assert.Equal(t, numDigits, numDigitsSecret, err)
}

func TestGetClues(t *testing.T) {
	g := newDefaultBagelsGenerator()

	testCases := []struct {
		guess    string
		secret   string
		expected string
	}{
		{"123", "123", "You got it!"},
		{"456", "123", "Bagels"},
		{"124", "123", "Fermi Fermi"},
		{"231", "123", "Pico Pico Pico"},
		{"132", "123", "Fermi Pico Pico"},
		{"156", "123", "Fermi"},
		{"456", "145", "Pico Pico"},
	}

	for _, tc := range testCases {
		clues := g.getClues(tc.guess, tc.secret)
		assert.Equal(t, tc.expected, clues)
	}
}

type testBagelsGenerator struct{}

func (t *testBagelsGenerator) getSecretNum() string {
	return "123"
}

func (t *testBagelsGenerator) getClues(guess, secret string) string {
	d := &defaultBagelsGenerator{}
	return d.getClues(guess, secret)
}

func TestSetStreams(t *testing.T) {
	cmd := &bagelsExecCommand{
		generator: &testBagelsGenerator{},
		stdin:     strings.NewReader(""),
		stdout:    &bytes.Buffer{},
		stderr:    &bytes.Buffer{},
	}

	var stdoutBuf bytes.Buffer
	input := "123\nn\n"
	cmd.SetStdin(strings.NewReader(input))
	cmd.SetStdout(&stdoutBuf)
	cmd.SetStderr(&bytes.Buffer{}) // not testing stderr here

	err := cmd.Run()
	assert.NoError(t, err)

	output := stdoutBuf.String()
	assert.Contains(t, output, "You got it!")
	assert.Contains(t, output, "Thanks for playing")
}
