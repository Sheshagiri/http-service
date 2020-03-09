docker-build:
	docker build -t http-service .

build:
	go build

run:
	go run main.go
