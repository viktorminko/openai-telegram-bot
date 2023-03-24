.PHONY: build run stop test lint

IMAGE_NAME := telegram-bot-api
IMAGE_TAG := telegram-bot-api

build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

run: build
	docker run -it --env-file .env $(IMAGE_NAME):$(IMAGE_TAG)

stop:
	docker stop $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)

test:
	go test -v ./...

lint:
	golangci-lint run
