###
 # @Author: F1
 # @Date: 2020-07-20 16:19:50
 # @LastEditors: F1
 # @LastEditTime: 2021-09-02 21:16:41
 # @Description: 用于生成protocbuf .pb.go文件
### 
export PATH=$PATH;
#v3.17.3
protoc resetdns.proto --go_out=.
protoc register.proto --go_out=.
protoc execsql.proto --go_out=.
protoc cmd.proto --go_out=.
protoc http.proto --go_out=.
protoc sendto.proto --go_out=.