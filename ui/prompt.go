package ui

import "github.com/charmbracelet/bubbles/viewport"

type PromptArea struct {
	viewport       viewport.Model
	content        string
	characterLimit int
	focused        bool
}
