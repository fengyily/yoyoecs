package protocols

import "yoyoecs/utils"

// Flag ８bit
//  ____________________________________________
// | DataType  | Encrytion |COMPRESS|ext|ext|ext|
// |-----------|-----------|--------|---|---|---|
// | 3 bit     | 1 bit     | 1 bit  | 1 | 1 | 1 |
// |-----------|-----------|--------|---|---|---|
// |[0] [1] [2]|    [3]    |   [4]  |[5]|[6]|[7]|
// |____________________________________________|
type Flag uint8

const (
	HEADER_LENGTH = 4

	HEADER_FLAG_DATA_TYPE_JSON   Flag = 1      // 0000 0001
	HEADER_FLAG_DATA_TYPE_PB     Flag = 2      // 0000 0010
	HEADER_FLAG_DATA_TYPE_STRING Flag = 3      // 0000 0011
	HEADER_FLAG_DATA_TYPE_EXT1   Flag = 4      // 0000 0100
	HEADER_FLAG_DATA_TYPE_EXT2   Flag = 5      // 0000 0101
	HEADER_FLAG_DATA_TYPE_EXT3   Flag = 6      // 0000 0110
	HEADER_FLAG_DATA_TYPE_EXT4   Flag = 7      // 0000 0111
	HEADER_FLAG_IS_ENCRYTION     Flag = 1 << 3 // 0000 1000 是否加密
	HEADER_FLAG_IS_COMPRESS      Flag = 1 << 4 // 0001 0000 是否开启了压缩
)

type Header struct {
	Cmd    Command
	Flag   uint8
	Length uint16
}

func (header *Header) ToBytes() []byte {
	var data []byte
	if header.Cmd == REQUEST_HEARTBEAT || header.Cmd == RESPONSE_HEARTBEAT || header.Cmd == RESPONSE_REGISTER_SUCCESS {
		data = make([]byte, 1)
		data[0] = byte(header.Cmd)
	} else {
		data = make([]byte, HEADER_LENGTH)
		data[0] = byte(header.Cmd)
		data[1] = header.Flag

		data = append(data, utils.Uint16ToBytes(header.Length)...)
	}
	return data
}

func LoadHeader(buffer *[]byte) (ok bool, header Header) {
	ok = true
	for total := len(*buffer); total > 0; total = len(*buffer) {
		i := 0
		header.Cmd = Command((*buffer)[i])
		if header.Cmd == REQUEST_HEARTBEAT ||
			header.Cmd == RESPONSE_REGISTER_SUCCESS ||
			header.Cmd == RESPONSE_HEARTBEAT {
			*buffer = (*buffer)[i+1:]
			return ok, header
		}

		//var length uint16
		// 正常包的处理逻辑，如果遇到不识别的包，重新定位头部位置。
		if !header.Cmd.IsCommandType((*buffer)[i]) {
			rIndex := 0
			for ; rIndex < total-i; rIndex++ {
				if header.Cmd.IsCommandType((*buffer)[i+rIndex]) {
					//fmt.Println("重新定位成功。丢弃　", i+rIndex, "字节。")
					break
				}
			}
			*buffer = (*buffer)[i+rIndex:]
			//fmt.Println("未识别的指令", rIndex, len(*buffer))
			continue
		}

		// 是否满足正常包的一个头部
		if total < i+HEADER_LENGTH {
			//fmt.Println("len(cs.Buffer) < i+3")
			return false, header
		}
		header.Flag = (*buffer)[1]
		header.Length = utils.BytesToUInt16((*buffer)[2:4])

		break
	}
	return true, header
}
