package videoplayer

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

type VideoData struct {
	Width    int
	Heigth   int
	Depth    int
	Fps      float64
	Duration string
}

type Video struct {
	stream       string
	Data         VideoData
	Framebuffer  []byte
	frameCounter int
	FrameMap     sync.Map
	pipe         io.ReadCloser
	cmd          *exec.Cmd
	EOF          bool
}

func NewVideo(file string, data VideoData) *Video {
	return &Video{
		Data:   data,
		stream: file,
	}
}

// TODO: allow to read data from the stream insted of the file
func (v *Video) init() error {
	//scale_factor := 8
	v.cleanup()
	// ffmpeg command to pipe video data to stdout in 8-bit RGBA format.
	cmd := exec.Command(
		"ffmpeg",
		"-i", v.stream,
		"-f", "image2pipe",
		"-loglevel", "quiet",
		"-pix_fmt", "rgba",
		"-vcodec", "rawvideo",
		"-",
	)

	v.cmd = cmd
	pipe, err := cmd.StdoutPipe()

	if err != nil {
		return err
	}
	v.pipe = pipe

	if err := cmd.Start(); err != nil {
		return err
	}

	if v.Framebuffer == nil {
		v.Framebuffer = make([]byte, v.Data.Width*v.Data.Heigth*v.Data.Depth)
	}
	return nil
}

func (v *Video) Read() bool {
	// If cmd is nil, video reading has not been initialized.
	if v.cmd == nil {
		if err := v.init(); err != nil {
			return false
		}
	}

	if _, err := io.ReadFull(v.pipe, v.Framebuffer); err != nil {
		fmt.Println("error:", err)

		if err.Error() == "EOF" {
			v.EOF = true
		}

		v.Close()
		return false
	}
	v.frameCounter++

	return true
}

func (v *Video) Close() {
	v.EOF = false
	v.frameCounter = 0
	if v.pipe != nil {
		v.pipe.Close()
	}

	if v.cmd != nil {
		v.cmd.Wait()
	}
}

func (v *Video) cleanup() {
	v.frameCounter = 0
	v.EOF = false
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if v.pipe != nil {
			v.pipe.Close()
		}
		if v.cmd != nil {
			v.cmd.Process.Kill()
		}
		os.Exit(1)
	}()
}
