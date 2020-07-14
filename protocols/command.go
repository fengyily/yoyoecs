package protocols

import "strconv"

//  Command 消息指令
//  __________________________________________________________
// | 1 byte  | 1 byte  | 2 byte  | length byte                |
// | ------  | ------  | ------- | ----------------           |
// | command | json,pb | length  |  body                      |
// | ------  | ------  | ------- | ----------------           |
// | [0]     | [0]     | [1][2]  | [3][][][][]...[length+3]   |
// | 0-255   | 0-255   | 0-65535 | length                     |
// |__________________________________________________________|
type Command byte

const (
	REQUEST_HEARTBEAT         Command = 0   // 心跳包
	RESPONSE_HEARTBEAT        Command = 100 // 心跳包响应
	REQUEST_REGISTER          Command = 1   // 边缘端向服务端注册
	RESPONSE_REGISTER_SUCCESS Command = 101 // 注册响应 成功
	RESPONSE_REGISTER_FAILED  Command = 201 // 注册响应 成功

	REQUEST_TRANS_SKU_DATA  Command = 2   // 传输SKU信息
	RESPONSE_TRANS_SKU_DATA Command = 102 // 传输SKU信息响应包

	REQUEST_TRANS_ITEM_DATA  Command = 3   // 传输匹配信息
	RESPONSE_TRANS_ITEM_DATA Command = 103 //传输匹配信息响应
	REQUEST_UPLOAD_SKU_DATA  Command = 4   // 边缘端传输SKU信息
	RESPONSE_UPLOAD_SKU_DATA Command = 104 // 传输SKU信息响应包

	REQUEST_PASSIVE_UPLOAD_SKU_DATA  Command = 5 // 被动上传
	RESPONSE_PASSIVE_UPLOAD_SKU_DATA Command = 105

	// 快捷键
	//
)

// 判断是否是定义的指令
func (cmd Command) IsCommandType(c byte) bool {
	return c == byte(RESPONSE_HEARTBEAT) ||
		c == byte(REQUEST_REGISTER) ||
		c == byte(RESPONSE_REGISTER_SUCCESS) ||
		c == byte(RESPONSE_REGISTER_FAILED) ||
		c == byte(REQUEST_TRANS_SKU_DATA) ||
		c == byte(RESPONSE_TRANS_SKU_DATA) ||
		c == byte(REQUEST_TRANS_ITEM_DATA) ||
		c == byte(RESPONSE_TRANS_ITEM_DATA) ||
		c == byte(REQUEST_UPLOAD_SKU_DATA) ||
		c == byte(RESPONSE_UPLOAD_SKU_DATA)
}

func (cmd Command) ToString() string {
	return strconv.Itoa(int(cmd))
}
