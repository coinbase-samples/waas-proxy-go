package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func HttpBadGateway(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadGateway)
}
func HttpOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func HttpBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func HttpGatewayTimeout(w http.ResponseWriter) {
	w.WriteHeader(http.StatusGatewayTimeout)
}

func HttpJsonContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func HttpWriteJsonResponseWithStatusOk(w http.ResponseWriter, body []byte) (err error) {
	return HttpWriteJsonResponseWithStatus(w, body, http.StatusOK)
}

func HttpWriteJsonResponseWithStatusCreated(w http.ResponseWriter, body []byte) (err error) {
	return HttpWriteJsonResponseWithStatus(w, body, http.StatusCreated)
}

func MarhsallAndWriteJsonResponseWithOk(w http.ResponseWriter, v any) error {
	return HttpMarhsallAndWriteJsonResponseWithStatus(w, v, http.StatusOK)
}

func HttpMarhsallAndWriteJsonResponseWithStatus(w http.ResponseWriter, v any, status int) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return HttpWriteJsonResponseWithStatus(w, body, status)
}

func HttpWriteJsonResponseWithStatus(w http.ResponseWriter, body []byte, status int) (err error) {
	HttpJsonContentType(w)
	w.WriteHeader(status)
	_, err = io.WriteString(w, string(body))

	if err != nil {
		err = fmt.Errorf("Unable to write json response body %w", err)
	}
	return
}
