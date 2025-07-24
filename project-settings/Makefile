CRIME_BINARY=crimeServiceApp
BROKER_BINARY=brokerServiceApp

up:
	@echo "Starting Docker images.."
	docker-compose up -d
	@echo "Docker images started!"


up_build: build_crime build_broker
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"


build: build_crime build_broker

down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"


// TODO: Add deployment GOARCH as amd64
build_crime:
	@echo "Building crime binary..."
	cd ../crime-service && env GOOS=linux CGO_ENABLED=0 go build -o ${CRIME_BINARY} ./cmd/api
	@echo "Done!"


// TODO: Add deployment GOARCH as amd64
build_broker:
	@echo "Building crime binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

