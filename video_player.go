package main

import (
	"fmt"

	ascii "github.com/AndriiPets/terminal_yt/image_manipulation"
)

type VideoPlayer struct {
	Video *Video
}

func NewVideoPlayer() *VideoPlayer {
	return &VideoPlayer{}
}

func (vp *VideoPlayer) LoadFromURL(url string) error {
	data, url := getStream(url)
	video := NewVideo(url, data)
	err := video.init()

	if err != nil {
		fmt.Println("video init error:", err)
		return err
	}
	vp.Video = video

	for video.Read() {
		frame, err := ascii.Byte2ascii2(video.framebuffer, data.width, data.heigth, ascii.AsciiTableSimple)

		if video.EOF {
			frame = "EOF"
		}

		if err != nil {
			fmt.Printf("unable to encode frame No:%d", video.frameCounter)
			return err
		}
		video.frameMap[video.frameCounter] = frame
	}

	return nil
}
