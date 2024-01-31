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

after cloning the repo, you can run the app right away with default configs
```
go run .
```
or build an executable with
```
go build .
```
to customize, make a copy of `defaultConf.go` and edit it to your liking.
to run your configured version
```
go run chronogopher.go <yourConfigFile>.go
```
or to build the executable
```
go run chronogopher.go <yourConfigFile>.go
```

# features planned to be implemented in the future

- scheule tasks to be done on a certain day or time
-- attach a glyph to days on the calendar if there's a task scheduled that day
-- do some kind of alarm thing
-- show a general/longterm todo list and a day specific todo list.

- colors adjustable while running

- different clock faces
-- large 5x6
-- tiny plain text
-- word clock like this https://youtu.be/SXYwSN6mX_Q?t=279

- pomodoro
