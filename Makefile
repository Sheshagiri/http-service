docker-build-master:
	docker build -t sheshagiri/http-service:master .

build:
	go build

run:
	go run main.go

docker-push-quay: docker-build-master
	docker tag sheshagiri/http-service:master sheshagiri0/http-service:0.0.2