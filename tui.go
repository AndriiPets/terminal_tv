package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	ascii "github.com/AndriiPets/terminal_yt/image_manipulation"
	videoplayer "github.com/AndriiPets/terminal_yt/video_player"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	default_w = 160
	default_h = 90
)

type App struct {
	Keys          KeyMap
	VideoPlayer   *videoplayer.VideoPlayer
	textInput     textinput.Model
	videoUrlInput string
	fps           float64

	fCounter      int
	currFrame     ascii.Frame
	prevFrame     ascii.Frame
	num_of_frames int

	pause bool

	err error

	isVideoStarted bool
}

type (
	ErrMsg  error
	TickMsg time.Time
)

func initModel(vp *videoplayer.VideoPlayer) *App {
	//Textimput
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	//VideoPlayer settings
	fps := 30.0

	app := &App{
		Keys: KeyMap{
			Quit:  key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
			Pause: key.NewBinding(key.WithKeys("spacebar", "p"), key.WithHelp("spacebar/p", "pause")),
			Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "enter the input")),
		},
		VideoPlayer: vp,
		fps:         fps,
		textInput:   ti,
	}

	return app
}

func RunTUI(vp *videoplayer.VideoPlayer) {

	p := tea.NewProgram(initModel(vp))
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

func (app *App) startStream() error {
	url := app.videoUrlInput
	err := app.VideoPlayer.LoadVideoMetadata(url)
	if err != nil {
		return err
	}

	durMilisecond, _ := strconv.Atoi(app.VideoPlayer.Video.Data.Duration)
	numFrames := (durMilisecond / 1000) * int(app.fps)
	app.num_of_frames = numFrames

	go app.VideoPlayer.StartStream()
	return nil
}

func (app *App) updateFrame() (tea.Model, tea.Cmd) {

	if app.VideoPlayer.Video == nil {
		app.currFrame = ascii.Frame{Content: "//////NO_VIDEO//////", Width: default_w, Height: default_h}
		return app, app.tick()
	}

	nextFrm, ok := app.VideoPlayer.Video.FrameMap.Load(app.fCounter)
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
	return tea.Batch(app.tick(), textinput.Blink)
}

func (app *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch {

		case key.Matches(msg, app.Keys.Quit):
			fmt.Println()
			return app, tea.Quit

		case key.Matches(msg, app.Keys.Pause):
			app.pause = !app.pause

		case key.Matches(msg, app.Keys.Enter):
			if !app.isVideoStarted {
				err := app.startStream()
				if err != nil {
					app.err = err
				} else {
					app.isVideoStarted = !app.isVideoStarted
				}
			}
		}

	case TickMsg:
		if !app.pause && app.VideoPlayer.Video != nil && app.isVideoStarted {
			app.fCounter++
		}

		if app.fCounter >= app.num_of_frames {
			app.fCounter = 0
		}

		return app.updateFrame()

	case tea.WindowSizeMsg:
		app.VideoPlayer.Width = msg.Width
		app.VideoPlayer.Heigth = msg.Height
	}

	app.textInput, cmd = app.textInput.Update(msg)
	app.videoUrlInput = app.textInput.Value()

	return app, cmd
}

func (app *App) View() string {

	var frame string

	if app.isVideoStarted {
		frame = app.currFrame.Content
		frame += fmt.Sprintf("counter: %d total: %d", app.fCounter, app.num_of_frames)

		//playback bar
		var sb strings.Builder
		vPercent := (app.fCounter * 100) / app.num_of_frames
		fPercent := (vPercent * app.currFrame.Width) / app.num_of_frames
		passed := strings.Repeat("#", fPercent)
		left := strings.Repeat(".", app.currFrame.Width-fPercent)
		sb.WriteString(passed)
		sb.WriteString(left)

		frame += fmt.Sprintf("\n%s", sb.String())
	} else {
		b := &strings.Builder{}
		b.WriteString("Paste url of the video below :")
		b.WriteString(app.textInput.View())

		frame = b.String()

		if app.err != nil {
			frame += fmt.Sprintf("\nError :%s", app.err.Error())
		}
	}

	return frame
}
