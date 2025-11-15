PROTO_DIR=proto
OUTDIR_PROTO=pb


generate-proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(OUTDIR_PROTO) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUTDIR_PROTO) --go-grpc_opt=paths=source_relative \
		$$(find $(PROTO_DIR) -name "*.proto")