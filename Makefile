.PHONY: run
run: build
	./tgbot

build:
	go build -o tgbot ./cmd/bot/main.go