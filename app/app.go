package app

import (
	"fmt"
	"strings"
	//	"os"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/v1s10n-1/sai/models"
)

var useHighPerformanceRenderer = false

type Program struct {
	models   models.Models
	selected map[int]struct{}
	cursor   int

	prompt   textarea.Model
	err      error
	viewport viewport.Model
}

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "|"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()
	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "|"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

type errMsg error

func NewApp() Program {
	ta := textarea.New()
	ta.Placeholder = "twój prompt tutaj"
	ta.Focus()

	return Program{
		models:   models.LoadModels(),
		selected: make(map[int]struct{}),
		prompt:   ta,
		err:      nil,
	}
}

func (m Program) Init() tea.Cmd {
	return textarea.Blink
}

func (m Program) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.models.Models)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
		m.viewport.YPosition = headerHeight
		m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
		m.viewport.SetContent(m.View())

		// This is only necessary for high performance rendering, which in
		// most cases you won't need.
		//
		// Render the viewport one line below the header.
		m.viewport.YPosition = headerHeight + 1
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.prompt, cmd = m.prompt.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Program) headerView() string {
	title := titleStyle.Render("GREAT LLM")
	line := strings.Repeat("=", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Program) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("-", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m Program) View() string {
	s := "Co ode mnie chcesz?\n\n"

	for i, Model := range m.models.Models {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, Model.Name)
	}
	s += "\nPress q to quit. \n"
	return s
}
