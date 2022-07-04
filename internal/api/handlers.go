package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

var ErrBlocked = errors.New("comment blocked by moderator")

type Comment struct {
	Text string // тело комментария
}

//ModerateHandler проверяет комментарий на содержание запрещенных слов.
//Если такие слова присутствуют, возвращает 400 ошибку.
func (m *Moderator) ModerateHandler(w http.ResponseWriter, r *http.Request) {
	var comment Comment

	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		m.writeResponseError(w, err, http.StatusBadRequest)
		return
	}

	words := strings.Fields(comment.Text)
	for _, word := range words {
		if _, ok := m.forbiddenWords[word]; ok {
			m.writeResponseError(w, ErrBlocked, http.StatusBadRequest)
			return
		}
	}

	w.Header().Add("Code", strconv.Itoa(http.StatusNoContent))
	w.WriteHeader(http.StatusNoContent)
}
