.PHONY: build build-codeprocessor build-all up up-codeprocessor up-all tests swagger down-all

SRC_DIR = ./src
SRC_MAIN_PATH = cmd/main.go
DOCS_DIR = $(SRC_DIR)/docs

CODEPROCESSOR_DIR=./codeprocessor
CODEPROCESSOR_COMPOSE=./codeprocessor/docker-compose.yml
CODEPROCESSOR_BUILD_CMD=make build-with-certs
CODEPROCESSOR_UP_CMD=make up-all
CODEPROCESSOR_DOWN_CMD=make down


build:
	@echo "Собираем API..."
	docker compose build

build-codeprocessor:
	@echo "Собираем сервис Codeprocessor..."
	cd ${CODEPROCESSOR_DIR} && ${CODEPROCESSOR_BUILD_CMD}

build-all: build build-codeprocessor
		@echo "Образы собраны"
up:
	@echo "Поднимаем API..."
	docker compose up -d

up-codeprocessor: 
	@echo "Поднимаем сервис Codeprocessor..."
	cd ${CODEPROCESSOR_DIR} && ${CODEPROCESSOR_UP_CMD}

up-all: up up-codeprocessor
		@echo "Сервисы запущены"

down-all:
		docker compose down
		cd ${CODEPROCESSOR_DIR} && ${CODEPROCESSOR_DOWN_CMD}

swagger:
	@echo "Генерация Swagger-документации..."
	swag init -d $(SRC_DIR) -g $(SRC_MAIN_PATH)  -o $(DOCS_DIR)
	@echo "Swagger-документация сгенерирована в $(DOCS_DIR)"