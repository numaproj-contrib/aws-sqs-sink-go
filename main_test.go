package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_awsSQSSink(t *testing.T) {
	t.Run("Generate SQS Client", func(t *testing.T) {
		ctx := context.TODO()

		t.Setenv("AWS_REGION", "us-east-1")
		sqsClient := awsSQSClient(ctx)

		assert.NotEmpty(t, sqsClient.logger)
		assert.NotEmpty(t, sqsClient.sqsClient)
	})
}
