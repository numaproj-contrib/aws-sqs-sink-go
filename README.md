# AWS SQS SINK

AWS SQS Sink for Numaflow implemented in Golang, which will push the vertex data to AWS SQS.

### Environment Variables

Specify the environment variables based on the supported envs specified here https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#specifying-credentials

### Example Pipeline Configuration

```yaml
    - name: aws-sink
      scale:
        min: 1
      sink:
        udsink:
          container:
            image: quay.io/numaio/numaflow-sink/aws-sqs-sink:v0.0.1
            env:
              - name: AWS_SQS_QUEUE_NAME
                value: "test"
              - name: AWS_REGION
                value: "us-east-1"
              - name: AWS_ACCESS_KEY_ID
                value: "testing" ## This can be passed as k8s secret as well
              - name: AWS_SECRET_ACCESS_KEY
                value: "testing" # ## This can be passed as k8s secret as well
```
    