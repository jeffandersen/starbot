build:
	GOOS=linux GOARCH=amd64 go build -v -o bin/starbot .

heroku: build
all: build
