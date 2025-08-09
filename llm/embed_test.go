package llm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Hassan-Ibrahim-1/research/command"
)

func TestEmbed(t *testing.T) {
	embedTag := "[embed]"
	// all 's' fields must have an [embed] tag
	tests := []struct {
		s        string
		data     string
		expected string
	}{
		{"hello, [embed]", "world", "hello, world"},
		{"foo [embed] baz", "bar", "foo bar baz"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			start := strings.Index(tt.s, embedTag)
			if start == -1 {
				t.Fatalf("No embed tag '[embed]' found in tt.s %q", tt.s)
			}
			rng := command.Range{
				Start: start,
				End:   start + len(embedTag),
			}

			s := string(embed([]byte(tt.s), rng, []byte(tt.data)))
			if s != tt.expected {
				t.Errorf(
					"bad embeded string. got=%q. expected=%q.",
					s,
					tt.expected,
				)
			}
		})
	}
}
