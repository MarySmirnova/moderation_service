package internal

import (
	"net/http"

	"github.com/MarySmirnova/moderation_service/internal/api"
	"github.com/MarySmirnova/moderation_service/internal/config"

	log "github.com/sirupsen/logrus"
)

type Application struct {
	cfg            config.Application
	forbiddenWords []string
}

func NewApplication(cfg config.Application) *Application {
	return &Application{
		cfg: cfg,
		forbiddenWords: []string{
			"qwerty",
			"йцукен",
			"zxvbnm",
		},
	}
}

func (a *Application) StartServer() {
	srv := api.NewModerator(a.cfg.Server, a.forbiddenWords)
	s := srv.GetHTTPServer()

	log.WithField("listen", s.Addr).Info("start server")

	err := s.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		log.WithError(err).Error("the channel raised an error")
		return
	}

	log.Info("server has been stoped")
}
