create-database:
	docker compose exec app sqlite3 data/data.db ".databases"

build-dev:
	docker compose build

up:
	docker compose up -d

dev: build-dev up create-database