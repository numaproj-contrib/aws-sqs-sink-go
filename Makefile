DOCKERIO_ORG=quay.io/numaio
PLATFORM=linux/x86_64
TARGET=sqs-sink

IMAGE_TAG=$(TAG)
ifeq ($(IMAGE_TAG),)
	IMAGE_TAG=dev
endif

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/aws-sqs-sink main.go

.PHONY: image
image: build
	docker buildx build -t "$(DOCKERIO_ORG)/numaproj-contrib/aws-sqs-sink-go:$(IMAGE_TAG)" --platform $(PLATFORM) --target $(TARGET) . --load

.PHONY: lint
lint:
	go mod tidy
	golangci-lint run --fix --verbose --concurrency 4 --timeout 5m

clean:
	-rm -rf ./dist
