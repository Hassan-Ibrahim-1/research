package llm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestEmbedUrl(t *testing.T) {
	serverResponse := "<p>Hello, World</p>"

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(serverResponse))
		}),
	)
	defer server.Close()

	embedTag := "[embed]"
	str := "test " + embedTag
	start := strings.Index(str, embedTag)
	rng := command.Range{
		Start: start,
		End:   start + len(embedTag),
	}

	embedded, err := embedURL([]byte(str), rng, server.URL)
	if err != nil {
		t.Errorf("Failed to embed: %s", err)
	}

	expected := "test " + serverResponse
	if s := string(embedded); s != expected {
		t.Errorf(
			"embedded string not equal to expected. got=%q. expected=%q",
			s,
			expected,
		)
	}
}
