package utils

import (
	"bytes"
	"compress/gzip"

	"github.com/sirupsen/logrus"
)

func Compress(gbk []byte) []byte {
	defer func() {
		println("Compress", recover())
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

func UnCompress(gbk []byte) []byte {
	defer func() {
		println("UnCompress", recover())
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
