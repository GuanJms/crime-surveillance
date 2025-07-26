CRIME_BINARY=crimeServiceApp
BROKER_BINARY=brokerServiceApp
AUTH_BINARY=authServiceApp

up:
	@echo "Starting Docker images.."
	docker-compose up -d
	@echo "Docker images started!"


up_build: build_crime build_broker build_auth
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"


build: build_crime build_broker build_auth

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

build_auth:
	@echo "Building crime binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"


generate_proto_crime:
	@echo "Generate crime server proto file"
	cd ../crime-service/proto && \
	protoc --go_out=. --go-grpc_out=. \
	--experimental_allow_proto3_optional \
	--proto_path=. crime.proto
	@echo "Done!"

	@echo "Generate broker-service crime proto file"
	cd ../broker-service/internal/crime_broker/proto && \
	protoc --go_out=. --go-grpc_out=. \
	--experimental_allow_proto3_optional \
	--proto_path=. crime.proto
	@echo "Done!"


generate_proto_auth:
	@echo "Generate authentication-service auth proto file"
	cd ../authentication-service/proto && \
	protoc --go_out=. --go-grpc_out=. \
	--experimental_allow_proto3_optional \
	--proto_path=. auth.proto
	@echo "Done!"

	@echo "Generate broker-service crime proto file"
	cd ../broker-service/internal/auth_broker/proto && \
	protoc --go_out=. --go-grpc_out=. \
	--experimental_allow_proto3_optional \
	--proto_path=. auth.proto
	@echo "Done!"

