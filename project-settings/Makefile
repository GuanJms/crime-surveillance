CRIME_BINARY=crimeServiceApp
BROKER_BINARY=brokerServiceApp
AUTH_BINARY=authServiceApp
PATROL_BINARY=patrolServiceApp

up:
	@echo "Starting Docker images.."
	docker-compose up -d
	@echo "Docker images started!"


up_build: prepare_frontend build_crime build_broker build_auth build_patrol build_frontend
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

up_build_amd64: prepare_frontend build_crime_amd64 build_broker_amd64 build_auth_amd64 build_patrol_amd64 build_frontend
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"


build: prepare_frontend build_crime build_broker build_auth build_patrol build_frontend

build_amd64: prepare_frontend build_crime_amd64 build_broker_amd64 build_auth_amd64 build_patrol_amd64 build_frontend_amd64

down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"


build_crime:
	@echo "Building crime binary..."
	cd ../crime-service && env GOOS=linux CGO_ENABLED=0 go build -o ${CRIME_BINARY} ./cmd/api
	@echo "Done!"

build_crime_amd64:
	@echo "Building crime binary..."
	cd ../crime-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${CRIME_BINARY} ./cmd/api
	@echo "Done!"



build_broker:
	@echo "Building crime binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

build_broker_amd64:
	@echo "Building crime binary..."
	cd ../broker-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

build_auth:
	@echo "Building crime binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

build_auth_amd64:
	@echo "Building crime binary..."
	cd ../authentication-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

build_patrol:
	@echo "Building patrol binary..."
	cd ../patrol-service && env GOOS=linux CGO_ENABLED=0 go build -o ${PATROL_BINARY} ./cmd/api
	@echo "Done!"

build_patrol_amd64:
	@echo "Building patrol binary..."
	cd ../patrol-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${PATROL_BINARY} ./cmd/api
	@echo "Done!"


generate_proto_broker_service:
	@echo "Generate broker-service crime proto file"
	cd ../broker-service/internal/crime_broker/proto && \
	protoc --go_out=. --go-grpc_out=. \
	--experimental_allow_proto3_optional \
	--proto_path=. crime.proto
	@echo "Done!"

	@echo "Generate broker-service auth proto file"
	cd ../broker-service/internal/auth_broker/proto && \
	protoc --go_out=. --go-grpc_out=. \
	--experimental_allow_proto3_optional \
	--proto_path=. auth.proto
	@echo "Done!"

	@echo "Generate broker-service patrol proto file"
	cd ../broker-service/internal/patrol_broker/proto && \
	protoc --go_out=. --go-grpc_out=. \
	--experimental_allow_proto3_optional \
	--proto_path=. patrol.proto
	@echo "Done!"


generate_proto_auth_service:
	@echo "Generate authentication-service auth proto file"
	cd ../authentication-service/proto && \
	protoc --go_out=. --go-grpc_out=. \
	--experimental_allow_proto3_optional \
	--proto_path=. auth.proto
	@echo "Done!"


generate_proto_crime_service:
	@echo "Generate crime-service crime proto file"
	cd ../crime-service/proto && \
	protoc --go_out=. --go-grpc_out=. \
	--experimental_allow_proto3_optional \
	--proto_path=. crime.proto
	@echo "Done!"

generate_proto_patrol_service:
	@echo "Generate patrol-service patrol proto file"
	cd ../patrol-service/proto && \
	protoc --go_out=. --go-grpc_out=. \
	--experimental_allow_proto3_optional \
	--proto_path=. patrol.proto
	@echo "Done!"

	@echo "Generate patrol-service crime proto file"
	cd ../patrol-service/data/crime_repo/proto && \
	protoc --go_out=. --go-grpc_out=. \
	--experimental_allow_proto3_optional \
	--proto_path=. crime.proto
	@echo "Done!"


generate_sqlc_patrol_service:
	@echo "Generate patrol-service sqlc file"
	cd ../patrol-service && \
	sqlc generate
	@echo "Done!"


frontend: prepare_frontend
	@echo "Building & starting frontend container (port 8081).."
	docker-compose up -d front-end
	@echo "Frontend container running!"

prepare_frontend:
	@echo "Copying openAPI.yaml to front-end..."
	cp ../api/openapi.yaml ../front-end/openapi.yaml
	@echo "Copying README pdf to front-end..."
	cp ../README.pdf ../front-end/README.pdf
	@echo "Done!"

build_frontend:
	@echo "Building front-end image..."
	docker-compose build front-end
	@echo "Done!"

build_frontend_amd64:
	@echo "Building front-end image..."
	docker-compose build --build-arg TARGETARCH=amd64 front-end
	@echo "Done!"

