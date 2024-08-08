package cmd

import (
	"fmt"

	"github.com/AndriiPets/terminal_yt/ui"
	videoplayer "github.com/AndriiPets/terminal_yt/video_player"
	"github.com/spf13/cobra"
)

var urlCmd = &cobra.Command{
	Use:     "url <URL>",
	Aliases: []string{"URL"},
	Short:   "Load video from URL",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vPlayer := videoplayer.NewVideoPlayer()
		err := vPlayer.LoadVideoMetadata(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		ui.RunTUI(vPlayer, false)
	},
}

func init() {
	rootCmd.AddCommand(urlCmd)
}
