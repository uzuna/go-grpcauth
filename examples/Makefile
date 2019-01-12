
setup:
	go get -u -v github.com/yoheimuta/protolint/cmd/pl

pb:
	protoc -I protos/ protos/world.proto --go_out=plugins=grpc:pb

pl:
	pl lint -fix .