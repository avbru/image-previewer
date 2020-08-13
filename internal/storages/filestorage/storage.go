package filestorage

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

type FileStore struct {
	path string
}

func New(path string) (fileStore *FileStore, err error) {
	err = os.MkdirAll(path, os.ModeDir)
	if err != nil {
		return nil, fmt.Errorf("stogare cann't create dir: %w", err)
	}

	return &FileStore{
		path: path,
	}, nil
}

func (s *FileStore) SaveFile(ctx context.Context, fName string, reader io.Reader) error {
	file, err := os.Create(s.path + "/" + fName)
	if err != nil {
		return fmt.Errorf("stogare cann't create file: %w", err)
	}
	defer fClose(file)

	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("stogare cann't read/save image to file: %w", err)
	}

	return nil
}

func (s *FileStore) ReadFile(ctx context.Context, fName string, w io.Writer) error {
	file, err := os.Open(s.path + "/" + fName)
	if err != nil {
		return fmt.Errorf("stogare ReadFile cann't open file: %w", err)
	}
	defer fClose(file)

	_, err = io.Copy(w, file)
	if err != nil {
		return fmt.Errorf("stogare ReadFile cann't read file: %w", err)
	}

	return nil
}

func (s *FileStore) DeleteFile(ctx context.Context, fName string) {
	removeFile(ctx, s.path+"/"+fName)
}

func (s *FileStore) Purge(ctx context.Context) error {
	files, err := ioutil.ReadDir(s.path)
	if err != nil {
		return err
	}

	for _, f := range files {
		// Remove only jpg files for a safety reason
		if !strings.HasSuffix(f.Name(), ".jpg") {
			continue
		}
		if err := os.Remove(s.path + "/" + f.Name()); err != nil {
			return err
		}
	}

	return nil
}

func (s *FileStore) Volume(ctx context.Context) (int64, error) {
	files, err := ioutil.ReadDir(s.path)
	if err != nil {
		return 0, err
	}
	var size int64
	for _, file := range files {
		if file.Mode().IsRegular() {
			size += file.Size()
		}
	}
	return size, nil
}

func removeFile(ctx context.Context, fName string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, err := os.Stat(fName)
			if os.IsNotExist(err) {
				return
			}
			if err := os.Remove(fName); err == nil {
				return
			}
		}
	}
}

func fClose(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Err(err)
	}
}
