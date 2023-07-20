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
package backup

import (
	"net/http"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/coinbase-samples/waas-proxy-go/cloud/aws/s3"
	"github.com/coinbase-samples/waas-proxy-go/utils"
	log "github.com/sirupsen/logrus"
)

type GetPresignedS3UrlResponse struct {
	Url string `json:"url"`
}

func GetPresignedS3Url(w http.ResponseWriter, r *http.Request) {

	s3Method := utils.HttpPathVarOrSendBadRequest(w, r, "s3Method")
	if len(s3Method) == 0 {
		return
	}

	objectKey := utils.HttpPathVarOrSendBadRequest(w, r, "objectKey")
	if len(objectKey) == 0 {
		return
	}

	ctx := r.Context()
	var httpReq *v4.PresignedHTTPRequest
	var err error
	if s3Method == "putObject" {
		httpReq, err = s3.GeneratePutObjectUrl(ctx, objectKey)
	} else {
		httpReq, err = s3.GenerateGetObjectUrl(ctx, objectKey)
	}
	if err != nil {
		utils.HttpBadGateway(w)
		log.Error(err)
		return
	}

	resp := &GetPresignedS3UrlResponse{
		Url: httpReq.URL,
	}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, resp); err != nil {
		log.Errorf("cannot marshal and write GetPresignedS3UrlResponse: %v", err)
		utils.HttpBadGateway(w)
	}
}
