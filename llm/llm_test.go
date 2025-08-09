package llm

import (
	"testing"
)

func TestConstructPrompt(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hey", "<User Prompt>\nHey\n</User Prompt>\n"},
		{
			"Hello, @text(World!)",
			"<User Prompt>\nHello, World!\n</User Prompt>\n",
		},
	}

	for _, tt := range tests {
		s := Session{}
		prompt, err := s.constructPrompt(tt.input)
		if err != nil {
			t.Errorf(
				"Failed to construct prompt with input %s, %v",
				tt.input,
				err,
			)
			continue
		}

		if prompt != tt.expected {
			t.Errorf("invalid prompt: got=%q, expected=%q", prompt, tt.expected)
		}
	}
}
