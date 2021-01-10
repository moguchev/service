package utils

import (
	"encoding/json"
	"net/http"

	"github.com/moguchev/service/pkg/logger"
)

// ErrorMessage - answer with error
type ErrorMessage struct {
	Message string `json:"error"`
}

// RespondWithError - answer with error log
func RespondWithError(w http.ResponseWriter, r *http.Request, code int, err error) {
	RespondWithJSON(w, r, code, ErrorMessage{Message: err.Error()})
}

// RespondWithJSON - http json respond
func RespondWithJSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.GetLogger(r.Context()).WithError(err).Error("encode")
	}
}
