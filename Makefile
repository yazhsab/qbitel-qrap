.PHONY: all setup build test lint clean migrate docker-deps docker-all docker-down

# ==== Setup ====
setup: setup-go setup-python setup-node
	@echo "Setup complete"

setup-go:
	go work sync

setup-python:
	cd ml && python3 -m venv .venv && .venv/bin/pip install -e ".[dev]"

setup-node:
	cd web && npm install

# ==== Build ====
build: build-go build-node

build-go:
	cd api && go build ./...

build-node:
	cd web && npm run build

# ==== Test ====
test: test-go test-python test-node

test-go:
	cd api && go test ./...
	cd shared/go && go test ./...

test-python:
	cd ml && .venv/bin/pytest tests/ -v

test-node:
	cd web && npx tsc --noEmit

# ==== Lint ====
lint: lint-go lint-python

lint-go:
	cd api && golangci-lint run ./...
	cd shared/go && golangci-lint run ./...

lint-python:
	cd ml && .venv/bin/ruff check .

# ==== Database Migrations ====
migrate:
	migrate -path db/migrations -database "$${QRAP_DATABASE_URL}" up

migrate-down:
	migrate -path db/migrations -database "$${QRAP_DATABASE_URL}" down 1

# ==== Docker ====
docker-deps:
	docker compose -f infra/docker/docker-compose.deps.yml up -d

docker-all:
	docker compose -f infra/docker/docker-compose.yml up -d

docker-down:
	docker compose -f infra/docker/docker-compose.yml down

docker-build:
	docker compose -f infra/docker/docker-compose.yml build

# ==== Clean ====
clean:
	go clean -cache
	rm -rf web/dist web/node_modules/.cache ml/.venv
