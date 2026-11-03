# Commands

## Start server

```bash
go run cmd/api/main.go
```

## Create Migration

```bash
migrate create -ext sql -dir migrations -seq create_users_table
```

## Docker

### Start Container

```bash
docker compose up -d
```

### Stop Container

```bash
docker compose down -v
```