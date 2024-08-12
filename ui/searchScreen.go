package ui

import (
	"fmt"

	videoplayer "github.com/AndriiPets/terminal_yt/video_player"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type SearchScreen struct {
	err     error
	spinner spinner.Model
	query   string
	keys    KeyMap
	vp      *videoplayer.VideoPlayer
}

func initSearchScreen(vp *videoplayer.VideoPlayer) SearchScreen {
	s := spinner.New()
	s.Spinner = spinner.Line

	return SearchScreen{
		spinner: s,
		vp:      vp,
		keys: KeyMap{
			Quit:  key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
			Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "enter the input")),
			Up:    key.NewBinding(key.WithKeys("up"), key.WithHelp("up", "up")),
			Down:  key.NewBinding(key.WithKeys("down"), key.WithHelp("down", "down")),
		},
	}
}

func (s *SearchScreen) Init() tea.Cmd {
	return s.spinner.Tick
}

func (s *SearchScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Quit):
			fmt.Println()
			return s, tea.Quit

		case key.Matches(msg, s.keys.Enter):
			videoScreen := initVideoScreen(s.vp, false)
			return NewRootScreen(s.vp, Video).SwitchScreen(&videoScreen)
		}

	case error:
		s.err = msg
		return s, nil

	default:
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd
	}
	return s, nil
}

func (s *SearchScreen) View() string {
	str := fmt.Sprintf("\n   %s This is search screen...\n\n", s.spinner.View())
	return str
}
