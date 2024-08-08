package main

import videoplayer "github.com/AndriiPets/terminal_yt/video_player"

func main() {
	vPlayer := videoplayer.NewVideoPlayer()
	//vPlayer.LoadVideoMetadata("https://www.youtube.com/watch?v=dQw4w9WgXcQ") //https://www.youtube.com/watch?v=dQw4w9WgXcQ
	//go vPlayer.StartStream()                                                 //https://www.youtube.com/watch?v=_C6PbG5cH14&list=LL&index=3

	RunTUI(vPlayer)
}
