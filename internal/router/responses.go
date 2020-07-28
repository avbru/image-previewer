package router

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

//func httpJSON(w http.ResponseWriter, v interface{}) {
//	data, err := json.Marshal(v)
//	if err != nil {
//		httpInternalServerError(w, "failed to encode response body", err)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	_, err = w.Write(data)
//	if err != nil {
//		log.Error().Msg("error writing response")
//	}
//}

func httpBadRequest(w http.ResponseWriter, msg string, err error) {
	httpError(w, http.StatusBadRequest, "400 bad request: "+msg, err)
	log.Warn().Msg(msg + err.Error())
}

func httpMethodNotAllowed(w http.ResponseWriter, msg string, err error) {
	httpError(w, http.StatusMethodNotAllowed, "405 method not allowed: "+msg, err)
}

func httpInternalServerError(w http.ResponseWriter, msg string, err error) {
	httpError(w, http.StatusInternalServerError, msg, err)
}

func httpRemoteServerError(w http.ResponseWriter, msg string, err error) {
	httpError(w, http.StatusUnprocessableEntity, msg, err)
}

func httpError(w http.ResponseWriter, httpStatus int, msg string, err error) {
	http.Error(w, msg, httpStatus)
	log.Warn().Msgf("%s: %v", msg, err)
}
