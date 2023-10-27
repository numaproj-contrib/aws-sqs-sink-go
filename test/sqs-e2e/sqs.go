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

package sqs_e2e

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"time"
)

type sqsClient struct {
	client    *sqs.Client
	queueName string
}

func (s *sqsClient) createSQSQueue(ctx context.Context) error {
	if _, err := s.client.CreateQueue(ctx, &sqs.CreateQueueInput{QueueName: aws.String(s.queueName)}); err != nil {
		return err
	}

	return nil
}

func awsSQSClient(region, accessKey, accessSecret, queueName string) *sqsClient {
	awsConfig := aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentialsProvider(accessKey, accessSecret, ""),
	}

	client := sqs.NewFromConfig(awsConfig, func(options *sqs.Options) {
		options.BaseEndpoint = aws.String("http://localhost:5000")
	})

	return &sqsClient{
		client:    client,
		queueName: queueName,
	}
}

func (s *sqsClient) isQueueContainMessages(ctx context.Context) (bool, error) {
	queueURL, err := s.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(s.queueName)})
	if err != nil {
		return false, err
	}

	for i := 0; i < 10; i++ {
		message, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{QueueUrl: queueURL.QueueUrl, MaxNumberOfMessages: 2})
		if err != nil {
			fmt.Println(err)
		}

		messageCount := 0
		for _, msg := range message.Messages {
			fmt.Println("message received:", *msg.Body)
			messageCount++
		}

		if messageCount > 0 {
			return true, nil
		}

		// wait for 2 second and poll the message from queue again
		time.Sleep(2 * time.Second)
	}

	return false, errors.New("retry exceeded, queue doesn't have any message")
}
