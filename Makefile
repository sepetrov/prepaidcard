.PHONY: \
	clean \
	config \
	config-dev \
	dev \
	doc \
	exec \
	help \
	install \
	logs \
	ps \
	query-db \
	query-testdb \
	start \
	stop \
	tail-logs \
	test \
	test-integration \
	test-unit \
	up \

.DEFAULT_GOAL:=help

# variables
BINARY:=prepaidcard
GOVERSION:=1.10
VERSION=unknown
PACKAGE:=github.com/sepetrov/prepaidcard

# target-specific variables
config doc up: VERSION:=$(shell git -C . describe --abbrev=0 --tags 2> /dev/null || git -C . rev-parse --short HEAD)
doc: DOC_PORT:=$(shell grep DOC_PORT .env 2> /dev/null | sed -e 's/DOC_PORT\s*=\s*\(.*\)/\1/g')
query-db query-testdb: DB_ROOT_PASSWORD:=$(shell grep DB_ROOT_PASSWORD .env 2> /dev/null | sed -e 's/DB_ROOT_PASSWORD\s*=\s*\(.*\)/\1/g')

##
##  * Host targets
##

clean:        ## Remove Docker containers and images
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml rm -fsv
	-docker rmi -f $(shell docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml images -q 2>/dev/null)

config:       ## Show Docker configuration for development mode
	@BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) VERSION=$(VERSION) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml config

config-dev:   ## Show Docker configuration for production mode
	@BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) VERSION=$(VERSION) docker-compose -p prepaidcard -f docker-compose.yml config

dev:          ## Build and start Docker conainer in development mode
	BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) VERSION=$(VERSION) docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml up --build -d --remove-orphans api
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml exec api sh -c "make install && $(BINARY)"

doc:          ## Build and start Swagger API container
	VERSION=$(VERSION) docker-compose -p prepaidcard -f docker-compose.yml up -d doc
	-[ -z "$(DOC_PORT)" ] || open http://localhost:$(DOC_PORT) 2> /dev/null

exec:         ## SSH into API container running in development mode
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml exec api sh

logs:         ## Show Docker container logs
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml logs

ps:           ## Show Docker container status
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml  ps

query-db:     ## Open MySQL client
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml exec db mysql -uroot -p'$(DB_ROOT_PASSWORD)' $(BINARY)

query-testdb: ## Open MySQL client to test databse
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml exec testdb mysql -uroot -p'$(DB_ROOT_PASSWORD)' $(BINARY)_test

start:        ## Start Docker containers in development mode
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml start

stop:         ## Stop Docker container
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml stop

tail-logs:    ## Show Docker container logs continuously
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml logs -f

up:           ## Build and start Docker containers in prodution mode
	BINARY=$(BINARY) GOVERSION=$(GOVERSION) PACKAGE=$(PACKAGE) VERSION=$(VERSION) docker-compose -p prepaidcard -f docker-compose.yml up --build --remove-orphans -d api

##
##  * Container targets
##

install:          ## Install application binary
	CGO_ENABLED=0 go install -a -ldflags '-s -w' -v $(PACKAGE)/cmd/$(BINARY)

test:             ## Run tests
	CGO_ENABLED=0 go test -a -ldflags '-s -w' -v $(PACKAGE)/...

test-integration: ## Run integration tests
	go test -tags=integration -v $(PACKAGE)/...

test-unit:        ## Run unit tests
	go test -tags=unit -v $(PACKAGE)/...


##
##  * Help
##

help:    ## Show this help message
	@echo
	@echo '  Usage:'
	@echo '    make <target>'
	@echo
	@echo '  Targets:'
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
	@echo

version: ## Print the version.
	@echo $(VERSION)