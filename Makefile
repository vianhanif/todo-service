
IMAGE := todo-api

# This version-strategy uses git tags to set the version string
VERSION := $(shell git describe --tags --always --dirty)

migrate:
	@go run -v ./databases/migrate.go

test:
	@go test -v -cover -p 1 ./...

container:
	@echo "building the image..."
	@docker build --label "version=$(VERSION)" -t $(IMAGE):$(VERSION) -f ./docker/Dockerfile .

push:
	@docker tag $(IMAGE):$(VERSION) localhost:5000/$(IMAGE):$(VERSION)
	@docker push $(IMAGE):$(VERSION)