/*
 * protocols 协议包
 * @Author: F1
 * @Date: 2020-07-14 21:16:18
 * @LastEditTime: 2020-11-10 20:01:20
 * @LastEditors: F1
 * @Description: 协议包中指令部份，目前支持0-255的指令定义
 *
 * @FilePath: /edge-cloud-socket/protocols/command.go
 */
package protocols

import (
	"reflect"
	"strconv"
)

/**
 * @Title: Command 消息指令定义
 * @Description:
 *
 * 	    消息指令位于头部（Header）中第一个字节，最多支持255种指令
 *  __________________________________________________________
 * | 1 byte  | 1 byte  | 2 byte  | length byte                |
 * | ------  | ------  | ------- | ----------------           |
 * | command | flag    | length  |  body                      |
 * | ------  | ------  | ------- | ----------------           |
 * | [0]     | [1]     | [2][3]  | [4][5][6][7][]	 		  |
 * | 0-255   | 0-255   | 0-65535 | length                     |
 * |__________________________________________________________|
 *
 * @Author: F1
 * @Date: 2020-07-21 11:02:42
 */
type Command byte

const (
	REQUEST_HEARTBEAT                Command = 0   // 心跳包
	RESPONSE_HEARTBEAT               Command = 100 // 心跳包响应
	REQUEST_REGISTER                 Command = 1   // 边缘端向服务端注册
	RESPONSE_REGISTER_SUCCESS        Command = 101 // 注册响应 成功
	RESPONSE_REGISTER_FAILED         Command = 201 // 注册响应 成功
	REQUEST_TRANS_SKU_DATA           Command = 2   // 传输SKU信息
	RESPONSE_TRANS_SKU_DATA          Command = 102 // 传输SKU信息响应包
	REQUEST_TRANS_ITEM_DATA          Command = 3   // 传输匹配信息
	RESPONSE_TRANS_ITEM_DATA         Command = 103 // 传输匹配信息响应
	REQUEST_UPLOAD_SKU_DATA          Command = 4   // 边缘端传输SKU信息
	RESPONSE_UPLOAD_SKU_DATA         Command = 104 // 传输SKU信息响应包
	REQUEST_PASSIVE_UPLOAD_SKU_DATA  Command = 5   // 被动上传
	RESPONSE_PASSIVE_UPLOAD_SKU_DATA Command = 105 // 被动上传响应包
	REQUEST_TRANS_YOYOINFO_DATA      Command = 6   // 云端下发Yoyo数据
	RESPONSE_TRANS_YOYOINFO_DATA     Command = 106 // 云端下发Yoyo数据响应
	REQUEST_EXEC_CMD                 Command = 110
	RESPONAE_EXEC_CMD_REPLY          Command = 210 //
	HTTP_REQUEST_CMD                 Command = 111 // 发起ＨＴＴＰ请求
	HTTP_REQUEST_REPLY               Command = 211 //　HTTP请求返回
	TARGET_CMD_SKU_DATA_DONE         Command = 200 // 标识ＳＫＵ数据已传输完毕
)

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
