package main

import (
	"fmt"
	"os"
	"time"
	"strings"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textinput"
)

// styles
var (
	// style for whole app
	appStyle = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(lipgloss.Color("2"))
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

func todoItemStyle(m model, i byte) lipgloss.Style {
	if i == m.todo.sel {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	}
	return lipgloss.NewStyle()
}

// CUSTOM TYPES //
// clockType types
type ClockType byte

const (
	H12 ClockType = iota
	H24
)
// keep this variable in sync with number of clock types available
// should always be numberOfClock - 1
const HIGHEST_CLOCK ClockType = 1

func incClk (c ClockType) ClockType {
	if c == HIGHEST_CLOCK {
		return 0
	}
	c++
	return c
}
func decClk (c ClockType) ClockType {
	if c == 0 {
		return HIGHEST_CLOCK
	}
	c--
	return c
}

// section type
type Section byte

const (
	 ClockSect Section = iota
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

// Todo list types
const TODO_LIST string = "../.cgtodo"
// where to look for the todolist
// can be written as an absolute path
// or relative to where it's being launched from

type Todo struct {
	tasks  []Task
	sel    byte // if my list get's longer than 256 items then I'll change this
}

type Task struct {
	task     string
	// subtasks [] string
	// alarm Time.time
}

// model
type model struct {
	dt        time.Time
	todo      Todo
	sel       Section
	clkTyp    ClockType
	todoInput textinput.Model
}

func initialModel() model {
	initTasks := fetchTasks()
	ti := textinput.New()
	ti.Blur()
	ti.Placeholder = "      what to do...      "
	ti.Width = 25

	return model {
		dt: time.Now(),
		todo: Todo {
			tasks: initTasks,
			sel: 0,
		},
		sel: 2,
		clkTyp: 0,
		todoInput: ti,
	}
}

func writeTasks (ts []Task) {
	os.Remove(TODO_LIST)
	fptr, _ := os.Create(TODO_LIST)
	fptr.WriteString(taskString(ts))
	fptr.Close()
	// TODO: actually check for and handle these errors
}

func taskString (ts []Task) string {
	tdl := ""
	for _, t := range ts {
		if t.task != "" {
			tdl += fmt.Sprint(t.task, "\n")
		}
	}
	return tdl
}

func styledTaskString (m model) string {
	ts := m.todo.tasks
	tdl := ""
	for i, t := range ts {
		if t.task != ""{
			if byte(i)==m.todo.sel {
				tdl += fmt.Sprint("⊢", t.task, "\n")
			} else {
				tdl += fmt.Sprint(" ", t.task, "\n")
			}
		}
	}
	return tdl
}

func fetchTasks () []Task {
	dat, err := os.ReadFile(TODO_LIST)
	var initTasks []Task
	if err!=nil {
		initTasks = make([]Task, 0)
	} else {
		initTasks = make([]Task, 0)
		for _,t:=range strings.Split(string(dat), "\n") {
			nt := Task{task: t}
			initTasks = append(initTasks, nt)
		}
	}
	return initTasks
}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
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
				case "ctrl+q":
					return m, tea.Quit
				case "tab":
					m.sel = incSel (m.sel)
					return m, nil
				case "shift+tab":
					m.sel = decSel (m.sel)
					return m, nil
				case "left":
					switch m.sel {
						case ClockSect:
							m.clkTyp = decClk(m.clkTyp)
							return m, nil
						case CalSect:
							return m, nil
						case TodoSect:
							return m, nil
					}
				case "right":
					switch m.sel {
						case ClockSect:
							m.clkTyp = incClk(m.clkTyp)
							return m, nil
						case CalSect:
							return m, nil
						case TodoSect:
							return m, nil
					}
				case "up":
					switch m.sel{
						case TodoSect:
							if int(m.todo.sel) == 0 {
							// TODO: figure out why I need to subtract two instead of one here and below
								m.todo.sel = byte(len(m.todo.tasks)-2)
								return m, nil
							} else {
								m.todo.sel--
								return m, nil
							}
					}
				case "down":
					switch m.sel{
						case TodoSect:
							if int(m.todo.sel) == len(m.todo.tasks)-2 {
								m.todo.sel = 0
								return m, nil
							} else {
								m.todo.sel++
								return m, nil
							}
					}
				case "enter":
					switch m.sel{
						case TodoSect:
							if m.todoInput.Focused() {
								// if enter is pressed while it's focused
								// update model and file todolist
								nt := m.todoInput.Value()
								m.todoInput.Reset()
								m.todo.tasks = append (m.todo.tasks, Task{task: nt})
								writeTasks(m.todo.tasks)
								// unfocus textinput
								m.todoInput.Blur()
							} else { // focus it
								foc := m.todoInput.Focus()
								return m, foc
							}
					}
				case "esc":
					if m.sel == TodoSect {
						m.todoInput.Reset()
						m.todoInput.Blur()
					}
					return m, nil
			}
		case TickMsg:
			if !slices.Equal(m.todo.tasks, fetchTasks()){
				m.todo.tasks = fetchTasks()
				if int(m.todo.sel) > len(m.todo.tasks){
					m.todo.sel--
				}
			}
			m.dt = time.Now()
			return m, doTick()
	}
	var cmd tea.Cmd
	m.todoInput, cmd = m.todoInput.Update(msg)
	return m, cmd
}


func (m model) View() string {
// creating the clock section
	clock := clockStyle(m).Render(genClock(m))
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
	tdl := styledTaskString(m)
	tdl = strings.TrimSpace(tdl)
	tdl = lipgloss.JoinVertical(.5, tdl)
	todo := todoStyle(m).Render(lipgloss.JoinVertical(.5,
		"todo",
		tdl,
		m.todoInput.View()))
// here's the actual view to be rendered
		return appStyle.Render(lipgloss.JoinVertical(0.5,
			clock,
			mdy,
			cal,
			todo))
}

func genClock (m model) string{
	hour, min, sec := m.dt.Clock()
	switch m.clkTyp {
		case H12:
			ampm := "am"
			if hour>12 {
				hour -= 12
				ampm = "pm"
			}
			if hour == 12 {
				ampm = "pm"
			}
			if hour == 0 {
				hour = 12
			}
			time := stringTime(hour, min, sec)
			return fmt.Sprintf("%s\n%s\n%s",
				timeTopLine(time),
				timeMidLine(time)+ampm,
				timeBotLine(time))
		case H24:
			time := stringTime(hour, min, sec)
			return fmt.Sprintf("%s\n%s\n%s",
				timeTopLine(time),
				timeMidLine(time),
				timeBotLine(time))
	}
	return ""
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
