package serilaizer

import (
	"github.com/jimyag/grpc-go/sample"
	"github.com/stretchr/testify/require"
	"testing"
)

//
// TestWriteProtobufBinaryFile
//  @Description: 测试将 message 写入二进制文件
//  @param t
//
func TestWriteProtobufBinaryFile(t *testing.T) {
	t.Parallel()
	binaryFile := "../tmp/laptop.bin"
	laptop := sample.NewLaptop()
	err := WriteProtobufBinaryFile(laptop, binaryFile)
	require.NoError(t, err)
}

//
// TestReadProtobufFormBinaryFile
//  @Description: 测试从二进制文件读取 message
//  @param t
//
func TestReadProtobufFormBinaryFile(t *testing.T) {
	t.Parallel()
	binaryFile := "../tmp/laptop.bin"
	laptop := sample.NewLaptop()
	err := ReadProtobufFormBinaryFile(binaryFile, laptop)
	require.NoError(t, err)
}

//
// TestWriteProtobufToJSONFile
//  @Description: 测试将 message 写入 JSON 文件
//  @param t
//
func TestWriteProtobufToJSONFile(t *testing.T) {
	t.Parallel()
	jsonFile := "../tmp/laptop.json"
	laptop := sample.NewLaptop()
	err := WriteProtobufToJSONFile(laptop, jsonFile)
	require.NoError(t, err)
}
