package resizer

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"

	"github.com/disintegration/imaging"
)

func Resize(width, height int, filter imaging.ResampleFilter, reader io.Reader) (buffer io.Reader, err error) {
	original, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("resize can't decode image %w", err)
	}

	resized := imaging.Fill(original, width, height, imaging.Center, filter)

	var buf bytes.Buffer

	if err := jpeg.Encode(bufio.NewWriter(&buf), resized, nil); err != nil {
		return nil, fmt.Errorf("resize can't encode image %w", err)
	}

	return &buf, nil
}
