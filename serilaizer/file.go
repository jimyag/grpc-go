package serilaizer

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
)

//
// WriteProtobufBinaryFile
//  @Description: 将message写入二进制文件
//  @param message
//  @param filename
//  @return error
//
func WriteProtobufBinaryFile(message proto.Message, filename string) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to binary:%w", err)
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write binary to file :%w", err)
	}
	return nil
}

//
// ReadProtobufFormBinaryFile
//  @Description: 从二进制文件中读取 proto message
//  @param filename
//  @param message
//  @return error
//
func ReadProtobufFormBinaryFile(filename string, message proto.Message) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read binary file :%w", err)
	}

	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("cannot unmarshal binary to proto message :%w", err)
	}
	return nil
}

//
// WriteProtobufToJSONFile
//  @Description: 将 proto message 写入到 JSON 文件
//  @param message
//  @param filename
//  @return error
//
func WriteProtobufToJSONFile(message proto.Message, filename string) error {
	data, err := ProtobufToJSON(message)
	if err != nil {
		return fmt.Errorf("cannot meshaler protobuf to json file :%w", err)
	}

	err = ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("cannot write json data to file :%w", err)
	}
	return nil
}
