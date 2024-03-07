LOCAL_BIN:=$(CURDIR)/bin
PACKAGE=cmd/url_generator/main.go

install-go-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v25.3/protoc-25.3-osx-aarch_64.zip
	unzip protoc-25.3-osx-aarch_64.zip
	mv readme.txt $(LOCAL_BIN)
	mv include $(LOCAL_BIN)
	rm -rf protoc-25.3-osx-aarch_64.zip


vendor-proto:
		mkdir -p vendor-proto
		@if [ ! -d vendor-proto/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor-proto/googleapis &&\
			mkdir -p  vendor-proto/google/ &&\
			mv vendor-proto/googleapis/google/api vendor-proto/google &&\
			rm -rf vendor-proto/googleapis ;\
		fi
		@if [ ! -d vendor-proto/google/protobuf ]; then\
			git clone https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
			mkdir -p  vendor-proto/google/protobuf &&\
			mv vendor-proto/protobuf/src/google/protobuf/*.proto vendor-proto/google/protobuf &&\
			rm -rf vendor-proto/protobuf ;\
		fi

generate:
	mkdir -p pkg/url_generator

	$(LOCAL_BIN)/protoc -I api/ -I vendor-proto \
	--go_out=pkg/url_generator --go_opt=paths=source_relative \
	--go-grpc_out=pkg/url_generator --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pkg/url_generator --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	api/url_generator.proto

build: bindir
	go build -o ${LOCAL_BIN}/url_generator ${PACKAGE}
bindir:
	mkdir -p ${LOCAL_BIN}
run:
	go run ${PACKAGE}