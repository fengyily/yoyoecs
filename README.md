
## ECS in scene
![Image text](https://raw.githubusercontent.com/fengyily/yoyoecs/master/resources/resource1.jpg)

## ECS 协议
> Edge Client Server 协议
### Header
```go
/**
 * @Title:协议头部，包含指令，标识及包体的长度
 * @Description:
 *      一般来讲一个完整的包是由Header|Body组成，特殊的包只有包头：
 *　 心跳包及心跳回复包，只有1字节，其它包允许只有包头，也就是Body为空，此时Header.Length为0
 *
 * @Author: F1
 * @Date: 2020-07-21 10:01:05
 * @LastEditors: F1
 * @Param:
 * @Return:
 *      []byte 长度为：protocols.HEADER_LENGTH
 */
type Header struct {
    Cmd    Command
    Flag   Flag
    Length uint32
}
```

### Command
```go
/**
 * @Title: Command 消息指令定义
 * @Description:
 *
 *      消息指令位于头部（Header）中第一个字节，最多支持255种指令
 *  __________________________________________________________________
 * | 1 byte  | 1 byte  | 4 byte          | length byte                |
 * | ------  | ------  | --------------- | ----------------           |
 * | command | flag    | length          |  body                      |
 * | ------  | ------  | --------------- | ----------------           |
 * | [0]     | [1]     | [2][3][4][5]    | [4][5][6][7][]             |
 * | 0-255   | 0-255   | 0-2^32          | length                     |
 * |__________________________________________________________________|
 *
 * @Author: F1
 * @Date: 2020-07-21 11:02:42
 */
type Command byte
```

### Flag
```go
/**
 * @Title: Flag 长度８位，占１字节
 * @Description:
 *
 *              DataType 前３位用来标识传输的数据类型，所以最多支持８种，
 *              Encrytion 第４位表示数据包是否加密
 *              COMPRESS 第５位表示数据包是否开启了压缩
 *              第６到８位为扩展备用位
 *                ____________________________________________
 *               | DataType  | Encrytion |COMPRESS|ext|ext|ext|
 *               |-----------|-----------|--------|---|---|---|
 *               | 3 bit     | 1 bit     | 1 bit  | 1 | 1 | 1 |
 *               |-----------|-----------|--------|---|---|---|
 *               |[0] [1] [2]|    [3]    |   [4]  |[5]|[6]|[7]|
 *               |____________________________________________|
 *
 * @Author: F1
 * @Date: 2020-07-21 10:55:39
 * @Param:
 * @Return:
 */
type Flag byte
```

### Body
> 支持 Protocol Bufffer & Json，在 Flag 中标识 Body 的类型、是否加密、是否对内容进行压缩


### 消息转发协议
> SendTo 协议核心代码
```go
if header.Flag&protocols.HEADER_FLAG_IS_COMPRESS > 0 {
        data = utils.UnCompress(data)
        header.Length = uint32(len(data))
}
if header.Cmd == protocols.REQUEST_SENDTO_CMD {
        if cs.OnSendToMessage != nil {
                // 消息转发
                sendTo := protoc.SendTo{}
                err := proto.Unmarshal(data, &sendTo)
                if err != nil {
                        fmt.Println("proto.Unmarshal(b, &sendTo)", err)
                        return
                }
                cs.OnSendToMessage(sendTo.CID, &sendTo, cs)
        }
} else {
        cs.OnRecvMessage(header, data, cs)
}
```


## Generate Protocal Buffer

go get -u github.com/golang/protobuf/proto

go get -u github.com/golang/protobuf/protoc-gen-go

cd examples/protobuf/

sh bash.sh 

