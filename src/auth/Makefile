
$(shell PATH=$PATH:$GOPATH/bin)
PROTO_DIR := ../../pb

.DEFAULT_GOAL := help

help: ## List targets & descriptions
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

genproto:
	protoc --go_out=plugins=grpc:proto -I $(PROTO_DIR) $(PROTO_DIR)/auth/auth.proto

