package main

import (
	"github.com/charmbracelet/lipgloss"
)

// CONSTANTS and vars FOR CONFIGURATION
// GLOBAL

// style for whole app
func appStyle (m model) lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(lipgloss.ANSIColor(2))
}

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
var currMonthDay = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
var othMonthDay = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))


func calStyle(m model) lipgloss.Style {
	if m.sel == CalSect {
		return lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	}
	return lipgloss.NewStyle().Padding(1)
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
