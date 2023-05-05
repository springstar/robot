package main

import (
	"fmt"
	"strings"
	"log"
	"os"
	"io/ioutil"
	"path/filepath"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/golang/protobuf/proto"
	"github.com/springstar/robot/pb"
	"github.com/gobuffalo/plush"
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

type ProtoGen struct {
	files []string
	mds map[int32]*desc.MessageDescriptor
	fids map[string][]*desc.FieldDescriptor
	id2names map[int32]string

}

func newProtoGen() *ProtoGen {
	return &ProtoGen{
		mds: make(map[int32]*desc.MessageDescriptor),
		id2names: make(map[int32]string),
		fids: make(map[string][]*desc.FieldDescriptor),
	}
}

func(g *ProtoGen) parse(path string) {
	files, err := ioutil.ReadDir(path)
	if (err != nil) {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		
		ext := filepath.Ext(f.Name())
		if strings.Compare(ext, ".proto") != 0 {
			continue
		}

		p := filepath.Join(path, f.Name())
		
		g.files = append(g.files, p)
	}

	var parser protoparse.Parser
	for _, f := range g.files {
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
					g.mds[ *mid.(*int32)] = md
					g.id2names[*mid.(*int32)] = md.GetName()
					var fids []*desc.FieldDescriptor
					fids = append(fids, md.GetFields()...)
					g.fids[md.GetName()] = fids
				}
			}
		}
	}

}

func (g *ProtoGen) generate() {
	content, err := ioutil.ReadFile("template/funcdecl.tpl")
	if (err != nil) {
		log.Fatal(err)
	}

	template := string(content[:])

	ctx := plush.NewContext()
	var protos []string
	for _, p := range g.id2names {
		protos = append(protos, p)
	}

	ctx.Set("package", "package main")
	ctx.Set("imports", func() []string {
		return []string{
			"github.com/springstar/protogen/pb",
			"github.com/jhump/protoreflect/desc",
			"github.com/jhump/protoreflect/dynamic",
		}
	})

	ctx.Set("names", protos)
	ctx.Set("decl", func(p string) string {
		var sb strings.Builder
		sb.WriteString("func parse")
		sb.WriteString(p)
		sb.WriteString("(id int32, bytes []byte)")
		sb.WriteString(" *")
		sb.WriteString(p)
		return sb.String()
		
		
	})

	ctx.Set("lbrack", func() string {
		return " {"
	})

	ctx.Set("rbrack", func() string {
		return "}"
	})

	var tm map[string][]Field = make(map[string][]Field)
	
	for name, fields := range g.fids {
		var fids []Field
		for _, field := range fields {
			t := typeMap[field.GetType().String()]
			fld := Field{
				typ: t,
				name: field.GetName(),
			}

			fids = append(fids, fld)
		}
		
		tm[name] = fids
	}

	ctx.Set("fields", func(name string) []string {
		var fns []string
		fields := g.fids[name]
		for _, f := range fields {
			fns = append(fns, f.GetName())
		}
		
		return fns
	})

	ctx.Set("params", func(msg string) string {					
		fields := tm[msg]
		var sb strings.Builder
		c := len(fields)
		var i int = 0
		for _, field := range fields {
			k := field.typ
			v := field.name
			sb.WriteString(v)
			sb.WriteString(" ")
			if k == "msg" || k == "enum" {
				var sbb strings.Builder
				if k == "msg" {
					sbb.WriteString("*pb.")
				} else {
					sbb.WriteString("pb.")
				}
				
				cname := strings.Title(v)
				sbb.WriteString(cname)
				sb.WriteString(sbb.String())

			} else {
				sb.WriteString(k)
			}

			if i < c-1 {
				sb.WriteString(", ")
			}
			i++
		}

		return sb.String()
	})

	s, err := plush.Render(template, ctx)
	if (err != nil) {
		log.Fatal(err)
	}
	
	write(s)	

}

func write(s string) {
    file, err := os.Create("serializer.go")
    if err != nil {
        return
    }
    defer file.Close()

    file.WriteString(s)
}





