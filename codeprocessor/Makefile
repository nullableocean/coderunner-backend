CERTS_DIR=./docker/certs
COMPOSE_FILE=docker-compose.yml

.PHONY: all certsgen up down clean build

build-with-certs: certsgen build

build:
	docker compose build

up-all:
	docker compose -f $(COMPOSE_FILE) up -d

up-docdemon:
	docker compose -f $(COMPOSE_FILE) up docdemon -d

up-processor:
	docker compose -f $(COMPOSE_FILE) up code-processor -d

restart-processor:
	docker compose -f $(COMPOSE_FILE) restart code-processor -d --no-deps --build

down:
	docker compose -f $(COMPOSE_FILE) down

certsgen: clean
	mkdir -p $(CERTS_DIR)
	@bash/certsgen.sh $(CERTS_DIR)

clean: down
	rm -rf $(CERTS_DIR)

run: build-with-certs up-all
