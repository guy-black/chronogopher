package main

import (
	"fmt"
	"os"
	"time"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var appStyle = lipgloss.NewStyle().Padding(1)

type model struct {
	dt time.Time
}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
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
		tea.EnterAltScreen,
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
	hour, min, sec := m.dt.Clock()
// some time later I'll make this conditional on a 24/12 hour clock option
// but I like 12 hour time so it's just like this for now
	if hour>12 {
		hour -= 12
	}
// I convert the time to a hh:mm:ss string to make it easier
// to convert into big text with box drawing unicode chars
	time := stringTime(hour, min, sec)

// creating the clock section
	clock := fmt.Sprintf("%s\n%s\n%s",
		timeTopLine(time),
		timeMidLine(time),
		timeBotLine(time))

// creating the date
	year, month, day := m.dt.Date()
	mdy := fmt.Sprintf("%s %d, %d",
		month.String(),
		day,
		year,)

// here's the actual view to be rendered
		return appStyle.Render(lipgloss.JoinVertical(0.5,
			clock,
			mdy,))
}

func timeTopLine (time string) string {
	time = strings.ReplaceAll(time, "1", OneTop)
	time = strings.ReplaceAll(time, "2", TwoTop)
	time = strings.ReplaceAll(time, "3", ThrTop)
	time = strings.ReplaceAll(time, "4", FouTop)
	time = strings.ReplaceAll(time, "5", FivTop)
	time = strings.ReplaceAll(time, "6", SixTop)
	time = strings.ReplaceAll(time, "7", SevTop)
	time = strings.ReplaceAll(time, "8", EigTop)
	time = strings.ReplaceAll(time, "9", NinTop)
	time = strings.ReplaceAll(time, "0", ZerTop)
	time = strings.ReplaceAll(time, ":", ColTop)
	return time
}

func timeMidLine (time string) string {
	time = strings.ReplaceAll(time, "1", OneMid)
	time = strings.ReplaceAll(time, "2", TwoMid)
	time = strings.ReplaceAll(time, "3", ThrMid)
	time = strings.ReplaceAll(time, "4", FouMid)
	time = strings.ReplaceAll(time, "5", FivMid)
	time = strings.ReplaceAll(time, "6", SixMid)
	time = strings.ReplaceAll(time, "7", SevMid)
	time = strings.ReplaceAll(time, "8", EigMid)
	time = strings.ReplaceAll(time, "9", NinMid)
	time = strings.ReplaceAll(time, "0", ZerMid)
	time = strings.ReplaceAll(time, ":", ColMid)
	return time
}

func timeBotLine (time string) string {
	time = strings.ReplaceAll(time, "1", OneBot)
	time = strings.ReplaceAll(time, "2", TwoBot)
	time = strings.ReplaceAll(time, "3", ThrBot)
	time = strings.ReplaceAll(time, "4", FouBot)
	time = strings.ReplaceAll(time, "5", FivBot)
	time = strings.ReplaceAll(time, "6", SixBot)
	time = strings.ReplaceAll(time, "7", SevBot)
	time = strings.ReplaceAll(time, "8", EigBot)
	time = strings.ReplaceAll(time, "9", NinBot)
	time = strings.ReplaceAll(time, "0", ZerBot)
	time = strings.ReplaceAll(time, ":", ColBot)
	return time
}

func stringTime (hour, min, sec int) string {
	var h,m,s string
	if hour<10 {
		h = fmt.Sprintf ("0%d", hour)
	} else {
		h = fmt.Sprint (hour)
	}
	if min<10 {
		m = fmt.Sprintf ("0%d", min)
	} else {
		m = fmt.Sprint (min)
	}
	if sec<10 {
		s = fmt.Sprintf ("0%d", sec)
	} else {
		s = fmt.Sprint (sec)
	}
	return fmt.Sprint(h,":",m,":",s)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("oooopsies: %v", err)
		os.Exit(1)
	}
}

// constants

const (
	OneTop = " ╻ "
	OneMid = " ┃ "
	OneBot = " ╹ "
	TwoTop = "╺━┓"
	TwoMid = "┏━┛"
	TwoBot = "┗━╸"
	ThrTop = "╺━┓"
	ThrMid = "╺━┫"
	ThrBot = "╺━┛"
	FouTop = "╻ ╻"
	FouMid = "┗━┫"
	FouBot = "  ╹"
	FivTop = "┏━╸"
	FivMid = "┗━┓"
	FivBot = "╺━┛"
	SixTop = "┏━╸"
	SixMid = "┣━┓"
	SixBot = "┗━┛"
	SevTop = "╺━┓"
	SevMid = "  ┃"
	SevBot = "  ╹"
	EigTop = "┏━┓"
	EigMid = "┣━┫"
	EigBot = "┗━┛"
	NinTop = "┏━┓"
	NinMid = "┗━┫"
	NinBot = "╺━┛"
	ZerTop = "┏━┓"
	ZerMid = "┃ ┃"
	ZerBot = "┗━┛"
	ColTop = "╻"
	ColMid = " "
	ColBot = "╹"
)
