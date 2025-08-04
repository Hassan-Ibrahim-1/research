package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Hassan-Ibrahim-1/research/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f := enableLogs()
		defer func() {
			_ = f.Close()
		}()
		log.SetOutput(f)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetOutput(io.Discard)
	}

	content, err := os.ReadFile("content.md")
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(
		ui.NewModel(string(content)),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Err: %v", err)
		os.Exit(1)
	}
}

func enableLogs() io.WriteCloser {
	f, err := tea.LogToFile("debug.log", "")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}

	return f
}
