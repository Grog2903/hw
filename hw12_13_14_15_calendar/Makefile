LOCAL_BIN:=$(CURDIR)/bin

PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc

DOCKER_IMG="calendar:develop"
LOCAL_MIGRATION_DIR="./migrations/"
LOCAL_MIGRATION_DSN="host=localhost port=54321 dbname=calendar user=calendar-user password=calendar-password sslmode=disable"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(LOCAL_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(LOCAL_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender
	go build -v -o $(LOCAL_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar_sheduler

run: build
	$(BIN) -config configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/storage/memory/ ./internal/api/event/

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

install-deps:

.PHONY: .bin-deps
.bin-deps:
	$(info Installing binary dependencies...)

	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest && \
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

.vendor-proto: .vendor-rm vendor-proto/google/api vendor-proto/google/protobuf vendor-proto/protoc-gen-openapiv2/options

.PHONY: .vendor-rm
.vendor-rm:
	rm -rf vendor-proto

vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor-proto/googleapis && \
 	cd vendor-proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor-proto/grpc-ecosystem && \
 	cd vendor-proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor-proto/protoc-gen-openapiv2
	mv vendor-proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor-proto/protoc-gen-openapiv2
	rm -rf vendor-proto/grpc-ecosystem

vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf

EVENT_PROTO_PATH:="api/event/v1"

PHONY: .protoc-generate
.protoc-generate:
	mkdir -p pkg/${EVENT_PROTO_PATH}
	mkdir -p api/openapiv2
	$(PROTOC) -I ${EVENT_PROTO_PATH} -I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out pkg/${EVENT_PROTO_PATH} --go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out pkg/${EVENT_PROTO_PATH} --go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway --grpc-gateway_out pkg/${EVENT_PROTO_PATH} --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
	--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 --openapiv2_out api/openapiv2 --openapiv2_opt logtostderr=true \
	api/event/v1/*.proto
	go mod tidy

PHONY: generate
generate: .bin-deps .vendor-proto .protoc-generate

.PHONY: fast-generate
fast-generate: .protoc-generate

migrate:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

migrate-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

up:
	docker-compose up -d --build

down:
	docker-compose down

restart: down up

integration-tests:
	set -e ;\
	docker-compose -f docker-compose.test.yml up --build  --force-recreate  --abort-on-container-exit;\
	docker-compose -f docker-compose.test.yml run integration_tests go test -tags=integration ./... ;\
	docker-compose -f docker-compose.test.yml down

test-cleanup:
	docker-compose -f docker-compose.test.yml down \
        --rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60; \
  	docker-compose rm -f