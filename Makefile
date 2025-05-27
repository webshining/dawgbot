discord:
	go run cmd/discord/main.go
telegram:
	go run cmd/telegram/main.go
pg_dump:
	mkdir -p ./data/backups/postgres && docker compose exec -T postgres pg_dump -U postgres postgres --no-owner \
	| gzip -9 > ./data/backups/postgres/backup-$(shell date +%Y-%m-%d_%H-%M-%S).sql.gz
rebuild:
	docker compose up -d --no-deps --force-recreate --build