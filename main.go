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
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
	"github.com/numaproj/numaflow/pkg/shared/logging"
	"go.uber.org/zap"
)

const (
	sqsQueueName   = "AWS_SQS_QUEUE_NAME"
	awsEndpointURL = "AWS_ENDPOINT_URL"
)

type awsSQSSink struct {
	logger    *zap.SugaredLogger
	sqsClient *sqs.Client
}

// awsSQSClient will generate the aws sqs client with default aws config.
func awsSQSClient(ctx context.Context) *awsSQSSink {
	// initialize the logger
	logger := logging.NewLogger().Named("aws-sqs-sink")

	// Load default configs for aws based on env variable provided based on
	// https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#specifying-credentials
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.Fatalln("failed loading aws config, err: ", err)
	}

	// generate the sqs client based on default values or if AWS_ENDPOINT_URL is passed as env.
	awsEndpoint := os.Getenv(awsEndpointURL)

	var client *sqs.Client
	if awsEndpoint != "" {
		client = sqs.NewFromConfig(cfg, func(options *sqs.Options) {
			options.BaseEndpoint = aws.String(awsEndpoint)
		})
	} else {
		client = sqs.NewFromConfig(cfg)
	}

	return &awsSQSSink{
		logger:    logger,
		sqsClient: client,
	}
}

// Sink will publish the vertex data to aws sqs sink
func (s *awsSQSSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	responses := sinksdk.ResponsesBuilder()

	// generate the queue url to publish data to queue via queue name.
	queueURL, err := s.sqsClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(os.Getenv(sqsQueueName))})
	if err != nil {
		s.logger.Fatalln("failed to generate SQS Queue url, err: ", err)
	}

	// generate message request entries for processing message in a batch
	var messageRequests []sqsTypes.SendMessageBatchRequestEntry
	for datum := range datumStreamCh {
		messageRequests = append(messageRequests, sqsTypes.SendMessageBatchRequestEntry{
			Id:          aws.String(datum.ID()),
			MessageBody: aws.String(string(datum.Value())),
		})
	}

	// send batch message to aws queue
	response, err := s.sqsClient.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
		Entries:  messageRequests,
		QueueUrl: queueURL.QueueUrl,
	})
	if err != nil {
		s.logger.Errorf("failed to push message %v", err)
	}

	// log the failure response and append that to responses object.
	if len(response.Failed) > 0 {
		s.logger.Error("failed to push message, err: %v", response.Failed)

		for _, fail := range response.Failed {
			responses = responses.Append(sinksdk.ResponseFailure(aws.ToString(fail.Id), "failed to push message"))
		}
	}

	// append the success response to responses object
	for _, success := range response.Successful {
		responses = responses.Append(sinksdk.ResponseOK(aws.ToString(success.Id)))
	}

	return responses
}

func main() {
	ctx := context.Background()

	// generate aws sqs queue client based on provided config.
	configs := awsSQSClient(ctx)

	// start a new sink server which will push data to aws sqs queue.
	if err := sinksdk.NewServer(configs).Start(ctx); err != nil {
		configs.logger.Fatalln("failed to start aws sqs sink server, err: %v", err)
	}
}
