/*
Copyright 2023 The Numaproj Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
	"github.com/numaproj/numaflow/pkg/shared/logging"
	"go.uber.org/zap"
)

type awsSQSSink struct {
	logger          *zap.SugaredLogger
	sqsClient       *sqs.Client
	queueName       string
	region          string
	awsAccessKey    string
	awsAccessSecret string
	awsBaseEndpoint string
}

// newAWSSQSSink will read the environment variable required to authenticate aws and publish data to SQS
func newAWSSQSSink() *awsSQSSink {
	logger := logging.NewLogger().Named("aws-sqs-sink")
	queueName, ok := os.LookupEnv("AWS_SQS_QUEUE_NAME")
	if !ok {
		logger.Fatalln("AWS_SQS_QUEUE_NAME not found")
	}
	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		logger.Fatalln("AWS_REGION not found")
	}
	accessKey, ok := os.LookupEnv("AWS_ACCESS_KEY")
	if !ok {
		logger.Fatalln("AWS_ACCESS_KEY not found")
	}
	accessSecret, ok := os.LookupEnv("AWS_ACCESS_SECRET")
	if !ok {
		logger.Fatalln("AWS_ACCESS_SECRET not found")
	}

	return &awsSQSSink{
		logger:          logger,
		queueName:       queueName,
		region:          region,
		awsAccessKey:    accessKey,
		awsAccessSecret: accessSecret,
		awsBaseEndpoint: os.Getenv("AWS_BASE_ENDPOINT"),
	}
}

// awsSQSClient will generate the aws sqs client using access_key, access_secret and region.
func (s *awsSQSSink) awsSQSClient() *sqs.Client {
	config := aws.Config{
		Region:      s.region,
		Credentials: credentials.NewStaticCredentialsProvider(s.awsAccessKey, s.awsAccessSecret, ""),
	}

	if s.awsBaseEndpoint != "" {
		return sqs.NewFromConfig(config, func(options *sqs.Options) {
			options.BaseEndpoint = aws.String(s.awsBaseEndpoint)
		})
	}

	return sqs.NewFromConfig(config)
}

// Sink will publish the vertex data to aws sqs sink
func (s *awsSQSSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	ok := sinksdk.ResponsesBuilder()

	// generate the queue url to publish data to queue via queue name.
	queueURL, err := s.sqsClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: &s.queueName})
	if err != nil {
		s.logger.Fatalln("failed to generate SQS Queue url, err: ", err)
	}

	// generate message request entries for processing message in a batch
	requestEntries := make([]sqsTypes.SendMessageBatchRequestEntry, len(datumStreamCh))
	count := 0
	for datum := range datumStreamCh {
		requestEntries[count] = sqsTypes.SendMessageBatchRequestEntry{
			Id:          aws.String(datum.ID()),
			MessageBody: aws.String(string(datum.Value())),
		}
		count++
	}

	// send batch message to aws queue
	response, err := s.sqsClient.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
		Entries:  requestEntries,
		QueueUrl: queueURL.QueueUrl,
	})
	if err != nil {
		s.logger.Errorf("failed to push message %v", err)
	} else {
		// check for response in case of partial failure, then log those responses otherwise return id of success request.
		if len(response.Failed) > 0 {
			s.logger.Error("failed to push message, err: %v", response.Failed)
		}
		if len(response.Successful) > 0 {
			for _, success := range response.Successful {
				ok = ok.Append(sinksdk.ResponseOK(aws.ToString(success.Id)))
			}
		}
	}

	return ok
}

func main() {
	// generate aws sqs sink configuration using aws_access_key, aws_access_secret, region and queue_name.
	config := newAWSSQSSink()

	// generate aws sqs client using region, aws_access_key and aws_access_secret, used for push data to queue.
	config.sqsClient = config.awsSQSClient()

	// start a new sink server which will push data to aws sqs queue.
	if err := sinksdk.NewServer(config).Start(context.Background()); err != nil {
		config.logger.Fatalln("failed to start aws sqs sink server, err: %v", err)
	}
}
