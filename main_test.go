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
		t.Setenv("AWS_SQS_QUEUE_NAME", "unit_test")
		sqsClient, err := newSQSSinkConfig(ctx)
		assert.NoError(t, err)

		assert.NotEmpty(t, sqsClient.sqsClient)
		assert.NotEmpty(t, sqsClient.queueURL)
	})
}
