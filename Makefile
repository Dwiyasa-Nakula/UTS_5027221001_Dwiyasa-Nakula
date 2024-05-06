gen:
	@protoc \
		--proto_path=protobuf "protobuf/musicplaylist.proto" \
		--go_out=backend/genproto/musicplaylist --go_opt=paths=source_relative \
	--go-grpc_out=backend/genproto/musicplaylist --go-grpc_opt=paths=source_relative

server:
	@go run grpc/server/main.go $(profile)

client:
	@go run grpc/client/main.go $(profile)