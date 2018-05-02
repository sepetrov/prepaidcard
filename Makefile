 .PHONY: clean config config-dev dev doc exec install logs ps query-db query-testdb start stop tail-logs test test-all up
.DEFAULT_GOAL:=up

# variables
BINARY:=prepaidcard
GOVERSION:=1.10
VERSION=unknown
PACKAGE:=github.com/sepetrov/prepaidcard

# target-specific variables
config doc up: VERSION:=$(shell git -C . describe --abbrev=0 --tags 2> /dev/null || git -C . rev-parse --short HEAD)
doc: DOC_PORT:=$(shell grep DOC_PORT .env 2> /dev/null | sed -e 's/DOC_PORT\s*=\s*\(.*\)/\1/g')
query-db query-testdb: DB_ROOT_PASSWORD:=$(shell grep DB_ROOT_PASSWORD .env 2> /dev/null | sed -e 's/DB_ROOT_PASSWORD\s*=\s*\(.*\)/\1/g')

# main targets
clean:
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml rm -fsv
	-docker rmi -f $(shell docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml images -q 2>/dev/null)

dev:
	BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) VERSION=$(VERSION) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml up --build -d --remove-orphans api
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml exec api sh -c "make install && $(BINARY)"

up:
	BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) VERSION=$(VERSION) docker-compose -p prepaidcard -f docker-compose.yml up --build --remove-orphans -d api

# helper targets for the host
config:
	@BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) VERSION=$(VERSION) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml config

config-dev:
	@BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) VERSION=$(VERSION) docker-compose -p prepaidcard -f docker-compose.yml config

doc:
	VERSION=$(VERSION) docker-compose -p prepaidcard -f docker-compose.yml up -d doc
	-[ -z "$(DOC_PORT)" ] || open http://localhost:$(DOC_PORT) 2> /dev/null

exec:
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml exec api sh

logs:
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml logs

ps:
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml  ps

query-db:
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml exec db mysql -uroot -p'$(DB_ROOT_PASSWORD)' $(BINARY)

query-testdb:
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml exec testdb mysql -uroot -p'$(DB_ROOT_PASSWORD)' $(BINARY)_test

start:
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml start

stop:
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml stop

tail-logs:
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml logs -f

# helper targets for the container
install:
	CGO_ENABLED=0 go install -a -ldflags '-s -w' -v $(PACKAGE)/cmd/$(BINARY)

test:
	CGO_ENABLED=0 go test -a -ldflags '-s -w' -v $(PACKAGE)/...

test-all:
	go test -tags=integration -v $(PACKAGE)/...
