// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpc/vector_clock.proto

package rpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type VC struct {
	Store map[string]uint64 `protobuf:"bytes,1,rep,name=store" json:"store,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
}

func (m *VC) Reset()                    { *m = VC{} }
func (m *VC) String() string            { return proto.CompactTextString(m) }
func (*VC) ProtoMessage()               {}
func (*VC) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *VC) GetStore() map[string]uint64 {
	if m != nil {
		return m.Store
	}
	return nil
}

func init() {
	proto.RegisterType((*VC)(nil), "rpc.VC")
}

func init() { proto.RegisterFile("rpc/vector_clock.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 143 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2b, 0x2a, 0x48, 0xd6,
	0x2f, 0x4b, 0x4d, 0x2e, 0xc9, 0x2f, 0x8a, 0x4f, 0xce, 0xc9, 0x4f, 0xce, 0xd6, 0x2b, 0x28, 0xca,
	0x2f, 0xc9, 0x17, 0x62, 0x2e, 0x2a, 0x48, 0x56, 0xca, 0xe0, 0x62, 0x0a, 0x73, 0x16, 0xd2, 0xe0,
	0x62, 0x2d, 0x2e, 0xc9, 0x2f, 0x4a, 0x95, 0x60, 0x54, 0x60, 0xd6, 0xe0, 0x36, 0x12, 0xd2, 0x2b,
	0x2a, 0x48, 0xd6, 0x0b, 0x73, 0xd6, 0x0b, 0x06, 0x09, 0xba, 0xe6, 0x95, 0x14, 0x55, 0x06, 0x41,
	0x14, 0x48, 0x59, 0x70, 0x71, 0x21, 0x04, 0x85, 0x04, 0xb8, 0x98, 0xb3, 0x53, 0x2b, 0x25, 0x18,
	0x15, 0x18, 0x35, 0x38, 0x83, 0x40, 0x4c, 0x21, 0x11, 0x2e, 0xd6, 0xb2, 0xc4, 0x9c, 0xd2, 0x54,
	0x09, 0x26, 0x05, 0x46, 0x0d, 0x96, 0x20, 0x08, 0xc7, 0x8a, 0xc9, 0x82, 0x31, 0x89, 0x0d, 0x6c,
	0xab, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x2f, 0x04, 0x60, 0x14, 0x8f, 0x00, 0x00, 0x00,
}
