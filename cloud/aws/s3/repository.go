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
package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/coinbase-samples/waas-proxy-go/config"
	log "github.com/sirupsen/logrus"
)

var repo *Repository

type Repository struct {
	App    *config.AppConfig
	Client *s3.Client
}

func InitRepo(ctx context.Context, a *config.AppConfig) {
	cfg, err := awsConfig.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatalf("unable to get aws config: %v", err)
	}

	repo = &Repository{
		App:    a,
		Client: s3.NewFromConfig(cfg),
	}
}

func GenerateGetObjectUrl(
	ctx context.Context,
	objectKey string,
) (*v4.PresignedHTTPRequest, error) {
	presignClient := s3.NewPresignClient(repo.Client)
	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(repo.App.BucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(repo.App.PresignedUrlExpiration * int64(time.Second))
	})
	if err != nil {
		return nil, fmt.Errorf(
			"generate presigned get object url failed for object: %s err: %w",
			objectKey,
			err,
		)
	}
	return request, nil
}

func GeneratePutObjectUrl(
	ctx context.Context,
	objectKey string,
) (*v4.PresignedHTTPRequest, error) {
	presignClient := s3.NewPresignClient(repo.Client)
	request, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(repo.App.BucketName),
		ContentType: aws.String("application/octet-stream"),
		Key:         aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(repo.App.PresignedUrlExpiration * int64(time.Second))
	})
	if err != nil {
		return nil, fmt.Errorf(
			"generate presigned put object url failed for object: %s err: %w",
			objectKey,
			err,
		)
	}
	return request, nil
}
