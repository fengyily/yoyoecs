// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.6.1
// source: sku.proto

package protoc

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Sku struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          int32   `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	SkuName     string  `protobuf:"bytes,2,opt,name=SkuName,proto3" json:"SkuName,omitempty"`
	Price       int64   `protobuf:"varint,3,opt,name=Price,proto3" json:"Price,omitempty"`
	PriceUnit   int64   `protobuf:"varint,5,opt,name=PriceUnit,proto3" json:"PriceUnit,omitempty"`
	PluCode     string  `protobuf:"bytes,6,opt,name=PluCode,proto3" json:"PluCode,omitempty"`
	Stock       float32 `protobuf:"fixed32,7,opt,name=Stock,proto3" json:"Stock,omitempty"`
	IsOn        int32   `protobuf:"varint,8,opt,name=IsOn,proto3" json:"IsOn,omitempty"`
	YoyoCode    int64   `protobuf:"varint,9,opt,name=YoyoCode,proto3" json:"YoyoCode,omitempty"`
	ImageUrl    string  `protobuf:"bytes,10,opt,name=ImageUrl,proto3" json:"ImageUrl,omitempty"`
	DisplayName string  `protobuf:"bytes,11,opt,name=DisplayName,proto3" json:"DisplayName,omitempty"`
	IsLock      int32   `protobuf:"varint,12,opt,name=IsLock,proto3" json:"IsLock,omitempty"`
}

func (x *Sku) Reset() {
	*x = Sku{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sku_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sku) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sku) ProtoMessage() {}

func (x *Sku) ProtoReflect() protoreflect.Message {
	mi := &file_sku_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sku.ProtoReflect.Descriptor instead.
func (*Sku) Descriptor() ([]byte, []int) {
	return file_sku_proto_rawDescGZIP(), []int{0}
}

func (x *Sku) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Sku) GetSkuName() string {
	if x != nil {
		return x.SkuName
	}
	return ""
}

func (x *Sku) GetPrice() int64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *Sku) GetPriceUnit() int64 {
	if x != nil {
		return x.PriceUnit
	}
	return 0
}

func (x *Sku) GetPluCode() string {
	if x != nil {
		return x.PluCode
	}
	return ""
}

func (x *Sku) GetStock() float32 {
	if x != nil {
		return x.Stock
	}
	return 0
}

func (x *Sku) GetIsOn() int32 {
	if x != nil {
		return x.IsOn
	}
	return 0
}

func (x *Sku) GetYoyoCode() int64 {
	if x != nil {
		return x.YoyoCode
	}
	return 0
}

func (x *Sku) GetImageUrl() string {
	if x != nil {
		return x.ImageUrl
	}
	return ""
}

func (x *Sku) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

func (x *Sku) GetIsLock() int32 {
	if x != nil {
		return x.IsLock
	}
	return 0
}

type SkuList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sku []*Sku `protobuf:"bytes,1,rep,name=sku,proto3" json:"sku,omitempty"`
}

func (x *SkuList) Reset() {
	*x = SkuList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sku_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SkuList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SkuList) ProtoMessage() {}

func (x *SkuList) ProtoReflect() protoreflect.Message {
	mi := &file_sku_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SkuList.ProtoReflect.Descriptor instead.
func (*SkuList) Descriptor() ([]byte, []int) {
	return file_sku_proto_rawDescGZIP(), []int{1}
}

func (x *SkuList) GetSku() []*Sku {
	if x != nil {
		return x.Sku
	}
	return nil
}

var File_sku_proto protoreflect.FileDescriptor

var file_sku_proto_rawDesc = []byte{
	0x0a, 0x09, 0x73, 0x6b, 0x75, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x22, 0x99, 0x02, 0x0a, 0x03, 0x53, 0x6b, 0x75, 0x12, 0x0e, 0x0a, 0x02, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x53,
	0x6b, 0x75, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x53, 0x6b,
	0x75, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x50,
	0x72, 0x69, 0x63, 0x65, 0x55, 0x6e, 0x69, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x50, 0x72, 0x69, 0x63, 0x65, 0x55, 0x6e, 0x69, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x50, 0x6c, 0x75,
	0x43, 0x6f, 0x64, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x50, 0x6c, 0x75, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x02, 0x52, 0x05, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x49, 0x73, 0x4f,
	0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x49, 0x73, 0x4f, 0x6e, 0x12, 0x1a, 0x0a,
	0x08, 0x59, 0x6f, 0x79, 0x6f, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x08, 0x59, 0x6f, 0x79, 0x6f, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x55, 0x72, 0x6c, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x20, 0x0a, 0x0b, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x44, 0x69, 0x73, 0x70,
	0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x49, 0x73, 0x4c, 0x6f, 0x63,
	0x6b, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x49, 0x73, 0x4c, 0x6f, 0x63, 0x6b, 0x22,
	0x28, 0x0a, 0x07, 0x53, 0x6b, 0x75, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x03, 0x73, 0x6b,
	0x75, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x2e, 0x53, 0x6b, 0x75, 0x52, 0x03, 0x73, 0x6b, 0x75, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_sku_proto_rawDescOnce sync.Once
	file_sku_proto_rawDescData = file_sku_proto_rawDesc
)

func file_sku_proto_rawDescGZIP() []byte {
	file_sku_proto_rawDescOnce.Do(func() {
		file_sku_proto_rawDescData = protoimpl.X.CompressGZIP(file_sku_proto_rawDescData)
	})
	return file_sku_proto_rawDescData
}

var file_sku_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_sku_proto_goTypes = []interface{}{
	(*Sku)(nil),     // 0: protoc.Sku
	(*SkuList)(nil), // 1: protoc.SkuList
}
var file_sku_proto_depIdxs = []int32{
	0, // 0: protoc.SkuList.sku:type_name -> protoc.Sku
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_sku_proto_init() }
func file_sku_proto_init() {
	if File_sku_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sku_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sku); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_sku_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SkuList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_sku_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_sku_proto_goTypes,
		DependencyIndexes: file_sku_proto_depIdxs,
		MessageInfos:      file_sku_proto_msgTypes,
	}.Build()
	File_sku_proto = out.File
	file_sku_proto_rawDesc = nil
	file_sku_proto_goTypes = nil
	file_sku_proto_depIdxs = nil
}
