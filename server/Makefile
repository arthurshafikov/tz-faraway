DOCKER_COMPOSE_FILE := "./deployments/docker-compose.yml"
APP_NAME := "tz-faraway-server"

run: 
	go run ./...

lint:
	golangci-lint run ./...

test:
	go test ./... -race

up:
	docker-compose -f ${DOCKER_COMPOSE_FILE} -p ${APP_NAME} up --build

down:
	docker-compose -f ${DOCKER_COMPOSE_FILE} -p ${APP_NAME} down --volumes
