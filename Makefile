include .env
export

up:
	docker compose up -d

run:
	go run . api

cronjob:
	go run . cronjob

redis:
	docker compose exec hubjob-redis redis-cli

build:
	go build .