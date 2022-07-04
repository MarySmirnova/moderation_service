package api

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MarySmirnova/moderation_service/internal/config"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

type ContextKey string

const ContextReqIDKey ContextKey = "request_id"

type Moderator struct {
	httpServer     *http.Server
	forbiddenWords map[string]struct{}
}

func NewModerator(cfg config.Server, forbiddenWords []string) *Moderator {
	m := &Moderator{
		forbiddenWords: map[string]struct{}{},
	}

	for _, word := range forbiddenWords {
		m.forbiddenWords[word] = struct{}{}
	}

	handler := mux.NewRouter()
	handler.Use(m.reqIDMiddleware, m.logMiddleware)
	handler.Name("moderate").Methods(http.MethodPost).Path("/moderate").HandlerFunc(m.ModerateHandler)

	m.httpServer = &http.Server{
		Addr:         cfg.Listen,
		Handler:      handler,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
	}

	return m
}

func (m *Moderator) GetHTTPServer() *http.Server {
	return m.httpServer
}

func (m *Moderator) reqIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqID int
		reqIDString := r.FormValue("request_id")

		if reqIDString == "" {
			reqID = m.generateReqID()
		}

		if reqIDString != "" {
			id, err := strconv.Atoi(reqIDString)
			if err != nil {
				m.writeResponseError(w, err, http.StatusBadRequest)
				return
			}
			reqID = id
		}

		ctx := context.WithValue(r.Context(), ContextReqIDKey, reqID)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Moderator) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			log.WithFields(log.Fields{
				"request_time": time.Now().Format("2006-01-02 15:04:05.000000"),
				"request_ip":   strings.TrimPrefix(strings.Split(r.RemoteAddr, ":")[1], "["),
				"code":         w.Header().Get("Code"),
				"request_id":   r.Context().Value(ContextReqIDKey),
			}).Info("news reader response")
		}()

		next.ServeHTTP(w, r)
	})
}

func (m *Moderator) generateReqID() int {
	max := 999999999999
	min := 100000

	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func (m *Moderator) getPageAndFilterParams(w http.ResponseWriter, r *http.Request) (int, string, error) {
	var page int
	filter := r.FormValue("filter")
	pageString := r.FormValue("page")
	if pageString == "" {
		page = 1
	}
	if pageString != "" {
		p, err := strconv.Atoi(pageString)
		if err != nil {
			return 0, "", err
		}
		page = p
	}

	return page, filter, nil
}

func (m *Moderator) writeResponseError(w http.ResponseWriter, err error, code int) {
	w.Header().Add("Code", strconv.Itoa(code))
	log.WithError(err).Error("api error")
	w.WriteHeader(code)
	_, _ = w.Write([]byte(err.Error()))
}
