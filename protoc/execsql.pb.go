// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: execsql.proto

package protoc

import (
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

// protoc: v3.6.1
// package protoc;
type ExecType int32

const (
	ExecType_SELECT ExecType = 0
	ExecType_INSERT ExecType = 1
	ExecType_EXEC   ExecType = 2
	ExecType_COUNT  ExecType = 3
	ExecType_ONE    ExecType = 4
)

// Enum value maps for ExecType.
var (
	ExecType_name = map[int32]string{
		0: "SELECT",
		1: "INSERT",
		2: "EXEC",
		3: "COUNT",
		4: "ONE",
	}
	ExecType_value = map[string]int32{
		"SELECT": 0,
		"INSERT": 1,
		"EXEC":   2,
		"COUNT":  3,
		"ONE":    4,
	}
)

func (x ExecType) Enum() *ExecType {
	p := new(ExecType)
	*p = x
	return p
}

func (x ExecType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ExecType) Descriptor() protoreflect.EnumDescriptor {
	return file_execsql_proto_enumTypes[0].Descriptor()
}

func (ExecType) Type() protoreflect.EnumType {
	return &file_execsql_proto_enumTypes[0]
}

func (x ExecType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ExecType.Descriptor instead.
func (ExecType) EnumDescriptor() ([]byte, []int) {
	return file_execsql_proto_rawDescGZIP(), []int{0}
}

type ExecSql struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SqlStr  string   `protobuf:"bytes,1,opt,name=SqlStr,proto3" json:"SqlStr,omitempty"`
	OptType ExecType `protobuf:"varint,2,opt,name=OptType,proto3,enum=ExecType" json:"OptType,omitempty"`
	Args    []string `protobuf:"bytes,3,rep,name=Args,proto3" json:"Args,omitempty"`
}

func (x *ExecSql) Reset() {
	*x = ExecSql{}
	if protoimpl.UnsafeEnabled {
		mi := &file_execsql_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecSql) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecSql) ProtoMessage() {}

func (x *ExecSql) ProtoReflect() protoreflect.Message {
	mi := &file_execsql_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecSql.ProtoReflect.Descriptor instead.
func (*ExecSql) Descriptor() ([]byte, []int) {
	return file_execsql_proto_rawDescGZIP(), []int{0}
}

func (x *ExecSql) GetSqlStr() string {
	if x != nil {
		return x.SqlStr
	}
	return ""
}

func (x *ExecSql) GetOptType() ExecType {
	if x != nil {
		return x.OptType
	}
	return ExecType_SELECT
}

func (x *ExecSql) GetArgs() []string {
	if x != nil {
		return x.Args
	}
	return nil
}

type ExecSqlForJSONReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result string `protobuf:"bytes,1,opt,name=Result,proto3" json:"Result,omitempty"`
}

func (x *ExecSqlForJSONReply) Reset() {
	*x = ExecSqlForJSONReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_execsql_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecSqlForJSONReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecSqlForJSONReply) ProtoMessage() {}

func (x *ExecSqlForJSONReply) ProtoReflect() protoreflect.Message {
	mi := &file_execsql_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecSqlForJSONReply.ProtoReflect.Descriptor instead.
func (*ExecSqlForJSONReply) Descriptor() ([]byte, []int) {
	return file_execsql_proto_rawDescGZIP(), []int{1}
}

func (x *ExecSqlForJSONReply) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

var File_execsql_proto protoreflect.FileDescriptor

var file_execsql_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x65, 0x78, 0x65, 0x63, 0x73, 0x71, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x5a, 0x0a, 0x07, 0x45, 0x78, 0x65, 0x63, 0x53, 0x71, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x71,
	0x6c, 0x53, 0x74, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x71, 0x6c, 0x53,
	0x74, 0x72, 0x12, 0x23, 0x0a, 0x07, 0x4f, 0x70, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x09, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x54, 0x79, 0x70, 0x65, 0x52, 0x07,
	0x4f, 0x70, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x41, 0x72, 0x67, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x41, 0x72, 0x67, 0x73, 0x22, 0x2d, 0x0a, 0x13, 0x45,
	0x78, 0x65, 0x63, 0x53, 0x71, 0x6c, 0x46, 0x6f, 0x72, 0x4a, 0x53, 0x4f, 0x4e, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x2a, 0x40, 0x0a, 0x08, 0x45, 0x78,
	0x65, 0x63, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x45, 0x4c, 0x45, 0x43, 0x54,
	0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x49, 0x4e, 0x53, 0x45, 0x52, 0x54, 0x10, 0x01, 0x12, 0x08,
	0x0a, 0x04, 0x45, 0x58, 0x45, 0x43, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x43, 0x4f, 0x55, 0x4e,
	0x54, 0x10, 0x03, 0x12, 0x07, 0x0a, 0x03, 0x4f, 0x4e, 0x45, 0x10, 0x04, 0x42, 0x15, 0x5a, 0x13,
	0x2e, 0x2e, 0x2f, 0x2e, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x3b, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_execsql_proto_rawDescOnce sync.Once
	file_execsql_proto_rawDescData = file_execsql_proto_rawDesc
)

func file_execsql_proto_rawDescGZIP() []byte {
	file_execsql_proto_rawDescOnce.Do(func() {
		file_execsql_proto_rawDescData = protoimpl.X.CompressGZIP(file_execsql_proto_rawDescData)
	})
	return file_execsql_proto_rawDescData
}

var file_execsql_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_execsql_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_execsql_proto_goTypes = []interface{}{
	(ExecType)(0),               // 0: ExecType
	(*ExecSql)(nil),             // 1: ExecSql
	(*ExecSqlForJSONReply)(nil), // 2: ExecSqlForJSONReply
}
var file_execsql_proto_depIdxs = []int32{
	0, // 0: ExecSql.OptType:type_name -> ExecType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_execsql_proto_init() }
func file_execsql_proto_init() {
	if File_execsql_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_execsql_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecSql); i {
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
		file_execsql_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecSqlForJSONReply); i {
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
			RawDescriptor: file_execsql_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_execsql_proto_goTypes,
		DependencyIndexes: file_execsql_proto_depIdxs,
		EnumInfos:         file_execsql_proto_enumTypes,
		MessageInfos:      file_execsql_proto_msgTypes,
	}.Build()
	File_execsql_proto = out.File
	file_execsql_proto_rawDesc = nil
	file_execsql_proto_goTypes = nil
	file_execsql_proto_depIdxs = nil
}