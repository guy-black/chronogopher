# Chronogopher

a project to learn Go as well as imperative programming in general.  Currently just has a very basic digital 12 hour clock, calendar, and non interactive todo list.
it is officially technically good enough to replace my side panel eww widget and the window of tty-clock I keep open on my main workspace
but still have a while a lot of work to do before I can feel complete and happy about it

roadmap to 1.0

- figure out how to style the numbers on the calndar differently.
` -- maybe I'm just missing something but I can't get them to work right
  -- want to keep the days of the current month but not today the same
  -- and make days of previous/next month a different color, and current day a diffferent color
  -- add a * after any day with todo list items attached to it
- neatly organize all the color values and anything else worth customizing into their owns costs in one section together
  -- go compiles fast enough that I can just treat it as an interpreted language and config by editing source for now
  -- later on would like to actually read in and parse a config file just for completions sake
- read todo list from config file and save to it on close
  -- save it on every update, expensive but guaratees you'll always have an updated todo list
  -- could also poll it on every tick to keep it in sync if you have multiple instances open
- select between clock, calendar, and todo list sections
- when clock selected change clock from 12 hour to 24 hour clock or maybe others too in the future
- when calendar selected move day cursor up, down, left, and right around calendar
 -- able to add, remove, or view todo tasts on certain days
 -- also just see the calendar for that month
 -- repeat same todo over regular period (daily, weekly, monthly, etc)
- when todo list selected able to add/edit/remove todolist items
 -- scroll up and down todo list if list items exceed given lengths
 -- also add sub items to todo list items
 -- show a general/longterm todo list and a day specific todo list
 -- add optional time and alarm to todo list items, when the time comes either
   -- reverse colors repeatedly until user stops it, and/or
   -- execute arbitrary command/script to get the user's attention.
