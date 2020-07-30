package storages

import (
	"context"
	"io"
)

type FileStore interface {
	SaveFile(ctx context.Context, fName string, reader io.Reader) error
	ReadFile(ctx context.Context, fName string, w io.Writer) error
	DeleteFile(ctx context.Context, fName string) error
	Purge(ctx context.Context)
}
