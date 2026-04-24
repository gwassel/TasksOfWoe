# TasksOfWoe Agent Instructions

## Run commands
- `make test` - runs lint then tests (required order: lint → test)
- `make lint` - golangci-lint with formatters: gci, gofumpt, goimports, golines
- `make fmt` - format code with golangci-lint
- Local dev: `make up` (or `make upd` for detached) - builds docker-compose with db → migrate → app
- `make down` - stop services
- `make shell` - enter app container

## Architecture
- Entry point: `cmd/service/main.go` (binary: `task-tracker`)
- Layered: handler → usecase → persistence → domain
- Domain packages: `internal/domain/{task,user,analytics,encoder}`
- Telegram bot using `telegram-bot-api/v5`

## Testing
- Uses `go.uber.org/mock` with `gomock.Controller` pattern
- Mock generation with `mockgen` (configured in go.mod toolchain)
- Test files: `*_test.go` alongside source files
- Most tests are unit tests with mocked repositories

## Database
- PostgreSQL 13, migrations in `./migrations/`
- Migrations run via migrate/migrate container in docker-compose
- Connection via env vars: DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT
- Schema: tasks table with encrypted content, users, analytics events

## Environment
- Required: `TELEGRAM_BOT_TOKEN`, `ENCRYPTION_KEY`
- Logs: `/var/log/task-tracker/` (mounted to `./logs/`)
- Go version: 1.25

## Deployment
- CI deploys on push to main: `make down` → git pull → export env vars → `make upd`
- Uses self-hosted runners
