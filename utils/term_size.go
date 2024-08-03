package utils

import (
	"fmt"

	"golang.org/x/term"
)

func GetTermSize() (int, int, error) {
	if term.IsTerminal(0) {
		w, h, err := term.GetSize(0)
		if err != nil {
			return 0, 0, fmt.Errorf("Unable to get window size:%e", err)
		}

		return w, h, nil
	}
	return 0, 0, fmt.Errorf("Not a terminal")
}
