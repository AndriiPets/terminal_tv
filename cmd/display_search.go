package cmd

import (
	"github.com/AndriiPets/terminal_yt/ui"
	videoplayer "github.com/AndriiPets/terminal_yt/video_player"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:     "search <Query>",
	Aliases: []string{"Search"},
	Short:   "Search videos from query",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vPlayer := videoplayer.NewVideoPlayer()
		ui.RunTUI(vPlayer, ui.Search)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
