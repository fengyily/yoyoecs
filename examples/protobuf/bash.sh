###
 # @Author: F1
 # @Date: 2020-07-20 16:19:50
 # @LastEditors: F1
 # @LastEditTime: 2020-07-22 10:06:28
 # @Description: 用于生成protocbuf .pb.go文件
### 
export PATH=$PATH;/home/f1/go/bin;
protoc --proto_path=examples/protobuf/ --go_out=protoc examples/protobuf/register.proto 
protoc --proto_path=examples/protobuf/ --go_out=protoc examples/protobuf/sku.proto 
protoc --proto_path=examples/protobuf/ --go_out=protoc examples/protobuf/item.proto 
