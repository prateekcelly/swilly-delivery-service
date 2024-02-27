APP_NAME=swilly-delivery-service
APP_EXECUTABLE="./out/$(APP_NAME)"

cp-config:
	cp application.yml.sample application.yml

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

test:
	@echo "Running Tests"
	@go clean -testcache
	@go test `go list ./...`

build:
	@echo "Building Executable"
	go build -o ${APP_EXECUTABLE}

docker.run:
	docker-compose -f docker/docker-compose.yml up -d

docker.stop:
	docker-compose -f docker/docker-compose.yml down --remove-orphans

start-server: build
	${APP_EXECUTABLE} server

start-worker: build
	${APP_EXECUTABLE} worker

doc:
	@swag i
