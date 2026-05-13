.PHONY: test build compose-up compose-down migrate run create-key protected-ping revoke-key

test:
	go test ./...

build:
	go build ./...

compose-up:
	./scripts/generate-dev-pepper.sh
	docker compose up --build

compose-down:
	docker compose down -v

migrate:
	go run ./cmd/apikey-migrate up

run:
	go run ./cmd/sample-service

create-key:
	./scripts/curl-create-key.sh

protected-ping:
	./scripts/curl-protected-ping.sh

revoke-key:
	./scripts/curl-revoke-key.sh
