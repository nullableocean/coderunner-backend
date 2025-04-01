CERTS_DIR=./docker/certs
COMPOSE_FILE=docker-compose.yml

.PHONY: all certsgen up down clean

up-all:
	docker compose -f $(COMPOSE_FILE) up -d --build

up-docdemon:
	docker compose -f $(COMPOSE_FILE) up docdemon -d --build

up-processor:
	docker compose -f $(COMPOSE_FILE) up code-processor -d --build

restart-processor:
	docker compose -f $(COMPOSE_FILE) restart code-processor -d --no-deps --build

down:
	docker compose -f $(COMPOSE_FILE) down

certsgen: clean
	mkdir -p $(CERTS_DIR)
	@bash/certsgen.sh $(CERTS_DIR)

clean: down
	rm -rf $(CERTS_DIR)

run: gencerts up-all
