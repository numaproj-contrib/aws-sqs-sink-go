# AWS SQS Sink for Numaflow

The AWS SQS Sink is a custom user-defined sink for [Numaflow](https://numaflow.numaproj.io/) that enables the integration of Amazon Simple Queue Service (SQS) as a sink within your Numaflow pipelines.

## Quick Start
This quick start guide will walk you through setting up an AWS SQS sink in a Numaflow pipeline.

### Prerequisites
* [Install Numaflow on your Kubernetes cluster](https://numaflow.numaproj.io/quick-start/)
* [AWS CLI configured with access to AWS SQS](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-quickstart.html)

### Step-by-step Guide

#### 1. Create an AWS SQS Queue

Using AWS CLI or the AWS Management Console, [create a new SQS queue](https://docs.aws.amazon.com/cli/latest/reference/sqs/create-queue.html)

#### 2. Deploy a Numaflow Pipeline with AWS SQS Sink

- Save the following Kubernetes manifest to a file (e.g., `sqs-sink-pipeline.yaml`)
- Modifying the AWS region, queue name accordingly
- Specify the AWS credentials using any supported approach https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#specifying-credentials


```yaml
apiVersion: numaflow.numaproj.io/v1alpha1
kind: Pipeline
metadata:
  name: simple-pipeline
spec:
  vertices:
    - name: in
      source:
        generator:
          # How many messages to generate in the duration.
          rpu: 1
          duration: 5s
          # Optional, size of each generated message, defaults to 10.
          msgSize: 10
    - name: aws-sink
      scale:
        min: 1
      sink:
        udsink:
          container:
            image: quay.io/numaio/numaproj-contrib/aws-sqs-sink-go:v0.0.1
            env:
              - name: AWS_SQS_QUEUE_NAME
                value: "test-queue"
              - name: AWS_REGION
                value: "us-east-1"
              - name: AWS_ACCESS_KEY_ID
                value: "testing" ## This can be passed as k8s secret as well
              - name: AWS_SECRET_ACCESS_KEY
                value: "testing" ## This can be passed as k8s secret as well
  edges:
    - from: in
      to: aws-sink
```

Then apply it to your cluster:
```bash
kubectl apply -f sqs-sink-pipeline.yaml
```

#### 5. Verify the Log Sink

Verify the message using aws cli to poll message, refer: https://docs.aws.amazon.com/cli/latest/reference/sqs/receive-message.html

#### 6. Clean up

To delete the Numaflow pipeline:
```bash
kubectl delete -f sqs-sink-pipeline.yaml
```

To delete the SQS queue:
```bash
aws sqs delete-queue --queue-url <YourQueueUrl>
```

Congratulations!!! You have successfully set up an AWS SQS sink in a Numaflow pipeline.

## Additional Resources

For more information on Numaflow and how to use it to process data in a Kubernetes-native way, visit the [Numaflow Documentation](https://numaflow.numaproj.io/). For AWS SQS specific configuration, refer to the [AWS SQS Documentation](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/welcome.html).