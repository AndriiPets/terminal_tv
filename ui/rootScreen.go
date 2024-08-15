package ui

import (
	videoplayer "github.com/AndriiPets/terminal_yt/video_player"
	tea "github.com/charmbracelet/bubbletea"
)

type ScreenType int

const (
	Search ScreenType = iota
	Video
)

type RootScreen struct {
	model tea.Model
}

func NewRootScreen(vp *videoplayer.VideoPlayer, sType ScreenType) RootScreen {
	var curr tea.Model

	switch sType {

	case Video:
		vid := initVideoScreen(vp, false)
		curr = &vid

	case Search:
		search := initSearchScreen(vp)
		curr = &search
	}

	return RootScreen{
		model: curr,
	}
}

func (s RootScreen) Init() tea.Cmd {
	return s.model.Init() // rest methods are just wrappers for the model's methods
}

func (s RootScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s.model.Update(msg)
}

func (s RootScreen) View() string {
	return s.model.View()
}

func (s RootScreen) SwitchScreen(model tea.Model) (tea.Model, tea.Cmd) {
	s.model = model
	return s.model, s.model.Init()
}
