.PHONY: infra-up infra-down api web test vet lint build-web check

infra-up:
	docker compose up -d

infra-down:
	docker compose down

api:
	go run ./cmd/api

web:
	cd apps/web && npm run dev

test:
	go test ./cmd/... ./internal/...

vet:
	go vet ./cmd/... ./internal/...

lint:
	cd apps/web && npm run lint

build-web:
	cd apps/web && npm run build

check: test vet lint build-web
