GOOGLEAPIS_DIR="./api/googleapis"

protoc --proto_path=./api --proto_path=./api/googleapis --go_out=plugins=grpc:./src/payment/proto-gen ./api/payment/*.proto