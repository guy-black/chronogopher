package main

import (
	"fmt"
	"os"
	"time"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// styles
var (
	// style for whole app
	appStyle = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(lipgloss.Color("2"))
	// clock
	// calendar
	calCurrDay = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	othMonthDay = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
)

func clockStyle(m model) lipgloss.Style {
	if m.sel == ClockSect {
		return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Foreground(lipgloss.Color("5"))
	}
	return lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color("5"))
}

func calStyle(m model) lipgloss.Style {
	if m.sel == CalSect {
		return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Foreground(lipgloss.Color("3"))
	}
	return lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color("3"))
}

func todoStyle(m model) lipgloss.Style {
	if m.sel == TodoSect {
		return lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	}
	return lipgloss.NewStyle().Padding(1)
}

type Section byte

const (
	 ClockSect = iota
	 CalSect
	 TodoSect
)


func incSel (s Section) Section{
	if s == 2 {
		return 0
	}
	s++
	return s
}

func decSel (s Section) Section{
	if s == 0 {
		return 2
	}
	s--
	return s
}

type model struct {
	dt  time.Time
	td  []string
	sel Section
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
		td: []string {
			"finish chronogopher",
			"work on teh",
			},
		sel: 2,
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
				case "tab":
					m.sel = incSel (m.sel)
					return m, nil
				case "shift+tab":
					m.sel = decSel (m.sel)
					return m, nil
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
	clock := clockStyle(m).Render(fmt.Sprintf("%s\n%s\n%s",
		timeTopLine(time),
		timeMidLine(time),
		timeBotLine(time)))
// creating the date
	year, month, day := m.dt.Date()
	wd := m.dt.Weekday()
	mdy := fmt.Sprintf("%s, %s %d, %d",
		wd.String(),
		month.String(),
		day,
		year,)
// creating the calendar
	calDays := genCalDays (true, month, day, wd)
	cal := calStyle(m).Render(lipgloss.JoinVertical(0.5,
		// TODO: make this reference a value from m when I implement moving months
		fmt.Sprintf("%s  %d\n", month.String(), year),
		"Sun Mon Tue Wed Thu Fri Sat",
		calDays,))
// creating the todo
	tdl := ""
	for _, t := range m.td {
		tdl += fmt.Sprint(t, "\n")
	}
	fntdl, _ := strings.CutSuffix(tdl, "\n")
	todo := todoStyle(m).Render(lipgloss.JoinVertical(.5,
		"todo",
		fntdl))
// here's the actual view to be rendered
		return appStyle.Render(lipgloss.JoinVertical(0.5,
			clock,
			mdy,
			cal,
			todo))
}

func genCalDays(leap bool, mon time.Month, day int, wd time.Weekday) string {
// first figure how many last days of last month to render
// get first day of month with this weekday as fwdm
	fwdm := day
	for fwdm > 7 {
		fwdm -= 7
	}
	// fwdmdeb := fwdm
// get weekday of first of the month as wdfm
	wdfm := wd
	// as long as fwdm > 1 decrement fwdm and go back one weekday from wdfm
	for fwdm > 1 {
		wdfm = backOneWeekday(wdfm)
		fwdm -= 1
	} // wdfm should now be the weekday for the first of the month
	// thus wdfm is also the number of days of the last month to show
	finStr := ""
	daysInPrevMon := daysInMonth(leap, backOneMonth(mon))
	daysToPrint := 42 // 6 weeks * 7 days
	for wdfm > 0 {
		finStr = fmt.Sprint(daysInPrevMon, " ") + finStr
		daysInPrevMon--
		wdfm--
		daysToPrint--
	}
	daysThisMon := daysInMonth(leap, mon)
	for i:=1; i<=daysThisMon; i++ {
		dayNum := ""
		if i==day{
			if i<10 {
				dayNum = fmt.Sprint("  ", i, " ")
			} else {
				dayNum = fmt.Sprint(" ", i, " ")
			}
		} else {
			if i<10 {
				dayNum = fmt.Sprint("  ", i, " ")
			} else {
				dayNum = fmt.Sprint(" ", i, " ")
			}
		}
		finStr += dayNum
		daysToPrint--
		if daysToPrint%7 == 0 && daysToPrint != 0 {
			finStr += "\n"
		}
	}
	// okay now just adding remaining days from next month
	// create new var dltp equal to daysToPrint so I can still keep track of if a
	// newline is needed by days left to print while counting up days to print
	dltp := daysToPrint
	for i:=1; i<=dltp; i++ {
		dayNum := ""
		if i<10 {
			dayNum = fmt.Sprint("  ", i, " ")
		} else {
			dayNum = fmt.Sprint(" ", i, " ")
		}
		finStr += dayNum
		daysToPrint--
		if daysToPrint%7 == 0 && daysToPrint != 0 {
			finStr += "\n"
		}
	} // 50% grug brain
	return finStr
}

func backOneWeekday(wd time.Weekday) time.Weekday{
	if wd == 0 {
		return 6
	} else {
		return wd - 1
	}
}

func backOneMonth(mon time.Month) time.Month{
	if mon == 1 {
		return 12
	} else {
		return mon - 1
	}
}

func daysInMonth(leap bool, mon time.Month) int {
	switch mon.String() {
		case "January"  : return 31
		case "February" :
			 if leap {
			 	return 29
			 } else {
			 	return 28
			 }
		case "March"    : return 31
		case "April"    : return 30
		case "May"      : return 31
		case "June"     : return 30
		case "July"     : return 31
		case "August"   : return 31
		case "September": return 30
		case "October"  : return 31
		case "November" : return 30
		case "December" : return 31
		default         : return 0 // shouldn't need this but here just in case
	}
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
