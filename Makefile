proto:
	rm -rf internal/pb
	mkdir -p internal/pb

	protoc \
		--go_out=internal/pb \
		--go-grpc_out=internal/pb \
		api/*.proto

