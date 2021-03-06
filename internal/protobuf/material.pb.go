// Code generated by protoc-gen-go. DO NOT EDIT.
// source: material.proto

package protobuf

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Material struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	VideoLength          int32    `protobuf:"varint,2,opt,name=videoLength,proto3" json:"videoLength,omitempty"`
	VideoSize            int32    `protobuf:"varint,3,opt,name=videoSize,proto3" json:"videoSize,omitempty"`
	Width                int32    `protobuf:"varint,4,opt,name=width,proto3" json:"width,omitempty"`
	Height               int32    `protobuf:"varint,5,opt,name=height,proto3" json:"height,omitempty"`
	BitRate              int32    `protobuf:"varint,6,opt,name=bitRate,proto3" json:"bitRate,omitempty"`
	IValue               int32    `protobuf:"varint,7,opt,name=iValue,proto3" json:"iValue,omitempty"`
	SValue               string   `protobuf:"bytes,8,opt,name=sValue,proto3" json:"sValue,omitempty"`
	FValue               float64  `protobuf:"fixed64,9,opt,name=fValue,proto3" json:"fValue,omitempty"`
	Resolution           string   `protobuf:"bytes,10,opt,name=resolution,proto3" json:"resolution,omitempty"`
	Mime                 string   `protobuf:"bytes,11,opt,name=mime,proto3" json:"mime,omitempty"`
	FMd5                 string   `protobuf:"bytes,12,opt,name=fMd5,proto3" json:"fMd5,omitempty"`
	Orientation          int32    `protobuf:"varint,13,opt,name=orientation,proto3" json:"orientation,omitempty"`
	Protocol             int32    `protobuf:"varint,14,opt,name=protocol,proto3" json:"protocol,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Material) Reset()         { *m = Material{} }
func (m *Material) String() string { return proto.CompactTextString(m) }
func (*Material) ProtoMessage()    {}
func (*Material) Descriptor() ([]byte, []int) {
	return fileDescriptor_154cda5588300428, []int{0}
}

func (m *Material) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Material.Unmarshal(m, b)
}
func (m *Material) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Material.Marshal(b, m, deterministic)
}
func (m *Material) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Material.Merge(m, src)
}
func (m *Material) XXX_Size() int {
	return xxx_messageInfo_Material.Size(m)
}
func (m *Material) XXX_DiscardUnknown() {
	xxx_messageInfo_Material.DiscardUnknown(m)
}

var xxx_messageInfo_Material proto.InternalMessageInfo

func (m *Material) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Material) GetVideoLength() int32 {
	if m != nil {
		return m.VideoLength
	}
	return 0
}

func (m *Material) GetVideoSize() int32 {
	if m != nil {
		return m.VideoSize
	}
	return 0
}

func (m *Material) GetWidth() int32 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *Material) GetHeight() int32 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *Material) GetBitRate() int32 {
	if m != nil {
		return m.BitRate
	}
	return 0
}

func (m *Material) GetIValue() int32 {
	if m != nil {
		return m.IValue
	}
	return 0
}

func (m *Material) GetSValue() string {
	if m != nil {
		return m.SValue
	}
	return ""
}

func (m *Material) GetFValue() float64 {
	if m != nil {
		return m.FValue
	}
	return 0
}

func (m *Material) GetResolution() string {
	if m != nil {
		return m.Resolution
	}
	return ""
}

func (m *Material) GetMime() string {
	if m != nil {
		return m.Mime
	}
	return ""
}

func (m *Material) GetFMd5() string {
	if m != nil {
		return m.FMd5
	}
	return ""
}

func (m *Material) GetOrientation() int32 {
	if m != nil {
		return m.Orientation
	}
	return 0
}

func (m *Material) GetProtocol() int32 {
	if m != nil {
		return m.Protocol
	}
	return 0
}

func init() {
	proto.RegisterType((*Material)(nil), "protobuf.Material")
}

func init() { proto.RegisterFile("material.proto", fileDescriptor_154cda5588300428) }

var fileDescriptor_154cda5588300428 = []byte{
	// 252 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0x49, 0xdb, 0xa4, 0xc9, 0x54, 0x8b, 0x0c, 0x22, 0x83, 0x88, 0x04, 0x4f, 0x39, 0x79,
	0x11, 0xdf, 0xc2, 0x5e, 0x22, 0x78, 0x4f, 0xcc, 0xa4, 0x19, 0x48, 0xb2, 0xb2, 0xdd, 0x28, 0xf8,
	0xd8, 0x3e, 0x81, 0xec, 0x6c, 0xb4, 0x39, 0xe5, 0xff, 0xbe, 0x9f, 0x09, 0xcb, 0x0f, 0xfb, 0xa1,
	0x72, 0x6c, 0xa5, 0xea, 0x1f, 0x3f, 0xac, 0x71, 0x06, 0x53, 0xfd, 0xd4, 0x53, 0xfb, 0xf0, 0xb3,
	0x82, 0xf4, 0x30, 0x97, 0x78, 0x05, 0xeb, 0xc9, 0xf6, 0x14, 0xe5, 0x51, 0x91, 0x95, 0x3e, 0x62,
	0x0e, 0xbb, 0x4f, 0x69, 0xd8, 0xbc, 0xf0, 0x78, 0x74, 0x1d, 0xad, 0xf2, 0xa8, 0x88, 0xcb, 0xa5,
	0xc2, 0x3b, 0xc8, 0x14, 0x5f, 0xe5, 0x9b, 0x69, 0xad, 0xfd, 0x59, 0xe0, 0x35, 0xc4, 0x5f, 0xd2,
	0xb8, 0x8e, 0x36, 0xda, 0x04, 0xc0, 0x1b, 0x48, 0x3a, 0x96, 0x63, 0xe7, 0x28, 0x56, 0x3d, 0x13,
	0x12, 0x6c, 0x6b, 0x71, 0x65, 0xe5, 0x98, 0x12, 0x2d, 0xfe, 0xd0, 0x5f, 0xc8, 0x5b, 0xd5, 0x4f,
	0x4c, 0xdb, 0x70, 0x11, 0xc8, 0xfb, 0x53, 0xf0, 0xa9, 0x3e, 0x7a, 0x26, 0xef, 0xdb, 0xe0, 0xb3,
	0x3c, 0x2a, 0xa2, 0x72, 0x26, 0xbc, 0x07, 0xb0, 0x7c, 0x32, 0xfd, 0xe4, 0xc4, 0x8c, 0x04, 0x7a,
	0xb3, 0x30, 0x88, 0xb0, 0x19, 0x64, 0x60, 0xda, 0x69, 0xa3, 0xd9, 0xbb, 0xf6, 0xd0, 0x3c, 0xd3,
	0x45, 0x70, 0x3e, 0xfb, 0x5d, 0x8c, 0x15, 0x1e, 0x5d, 0xa5, 0x3f, 0xba, 0x0c, 0xbb, 0x2c, 0x14,
	0xde, 0x42, 0x18, 0xf9, 0xdd, 0xf4, 0xb4, 0xd7, 0xfa, 0x9f, 0xeb, 0x44, 0xd3, 0xd3, 0x6f, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x23, 0xb4, 0xdb, 0x17, 0x97, 0x01, 0x00, 0x00,
}
