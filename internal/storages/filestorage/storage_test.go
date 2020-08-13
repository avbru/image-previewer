package filestorage

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var dir = "images_test"

func Test_SaveFile(t *testing.T) {
	store, _ := New(dir)
	err := store.SaveFile(context.Background(), "image.jpg", strings.NewReader("some image"))
	require.NoError(t, err)
	require.FileExists(t, dir+"/image.jpg")
	clean(t)
}

func Test_ReadFile(t *testing.T) {
	store, _ := New(dir)
	err := store.SaveFile(context.Background(), "image.jpg", strings.NewReader("some image"))
	require.NoError(t, err)
	var buf bytes.Buffer
	err = store.ReadFile(context.Background(), "image.jpg", bufio.NewWriter(&buf))
	require.NoError(t, err)
	clean(t)
}

func Test_DeleteFile(t *testing.T) {
	store, _ := New(dir)
	err := store.SaveFile(context.Background(), "image.jpg", strings.NewReader("some image"))
	require.NoError(t, err)
	// Emulate slow reader...
	go func() {
		file, _ := os.Open(dir + "/image.jpg")
		time.Sleep(time.Millisecond * 200)
		_ = file.Close()
	}()
	require.FileExists(t, dir+"/image.jpg")
	store.DeleteFile(context.Background(), "image.jpg")
	//
	// Check file deleted
	//time.Sleep(time.Millisecond * 400)
	require.NoFileExists(t, dir+"/image.jpg")
	clean(t)
}

func Test_Volume(t *testing.T) {
	store, _ := New(dir)
	_ = store.SaveFile(context.Background(), "image.jpg", strings.NewReader("some image"))
	_ = store.SaveFile(context.Background(), "image2.jpg", strings.NewReader("some image"))
	volume, err := store.Volume(context.Background())
	require.NoError(t, err)
	require.Equal(t, int64(20), volume)
	clean(t)
}

func clean(t *testing.T) {
	err := os.RemoveAll(dir)
	require.NoError(t, err)
	require.NoDirExists(t, dir)
}
