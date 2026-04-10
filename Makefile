.PHONY: up seed swagger

# Сгенерировать Swagger-документацию
swagger:
	swag init -g cmd/main.go -o docs

# Поднять всё через docker-compose
up:
	docker-compose up --build -d

# Наполнить БД тестовыми данными
seed:
	go run cmd/seed/main.go