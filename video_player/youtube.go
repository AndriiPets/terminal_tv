package videoplayer

import (
	"io"
	"os"

	"github.com/kkdai/youtube/v2"
)

func getStream(url string) (VideoData, string, error) {
	client := youtube.Client{}

	video, err := client.GetVideo(url)
	if err != nil {
		return VideoData{}, "", err
	}

	formats := video.Formats.Quality("medium")
	format := &formats[0]
	stream, err := client.GetStreamURL(video, format)
	if err != nil {
		return VideoData{}, "", err
	}
	vData := VideoData{
		Width:    format.Width,
		Heigth:   format.Height,
		Depth:    4, //idk this value was in the example
		Fps:      float64(format.FPS),
		Duration: format.ApproxDurationMs,
	}

	return vData, stream, nil

}

func download_video(url string) VideoData {
	client := youtube.Client{}

	video, err := client.GetVideo(url)
	if err != nil {
		panic(err)
	}

	// Typically youtube only provides separate streams for video and audio.
	// If you want audio and video combined, take a look a the downloader package.
	formats := video.Formats.Quality("medium")
	format := &formats[0]
	stream, _, err := client.GetStream(video, format)

	vData := VideoData{
		Width:  format.Width,
		Heigth: format.Height,
		Depth:  4, //idk this value was in the example
		Fps:    float64(format.FPS),
	}
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	file, err := os.Create("video.mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}

	return vData
}
