
$(shell PATH=$PATH:$GOPATH/bin)
PROTO_DIR := ../../pb

.DEFAULT_GOAL := help

help: ## List targets & descriptions
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

genproto:
	protoc --go_out=plugins=grpc:proto -I $(PROTO_DIR) $(PROTO_DIR)/health.proto

gateway:	
	if [ ! -d "proto-gen" ]; then \
			mkdir proto-gen; \
	fi

	protoc -I../../proto/api -I../../proto/third_party/googleapis --go_out=plugins=grpc:./proto-gen ../../proto/api/payment/*.proto
	protoc -I../../proto/api -I../../proto/third_party/googleapis --go_out=plugins=grpc:./proto-gen ../../proto/api/payment/*.proto
	
	protoc -I/usr/local/include -I. \
		-I../../proto/api \
		-I../../proto/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true,allow_repeated_fields_in_body=true:proto-gen \
		../../proto/api/**/*.proto

	