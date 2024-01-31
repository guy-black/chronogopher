package main

import (
	"github.com/charmbracelet/lipgloss"
)

// CONSTANTS and vars FOR CONFIGURATION

// GLOBAL
var appStyle = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(lipgloss.ANSIColor(2))
// conrols
const (
	// to close chronogopher
	QUIT string = "ctrl+q"
	// to cycle forwards through the sections
	CYCLE_SECTS string = "tab"
	// to cycle backwards through the sections
	REVCYCLE_SECTS string = "shift+tab"
)

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
// conrols
const (
	// to cycle back through clocks
	// to cycle foward through clocks
	PREV_CLOCK string = "left"
	NEXT_CLOCK string = "right"
)

// CALENDAR

// style for current calendar day
var calCurrDay = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
// style for calendar days of the selected month
var currMonthDay = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
// style's for other days from the next or previous month
var othMonthDay = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
func calStyle(m model) lipgloss.Style {
	if m.sel == CalSect {
		return lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	}
	return lipgloss.NewStyle().Padding(1)
}
//controls
const (
	// to go back to current day
	CAL_TODAY string = " "
	// to go to the previous day
	CAL_PREV_DAY string = "left"
	// to go to the next day
	CAL_NEXT_DAY string = "right"
	// to go back one week
	CAL_PREV_WEEK string = "up"
	// to go foward one week
	CAL_NEXT_WEEK string = "down"
	// to go back one month
	CAL_PREV_MON string = "ctrl+left"
	// to go forward one month
	CAL_NEXT_MON string = "ctrl+right"
)

//TODO
const(
	// where to look for the todolist
	// can be written as an absolute path
	// or relative to where it's being launched from
	// if this file does not exist, it will be created
	TODO_LIST string = ".cgtodo"
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
// controls
const (
	// to go up
	TD_UP string = "up"
	// to go down
	TD_DOWN string = "down"
	// if text input is focused and non empty add contents as new todo item below
	// the selected todo item else, focuses the text input
	TD_NEW_ADD string = "enter"
	// if text input is focused and non empty add contents as new todo item in place of the
	// selected task else, focuses the text input and prefill it with currently selected task
	TD_COPY_REPL string = "alt+enter"
	// unfocuses the textinput without changing todolist
	TD_CANCEL string = "esc"
	// delete selected task
	TD_DELETE string = "delete"
)


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
