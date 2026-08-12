.PHONY: infra-up infra-down migrate-up migrate-down api web test vet lint build-web check production-config production-build production-up production-down production-status production-logs production-smoke production-verify

infra-up:
	docker compose up -d

infra-down:
	docker compose down

migrate-up:
	migrate -path migrations -database "$$DATABASE_URL" up

migrate-down:
	migrate -path migrations -database "$$DATABASE_URL" down 1

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

production-config:
	docker compose --env-file .env.production -f docker-compose.production.yml config >/dev/null

production-build:
	docker compose --env-file .env.production -f docker-compose.production.yml build

production-up:
	docker compose --env-file .env.production -f docker-compose.production.yml up -d

production-down:
	docker compose --env-file .env.production -f docker-compose.production.yml down

production-status:
	docker compose --env-file .env.production -f docker-compose.production.yml ps

production-logs:
	docker compose --env-file .env.production -f docker-compose.production.yml logs

production-smoke:
	bash scripts/production-smoke-test.sh

production-verify:
	bash scripts/verify-production.sh
