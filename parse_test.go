package main

import (
	"fmt"
	"testing"
	_ "github.com/golang/protobuf/protoc-gen-go/descriptor"
	_ "github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
)

func TestParse(t *testing.T) {
	var parser protoparse.Parser
	descs, err := parser.ParseFiles("msg/protocol/test.proto")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, desc := range descs {
		fmt.Print(desc.GetName())
		for _, msgDesc := range desc.GetMessageTypes() {
			fmt.Println("\t", msgDesc.GetName())
			for _, fieldDesc := range msgDesc.GetFields() {
				fmt.Println("\t", fieldDesc.GetType().String(), fieldDesc.GetName())
			}
		}
	} 
}