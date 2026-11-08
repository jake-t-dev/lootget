.PHONY: run build clean test
COMPOSE_FILE=docker-compose.yml

run:
	go run cmd/main.go

build:
	go build -o bin/app cmd/main.go

clean:
	rm -rf bin/

up:
	docker-compose -f $(COMPOSE_FILE) up -d --build

down:
	docker-compose -f $(COMPOSE_FILE) down