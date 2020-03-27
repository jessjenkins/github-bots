SHELL=bash

VERSION:=$(shell git describe --tags --dirty)

.PHONY: docker
docker:
	 docker build --tag github-bots:$(VERSION) .

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: run
run: docker
	docker run --rm --env-file .env -p 8085:8085 --name github-bots github-bots:$(VERSION)