package config

import "time"

type Server struct {
	Listen       string        `env:"API_MODER_LISTEN" envDefault:":8083"`
	ReadTimeout  time.Duration `env:"API_MODER_READ_TIMEOUT" envDefault:"30s"`
	WriteTimeout time.Duration `env:"API_MODER_WRITE_TIMEOUT" envDefault:"30s"`
}
