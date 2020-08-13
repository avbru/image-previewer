package storages

import (
	"context"
	"io"
)

type FileStorage interface {
	SaveFile(ctx context.Context, fName string, reader io.Reader) error
	ReadFile(ctx context.Context, fName string, w io.Writer) error
	DeleteFile(ctx context.Context, fName string)
	Purge(ctx context.Context) error
	Volume(ctx context.Context) (int64, error)
}
