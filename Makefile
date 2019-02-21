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

#
# Variables
#

# The name of the binary.
BINARY:=prepaidcard

# The Go version.
GOVERSION:=1.10

# The name of the package.
PACKAGE:=github.com/sepetrov/prepaidcard

# The application version.
#
# The value is automatically generated using Git. You should never override it manually.
#
# Falls back gracefully to one of the values:
# - [tag] if HEAD is tagged
# - [tag]-[number of commits after the last tag]-g[short git commit hash] if HEAD has an offset from a tag
# - [git commit hash]
# - "unknown"
VERSION:=$(shell \
	git -C . describe --tags 2> /dev/null || \
	git -C . rev-parse --short HEAD 2> /dev/null || \
	echo "unknown" \
)

# Override the variables with the environment variables from .env if it exists.
-include .env

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
	CGO_ENABLED=0 go install -a -ldflags "-s -w -X '$(PACKAGE)/pkg/api/api.Version=$(VERSION)'" -v $(PACKAGE)

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