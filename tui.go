package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	ascii "github.com/AndriiPets/terminal_yt/image_manipulation"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	default_w = 160
	default_h = 90
)

type App struct {
	Keys        KeyMap
	VideoPlayer *VideoPlayer
	fps         float64

	fCounter      int
	currFrame     ascii.Frame
	prevFrame     ascii.Frame
	num_of_frames int

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

	durMilisecond, _ := strconv.Atoi(vp.Video.data.duration)
	app.num_of_frames = (durMilisecond / 1000) * int(app.fps)

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

	if app.VideoPlayer.Video == nil {
		app.currFrame = ascii.Frame{Content: "//////NO_VIDEO//////", Width: default_w, Height: default_h}
		return app, app.tick()
	}

	nextFrm, ok := app.VideoPlayer.Video.frameMap.Load(app.fCounter)
	if !ok {
		if len(app.prevFrame.Content) == 0 {
			app.currFrame = ascii.Frame{Content: "//////LOADING_VIDEO//////", Width: default_w, Height: default_h}
		} else {
			app.currFrame = app.prevFrame
		}
		return app, app.tick()
	}

	nextFrame := nextFrm.(ascii.Frame)

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
		if !app.pause && app.VideoPlayer.Video != nil {
			app.fCounter++
		}

		if app.fCounter >= app.num_of_frames {
			app.fCounter = 0
		}

		return app.updateFrame()
	}

	return app, nil
}

func (app *App) View() string {

	frame := app.currFrame.Content
	frame += fmt.Sprintf("counter: %d total: %d", app.fCounter, app.num_of_frames)

	//playback bar
	var sb strings.Builder
	vPercent := (app.fCounter * 100) / app.num_of_frames
	//fPercent := (vPercent * app.currFrame.Width) / app.num_of_frames
	passed := strings.Repeat("#", vPercent)
	left := strings.Repeat(".", 100-vPercent)
	sb.WriteString(passed)
	sb.WriteString(left)

	frame += fmt.Sprintf("\n%s", sb.String())

	return frame
}
