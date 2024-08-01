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

func (vp *VideoPlayer) LoadFromURL(url string) error {
	data, url := getStream(url)
	video := NewVideo(url, data)
	err := video.init()

	if err != nil {
		fmt.Println("video init error:", err)
		return err
	}
	vp.Video = video

	var wg sync.WaitGroup

	for video.Read() {
		wg.Add(1)

		go func(video *Video) {
			defer wg.Done()
			frameBuilder, _ := ascii.Byte2ascii2(video.framebuffer, data.width, data.heigth, ascii.AsciiTableSimple)
			frame := frameBuilder.String()

			if video.EOF {
				frame = "EOF"
			}

			// if err != nil {
			// 	fmt.Printf("unable to encode frame No:%d", video.frameCounter)
			// 	return err
			// }
			video.frameMap.Store(video.frameCounter, frame)
		}(video)
	}

	wg.Wait()

	return nil
}
