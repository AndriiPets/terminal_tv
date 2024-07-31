package ascii

type Pixel struct {
	r uint8
	g uint8
	b uint8
}

func byte_to_ascii(rawFrame []byte, h, w int) string {
	var i int
	var ascii_string string

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			index := (y*w + x) * 4
			p := rawFrame[index : index+4]
			pix := Pixel{
				r: p[0],
				g: p[1],
				b: p[2],
			}
			i += 4
			char := getAsciiChar(pix, AsciiTableSimple, true)
			ascii_string += char
		}
		ascii_string += "\n"
	}
	return ascii_string
}

func getAsciiChar(pix Pixel, code string, invert bool) string {
	var charCode int
	brightness := float64(pix.r + pix.g + pix.b/3)

	if invert {
		charCode = int(mapRange(brightness, 0, 225, 0, float64(len(code)-1)))
	} else {
		//charCode = int(mapRange(brightness, 0, 225, float64(len(code)), 0))
		charCode = int(brightness / 255 * float64(len(code)-1))
	}
	return string(code[charCode])
}

func mapRange(value, minInput, maxInput, minOutput, maxOutput float64) float64 {
	//Clamping value to the range
	if value < minInput {
		value = minInput
	} else if value > maxInput {
		value = maxInput
	}

	proportion := (value - minInput) / (maxInput - minInput)

	// Map the proportion to the output range
	newValue := proportion*(maxOutput-minOutput) + minOutput

	return newValue
}
