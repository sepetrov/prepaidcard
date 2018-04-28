 .PHONY: clean dev doc exec install logs ps start stop tail-logs test test-all up
.DEFAULT_GOAL:=up

# variables
BINARY:=prepaidcard
PACKAGE:=github.com/sepetrov/prepaidcard
VERSION:=$(shell git -C . describe --abbrev=0 --tags 2> /dev/null || git -C . rev-parse --short HEAD)

# main targets
clean:
	BINARY=$(BINARY) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml rm -fsv
	-BINARY=$(BINARY) PACKAGE=$(PACKAGE) docker rmi -f $(shell PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml images -q 2>/dev/null)

dev:
	BINARY=$(BINARY) PACKAGE=$(PACKAGE) VERSION="$(VERSION)" docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml up --build --remove-orphans -d api

up:
	BINARY=$(BINARY) PACKAGE=$(PACKAGE) VERSION="$(VERSION)" docker-compose -p prepaidcard -f docker-compose.yml up --build --remove-orphans -d api

# helper targets for the host
config:
	@BINARY=$(BINARY) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml config

doc:
	BINARY=$(BINARY) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml up -d doc

exec:
	BINARY=$(BINARY) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml exec api sh

logs:
	BINARY=$(BINARY) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml logs

ps:
	BINARY=$(BINARY) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml  ps

start:
	BINARY=$(BINARY) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml start

stop:
	BINARY=$(BINARY) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml stop

tail-logs:
	BINARY=$(BINARY) PACKAGE=$(PACKAGE) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml logs -f

# helper targets for the container
install:
	go install -v $(PACKAGE)/cmd/$(BINARY)

test:
	go test -v $(PACKAGE)/...

test-all:
	TEST_INTEGRATION=1 go test -v $(PACKAGE)/...