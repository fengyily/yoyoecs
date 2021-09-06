/*
 * protocols 协议包
 * @Author: F1
 * @Date: 2020-07-14 21:16:18
 * @LastEditTime: 2021-09-06 19:58:53
 * @LastEditors: F1
 * @Description: 协议包中指令部份，目前支持0-255的指令定义
 *
 * @FilePath: /yoyoecs/protocols/command.go
 */
package protocols

import (
	"reflect"
	"strconv"

	"google.golang.org/protobuf/reflect/protoreflect"
)

/**
 * @Title: Command 消息指令定义
 * @Description:
 *
 * 	    消息指令位于头部（Header）中第一个字节，最多支持255种指令
 *  __________________________________________________________________
 * | 1 byte  | 1 byte  | 4 byte          | length byte                |
 * | ------  | ------  | --------------- | ----------------           |
 * | command | flag    | length          |  body                      |
 * | ------  | ------  | --------------- | ----------------           |
 * | [0]     | [1]     | [2][3][4][5]    | [4][5][6][7][]	 		  |
 * | 0-255   | 0-255   | 0-2^32          | length                     |
 * |__________________________________________________________________|
 *
 * @Author: F1
 * @Date: 2020-07-21 11:02:42
 */
type Command byte

const (
	REQUEST_HEARTBEAT         Command = 0   // 心跳包
	RESPONSE_HEARTBEAT        Command = 100 // 心跳包响应
	REQUEST_REGISTER          Command = 1   // 边缘端向服务端注册
	RESPONSE_REGISTER_SUCCESS Command = 101 // 注册响应 成功
	RESPONSE_REGISTER_FAILED  Command = 201 // 注册响应 成功
	REQUEST_SENDTO_CMD        Command = 102 // 1 - 1 发消息
	REQUEST_SENDTO_REPLY      Command = 202 // 1 - 1 消息响应
	REQUEST_CAST_MSG_CMD      Command = 103 // 广播消息
	REQUEST_CAST_MSG_REPLY    Command = 203 // 广播响应
	REQUEST_EXEC_CMD          Command = 110 // Shell 脚本执行
	RESPONAE_EXEC_CMD_REPLY   Command = 210 // Shell 脚本执行响应
	HTTP_REQUEST_CMD          Command = 111 // 发起ＨＴＴＰ请求
	HTTP_REQUEST_REPLY        Command = 211 //　HTTP请求返回
	SQL_REQUEST_CMD           Command = 112 // SQL 请求
	SQL_REQUEST_REPLY         Command = 212 // SQL 请求响应
	RESET_DB_DNS_CMD          Command = 113 // 重置 DNS
	RESET_DB_DNS_REPLY        Command = 213 // 重置 DNS 响应
)

// Enum value maps for Command.
var (
	Command_name = map[int32]string{
		0:   "REQUEST_HEARTBEAT",
		100: "RESPONSE_HEARTBEAT",
		1:   "REQUEST_REGISTER",
		101: "RESPONSE_REGISTER_SUCCESS",
		201: "RESPONSE_REGISTER_FAILED",
		102: "REQUEST_SENDTO_CMD",
		202: "REQUEST_SENDTO_REPLY",
		103: "REQUEST_CAST_MSG_CMD",
		203: "REQUEST_CAST_MSG_REPLY",
		110: "REQUEST_EXEC_CMD",
		210: "RESPONAE_EXEC_CMD_REPLY",
		111: "HTTP_REQUEST_CMD",
		211: "HTTP_REQUEST_REPLY",
		112: "SQL_REQUEST_CMD",
		212: "SQL_REQUEST_REPLY",
		113: "RESET_DB_DNS_CMD",
		213: "RESET_DB_DNS_REPLY",
	}
)

func (x Command) String() string {
	return Command_name[int32(x.Number())]
}

func (x Command) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

/**
 * @Title: IsCommandType
 * @Description:
 *
 * 			说明：判断字节流中的当前位是否是定义的指令类型
 *
 * @Author: F1
 * @Date: 2020-07-21 11:05:14
 * @Param:
 * 		c byte 字节流中的第一位
 * @Return:
 *		ok bool 是否是指令类型
 */
func (cmd Command) IsCommandType(c byte) bool {
	tmp := Command(c)
	return reflect.TypeOf(tmp) == reflect.TypeOf(RESPONSE_HEARTBEAT)
}

/**
* @Title: ToString
* @Description:
* @Author: F1
* @Date: 2020-07-21 11:07:06
 * @Return: string
*/
func (cmd Command) ToString() string {
	return strconv.Itoa(int(cmd))
}
