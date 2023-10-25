# AWS SQS SINK

AWS SQS Sink for Numaflow implemented in Golang, which will push the vertex data to AWS SQS.

### Environment Variables
	AWS_SQS_QUEUE_NAME      : AWS SQS Queue name for push the vertex data
	AWS_REGION              : Region for AWS SQS 

### Secret Variables
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: aws-secret
type: Opaque
data:
  AWS_ACCESS_KEY: "<BASE_64_OF_AWS_ACCESS_KEY>" # Replace with aws access key in base64 format
  AWS_ACCESS_SECRET: "<BASE_64_AWS_ACCESS_SECRET>" # Replace with aws access secret in base64 format
```

### Example Pipeline Configuration

```yaml
    - name: aws-sink
      scale:
        min: 1
      sink:
        udsink:
          container:
            image: quay.io/numaio/numaflow-sink/aws-sqs-sink:v0.0.3
            env:
              - name: AWS_SQS_QUEUE_NAME
                value: "test-queue" # Replace with aws sqs queue name
              - name: AWS_REGION
                value: "eu-north-1" # Replace with aws sqs queue region
            envFrom:
              - secretRef:
                  name: aws-secret # Replace with k8s secret name created earlier
```
    