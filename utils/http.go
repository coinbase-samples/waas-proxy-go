package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func HttpReadBodyOrSendGatewayTimeout(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		HttpGatewayTimeout(w)
	}

	return body, err
}

func HttpPathVarOrSendBadRequest(w http.ResponseWriter, r *http.Request, name string) string {

	vars := mux.Vars(r)

	v, found := vars[name]

	if !found {
		HttpBadRequest(w)
	}

	return v
}

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

func HttpMarshalAndWriteJsonResponseWithOk(w http.ResponseWriter, v any) error {
	return HttpMarshalAndWriteJsonResponseWithStatus(w, v, http.StatusOK)
}

func HttpMarshalAndWriteJsonResponseWithStatus(w http.ResponseWriter, v any, status int) error {
	protoArg, ok := v.(proto.Message)
	// TODO: revisit when I can only use proto.Message
	if !ok {
		body, err := json.Marshal(v)
		if err != nil {
			return err
		}

		return HttpWriteJsonResponseWithStatus(w, body, status)
	}
	body, err := protojson.Marshal(protoArg)
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
