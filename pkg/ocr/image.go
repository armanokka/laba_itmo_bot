package ocr

import (
	"errors"
	"github.com/fogleman/gg"
	"os"
	"strings"
)

func WriteTextOnImage(text, fontFile, pngFilename string) error {
	text = strings.ReplaceAll(text, "\n", " ")
	var (
		width          int     = 1920
		height         int     = 1920
		leftIdentation float64 = 10
	)
	if err := os.Remove(pngFilename); !errors.Is(err, os.ErrNotExist) {
		return err
	}

	dc := gg.NewContext(width, height)
	if err := dc.LoadFontFace(fontFile, 28); err != nil {
		return err
	}
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	strs := dc.WordWrap(text, float64(width))
	for i, s := range strs {
		y := float64((i + 1) * 50)
		dc.DrawString(s, leftIdentation, y)
	}

	return dc.SavePNG(pngFilename)
}
