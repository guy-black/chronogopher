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
	// date  date
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
	initTasks := make([]Task, 0)
	if err!=nil {
		// TODO, check check if this can produce errors besides file not found
		// if it can than check for them and handle them.  file not found already
		// solved with writeTask and by returning the empty []Task
	} else {
		for _,t:=range strings.Split(string(dat), "\n") {
			if t != "" {
				nt := Task{task: t}
 				initTasks = append(initTasks, nt)
			}
		}
	}
	return initTasks
}

// date type

type date struct {
	year  int
	month time.Month
	day   int
}

func backOneDay (d date) date {
	if d.day > 1 {
	// easy solution just have to change day
		d.day--
	} else {// otherwise day was 1 and need to shift month and maybe year too
		d.month = backOneMonth(d.month)
		d.day = daysInMonth(d.year%4==0, d.month)
		if d.month.String() == "December"{
			d.year--
		}
	}
	return d
}

func forwardOneDay (d date) date {
	if d.day < daysInMonth(d.year%4==0, d.month) {
		// easy case just change day
		d.day++
	} else {
		// otherwise have to change month too
		d.day = 1
		d.month = forwardOneMonth(d.month)
		if d.month.String() == "January" {
			d.year++
		}
	}
	return d
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
	selDate    date
}

func initialModel() model {
	initTasks := fetchTasks()
	ti := textinput.New()
	ti.Blur()
	ti.Placeholder = TODO_PLACEHOLDER
	ti.Width = 25
	year, month, day := time.Now().Date()
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
		selDate: date {
			year: year,
			month: month,
			day: day,
		},
	}
}

type SyncMsg []Task

func doSync() tea.Cmd {
	return tea.Every(TODO_SYNC_FREQ, func(_ time.Time) tea.Msg {
		return SyncMsg (fetchTasks())
	})
}

