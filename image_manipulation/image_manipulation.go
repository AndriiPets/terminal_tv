package ascii

import (
	"image"
	"image/draw"
	"math"
	"strings"

	"github.com/disintegration/imaging"
)

var (
	AsciiTableSimple   = " .:-=+*#%@"
	AsciiTableDetailed = " .'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	CODES              = "Ñ@#W$9876543210?!abc;:+=-,._          "
	CHARS              = []byte(" .,:;i1tfLCG08@")
)

type Frame struct {
	Content string
	Width   int
	Height  int
}

func Byte2ascii2(raw []byte, w, h, termW, termH int, code string) (Frame, error) {
	// imgData, _, err := image.Decode(bytes.NewReader(raw))
	// if err != nil {
	// 	return "", fmt.Errorf("can't decode frame: %v", err)
	// }
	//Convert raw frame data to image.RGBA
	imgRect := image.Rect(0, 0, w, h)
	imgData := image.RGBA{
		Pix:    raw,
		Stride: 4 * imgRect.Dx(),
		Rect:   imgRect,
	}

	var sb strings.Builder
	//TODO: resize
	smallImg, err := resizeImage(&imgData, termW, termH, 2)
	if err != nil {
		return Frame{}, nil
	}

	//Create new imaga with rezized proportions, draw in original img
	rect := smallImg.Bounds()
	rgba := image.NewNRGBA(rect)
	draw.Draw(rgba, rect, smallImg, rect.Min, draw.Src)
	imgW, imgH := rect.Max.X, rect.Max.Y

	//extract color data
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			index := (y*imgW + x) * 4
			pix := rgba.Pix[index : index+4]
			r := pix[0]
			g := pix[1]
			b := pix[3]

			brightness := float64(r + g + b/3)
			charCode := int(brightness / 255 * float64(len(code)-1))
			sb.WriteString(string(code[charCode]))
		}
		sb.WriteString("\n")
	}

	return Frame{sb.String(), rect.Dx(), rect.Dy()}, nil
}

func resizeImage(img *image.RGBA, termW, termH, scale int) (image.Image, error) {
	imgW := float64(img.Bounds().Dx())
	imgH := float64(img.Bounds().Dy())
	aspect := imgW / imgH
	//termW, termH, _ := utils.GetTermSize()

	width := math.Min(float64(termH*scale), float64(termW))
	height := width / aspect

	smallImg := imaging.Resize(img, int(width), int(height), imaging.NearestNeighbor)

	return smallImg, nil
}
