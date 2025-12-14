PROTO_SRC=./pkg/proto
PROTO_OUT=./internal/pb

migrate:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down


generate-proto:
	protoc \
		--proto_path=$(PROTO_SRC) \
		--proto_path=$(PROTO_SRC)/product \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_SRC)/*.proto \
		$(PROTO_SRC)/product/*.proto \