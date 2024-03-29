# Makefile for OCPP Example project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD)fmt

# Folders definitions
BINARY_FOLDER=bin

# Get a short hash of the git had for building images.
VERSION = 1.0
TAG = $$(git rev-parse --short HEAD)
APP_NAME=ocpp_example
IMAGE_NAME = "${APP_NAME}"
DOCKERFILE = "Dockerfile"
PROJECT_PATH = $(shell pwd)

define docker-build =
	docker build -t ${IMAGE_NAME}:${VERSION} -f dev/${DOCKERFILE} .
	docker run --rm \
		--name toolbox-ocppexample-tests \
		-v ${PROJECT_PATH}:/go/ocppexample \
		${IMAGE_NAME}:${VERSION} \
		make -C /go/ocppexample $(1)
endef

define docker-run =
	docker run --rm \
		--name toolbox-ocppexample-tests \
		-v ${PROJECT_PATH}:/go/ocppexample \
		-v ${PROJECT_PATH}/logs:/tmp/logs \
		-v ${PROJECT_PATH}/example/configs.json:/tmp/configs.json \
		-p "9033:8080" \
		${IMAGE_NAME}:${VERSION} \
		$(1)
endef


clean:
	$(GOCLEAN)
	rm -rf $(BINARY_FOLDER)

modinit:
	$(GOMOD) init github.com/CoderSergiy/ocpp16-go

deps:
	export GO111MODULE=on
	$(GOGET) github.com/julienschmidt/httprouter
	$(GOGET) github.com/gorilla/websocket
	$(GOGET) github.com/CoderSergiy/golib/logging
	$(GOGET) github.com/CoderSergiy/golib/timelib
	$(GOGET) github.com/CoderSergiy/golib/tools
	$(GOGET) github.com/google/uuid

fmts:
	$(GOFMT) -s -d ./core/*.go
	$(GOFMT) -s -d ./messages/*.go
	$(GOFMT) -s -d ./example/*.go
	$(GOFMT) -s -d *.go

build:
	$(GOBUILD) -o $(BINARY_FOLDER)/server -v ./server.go

depsupdate:
	$(GOGET) -v -t ./...

test:
	@CGO_ENABLED=0 $(GOTEST) -v ./...

buildall: depsupdate fmts test build

# Command using Docker container

dockerclean:
	$(call docker-build, "clean")

dockerprojectsetup:
	# Call once when setup project
	$(call docker-build, "modinit")

dockerbuild:
	$(call docker-build, "buildall")

dockerserverrun:
	$(call docker-build, "build")
	$(call docker-run, 	 "$(BINARY_FOLDER)/server")
