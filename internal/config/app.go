package config

type Application struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"`
	Server   Server
}
