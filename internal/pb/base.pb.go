// TODO message结构的属性名请使用 驼峰小写字母开始
// TODO enum枚举的属性名请使用 驼峰大写字母开始

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v5.28.3
// source: base.proto

package pb

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

// 响应错误
// 当客户端发送Request类型消息时，如果有异常则返回S2CError
// 踢下线协议也会用到该结构
// push route:error
type S2CError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`      // 错误码
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"` // 错误内容
}

func (x *S2CError) Reset() {
	*x = S2CError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_base_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *S2CError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*S2CError) ProtoMessage() {}

func (x *S2CError) ProtoReflect() protoreflect.Message {
	mi := &file_base_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use S2CError.ProtoReflect.Descriptor instead.
func (*S2CError) Descriptor() ([]byte, []int) {
	return file_base_proto_rawDescGZIP(), []int{0}
}

func (x *S2CError) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *S2CError) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// 资源变更
// push route:resUpdate
type S2CResUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResInfos []*ResInfo `protobuf:"bytes,1,rep,name=resInfos,proto3" json:"resInfos,omitempty"` //玩家资源集合
}

func (x *S2CResUpdate) Reset() {
	*x = S2CResUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_base_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *S2CResUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*S2CResUpdate) ProtoMessage() {}

func (x *S2CResUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_base_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use S2CResUpdate.ProtoReflect.Descriptor instead.
func (*S2CResUpdate) Descriptor() ([]byte, []int) {
	return file_base_proto_rawDescGZIP(), []int{1}
}

func (x *S2CResUpdate) GetResInfos() []*ResInfo {
	if x != nil {
		return x.ResInfos
	}
	return nil
}

var File_base_proto protoreflect.FileDescriptor

var file_base_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62,
	0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x38,
	0x0a, 0x08, 0x53, 0x32, 0x43, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x37, 0x0a, 0x0c, 0x53, 0x32, 0x43, 0x52,
	0x65, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x27, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x49,
	0x6e, 0x66, 0x6f, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x62, 0x2e,
	0x52, 0x65, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x72, 0x65, 0x73, 0x49, 0x6e, 0x66, 0x6f,
	0x73, 0x42, 0x05, 0x5a, 0x03, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_base_proto_rawDescOnce sync.Once
	file_base_proto_rawDescData = file_base_proto_rawDesc
)

func file_base_proto_rawDescGZIP() []byte {
	file_base_proto_rawDescOnce.Do(func() {
		file_base_proto_rawDescData = protoimpl.X.CompressGZIP(file_base_proto_rawDescData)
	})
	return file_base_proto_rawDescData
}

var file_base_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_base_proto_goTypes = []interface{}{
	(*S2CError)(nil),     // 0: pb.S2CError
	(*S2CResUpdate)(nil), // 1: pb.S2CResUpdate
	(*ResInfo)(nil),      // 2: pb.ResInfo
}
var file_base_proto_depIdxs = []int32{
	2, // 0: pb.S2CResUpdate.resInfos:type_name -> pb.ResInfo
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_base_proto_init() }
func file_base_proto_init() {
	if File_base_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_base_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*S2CError); i {
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
		file_base_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*S2CResUpdate); i {
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
			RawDescriptor: file_base_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_base_proto_goTypes,
		DependencyIndexes: file_base_proto_depIdxs,
		MessageInfos:      file_base_proto_msgTypes,
	}.Build()
	File_base_proto = out.File
	file_base_proto_rawDesc = nil
	file_base_proto_goTypes = nil
	file_base_proto_depIdxs = nil
}
