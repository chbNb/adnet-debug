// Code generated by protoc-gen-go. DO NOT EDIT.
// source: creative.proto

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

type Creative struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	VideoLength          int32    `protobuf:"varint,2,opt,name=videoLength,proto3" json:"videoLength,omitempty"`
	VideoSize            int32    `protobuf:"varint,3,opt,name=videoSize,proto3" json:"videoSize,omitempty"`
	VideoResolution      string   `protobuf:"bytes,4,opt,name=videoResolution,proto3" json:"videoResolution,omitempty"`
	Width                int32    `protobuf:"varint,5,opt,name=width,proto3" json:"width,omitempty"`
	Height               int32    `protobuf:"varint,6,opt,name=height,proto3" json:"height,omitempty"`
	WatchMile            int32    `protobuf:"varint,7,opt,name=watchMile,proto3" json:"watchMile,omitempty"`
	BitRate              int32    `protobuf:"varint,8,opt,name=bitRate,proto3" json:"bitRate,omitempty"`
	IValue               int32    `protobuf:"varint,9,opt,name=iValue,proto3" json:"iValue,omitempty"`
	SValue               string   `protobuf:"bytes,10,opt,name=sValue,proto3" json:"sValue,omitempty"`
	FValue               float64  `protobuf:"fixed64,11,opt,name=fValue,proto3" json:"fValue,omitempty"`
	Resolution           string   `protobuf:"bytes,12,opt,name=resolution,proto3" json:"resolution,omitempty"`
	Mime                 string   `protobuf:"bytes,13,opt,name=mime,proto3" json:"mime,omitempty"`
	AdvCreativeId        string   `protobuf:"bytes,14,opt,name=advCreativeId,proto3" json:"advCreativeId,omitempty"`
	CreativeId           int64    `protobuf:"varint,15,opt,name=creativeId,proto3" json:"creativeId,omitempty"`
	FMd5                 string   `protobuf:"bytes,16,opt,name=fMd5,proto3" json:"fMd5,omitempty"`
	Source               int32    `protobuf:"varint,17,opt,name=source,proto3" json:"source,omitempty"`
	Orientation          int32    `protobuf:"varint,18,opt,name=orientation,proto3" json:"orientation,omitempty"`
	Protocal             int32    `protobuf:"varint,19,opt,name=protocal,proto3" json:"protocal,omitempty"`
	Cname                string   `protobuf:"bytes,20,opt,name=cname,proto3" json:"cname,omitempty"`
	UniqCId              int64    `protobuf:"varint,21,opt,name=uniqCId,proto3" json:"uniqCId,omitempty"`
	CsetId               int64    `protobuf:"varint,22,opt,name=csetId,proto3" json:"csetId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Creative) Reset()         { *m = Creative{} }
func (m *Creative) String() string { return proto.CompactTextString(m) }
func (*Creative) ProtoMessage()    {}
func (*Creative) Descriptor() ([]byte, []int) {
	return fileDescriptor_6db1d9509772806d, []int{0}
}

func (m *Creative) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Creative.Unmarshal(m, b)
}
func (m *Creative) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Creative.Marshal(b, m, deterministic)
}
func (m *Creative) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Creative.Merge(m, src)
}
func (m *Creative) XXX_Size() int {
	return xxx_messageInfo_Creative.Size(m)
}
func (m *Creative) XXX_DiscardUnknown() {
	xxx_messageInfo_Creative.DiscardUnknown(m)
}

var xxx_messageInfo_Creative proto.InternalMessageInfo

func (m *Creative) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Creative) GetVideoLength() int32 {
	if m != nil {
		return m.VideoLength
	}
	return 0
}

func (m *Creative) GetVideoSize() int32 {
	if m != nil {
		return m.VideoSize
	}
	return 0
}

func (m *Creative) GetVideoResolution() string {
	if m != nil {
		return m.VideoResolution
	}
	return ""
}

func (m *Creative) GetWidth() int32 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *Creative) GetHeight() int32 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *Creative) GetWatchMile() int32 {
	if m != nil {
		return m.WatchMile
	}
	return 0
}

func (m *Creative) GetBitRate() int32 {
	if m != nil {
		return m.BitRate
	}
	return 0
}

func (m *Creative) GetIValue() int32 {
	if m != nil {
		return m.IValue
	}
	return 0
}

func (m *Creative) GetSValue() string {
	if m != nil {
		return m.SValue
	}
	return ""
}

func (m *Creative) GetFValue() float64 {
	if m != nil {
		return m.FValue
	}
	return 0
}

func (m *Creative) GetResolution() string {
	if m != nil {
		return m.Resolution
	}
	return ""
}

func (m *Creative) GetMime() string {
	if m != nil {
		return m.Mime
	}
	return ""
}

func (m *Creative) GetAdvCreativeId() string {
	if m != nil {
		return m.AdvCreativeId
	}
	return ""
}

func (m *Creative) GetCreativeId() int64 {
	if m != nil {
		return m.CreativeId
	}
	return 0
}

func (m *Creative) GetFMd5() string {
	if m != nil {
		return m.FMd5
	}
	return ""
}

func (m *Creative) GetSource() int32 {
	if m != nil {
		return m.Source
	}
	return 0
}

func (m *Creative) GetOrientation() int32 {
	if m != nil {
		return m.Orientation
	}
	return 0
}

func (m *Creative) GetProtocal() int32 {
	if m != nil {
		return m.Protocal
	}
	return 0
}

func (m *Creative) GetCname() string {
	if m != nil {
		return m.Cname
	}
	return ""
}

func (m *Creative) GetUniqCId() int64 {
	if m != nil {
		return m.UniqCId
	}
	return 0
}

func (m *Creative) GetCsetId() int64 {
	if m != nil {
		return m.CsetId
	}
	return 0
}

func init() {
	proto.RegisterType((*Creative)(nil), "protobuf.Creative")
}

func init() { proto.RegisterFile("creative.proto", fileDescriptor_6db1d9509772806d) }

var fileDescriptor_6db1d9509772806d = []byte{
	// 353 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x92, 0xcd, 0x4e, 0xf2, 0x40,
	0x14, 0x86, 0xd3, 0x8f, 0xff, 0xc3, 0xc7, 0x8f, 0x47, 0x24, 0x27, 0xc6, 0x98, 0xc6, 0xb8, 0xe8,
	0xca, 0x8d, 0xf1, 0x0a, 0x58, 0x91, 0xc8, 0xa6, 0x26, 0xee, 0x87, 0xf6, 0x40, 0x27, 0x29, 0x1d,
	0x2d, 0x53, 0x48, 0xbc, 0x5a, 0x2f, 0xc5, 0xcc, 0x99, 0x02, 0xd5, 0x15, 0xf3, 0x3c, 0x2f, 0xbc,
	0x73, 0x66, 0x06, 0x18, 0x27, 0x25, 0x2b, 0xab, 0x0f, 0xfc, 0xf4, 0x51, 0x1a, 0x6b, 0xb0, 0x2f,
	0x1f, 0xeb, 0x6a, 0xf3, 0xf0, 0xdd, 0x86, 0xfe, 0xa2, 0x0e, 0x71, 0x0a, 0xad, 0xaa, 0xcc, 0x29,
	0x08, 0x83, 0x68, 0x10, 0xbb, 0x25, 0x86, 0x30, 0x3c, 0xe8, 0x94, 0xcd, 0x2b, 0x17, 0x5b, 0x9b,
	0xd1, 0xbf, 0x30, 0x88, 0x3a, 0x71, 0x53, 0xe1, 0x1d, 0x0c, 0x04, 0xdf, 0xf4, 0x17, 0x53, 0x4b,
	0xf2, 0x8b, 0xc0, 0x08, 0x26, 0x02, 0x31, 0xef, 0x4d, 0x5e, 0x59, 0x6d, 0x0a, 0x6a, 0x4b, 0xfb,
	0x5f, 0x8d, 0x33, 0xe8, 0x1c, 0x75, 0x6a, 0x33, 0xea, 0x48, 0x87, 0x07, 0x9c, 0x43, 0x37, 0x63,
	0xbd, 0xcd, 0x2c, 0x75, 0x45, 0xd7, 0xe4, 0x76, 0x3d, 0x2a, 0x9b, 0x64, 0x2b, 0x9d, 0x33, 0xf5,
	0xfc, 0xae, 0x67, 0x81, 0x04, 0xbd, 0xb5, 0xb6, 0xb1, 0xb2, 0x4c, 0x7d, 0xc9, 0x4e, 0xe8, 0xfa,
	0xf4, 0xbb, 0xca, 0x2b, 0xa6, 0x81, 0xef, 0xf3, 0xe4, 0xfc, 0xde, 0x7b, 0x90, 0xf1, 0x6a, 0x72,
	0x7e, 0xe3, 0xfd, 0x30, 0x0c, 0xa2, 0x20, 0xae, 0x09, 0xef, 0x01, 0xca, 0xcb, 0x91, 0xfe, 0xcb,
	0x6f, 0x1a, 0x06, 0x11, 0xda, 0x3b, 0xbd, 0x63, 0x1a, 0x49, 0x22, 0x6b, 0x7c, 0x84, 0x91, 0x4a,
	0x0f, 0xa7, 0xcb, 0x5e, 0xa6, 0x34, 0x96, 0xf0, 0xb7, 0x74, 0xcd, 0xc9, 0xe5, 0x2b, 0x93, 0x30,
	0x88, 0x5a, 0x71, 0xc3, 0xb8, 0xe6, 0xcd, 0x2a, 0x7d, 0xa1, 0xa9, 0x6f, 0x76, 0x6b, 0x99, 0xde,
	0x54, 0x65, 0xc2, 0x74, 0xe5, 0x4f, 0xe5, 0xc9, 0xbd, 0x9e, 0x29, 0x35, 0x17, 0x56, 0xc9, 0x98,
	0xe8, 0x5f, 0xaf, 0xa1, 0xf0, 0x16, 0xfc, 0x5f, 0x21, 0x51, 0x39, 0x5d, 0x4b, 0x7c, 0x66, 0xf7,
	0x22, 0x49, 0xa1, 0x76, 0x4c, 0x33, 0xd9, 0xca, 0x83, 0xbb, 0xdb, 0xaa, 0xd0, 0x9f, 0x8b, 0x65,
	0x4a, 0x37, 0x32, 0xdc, 0x09, 0xdd, 0x14, 0xc9, 0x9e, 0xed, 0x32, 0xa5, 0xb9, 0x04, 0x35, 0xad,
	0xbb, 0xd2, 0xf8, 0xfc, 0x13, 0x00, 0x00, 0xff, 0xff, 0x3d, 0x81, 0x14, 0x08, 0x85, 0x02, 0x00,
	0x00,
}