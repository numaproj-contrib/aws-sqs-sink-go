package sqs_e2e

import (
	"context"
	"testing"
	"time"

	. "github.com/numaproj-contrib/aws-sqs-sink-go/test/fixtures"
	"github.com/stretchr/testify/suite"
)

type AWSSQSSuite struct {
	E2ESuite
}

func (a *AWSSQSSuite) TestAWSSQSSinkPipeline() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	a.T().Log("port forwarding moto service")
	stopPortForward := a.StartPortForward("moto-0", 5000)

	client := awsSQSClient("us-east-1", "testing", "testing", "testing")

	err := client.createSQSQueue(ctx)
	a.NoError(err)
	a.T().Log("sqs queue is created!!!")

	w := a.Given().Pipeline("@testdata/aws-sqs-queue.yaml").
		When().
		CreatePipelineAndWait()
	defer w.DeletePipelineAndWait()

	// wait for all the pods to come up
	w.Expect().VertexPodsRunning()

	// Check for message in aws queue, it will retry for 10 times with 2-second sleep,
	// if message is available then it will return true otherwise false
	containMsg, err := client.isQueueContainMessages(ctx)
	a.NoError(err)

	a.True(containMsg)

	// stop the port forwarding of moto service
	stopPortForward()
}

func TestAWSSQSSinkPipelineSuite(t *testing.T) {
	suite.Run(t, new(AWSSQSSuite))
}
