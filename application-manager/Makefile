install:
	go get -u

build:
	go build -o bin/appmanager/main cmd/appmanager/main.go

run:
	go run cmd/appmanager/main.go

compile:
	GOOS=darwin GOARCH=amd64 go build -o bin/application-server-mac cmd/appmanager/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/application-server cmd/appmanager/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/application-server.exe cmd/appmanager/main.go
