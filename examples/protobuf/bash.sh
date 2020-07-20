export PATH=$PATH;/home/f1/go/bin;
protoc --proto_path=examples/protobuf/ --go_out=examples/protocols/ examples/protobuf/register.proto 
