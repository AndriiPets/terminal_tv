package ui

import (
	"fmt"
	"strings"

	"github.com/AndriiPets/terminal_yt/utils"
	videoplayer "github.com/AndriiPets/terminal_yt/video_player"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type SearchScreen struct {
	err     error
	spinner spinner.Model
	query   string
	keys    KeyMap
	vp      *videoplayer.VideoPlayer
	pending bool
	list    list.Model
}

type statusMsg int

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

func initSearchScreen(vp *videoplayer.VideoPlayer) SearchScreen {
	s := spinner.New()
	s.Spinner = spinner.Line

	return SearchScreen{
		spinner: s,
		vp:      vp,
		pending: true,
		keys: KeyMap{
			Quit:  key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
			Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "enter the input")),
			Up:    key.NewBinding(key.WithKeys("up"), key.WithHelp("up", "up")),
			Down:  key.NewBinding(key.WithKeys("down"), key.WithHelp("down", "down")),
		},
	}
}

func (s *SearchScreen) loadSearchResults() tea.Msg {
	results, err := utils.SearchYT(s.query)
	if err != nil {
		return errMsg{error: err}
	}

	var items []list.Item

	for _, res := range results {
		items = append(items, res)
	}
	s.list.Title = "Here is what i found"

	s.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	s.pending = false
	return statusMsg(1)
}

func (s *SearchScreen) Init() tea.Cmd {

	return tea.Batch(s.spinner.Tick, s.loadSearchResults)
}

func (s *SearchScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Quit):
			fmt.Println()
			return s, tea.Quit

		case key.Matches(msg, s.keys.Enter):
			if !s.pending {

				i := s.list.SelectedItem().(utils.SearchResult)
				err := s.vp.LoadVideoMetadata(i.Description())
				if err != nil {
					s.err = err
				}

				videoScreen := initVideoScreen(s.vp, false)
				return NewRootScreen(s.vp, Video).SwitchScreen(&videoScreen)
			}
		}

	case tea.WindowSizeMsg:
		if !s.pending {
			h, v := docStyle.GetFrameSize()
			s.list.SetSize(msg.Width-h, msg.Height-v)
		}

	case errMsg:
		s.err = msg
		return s, nil

	case spinner.TickMsg:
		s.spinner, cmd = s.spinner.Update(msg)
		cmds = append(cmds, cmd)
		return s, tea.Batch(cmds...)
	}

	if !s.pending {
		s.list, cmd = s.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return s, tea.Batch(cmds...)
}

func (s *SearchScreen) View() string {
	var sb strings.Builder
	if s.pending {
		sb.WriteString(fmt.Sprintf("\n   %s Searching for videos...\n\n", s.spinner.View()))
	} else {
		sb.WriteString(docStyle.Render(s.list.View()))
	}
	if s.err != nil {
		sb.WriteString(s.err.Error())
	}

	return sb.String()
}