func writeTasksCmd(tasks []Task) tea.Cmd {
	return func () tea.Msg {
		writeTasks(tasks)
		return nil
	}
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
		doSync(),
		doTick(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type){
		// todo on every second
		case SyncMsg:
			if !slices.Equal(m.todo.tasks, msg){
				m.todo.tasks = msg
				for int(m.todo.sel) >= len(m.todo.tasks)-1{
					m.todo.sel--
				}
			}
			return m, doSync()
		case TickMsg:
			// before updating dt get previous date
			py, pm, pd := m.dt.Date()
			now := time.Time(msg) // why do I need to cast TickMsg to time.Time but I
			cy, cm, cd := now.Date() // don't have to cast SyncMsg to []Task
			// check if the selDate is the same as the previous date
			if py == m.selDate.year && pm == m.selDate.month && pd == m.selDate.day {
				// selDate was equal to current day, so it should day in sync
				// if previous year/month/day isn't equal to current, update it
				if py != cy { m.selDate.year = cy}
				if pm != cm { m.selDate.month = cm}
				if pd != cd { m.selDate.day = cd}
			}
			m.dt = now
			return m, doTick()
		case tea.KeyMsg:
			// reacting to keypresses
			switch msg.String() {
				// global key press actions!!!
				case QUIT:
					return  m, tea.Quit
				case CYCLE_SECTS:
					m.sel = incSel (m.sel)
					return m, nil
				case REVCYCLE_SECTS:
					m.sel = decSel (m.sel)
					return m, nil
			}
			// section specific keypresses
			switch m.sel {
				case ClockSect:
					switch msg.String(){

// clock specific keypresses

						case PREV_CLOCK:
							m.clkTyp = decClk(m.clkTyp)
						case NEXT_CLOCK :
							m.clkTyp = incClk(m.clkTyp)
					}
				case CalSect:
					switch msg.String(){

// cal specific keypresses

						case CAL_TODAY: // going back to current day
							yr,mon,d := time.Now().Date()
							m.selDate.year = yr
							m.selDate.month = mon
							m.selDate.day = d
							return m, nil
						case CAL_PREV_DAY: // going back one day
							m.selDate = backOneDay(m.selDate)
							return m,nil
						case CAL_NEXT_DAY: // going foward one day
							m.selDate = forwardOneDay(m.selDate)
							return m,nil
						case CAL_PREV_WEEK: // going back 7 days
							for i:=0; i<7; i++ {
								m.selDate = backOneDay(m.selDate)
							}
							return m,nil
						case CAL_NEXT_WEEK: // going forward 7 days
							for i:=0; i<7; i++ {
								m.selDate = forwardOneDay(m.selDate)
							}
							return m,nil
						case CAL_PREV_MON: // going back one month
							m.selDate.month = backOneMonth (m.selDate.month)
							if m.selDate.month.String() == "December" {
								m.selDate.year--
							}
							dim := daysInMonth(m.selDate.year%4==0, m.selDate.month)
							if m.selDate.day > dim {
								m.selDate.day = dim
							}
							return m,nil
						case CAL_NEXT_MON: // going fotward one month
							m.selDate.month = forwardOneMonth (m.selDate.month)
							if m.selDate.month.String() == "January" {
								m.selDate.year++
							}
							dim := daysInMonth(m.selDate.year%4==0, m.selDate.month)
							if m.selDate.day > dim {
								m.selDate.day = dim
							}
							return m,nil
					}
				case TodoSect:
					switch msg.String(){

// todo specific keypresses

						case TD_UP:
							if int(m.todo.sel) == 0 {
								m.todo.sel = byte(len(m.todo.tasks)-1)
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
						case TD_DOWN:
							if m.todo.sel == byte(len(m.todo.tasks)-1) {
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
						case TD_NEW_ADD:
							if m.todoInput.Focused() {
								// if enter is pressed while it's focused
								// update model and file todolist
								nt := m.todoInput.Value()
								m.todoInput.Reset()
								if nt != "" {
									if len(m.todo.tasks) == 0 {
										m.todo.tasks = [] Task {Task{task: nt}}
									}else{
										m.todo.tasks = slices.Insert (m.todo.tasks,
											int(m.todo.sel)+1 ,
											Task{task: nt})
									}
								}
								// unfocus textinput
								m.todoInput.Blur()
								return m, writeTasksCmd(m.todo.tasks)
							} else { // focus it
								foc := m.todoInput.Focus()
								return m, foc
							}
						case TD_COPY_REPL:
							if m.todoInput.Focused() {
								if int(m.todo.sel) >= len(m.todo.tasks){
									// if sel is greater than the length of the slice of tasks
									// should only happen on an empty list
									// eitherway just do nothing to avoid trying to read from empty slice
									return m, nil
								}
								// check that it's not blank
								nt := m.todoInput.Value()
								m.todoInput.Reset()
								if nt!="" {
									selint:=int(m.todo.sel)
									rem:=slices.Delete(m.todo.tasks, selint, selint+1)
									ins:=slices.Insert(rem, selint, Task{task: nt})
									m.todo.tasks = ins
								}
								m.todoInput.Blur()
								return m, writeTasksCmd(m.todo.tasks)
							} else {
								foc := m.todoInput.Focus()
								if int(m.todo.sel) < len(m.todo.tasks) {
									txt := m.todo.tasks[int(m.todo.sel)].task
									m.todoInput.SetValue(txt)
								}
								return m,foc
							}
						case TD_CANCEL:
							m.todoInput.Reset()
							m.todoInput.Blur()
							return m, nil
						case TD_DELETE:
							selint := int(m.todo.sel)
							m.todo.tasks = slices.Delete(m.todo.tasks, selint, selint+1)
							if selint >= len(m.todo.tasks)-1 && selint >= 0 {
								m.todo.sel--
							}
							return m, writeTasksCmd(m.todo.tasks)
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
	leap := m.selDate.year%4 == 0
	selMon := m.selDate.month
	selDay := m.selDate.day
	selWd := time.Date(m.selDate.year,
		selMon,
		selDay,
		1,1,1,1,time.Local).Weekday()

	calDays := genCalDays (leap, selMon, selDay, selWd)
	cal := calStyle(m).Render(lipgloss.JoinVertical(0.5,
		fmt.Sprintf("%s  %d\n", selMon.String(), m.selDate.year),
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
		// printing previous month days
		finStr = othMonthDay.Render(fmt.Sprint(" ", daysInPrevMon, " ")) + finStr
		daysInPrevMon--
		wdfm--
		daysToPrint--
	}
	daysThisMon := daysInMonth(leap, mon)
	for i:=1; i<=daysThisMon; i++ {
		if i==day{
			if i<10 {
				// printing selected day
				finStr += calCurrDay.Render(fmt.Sprint("  ", i, " "))
			} else {
				finStr += calCurrDay.Render(fmt.Sprint(" ", i, " "))
			}
		} else {
			if i<10 {
				// printing days of the selected month
				finStr += currMonthDay.Render(fmt.Sprint("  ", i, " "))
			} else {
				finStr += currMonthDay.Render(fmt.Sprint(" ", i, " "))
			}
		}
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
		if i<10 {
			finStr += othMonthDay.Render(fmt.Sprint("  ", i, " "))
		} else {
			finStr += othMonthDay.Render(fmt.Sprint(" ", i, " "))
		}
		daysToPrint--
		if daysToPrint%7 == 0 && daysToPrint != 0 {
			finStr += "\n"
		}
	}
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

func forwardOneMonth(mon time.Month) time.Month{
	if mon == 12 {
		return 1
	} else {
		return mon + 1
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
