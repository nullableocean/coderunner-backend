SRC_DIR = ./src
SRC_MAIN_PATH = cmd/main.go
DOCS_DIR = $(SRC_DIR)/docs

swagger:
	@echo "Генерация Swagger-документации..."
	swag init -d $(SRC_DIR) -g $(SRC_MAIN_PATH)  -o $(DOCS_DIR)
	@echo "Swagger-документация сгенерирована в $(DOCS_DIR)"