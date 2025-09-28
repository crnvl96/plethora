package ideas

import (
	"bytes"
	"strings"
	"testing"

	"github.com/crnvl96/plethora/internal/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

type mockBagelsGenerator struct {
	mock.Mock
}

func (m *mockBagelsGenerator) getSecretNum() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockBagelsGenerator) getClues(guess, secret string) string {
	args := m.Called(guess, secret)
	return args.String(0)
}

func TestRunWin(t *testing.T) {
	mockGen := &mockBagelsGenerator{}
	mockGen.On("getSecretNum").Return("123")
	mockGen.On("getClues", "123", "123").Return("You got it!")

	cmd := &bagelsExecCommand{
		generator: mockGen,
		stdin:     strings.NewReader(""),
		stdout:    &bytes.Buffer{},
		stderr:    &bytes.Buffer{},
	}

	var stdoutBuf bytes.Buffer

	cmd.SetStdin(strings.NewReader("123\nn\n"))
	cmd.SetStdout(&stdoutBuf)
	cmd.SetStderr(&bytes.Buffer{})

	err := cmd.Run()
	assert.NoError(t, err)

	output := stdoutBuf.String()
	assert.Contains(t, output, "You got it!")
	assert.Contains(t, output, "Thanks for playing")

	mockGen.AssertExpectations(t)
}

func TestRunOutOfGuesses(t *testing.T) {
	mockGen := &mockBagelsGenerator{}
	mockGen.On("getSecretNum").Return("123")
	mockGen.On("getClues", mock.Anything, mock.Anything).Return("Fermi")

	cmd := &bagelsExecCommand{
		generator: mockGen,
		stdin:     strings.NewReader(""),
		stdout:    &bytes.Buffer{},
		stderr:    &bytes.Buffer{},
	}

	var stdoutBuf bytes.Buffer

	cmd.SetStdin(strings.NewReader("456\n456\n456\n456\n456\n456\n456\n456\n456\n456\nn\n"))
	cmd.SetStdout(&stdoutBuf)
	cmd.SetStderr(&bytes.Buffer{})

	err := cmd.Run()
	assert.NoError(t, err)

	output := stdoutBuf.String()
	assert.Contains(t, output, "You ran out of guesses")
	assert.Contains(t, output, "The answer was 123")
	assert.Contains(t, output, "Thanks for playing")

	mockGen.AssertExpectations(t)
}

func TestRunError(t *testing.T) {
	mockGen := &mockBagelsGenerator{}
	mockGen.On("getSecretNum").Return("123")

	cmd := &bagelsExecCommand{
		generator: mockGen,
		stdin:     strings.NewReader(""),
		stdout:    &bytes.Buffer{},
		stderr:    &bytes.Buffer{},
	}

	// We don't set the SetStdin so that an error is triggered
	err := cmd.Run()
	assert.Error(t, err)
	assert.Equal(t, "failed to read input", err.Error())

	callback := newBagelsExecCallback()
	msg := callback.OnErr(err)
	expectedMsg := ui.DoneMsg{Err: err}
	assert.Equal(t, expectedMsg, msg)

	mockGen.AssertExpectations(t)
}
