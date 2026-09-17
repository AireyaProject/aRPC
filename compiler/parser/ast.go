package parser

import (
	"bytes"
)

type Node interface {
	String() string
}

type File struct {
	Namespace string
	Enums     []*Enum
	Structs   []*Struct
	Services  []*Service
}

func (f *File) String() string {
	var buf bytes.Buffer
	buf.WriteString("namespace " + f.Namespace + ";\n")
	for _, e := range f.Enums {
		buf.WriteString(e.String())
	}
	for _, s := range f.Structs {
		buf.WriteString(s.String())
	}
	for _, s := range f.Services {
		buf.WriteString(s.String())
	}
	return buf.String()
}

type Enum struct {
	Name   string
	Values []*EnumValue
}

func (e *Enum) String() string {
	var buf bytes.Buffer
	buf.WriteString("enum " + e.Name + " {\n")
	for _, v := range e.Values {
		buf.WriteString("  " + v.Name + " = " + v.Value + ";\n")
	}
	buf.WriteString("}\n")
	return buf.String()
}

type EnumValue struct {
	Name  string
	Value string
}

type Struct struct {
	Name   string
	Fields []*Field
}

func (s *Struct) String() string {
	var buf bytes.Buffer
	buf.WriteString("struct " + s.Name + " {\n")
	for _, f := range s.Fields {
		buf.WriteString("  " + f.Tag + ": " + f.Type + " " + f.Name + ";\n")
	}
	buf.WriteString("}\n")
	return buf.String()
}

type Field struct {
	Tag  string
	Type string // e.g. "int64", "list<string>"
	Name string
}

type Service struct {
	Name string
	RPCs []*RPC
}

func (s *Service) String() string {
	var buf bytes.Buffer
	buf.WriteString("service " + s.Name + " {\n")
	for _, r := range s.RPCs {
		buf.WriteString("  " + r.String() + "\n")
	}
	buf.WriteString("}\n")
	return buf.String()
}

type RPC struct {
	Name       string
	ReqType    string
	ReqStream  bool
	RespType   string
	RespStream bool
}

func (r *RPC) String() string {
	req := r.ReqType
	if r.ReqStream {
		req = "stream(" + req + ")"
	}
	resp := r.RespType
	if r.RespStream {
		resp = "stream(" + resp + ")"
	}
	return "rpc " + r.Name + "(" + req + ") -> (" + resp + ");"
}
