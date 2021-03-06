// Code generated by protoc-gen-go.
// source: pool/program.proto
// DO NOT EDIT!

/*
Package pool is a generated protocol buffer package.

It is generated from these files:
	pool/program.proto

It has these top-level messages:
	Program
*/
package pool

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Program struct {
	Name       string   `protobuf:"bytes,1,opt,name=Name" json:"Name,omitempty"`
	Location   string   `protobuf:"bytes,2,opt,name=Location" json:"Location,omitempty"`
	UI         bool     `protobuf:"varint,3,opt,name=UI" json:"UI,omitempty"`
	Key        []byte   `protobuf:"bytes,4,opt,name=Key,proto3" json:"Key,omitempty"`
	Port32     uint32   `protobuf:"varint,5,opt,name=Port32" json:"Port32,omitempty"`
	Start      bool     `protobuf:"varint,6,opt,name=Start" json:"Start,omitempty"`
	Implements []string `protobuf:"bytes,7,rep,name=Implements" json:"Implements,omitempty"`
}

func (m *Program) Reset()                    { *m = Program{} }
func (m *Program) String() string            { return proto.CompactTextString(m) }
func (*Program) ProtoMessage()               {}
func (*Program) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Program) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Program) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

func (m *Program) GetUI() bool {
	if m != nil {
		return m.UI
	}
	return false
}

func (m *Program) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Program) GetPort32() uint32 {
	if m != nil {
		return m.Port32
	}
	return 0
}

func (m *Program) GetStart() bool {
	if m != nil {
		return m.Start
	}
	return false
}

func (m *Program) GetImplements() []string {
	if m != nil {
		return m.Implements
	}
	return nil
}

func init() {
	proto.RegisterType((*Program)(nil), "pool.Program")
}

func init() { proto.RegisterFile("pool/program.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 182 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0xce, 0xcf, 0x8a, 0xc2, 0x30,
	0x10, 0xc7, 0x71, 0xd2, 0xff, 0x1d, 0x76, 0x97, 0x65, 0x58, 0x96, 0xe0, 0x41, 0x82, 0xa7, 0x9c,
	0x14, 0xec, 0x53, 0x14, 0x45, 0x4a, 0xa4, 0x0f, 0x10, 0x25, 0x88, 0xd0, 0x74, 0x42, 0xcc, 0xc5,
	0x47, 0xf2, 0x2d, 0xa5, 0xa9, 0x88, 0xb7, 0xdf, 0xf7, 0x03, 0x03, 0x03, 0xe8, 0x88, 0x86, 0x8d,
	0xf3, 0x74, 0xf1, 0xda, 0xae, 0x9d, 0xa7, 0x40, 0x98, 0x4d, 0xb6, 0x7a, 0x30, 0x28, 0xbb, 0xd9,
	0x11, 0x21, 0x3b, 0x68, 0x6b, 0x38, 0x13, 0x4c, 0xd6, 0x2a, 0x6e, 0x5c, 0x40, 0xb5, 0xa7, 0xb3,
	0x0e, 0x57, 0x1a, 0x79, 0x12, 0xfd, 0xdd, 0xf8, 0x03, 0x49, 0xdf, 0xf2, 0x54, 0x30, 0x59, 0xa9,
	0xa4, 0x6f, 0xf1, 0x17, 0xd2, 0x9d, 0xb9, 0xf3, 0x4c, 0x30, 0xf9, 0xa5, 0xa6, 0x89, 0xff, 0x50,
	0x74, 0xe4, 0x43, 0xb3, 0xe5, 0xb9, 0x60, 0xf2, 0x5b, 0xbd, 0x0a, 0xff, 0x20, 0x3f, 0x06, 0xed,
	0x03, 0x2f, 0xe2, 0xf1, 0x1c, 0xb8, 0x04, 0x68, 0xad, 0x1b, 0x8c, 0x35, 0x63, 0xb8, 0xf1, 0x52,
	0xa4, 0xb2, 0x56, 0x1f, 0x72, 0x2a, 0xe2, 0xe3, 0xcd, 0x33, 0x00, 0x00, 0xff, 0xff, 0xd2, 0x7c,
	0x61, 0xca, 0xce, 0x00, 0x00, 0x00,
}
