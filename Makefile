docker-build-master:
	docker build -t sheshagiri/http-service:master .

build:
	go build

run:
	go run main.go
