package cacheservice

import (
	"errors"

	"github.com/disintegration/imaging"
)

type options struct {
	cacheSize    int
	maxImageSize int64
	path         string
	filter       imaging.ResampleFilter
}

type Option interface {
	apply(*options) error
}

// Cache size option.
type cacheSizeOption int

func (o cacheSizeOption) apply(opts *options) error {
	if int(o) <= 0 {
		return errors.New("provided cache size less than 1")
	}
	opts.cacheSize = int(o)
	return nil
}

func WithCacheSize(s uint64) Option {
	return cacheSizeOption(s)
}

// max remote image size option.
type maxImageSizeOption uint

func (c maxImageSizeOption) apply(opts *options) error {
	opt := int64(c)
	if opt <= 0 {
		return errors.New("max image size can't be 0 or negative")
	}
	opts.maxImageSize = int64(c)
	return nil
}

func WithMaxImageSize(s int64) Option {
	return maxImageSizeOption(s)
}

// Filestorage path option.
type pathOption string

func (c pathOption) apply(opts *options) error {
	if string(c) == "" {
		return errors.New("path can not be empty")
	}
	opts.path = string(c)
	return nil
}

func WithPath(s string) Option {
	return pathOption(s)
}

////Resizer Filter option.
type filterOption string

func (c filterOption) apply(opts *options) error {
	filter, ok := filters[string(c)]
	if !ok {
		return errors.New("filter unknown")
	}
	opts.filter = filter
	return nil
}

func WithFilter(s string) Option {
	return filterOption(s)
}

var filters = map[string]imaging.ResampleFilter{
	"NearestNeighbor":   imaging.NearestNeighbor,
	"Box":               imaging.Box,
	"Linear":            imaging.Linear,
	"Hermite":           imaging.Hermite,
	"MitchellNetravali": imaging.MitchellNetravali,
	"CatmullRom":        imaging.CatmullRom,
	"BSpline":           imaging.BSpline,
	"Gaussian":          imaging.Gaussian,
	"Bartlett":          imaging.Bartlett,
	"Lanczos":           imaging.Lanczos,
	"Hann":              imaging.Hann,
	"Hamming":           imaging.Hamming,
	"Blackman":          imaging.Blackman,
	"Welch":             imaging.Welch,
	"Cosine":            imaging.Cosine,
}
