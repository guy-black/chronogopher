package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	dt time.Time
}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func initialModel() model {
	return model {
		dt: time.Now(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Sequence(
		tea.SetWindowTitle("chronogopher"),
		doTick(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
				case "q":
					return m, tea.Quit
			}
		case TickMsg:
			m.dt = time.Now()
			return m, doTick()
	}
	return m, nil
}

func (m model) View() string {
	return m.dt.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("oooopsies: %v", err)
		os.Exit(1)
	}
}
