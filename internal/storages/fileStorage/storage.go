package filestorage

import (
	"context"
	"io"
	"io/ioutil"
	"os"
)

type FileStore struct {
	path string
}

func New(path string) (fileStore *FileStore, err error) {
	// TODO insure dir
	return &FileStore{
		path: path,
	}, nil
}

func (s *FileStore) SaveFile(ctx context.Context, fName string, reader io.Reader) error {
	file, err := os.Open(s.path + fName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	return nil
}

func (s *FileStore) ReadFile(ctx context.Context, fName string, w io.Writer) error {
	file, err := os.Open(s.path + fName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		return err
	}

	return nil
}

func (s *FileStore) DeleteFile(ctx context.Context, fName string) {
	go removeFile(ctx, s.path+fName)
}

func (s *FileStore) Purge(ctx context.Context) error {
	files, err := ioutil.ReadDir(s.path)
	if err != nil {
		return err
	}
	for _, v := range files {
		go removeFile(ctx, s.path+v.Name())
	}
	return nil
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
