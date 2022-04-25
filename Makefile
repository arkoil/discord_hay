APP_NAME = bot

build:
	go build -o ./bin/${APP_NAME} cmd/bot/main.go