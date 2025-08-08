package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Hassan-Ibrahim-1/research/llm"
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

	s := llm.NewSession("mistral")

	m := ui.New(&s)

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Err: %v", err)
		os.Exit(1)
	}

	// prompt := "here is some code: pub fn main() !void {return error.Oops;}"
	// response, err := s.SendPrompt(prompt)
	// if err != nil {
	// 	fmt.Println("Prompt failed:", err)
	// 	return
	// }
	//
	// fmt.Printf("user: %s\nmistral: ", prompt)
	// for str := range response {
	// 	fmt.Print(str)
	// }
	// fmt.Print("\n")
	//
	// prompt = "what language was the code i gave you written in?"
	// response, err = s.SendPrompt(prompt)
	// if err != nil {
	// 	fmt.Println("Prompt failed:", err)
	// 	return
	// }
	//
	// fmt.Printf("user: %s\nmistral: ", prompt)
	// for str := range response {
	// 	fmt.Print(str)
	// }
	// fmt.Print("\n")
}

func enableLogs() io.WriteCloser {
	f, err := tea.LogToFile("debug.log", "")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}

	return f
}
