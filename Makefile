
build:
	@docker build -t otus-mon:v1 .
	@go build -o ./bin/client ./cmd/client
	@go build -o ./bin/server ./cmd/server


run: build
	@./bin/server -d

proto: 
	# rm -rf internal/pb
	# mkdir -p internal/pb

	@protoc \
		--go_out=internal/pb \
		--go-grpc_out=internal/pb \
		api/*.proto

