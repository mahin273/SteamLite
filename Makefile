.PHONY: all test run-auth clean infra-up infra-down infra-logs infra-status infra-clean

all: test

test:
	go test -v streamlite.local/...

run-auth:
	go run ./services/auth/cmd/server/main.go

infra-up:
	docker compose -f deploy/docker/docker-compose.yml --env-file .env up -d

infra-down:
	docker compose -f deploy/docker/docker-compose.yml --env-file .env down

infra-logs:
	docker compose -f deploy/docker/docker-compose.yml --env-file .env logs -f

infra-status:
	docker compose -f deploy/docker/docker-compose.yml --env-file .env ps

infra-clean:
	docker compose -f deploy/docker/docker-compose.yml --env-file .env down -v

clean:
	rm -rf bin/ tmp/ scratch/

