package main

import (
	"fmt"
	"os"
	"strings"
	"log"
	"io/ioutil"
	"path/filepath"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/golang/protobuf/proto"
	"github.com/springstar/robot/pb"
	"github.com/springstar/robot/msg"
)

var typeMap map[string]string = map[string]string{
	"TYPE_STRING": "string",
	"TYPE_INT32": "int32",
	"TYPE_INT": "int",
	"TYPE_INT64": "int64",
	"TYPE_MESSAGE": "msg",
	"TYPE_ENUM": "enum",
}

type Field struct {
	typ string
	name string
}

type DescriptorGen struct {
	mds map[int32]*desc.MessageDescriptor
	id2names map[int32]string

}

func newDescriptorGen() *DescriptorGen {
	return &DescriptorGen{
		mds: make(map[int32]*desc.MessageDescriptor),
		id2names: make(map[int32]string),
	}
}

func(g *DescriptorGen) parse(path string) {
	if _, err := os.Stat("msg/serializer.go"); err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(path)
	if (err != nil) {
		log.Fatal(err)
	}

	var pfiles []string

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		
		ext := filepath.Ext(f.Name())
		if strings.Compare(ext, ".proto") != 0 {
			continue
		}

		p := filepath.Join(path, f.Name())
		pfiles = append(pfiles, p)
	}

	var parser protoparse.Parser
	for _, f := range pfiles {
		fds, err := parser.ParseFiles(f)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		for _, fd := range fds {
			mds := fd.GetMessageTypes()
			for _, md := range mds {
				options := md.GetOptions()
				mid, _ := proto.GetExtension(options, pb.E_Msgid)
				if (mid != nil) {
					msg.AddDescriptor(*mid.(*int32), md)
					g.id2names[*mid.(*int32)] = md.GetName()
				}
			}
		}
	}

}







