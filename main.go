package main

import (
	//	"fmt"
	"log"

	//	"github.com/charmbracelet/bubbles/textarea"
	//"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/v1s10n-1/sai/app"
	//	"github.com/v1s10n-1/sai/common"
)

func main() {
	//	model, err = viewport.NewExample()
	p := tea.NewProgram(app.NewApp(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
