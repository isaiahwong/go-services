GOOGLEAPIS_DIR="./api/googleapis"

protoc --proto_path=./proto/api --proto_path=./proto/third_party/googleapis --go_out=plugins=grpc:./src/payment/proto-gen ./proto/api/payment/*.proto