package videoplayer

import (
	"fmt"
	"sync"

	ascii "github.com/AndriiPets/terminal_yt/image_manipulation"
	"github.com/AndriiPets/terminal_yt/utils"
)

type VideoPlayer struct {
	Video  *Video
	Width  int
	Heigth int
}

func NewVideoPlayer() *VideoPlayer {
	termW, termH, _ := utils.GetTermSize()
	return &VideoPlayer{
		Width:  termW,
		Heigth: termH,
	}
}

func (vp *VideoPlayer) LoadVideoMetadata(url string) error {
	data, url, err := getStream(url)
	if err != nil {
		return err
	}
	video := NewVideo(url, data)

	vp.Video = video

	return nil
}

func (vp *VideoPlayer) StartStream() error {
	video := vp.Video
	data := video.Data

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
			frame, _ := ascii.Byte2ascii2(video.Framebuffer, data.Width, data.Heigth, vp.Width, vp.Heigth, ascii.AsciiTableSimple)

			video.FrameMap.Store(video.frameCounter, frame)
		}(video)
	}

	wg.Wait()

	return nil
}
