package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
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
	model string

	mu       sync.Mutex
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

func (s *Session) SendPrompt(prompt string) (<-chan string, error) {
	request := Request{
		Model:  s.model,
		Prompt: s.constructPrompt(prompt),
		Stream: true,
	}

	requestJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		OLLAMA_GENERATE_URL,
		"application/json",
		bytes.NewReader(requestJson),
	)
	if err != nil {
		return nil, err
	}

	log.Println("got a response", resp)

	ch := make(chan string)

	go func() {
		defer resp.Body.Close()
		defer close(ch)

		var fullResponse strings.Builder

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Bytes()
			if line == nil {
				continue
			}

			var partialResponse Response
			err = json.Unmarshal(line, &partialResponse)
			if err != nil {
				log.Println("Session err:", err)
				continue
			}

			if resp := partialResponse.Response; resp != "" {
				fullResponse.WriteString(resp)
				log.Println("partial response:", resp)
				ch <- resp
			}

			if partialResponse.Done {
				break
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println("error reading streaming response:", err)
		} else {
			s.addMessage(newMessage(prompt, fullResponse.String()))
		}
	}()

	return ch, nil
}

func (s *Session) addMessage(msg message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, msg)
}
