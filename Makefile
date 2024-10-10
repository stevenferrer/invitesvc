GOPATH ?= $(shell go env GOPATH)
GOBIN ?= $(GOPATH)/bin
GOOS ?=linux"
GOARCH ?=amd64

IMAGE_TAG=0.1.0
IMAGE_NAME=invitesvc

.PHONY: build
build:
	go build -v -ldflags "-w -s" -o ./cmd/invitesvc ./cmd/invitesvc

.PHONY: test
test:
	go test -v -cover -race ./...

.PHONY: postgres
postgres:
	docker rm -f invite-postgres || true
	docker run --name invite-postgres -e POSTGRES_PASSWORD=postgres -d --rm -p 5432:5432 postgres:13
	docker exec -it invite-postgres bash -c 'while ! pg_isready; do sleep 1; done;'

.PHONY: build-image
build-image:
	docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .