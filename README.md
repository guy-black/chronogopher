# Chronogopher

a project to learn Go as well as imperative programming in general.  Currently just has a very basic digital 12 hour clock, calendar, and non interactive todo list.
it is officially technically good enough to replace my side panel eww widget and the window of tty-clock I keep open on my main workspace
but still have a while a lot of work to do before I can feel complete and happy about it

![screenshot](https://raw.githubusercontent.com/guy-black/chronogopher/main/screenshot.png)

# default controls:

## appwide:
- ctrl+q to quit
- tab to go down one section
- shift-tab to go up one section
## clock
- left/right: cycle through clocks.  so far only 12hour and 24 hour digital clocks available
## Calendar
- left/right go back/forwards 1 day
- down/up go back/forwards 1 week
- ctrl+left/right go back/forwards 1 month
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

# setup and install

after cloning the repo copy or rename `conf.go.def` to `<whateverYouWantToNameYouCconfig>.go
you can run the app right away with default configs
```
go run .
```
or build an executable with
```
go build .
```
to customize, just read through your copy of the config and edit values, it is
quite heavily commented, but if anything is unclear please open an issue so I
clarify and improve the comments.  Once you're ready, if you only have one config
then you can build or run chronogopher with the same commands as above.  If you
have multiple config files you want to choose from then build or run with
```
go run chronogopher.go <configFileToUse>.go
```
or to build the executable
```
go run chronogopher.go <configFileToUse>.go
```
# features planned to be implemented in the future

- properly sync todolist with Tea.Msg
-- new const SYNC_FREQ of type Time.Duration for how often to sync
-- another function like doTick called doSync, but tea.Time will take SYNC_FREQ and it's
  function will ignore the time and return a SyncMsg carrying the current list of tasks in
  the TODO_LIST file as a []Task.  Either do one of two things in the update function whenever a SyncMsg tsks comes in
    - either always replace m.todo.tasks with the tsks from SyncMsg and decrease m.todo.sel if need be, or
    - check if m.todo.taks == tsks.  If so do nothing, if not replace replace m.todo.tasks with tsks and update m.todo.sel if need be
    - figure out how to test which is more efficient, but it's probably a miniscule difference it'd be a decent learning experience

- figure out how to make app take up entire window available
   -- for now create new consts MIN_APP_WIDTH and MIN_APP_HEIGHT, if the window is not atleast that big render an error
      instead of the app, like btop does.  Later on write funtion to calculate those values based on default clock and cal
   -- if toggling between clocks and cals, skip over options that won't fit in currently allotted space
   -- if a todolist item is too long, break it into multiple lines with some visual indicator that these lines are one task
      -- perhaps 4 new consts with defaults:
           START_TODO_ITEM="╮",
           MID_TODO_ITEM= "│",
           END_TODO_ITEM="╯",
           TODO_ITEM_BOOKEND="─"
           MULTILINE_TODO_INDICATOR_ON_RIGHT="true"
      -- which would render
         super duper overly long todo list item that takes too much space
      -- broken up into
         super duper
         overly long todo
         list item that
         takes too much
         space
      -- as
           super duper ──╮
         overy long todo │
          list item that │
          takes too much │
               apace ────╯
      -- algorithm to split todo item will try to take as many space seperated words as will fit in the allotted width-1
         - if a word is just to big it will be split and hyphenated
   -- every component is centered in as much horizontal space it is allotted
   -- cal and clock only take as much vertical space as needed, todolist will take as much vertial space as is available

- allow for user defined calendars and date formats similar to clock

- toggle between vertical alignment, horizontal alignment, or some combination (clock/cal sidebyside over todolist,
  clock on top over sidebyside cal and todo list, etc)

- schedule tasks to be done on a certain day or time
-- attach a glyph to days on the calendar if there's a task scheduled that day
-- do some kind of alarm thing
-- show a general/longterm todo list and a day specific todo list.

- colors adjustable while running

- different clock faces
-- large 5x6
-- tiny plain text
-- word clock like this https://youtu.be/SXYwSN6mX_Q?t=279

- pomodoro
