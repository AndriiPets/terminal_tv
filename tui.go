package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	Keys        KeyMap
	VideoPlayer *VideoPlayer
	fps         float64

	fCounter  int
	currFrame string
	prevFrame string

	pause bool
}

type TickMsg time.Time

func RunTUI(vp *VideoPlayer) {
	app := &App{
		Keys: KeyMap{
			Quit:  key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
			Pause: key.NewBinding(key.WithKeys("spacebar", "p"), key.WithHelp("spacebar/p", "pause")),
		},
		VideoPlayer: vp,
		fps:         30,
	}

	p := tea.NewProgram(app)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (app *App) tick() tea.Cmd {
	return tea.Tick(time.Second/time.Duration(app.fps), func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (app *App) updateFrame() (tea.Model, tea.Cmd) {
	nextFrame, ok := app.VideoPlayer.Video.frameMap[app.fCounter]
	if !ok {
		if len(app.prevFrame) == 0 {
			app.currFrame = "//////LOADING_VIDEO//////"
		} else {
			app.currFrame = app.prevFrame
		}
		return app, app.tick()
	}

	app.currFrame = nextFrame
	app.prevFrame = app.currFrame

	return app, app.tick()
}

func (app *App) Init() tea.Cmd {
	return app.tick()
}

func (app *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch {

		case key.Matches(msg, app.Keys.Quit):
			fmt.Println()
			return app, tea.Quit

		case key.Matches(msg, app.Keys.Pause):
			app.pause = !app.pause
		}

	case TickMsg:
		if !app.pause {
			app.fCounter++
		}

		if app.currFrame == "EOF" {
			app.fCounter = 0
		}

		return app.updateFrame()
	}

	return app, nil
}

func (app *App) View() string {

	frame := app.currFrame
	frame += fmt.Sprintf("counter: %d", app.fCounter)

	return frame
}
