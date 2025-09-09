PROTO_DIR=./shared/proto
OUT_DIR=./shared/proto

.PHONY: proto-user
proto-user:
	protoc -I=$(PROTO_DIR) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/user/user.proto

.PHONY: proto-order
proto-order:
	protoc -I=$(PROTO_DIR) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/order/order.proto

.PHONY: proto-review
proto-review:
	protoc -I=$(PROTO_DIR) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/review/review.proto

.PHONY: proto-listing
proto-listing:
	protoc -I=$(PROTO_DIR) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/listing/listing.proto

.PHONY: proto
proto: proto-user proto-order proto-review proto-listing
	@echo "All proto files generated!"
