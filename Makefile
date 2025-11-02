NAME := front-service

build-linux:
	go mod download
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./build/${NAME} ./cmd/${NAME}/main.go

