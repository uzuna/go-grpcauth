
PROTO_DIR=../protos
WORK_DIR=./src/client

proto:
	mkdir -p WORK_DIR
	protoc -I=${PROTO_DIR} world.proto \
	--js_out=import_style=commonjs:\
	${WORK_DIR} \
	--grpc-web_out=import_style=commonjs,mode=grpcwebtext:\
	${WORK_DIR}