package handlers

import (
	"fmt"
	"io"
	"net/http"
)

func httpCreated(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}

func httpBadGateway(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadGateway)
}
func httpOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func httpBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func httpGatewayTimeout(w http.ResponseWriter) {
	w.WriteHeader(http.StatusGatewayTimeout)
}

func jsonContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func writeJsonResponseWithStatusOk(w http.ResponseWriter, body []byte) (err error) {
	return writeJsonResponseWithStatus(w, body, http.StatusOK)
}

func writeJsonResponseWithStatusCreated(w http.ResponseWriter, body []byte) (err error) {
	return writeJsonResponseWithStatus(w, body, http.StatusCreated)
}

func writeJsonResponseWithStatus(w http.ResponseWriter, body []byte, status int) (err error) {
	jsonContentType(w)
	w.WriteHeader(status)
	_, err = io.WriteString(w, string(body))

	if err != nil {
		err = fmt.Errorf("Unable to write json response body %w", err)
	}
	return
}
