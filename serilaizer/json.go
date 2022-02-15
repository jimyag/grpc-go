package serilaizer

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

//
// ProtobufToJSON
//  @Description: 将 message 转换为 JSON
//  @param message
//  @return string
//  @return error
//
func ProtobufToJSON(message proto.Message) (string, error) {
	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true, Indent: "  ", OrigName: true,
	}

	return marshaler.MarshalToString(message)
}
