
build:
	@docker build -t otus-mon:v1 .
	@go build -o ./bin/client ./cmd/client
	@go build -o ./bin/server ./cmd/server


integration:
	@echo "Start server"
	@OTUS_MOD_START=cpu,loadavg go run ./cmd/server &
	@echo "Start test"
	@go test -v -count=1 ./cmd/client/

proto: 
	# rm -rf internal/pb
	# mkdir -p internal/pb

	@protoc \
		--go_out=internal/pb \
		--go-grpc_out=internal/pb \
		api/*.proto

