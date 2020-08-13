package config

type Config struct {
	Port         uint16 `env:"PREVIEWER_PORT" envDefault:"80"`
	CacheSize    uint64 `env:"PREVIEWER_CACHE_SIZE" envDefault:"100"`
	ImageDir     string `env:"PREVIEWER_IMAGE_DIR" envDefault:"images"`
	MaxImageSize int64  `env:"PREVIEWER_MAX_IMAGE_SIZE" envDefault:"5242880"`
	Filter       string `env:"PREVIEWER_FILTER" envDefault:"Linear"`
}
