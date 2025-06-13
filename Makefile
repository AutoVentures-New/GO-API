include .env
export

up:
	docker compose up -d

run:
	go run . api

cronjob:
	go run . cronjob

build:
	go build .