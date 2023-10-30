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

# E2E Tests

install-numaflow:
	kubectl create ns numaflow-system
	kubectl apply -n numaflow-system -f https://raw.githubusercontent.com/numaproj/numaflow/stable/config/install.yaml

test-e2e:
	go generate $(shell find ./test/sqs-e2e* -name '*.go')
	go test -v -timeout 15m -count 1 --tags test -p 1 ./test/sqs-e2e*