#!/bin/bash

base_filename="grpc-gateway"
proto_file=$base_filename.proto

# generate stubbed code
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  -I$GOPATH/src/github.com/gogo/protobuf/gogoproto/ \
  --gogo_out=plugins=grpc:models \
  $proto_file

# generate reverse proxy
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:models \
  $proto_file

# generate mocks
mockgen -source=models/$base_filename.pb.go -package=models -destination=models/$base_filename.pb.mock.go