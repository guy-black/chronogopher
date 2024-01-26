# Chronogopher

a project to learn Go as well as imperative programming in general.  Currently just has a very basic digital 12 hour clock, calendar, and non interactive todo list.
it is officially technically good enough to replace my side panel eww widget and the window of tty-clock I keep open on my main workspace
but still have a while a lot of work to do before I can feel complete and happy about it


# controls:

## appwide:
- ctrl+q to quit
- tab to go down one section
- shift-tab to go up one section
## clock
- left/right: cycle through clocks.  so far only 12hour and 24 hour digital clocks available
## todo list
- up/down: scroll through todo items
- delete: delete selected todo item
- esc: if textfield is focused, unfocus it and throw it's contents away
- enter:
  - if text field is active and typed in, insert that text as a todo item directly below the selected item
  - if text field is not active, make it active with an empty field
- alt+enter:
  - if text field is active and typed in, replace selected todo item with typed text
  - if text field is not active, make it active with the text of the selected todo item prefilled in


# roadmap to 1.0

- figure out how to style the numbers on the calendar differently.
` -- maybe I'm just missing something but I can't get them to work right
  -- want to keep the days of the current month but not today the same
  -- and make days of previous/next month a different color, and current day a diffferent color
  -- add a * after any day with todo list items attached to it
- neatly organize all the color values and anything else worth customizing into their owns consts in one section together
  -- go compiles fast enough that I can just treat it as an interpreted language and config by editing source for now
  -- later on would like to actually read in and parse a config file just for completions sake
- when calendar selected move day cursor up, down, left, and right around calendar
 -- able to add, remove, or view todo tasks on certain days
 -- also just see the calendar for that month
 -- repeat same todo over regular period (daily, weekly, monthly, etc)
 -- show a general/longterm todo list and a day specific todo list
 -- add optional time and alarm to todo list items, when the time comes either
   -- reverse colors repeatedly until user stops it, and/or
   -- execute arbitrary command/script to get the user's attention.
