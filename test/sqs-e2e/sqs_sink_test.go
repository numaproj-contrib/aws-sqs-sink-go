//go:build test

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

	client, err := awsSQSClient(ctx)
	a.NoError(err)

	err = client.createSQSQueue(ctx, "testing")
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
	containMsg, err := client.isQueueContainMessages(ctx, "testing")
	a.NoError(err)

	a.True(containMsg)

	// stop the port forwarding of moto service
	stopPortForward()
}

func TestAWSSQSSinkPipelineSuite(t *testing.T) {
	suite.Run(t, new(AWSSQSSuite))
}
