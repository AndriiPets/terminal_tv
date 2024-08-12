package ui

import (
	"fmt"
	"os"
	"time"

	videoplayer "github.com/AndriiPets/terminal_yt/video_player"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	default_w = 160
	default_h = 90
)

type (
	ErrMsg  error
	TickMsg time.Time
)

func RunTUI(vp *videoplayer.VideoPlayer, sType ScreenType) {

	p := tea.NewProgram(NewRootScreen(vp, sType))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
