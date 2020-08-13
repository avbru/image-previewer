package resizer

import (
	"image"
	"os"
	"testing"

	"github.com/disintegration/imaging"

	"github.com/stretchr/testify/require"
)

var resizeTests = []struct {
	file      string
	width     int
	height    int
	wantError bool
}{
	{
		file:      "gopher1.jpg",
		width:     100,
		height:    100,
		wantError: false,
	},
	{
		file:      "gopher1.jpg",
		width:     500,
		height:    500,
		wantError: false,
	},
	{
		file:      "1.txt",
		width:     100,
		height:    100,
		wantError: true,
	},
}

func Test_Resizer(t *testing.T) {
	for _, tt := range resizeTests {
		tCase := tt
		t.Run(tCase.file, func(t *testing.T) {
			file, err := os.Open("../../test/images/" + tCase.file)
			require.NoError(t, err)
			defer file.Close()

			buf, err := Resize(tCase.width, tCase.height, imaging.Linear, file)
			if tCase.wantError == true {
				return
			}
			require.NoError(t, err)

			img, _, err := image.Decode(buf)
			require.NoError(t, err)
			require.Equal(t, tCase.width, img.Bounds().Max.X)
			require.Equal(t, tCase.height, img.Bounds().Max.Y)
		})
	}
}
