
build-linux:
	@docker build -t otus-mon:v1 .
	@go build -tags=linux -o ./bin/client ./cmd/client
	@go build -tags=linux -o ./bin/server ./cmd/server

build-windows:
	@GOOS=windows go build -tags=windows -o ./bin/client.exe ./cmd/client
	@GOOS=windows go build -tags=windows -o ./bin/server.exe ./cmd/server

integration:
	@echo "Start server"
	@OTUS_MOD_START=cpu,loadavg,netstat go run -tags=linux ./cmd/server &
	@echo "Start test"
	@go test -v -count=1 ./cmd/client/

proto: 
	# rm -rf internal/pb
	# mkdir -p internal/pb

	@protoc \
		--go_out=internal/pb \
		--go-grpc_out=internal/pb \
		api/*.proto

