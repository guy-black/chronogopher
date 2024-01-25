# Chronogopher

a project to learn Go as well as imperative programming in general.  Currently just has a very basic digital 12 hour clock, calendar, and non interactive todo list.
it is officially technically good enough to replace my side panel eww widget and the window of tty-clock I keep open on my main workspace
but still have a while a lot of work to do before I can feel complete and happy about it

roadmap to 1.0

- figure out how to style the numbers on the calendar differently.
` -- maybe I'm just missing something but I can't get them to work right
  -- want to keep the days of the current month but not today the same
  -- and make days of previous/next month a different color, and current day a diffferent color
  -- add a * after any day with todo list items attached to it
- neatly organize all the color values and anything else worth customizing into their owns consts in one section together
  -- go compiles fast enough that I can just treat it as an interpreted language and config by editing source for now
  -- later on would like to actually read in and parse a config file just for completions sake
- when calendar selected move day cursor up, down, left, and right around calendar
 -- able to add, remove, or view todo tasts on certain days
 -- also just see the calendar for that month
 -- repeat same todo over regular period (daily, weekly, monthly, etc)
- when todo list selected able to edit todolist items
 -- scroll up and down todo list if list items exceed given lengths
 -- also add sub items to todo list items
 -- show a general/longterm todo list and a day specific todo list
 -- add optional time and alarm to todo list items, when the time comes either
   -- reverse colors repeatedly until user stops it, and/or
   -- execute arbitrary command/script to get the user's attention.
