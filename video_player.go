package main

import (
	"fmt"
	"sync"

	ascii "github.com/AndriiPets/terminal_yt/image_manipulation"
)

type VideoPlayer struct {
	Video *Video
}

func NewVideoPlayer() *VideoPlayer {
	return &VideoPlayer{}
}

func (vp *VideoPlayer) LoadVideoMetadata(url string) error {
	data, url := getStream(url)
	video := NewVideo(url, data)

	vp.Video = video

	return nil
}

func (vp *VideoPlayer) StartStream() error {
	video := vp.Video
	data := video.data

	err := video.init()

	if err != nil {
		fmt.Println("video init error:", err)
		return err
	}

	var wg sync.WaitGroup

	for video.Read() {
		wg.Add(1)

		go func(video *Video) {
			defer wg.Done()
			frame, _ := ascii.Byte2ascii2(video.framebuffer, data.width, data.heigth, ascii.AsciiTableSimple)

			video.frameMap.Store(video.frameCounter, frame)
		}(video)
	}

	wg.Wait()

	return nil
}
