up:
	mkdir -p my-dynamodb-data
	docker compose up -d
	sleep 5
	go run cmd/bootstrapinfra/main.go
.PHONY: up

down:
	docker compose stop
.PHONY: down

down-prune:
	docker compose exec redis redis-cli FLUSHALL
	docker compose stop
	docker compose rm -f -s -v
	rm my-dynamodb-data/*
.PHONY: down-prune

docker-infra:
	docker compose up redis dynamodb -d
	sleep 5
	go run cmd/bootstrapinfra/main.go
.PHONY: docker-infra

docker-bash:
	docker compose exec gameserver /bin/bash
.PHONY: docker-bash

docker-run-client:
	docker compose run -e REMOTE_ADDR=gameserver:50051 --rm gameserver go run cmd/client/main.go
.PHONY: docker-run-client

run-client:
	go run cmd/client/main.go
.PHONY: run-client

docker-redis:
	docker compose exec redis redis-cli
.PHONY: docker-redis

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	internal/ports/grpc/proto/battle/battle.proto \
	internal/ports/grpc/proto/leaderboard/leaderboard.proto \
	internal/ports/grpc/proto/player/player.proto 
.PHONY: proto

build:
	docker compose build
.PHONY: build

test:
	go test -v ./... -short
.PHONY: test