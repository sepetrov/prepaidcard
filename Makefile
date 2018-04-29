 .PHONY: clean dev doc exec install logs ps start stop tail-logs test test-all up
.DEFAULT_GOAL:=up

# variables
BINARY:=prepaidcard
GOVERSION:=1.10
PACKAGE:=github.com/sepetrov/prepaidcard
VERSION:=$(shell git -C . describe --abbrev=0 --tags 2> /dev/null || git -C . rev-parse --short HEAD)

# main targets
clean:
	BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml rm -fsv
	-BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) docker rmi -f $(shell PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml images -q 2>/dev/null)

dev:
	BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) VERSION="$(VERSION)" docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml up --build -d --remove-orphans api

up:
	BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) VERSION="$(VERSION)" docker-compose -p prepaidcard -f docker-compose.yml up --build --remove-orphans -d api

# helper targets for the host
config:
	@BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml config

doc:
	BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml up -d doc

exec:
	BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml exec api sh

logs:
	BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml logs

ps:
	BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml  ps

start:
	BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml start

stop:
	BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml stop

tail-logs:
	BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml logs -f

# helper targets for the container
install:
	CGO_ENABLED=0 GOOS=linux go install -a -v $(PACKAGE)/cmd/$(BINARY)

test:
	go test -v $(PACKAGE)/...

test-all:
	TEST_INTEGRATION=1 go test -v $(PACKAGE)/...