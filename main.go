package main

import (
	"fmt"
	"os"
	"time"
	"strings"
	"slices"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

// CUSTOM TYPES //
// clockType types
type ClockType byte

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

type Todo struct {
	tasks  []Task
	sel    byte // if my list get's longer than 256 items then I'll change this
}

type Task struct {
	task     string
	// alarm Time.time
	// date Time.time.Date()
}

// model
type model struct {
	dt         time.Time
	todo       Todo
	sel        Section
	clkTyp     ClockType
	todoInput  textinput.Model
	vpStart    int
	vpEnd      int
}

func initialModel() model {
	initTasks := fetchTasks()
	ti := textinput.New()
	ti.Blur()
	ti.Placeholder = TODO_PLACEHOLDER
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
		vpStart: 0,
		vpEnd: TODO_VP_LEN - 1,
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
		if t.task != "" && i>=m.vpStart && i <= m.vpEnd{
			if byte(i)==m.todo.sel {
				tdl += fmt.Sprint(TODO_SEL_PREF, t.task, "\n")
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
	switch msg := msg.(type){
		// todo on every second
		case TickMsg:
			if !slices.Equal(m.todo.tasks, fetchTasks()){
				m.todo.tasks = fetchTasks()
				for int(m.todo.sel) >= len(m.todo.tasks)-1{
					m.todo.sel--
				}
			}
			m.dt = time.Now()
			return m, doTick()
		case tea.KeyMsg:
			// reacting to keypresses
			switch msg.String() {
				// global key press actions!!!
				case "ctrl+q":
					return  m, tea.Quit
				case "tab":
					m.sel = incSel (m.sel)
					return m, nil
				case "shift+tab":
					m.sel = decSel (m.sel)
					return m, nil
			}
			// section specific keypresses
			switch m.sel {
				case ClockSect:
					switch msg.String(){

// clock specific keypresses

						case "left":
							m.clkTyp = decClk(m.clkTyp)
						case "right" :
							m.clkTyp = incClk(m.clkTyp)
					}
				case CalSect:
					switch msg.String(){

// cal specific keypresses

					}
				case TodoSect:
					switch msg.String(){

// todo specific keypresses

						case "up":
							if int(m.todo.sel) == 0 {
							// TODO: figure out why I need to subtract two instead of one here and below
								m.todo.sel = byte(len(m.todo.tasks)-2)
								if m.todo.sel > byte(m.vpEnd) {
									se := int(m.todo.sel) - m.vpEnd
									m.vpStart += se
									m.vpEnd += se
								}
								return m, nil
							} else {
								m.todo.sel--
								if m.todo.sel < byte (m.vpStart) {
									m.vpStart--
									m.vpEnd--
								}
								return m, nil
							}
						case "down":
							if m.todo.sel == byte(len(m.todo.tasks)-2) {
							// TODO: figure out why I need to subtract two instead of one here and below
								m.todo.sel = 0
								m.vpStart = 0
								m.vpEnd = 14 // basically resetting the viewport
								return m, nil
							} else {
								m.todo.sel++
								if int(m.todo.sel) > m.vpEnd {
									m.vpStart++
									m.vpEnd++
								}
								return m, nil
							}
						case "enter":
							if m.todoInput.Focused() {
								// if enter is pressed while it's focused
								// update model and file todolist
								nt := m.todoInput.Value()
								m.todoInput.Reset()
								m.todo.tasks = slices.Insert (m.todo.tasks,
									int(m.todo.sel)+1 ,
									Task{task: nt})
								writeTasks(m.todo.tasks)
								// unfocus textinput
								m.todoInput.Blur()
							} else { // focus it
								foc := m.todoInput.Focus()
								return m, foc
							}
						case "alt+enter":
							if m.todoInput.Focused() {
								// check that it's not blank
								nt := m.todoInput.Value()
								m.todoInput.Reset()
								if nt!="" {
									rem:=slices.Delete(m.todo.tasks, int(m.todo.sel), int(m.todo.sel)+1)
									ins:=slices.Insert(rem, int(m.todo.sel), Task{task: nt})
									m.todo.tasks = ins
									writeTasks(m.todo.tasks)
								}
								m.todoInput.Blur()
							} else {
								foc := m.todoInput.Focus()
								txt := m.todo.tasks[int(m.todo.sel)].task
								m.todoInput.SetValue(txt)
								return m,foc
							}
						case "esc":
							m.todoInput.Reset()
							m.todoInput.Blur()
							return m, nil
						case "delete":
							selint := int(m.todo.sel)
							m.todo.tasks = slices.Delete(m.todo.tasks, selint, selint+1)
							writeTasks(m.todo.tasks)
							if selint >= len(m.todo.tasks)-1 {
								m.todo.sel--
							}
							return m, nil
					}
			}
	}
	// handle bubbles
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
		TODO_LABEL,
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
			top := stringTime(hour, min, sec)
			mid := top // all three of these
			bot := top // should be equal
			for i,v := range digi3x3 {
				top = strings.ReplaceAll(top, i, v["top"])
				mid = strings.ReplaceAll(mid, i, v["mid"])
				bot = strings.ReplaceAll(bot, i, v["bot"])
			}
			return fmt.Sprintf("%s\n%s\n%s", top, mid+ampm, bot)
		case H24:
			top := stringTime(hour, min, sec)
			mid := top // all three of these
			bot := top // should be equal
			for i,v := range digi3x3 {
				top = strings.ReplaceAll(top, i, v["top"])
				mid = strings.ReplaceAll(mid, i, v["mid"])
				bot = strings.ReplaceAll(bot, i, v["bot"])
			}
			return fmt.Sprintf("%s\n%s\n%s", top, mid, bot)
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

// CONSTANTS and vars FOR CONFIGURATION
// GLOBAL

// style for whole app
var appStyle = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(lipgloss.Color("2"))

// CLOCK

// keep this variable should always be exactly 1 less than the amount of
// ClockTypes you have.  for two ClockTypes HIGHEST_CLOCK is 1.  if you add a
// ClockType then increment it, if you take one away, decrement it
// you must also update the genClock function just under the view to handle
// you're new ClockType.
const HIGHEST_CLOCK ClockType = 1
const (
	H12 ClockType = iota
	H24
)

func clockStyle(m model) lipgloss.Style {
	if m.sel == ClockSect {
		return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Foreground(lipgloss.Color("5"))
	}
	return lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color("5"))
}

// CALENDAR

var calCurrDay = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
var othMonthDay = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))

func calStyle(m model) lipgloss.Style {
	if m.sel == CalSect {
		return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Foreground(lipgloss.Color("3"))
	}
	return lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color("3"))
}

