include .env
export

up:
	docker compose up -d

run:
	go run .
