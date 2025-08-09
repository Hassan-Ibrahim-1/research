package embed

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"

	"github.com/Hassan-Ibrahim-1/research/command"
)

func Embed(str []byte, rng command.Range, data []byte) []byte {
	s := slices.Delete(str, rng.Start, rng.End)
	s = slices.Insert(s, rng.Start, data...)
	return s
}

// gets data from the url via http.Get
func URL(str []byte, rng command.Range, url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %w", err)
	}

	return Embed(str, rng, body), nil
}

func File(str []byte, rng command.Range, path string) ([]byte, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from file %s: %w", path, err)
	}
	return Embed(str, rng, contents), nil
}