//TODO
const(
	// where to look for the todolist
	// can be written as an absolute path
	// or relative to where it's being launched from
	TODO_LIST string = "../.cgtodo"
	// how much vertical space for the todo section to take
	// this should be equal to
	//   TODO_VP_LEN
	// + The number of lines in your TODO_LABEL
	// + 1 for the text input line
	// eg. TODO_VP_LEN = 15 + 1 line label + 1 = 17
	TODO_HEIGHT int = 17
	// the number of tasks to be visible at a time
	TODO_VP_LEN int = 15
	// label for the todo section
	TODO_LABEL string = "todo"
	// prefix for selected todo task
	TODO_SEL_PREF string = "⊢"
	// placeholder text for the textinpu
	TODO_PLACEHOLDER string = "      what to do...      "
)

func todoStyle(m model) lipgloss.Style {
	if m.sel == TodoSect {
		return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Height(TODO_HEIGHT)
	}
	return lipgloss.NewStyle().Padding(1).Height(17)
}


// misc
var digi3x3 = map[string]map[string]string{
	"1": {
		"top": " ╻ ",
		"mid": " ┃ ",
		"bot": " ╹ ",
	},
	"2": {
		"top": "╺━┓",
		"mid": "┏━┛",
		"bot": "┗━╸",
	},
	"3": {
		"top": "╺━┓",
		"mid": "╺━┫",
		"bot": "╺━┛",
	},
	"4": {
		"top": "╻ ╻",
		"mid": "┗━┫",
		"bot": "  ╹",
	},
	"5": {
		"top": "┏━╸",
		"mid": "┗━┓",
		"bot": "╺━┛",
	},
	"6": {
		"top": "┏━╸",
		"mid": "┣━┓",
		"bot": "┗━┛",
	},
	"7": {
		"top": "╺━┓",
		"mid": "  ┃",
		"bot": "  ╹",
	},
	"8": {
		"top": "┏━┓",
		"mid": "┣━┫",
		"bot": "┗━┛",
	},
	"9": {
		"top": "┏━┓",
		"mid": "┗━┫",
		"bot": "╺━┛",
	},
	"0": {
		"top": "┏━┓",
		"mid": "┃ ┃",
		"bot": "┗━┛",
	},
	":": {
		"top": "╻",
		"mid": " ",
		"bot": "╹",
	},
}
