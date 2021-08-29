/*
 * @Author: F1
 * @Date: 2020-07-17 18:07:04
 * @LastEditors: F1
 * @LastEditTime: 2021-08-29 10:40:13
 * @Description: file content
 */
package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/fengyily/yoyoecs/protoc"
	"github.com/fengyily/yoyoecs/utils"
	"google.golang.org/protobuf/proto"
)

func TestCompress(t *testing.T) {
	befor := []byte("这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的")
	after := utils.Compress(befor)
	println("befor", len(befor))

	println("after", len(after), string(after))

	out := utils.UnCompress(after)

	println("out", string(out))
}

func TestPB(t *testing.T) {
	reg := &protoc.Register{}

	reg.SN = "SN00112343"
	reg.CompanyID = 12345674879
	reg.ShopCode = "这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的"
	reg.IP = "127.0.0.1"
	pbb, _ := proto.Marshal(reg)
	println("pb序列化之后的长度：", len(pbb))

	outreg := &protoc.Register{}
	proto.Unmarshal(pbb, outreg)
	println(outreg.String())

	fmt.Println("pb=", reg.String())
	if outreg.SN == reg.SN && outreg.CompanyID == reg.CompanyID && outreg.ShopCode == reg.ShopCode {
		t.Log("TestPB测试通过")
	} else {
		t.Error("TestPB测试失败，因为输入与输出不符。")
	}
}

func TestPBCompress(t *testing.T) {
	start := time.Now()
	var after []byte
	var pbb []byte

	after = nil
	println("after", len(after))
	// regs := protoc.Registers{
	// 	Register: make([]*protoc.Register, 0),
	// }

	//pbb, _ = proto.Marshal(&regs)
	//pbb, _ = proto.Marshal(&yoyos)
	after = utils.Compress(pbb)
	println("befor", len(pbb), "500 tims Compress using time", time.Since(start).Microseconds(), "millis")

	println("after", len(after), string(after))

	start2 := time.Now()
	fmt.Println(time.Now())
	out := utils.UnCompress(after)
	fmt.Println(time.Now())
	println("unCompress using time= ", time.Since(start2).Microseconds())
	outreg := &protoc.Registers{}
	proto.Unmarshal(out, outreg)

	println(outreg.String())

	// fmt.Println("pb=", reg.String())
	// if outreg.SN == reg.SN && outreg.CompanyID == reg.CompanyID && outreg.ShopCode == reg.ShopCode {
	// 	t.Log("TestPBCompress测试通过")
	// } else {
	// 	t.Error("TestPBCompress测试失败，因为输入与输出不符。")
	// }

}
