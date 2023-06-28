/**
 * Copyright 2023 Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type httpRequestPageInfo struct {
	Size  int32
	Token string
}

const (
	httpRequestPageInfoTokenParam = "pageToken"
	httpRequestPageInfoSizeParam  = "pageSize"
)

func (h httpRequestPageInfo) Passed() bool {
	return h.Size > 0 || len(h.Token) > 0
}

func HttpRequestPageInfo(r *http.Request) (pageInfo httpRequestPageInfo, err error) {

	query := r.URL.Query()
	pageInfo = httpRequestPageInfo{Token: query.Get(httpRequestPageInfoTokenParam)}

	pageSize := query.Get(httpRequestPageInfoSizeParam)
	if len(pageSize) > 0 {
		var i int
		if i, err = strconv.Atoi(pageSize); err == nil && i > 0 {
			pageInfo.Size = int32(i)
		}
	}

	return
}

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
		err = fmt.Errorf("unable to write json response body %w", err)
	}
	return
}
