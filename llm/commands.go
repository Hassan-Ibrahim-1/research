package llm

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/Hassan-Ibrahim-1/research/command"
)

func (s *Session) attachFile(
	cmd command.Command,
	prompt []byte,
) ([]byte, error) {
	if len(cmd.Arguments) == 0 {
		return prompt, nil
	}

	var fileContents bytes.Buffer
	for _, file := range cmd.Arguments {
		b, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		_, err = fileContents.WriteString(file + ":\n")
		if err != nil {
			return nil, err
		}

		_, err = fileContents.Write(b)
		if err != nil {
			return nil, err
		}

		err = fileContents.WriteByte('\n')
		if err != nil {
			return nil, err
		}
	}

	embedded := embed(fileContents.Bytes(), cmd.Loc, prompt)
	return embedded, nil
}

func (s *Session) attachLink(
	cmd command.Command,
	prompt []byte,
) ([]byte, error) {
	if len(cmd.Arguments) == 0 {
		return prompt, nil
	}

	var urlContents bytes.Buffer
	for _, url := range cmd.Arguments {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		_, err = urlContents.WriteString(url + "\n")
		if err != nil {
			return nil, err
		}

		_, err = urlContents.Write(b)
		if err != nil {
			return nil, err
		}

		err = urlContents.WriteByte('\n')
		if err != nil {
			return nil, err
		}

		resp.Body.Close()
	}

	embedded := embed(urlContents.Bytes(), cmd.Loc, prompt)
	return embedded, nil
}
