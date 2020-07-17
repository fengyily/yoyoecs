package test

import (
	"testing"

	"github.com/fengyily/yoyoecs/utils"
)

func TestCompress(t *testing.T) {

	befor := []byte("这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的")
	after := utils.Compress(befor)
	println("befor", len(befor))

	println("after", len(after), string(after))

	out := utils.UnCompress(after)

	println("out", string(out))
}
