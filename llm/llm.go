package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type Response struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

const (
	OLLAMA_GENERATE_URL = "http://localhost:11434/api/generate"
)

// each message consists of one user message and one llm message
type message struct {
	userPrompt  string
	llmResponse string
}

func newMessage(prompt, response string) message {
	return message{
		userPrompt:  prompt,
		llmResponse: response,
	}
}

func (m message) String() string {
	return fmt.Sprintf(
		`<User Message>
        %s
        </User Message>
        <Assistant Response>
        %s
        </Assistant Response>
        `, m.userPrompt, m.llmResponse,
	)
}

type Session struct {
	model    string
	messages []message
}

func NewSession(model string) Session {
	return Session{
		model: model,
	}
}

func (s *Session) constructPrompt(prompt string) string {
	b := strings.Builder{}
	for _, msg := range s.messages {
		b.WriteString(msg.String())
	}
	b.WriteString(fmt.Sprintf("<User Prompt>\n%s\n</User Prompt>\n", prompt))
	return b.String()
}

func (s *Session) SendPrompt(prompt string) (string, error) {
	request := Request{
		Model:  s.model,
		Prompt: s.constructPrompt(prompt),
		Stream: false,
	}

	requestJson, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(
		OLLAMA_GENERATE_URL,
		"application/json",
		bytes.NewReader(requestJson),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var modelResponse Response
	err = json.NewDecoder(resp.Body).Decode(&modelResponse)
	if err != nil {
		return "", err
	}

	if !modelResponse.Done {
		return "", fmt.Errorf(
			"model response not done, only %q sent back",
			modelResponse.Response,
		)
	}

	s.messages = append(s.messages, newMessage(prompt, modelResponse.Response))

	return modelResponse.Response, nil
}
