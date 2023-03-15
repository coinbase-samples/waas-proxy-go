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

func jsonContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func writeJsonResponsWithStatusCreated(w http.ResponseWriter, body string) (err error) {
	jsonContentType(w)
	httpCreated(w)
	_, err = io.WriteString(w, body)

	if err != nil {
		err = fmt.Errorf("Unable to write response body %w", err)
	}
	return
}
