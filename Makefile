.PHONY: clean dev doc logs ps start stop tail-logs up
.DEFAULT_GOAL := up

# main targets
clean:
	docker-compose -p prepaidcard -f docker-compose.yml kill
	docker-compose -p prepaidcard -f docker-compose.yml rm -vf
	docker rmi -f $(shell docker-compose -p prepaidcard -f docker-compose.yml images -q 2>/dev/null) 2>/dev/null || true

dev:
	docker-compose -p prepaidcard -f docker-compose.yml -f docker-compose.override.yml up --remove-orphans -d doc

up:
	docker-compose -p prepaidcard -f docker-compose.yml up --remove-orphans -d doc

# helper targets
doc:
	docker-compose -p prepaidcard -f docker-compose.yml up -d doc

logs:
	docker-compose -p prepaidcard -f docker-compose.yml logs

ps:
	docker-compose -p prepaidcard -f docker-compose.yml ps

start:
	docker-compose -p prepaidcard -f docker-compose.yml start

stop:
	docker-compose -p prepaidcard -f docker-compose.yml stop

tail-logs:
	docker-compose -p prepaidcard -f docker-compose.yml logs -f
