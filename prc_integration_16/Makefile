.PHONY: up down test-compose logs all-tests

up:
	docker-compose up -d

down:
	docker-compose down

test-compose:
	go test -v -run DockerCompose ./integration -db-dsn="postgres://test:test@localhost:54321/notes_test?sslmode=disable"

logs:
	docker-compose logs -f

all-tests:
	docker-compose up -d
	go test -v -run DockerCompose ./integration -db-dsn="postgres://test:test@localhost:54321/notes_test?sslmode=disable"
	docker-compose down