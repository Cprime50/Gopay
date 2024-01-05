.PHONY: build docker-up docker-down run stop

# Set the name of your Go application binary
APP_BINARY_NAME := Gopay

build:
	# Build the Go application
	go build -o $(APP_BINARY_NAME) ./server

docker-up:
	# Run Docker Compose to start your services
	docker-compose up -d

docker-down:
	# Stop and remove the Docker Compose services
	docker-compose down

run: build docker-up
	# Run the Go application
	./$(APP_BINARY_NAME)

stop: docker-down
	# Stop Docker Compose services
	docker-compose stop
	# Remove Docker Compose services
	docker-compose rm -f

start: build docker-up run

.PHONY: start

