package resizer

import (
	"bufio"
	"bytes"
	"image"
	"image/jpeg"
	"io"

	"github.com/disintegration/imaging"
)

func Resize(width, height int, reader io.Reader) (buffer io.Reader, err error) {
	original, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	resized := imaging.Fill(original, width, height, imaging.Center, imaging.Linear)

	var buf bytes.Buffer

	if err := jpeg.Encode(bufio.NewWriter(&buf), resized, nil); err != nil {
		return nil, err
	}

	return &buf, nil
}
