package ideas

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestRun(t *testing.T) {
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
