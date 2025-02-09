include .env
export

up:
	docker compose up -d

run:
	go run .

redis:
	docker compose exec hubjob-redis redis-cli
