/*
 * @Author: F1
 * @Date: 2020-07-17 17:52:27
 * @LastEditors: F1
 * @LastEditTime: 2020-07-21 11:15:47
 * @Description:
 *
 *				工具包中，压缩解压缩的相关的功能，采用Gzip算法
 */
package utils

import (
	"bytes"
	"compress/gzip"

	"github.com/sirupsen/logrus"
)

/**
 * @Title: Compress
 * @Description:
 *
 * 				压缩方法
 *
 * @Author: F1
 * @Date: 2020-07-21 11:13:56
 * @Param:[]byte
 * @Return:[]byte
 */
func Compress(gbk []byte) []byte {
	defer func() {
		if err := recover(); err != nil {
			println("Compress", err)
		}
	}()
	//使用gzip压缩
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	_, err := zw.Write(gbk)
	zw.Flush()
	if err != nil {
		logrus.Fatal(err)
	}

	if err := zw.Close(); err != nil {
		logrus.Fatal(err)
	}

	return buf.Bytes()
}

/**
 * @Title: UnCompress
 * @Description:
 *
 * 				解压缩方法，同样在该应用场景下，限制了包的大小
 *
 * @Author: F1
 * @Date: 2020-07-21 11:13:56
 * @Param:[]byte
 * @Return:[]byte
 */
func UnCompress(gbk []byte) []byte {
	defer func() {
		if err := recover(); err != nil {
			println("UnCompress", err)
		}
	}()
	buf := bytes.NewBuffer(gbk)
	zr, err := gzip.NewReader(buf)

	if err != nil {
		println(err)
	}
	out := make([]byte, 2<<16)
	l, err := zr.Read(out)

	if err := zr.Close(); err != nil {
		println(err)
	}

	return out[:l]
}
