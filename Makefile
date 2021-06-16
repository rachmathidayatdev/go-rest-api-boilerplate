proto-build:
	@protoc --go_out=plugins=grpc:. *.proto

run-docs-generate:
	@swag init
	
run-service-book:
	@go run main.go book

run-service-book-migrate:
	@go run main.go book db-migrate

run-service-book-rollback:
	@go run main.go book db-rollback

# Docker
GO_BUILD_ENV := GOVERSION=1.13 CGO_ENABLED=0 GOOS=linux GOARCH=amd64
DOCKER_BUILD=$(shell pwd)/.docker_build
DOCKER_CMD=$(DOCKER_BUILD)/go-rest-api-boilerplate

$(DOCKER_CMD): clean
	mkdir -p $(DOCKER_BUILD)
	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

clean:
	rm -rf $(DOCKER_BUILD)

heroku: $(DOCKER_CMD)
	heroku container:push web