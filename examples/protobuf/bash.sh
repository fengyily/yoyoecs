###
 # @Author: F1
 # @Date: 2020-07-20 16:19:50
 # @LastEditors: F1
 # @LastEditTime: 2020-07-30 15:14:07
 # @Description: 用于生成protocbuf .pb.go文件
### 
export PATH=$PATH;/home/f1/go/bin;
protoc --proto_path=. --go_out=../../protoc register.proto 
protoc --proto_path=. --go_out=../../protoc sku.proto 
protoc --proto_path=. --go_out=../../protoc item.proto 
protoc --proto_path=. --go_out=../../protoc yoyoinfo.proto 
protoc --proto_path=. --go_out=../../protoc cmd.proto 
protoc --proto_path=. --go_out=../../protoc http.proto 