// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.20.3
// source: proto/human.proto

package pbs

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

type HumanRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AdminID int64 `protobuf:"varint,11,opt,name=AdminID,proto3" json:"AdminID,omitempty"` //管理员ID
}

func (x *HumanRequest) Reset() {
	*x = HumanRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_human_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HumanRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HumanRequest) ProtoMessage() {}

func (x *HumanRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_human_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HumanRequest.ProtoReflect.Descriptor instead.
func (*HumanRequest) Descriptor() ([]byte, []int) {
	return file_proto_human_proto_rawDescGZIP(), []int{0}
}

func (x *HumanRequest) GetAdminID() int64 {
	if x != nil {
		return x.AdminID
	}
	return 0
}

type HumanResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AdminID   int64  `protobuf:"varint,10,opt,name=AdminID,proto3" json:"AdminID,omitempty"`    //管理员ID
	AdminName string `protobuf:"bytes,11,opt,name=AdminName,proto3" json:"AdminName,omitempty"` //管理员名字
}

func (x *HumanResponse) Reset() {
	*x = HumanResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_human_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HumanResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HumanResponse) ProtoMessage() {}

func (x *HumanResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_human_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HumanResponse.ProtoReflect.Descriptor instead.
func (*HumanResponse) Descriptor() ([]byte, []int) {
	return file_proto_human_proto_rawDescGZIP(), []int{1}
}

func (x *HumanResponse) GetAdminID() int64 {
	if x != nil {
		return x.AdminID
	}
	return 0
}

func (x *HumanResponse) GetAdminName() string {
	if x != nil {
		return x.AdminName
	}
	return ""
}

var File_proto_human_proto protoreflect.FileDescriptor

var file_proto_human_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x68, 0x75, 0x6d, 0x61, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x03, 0x70, 0x62, 0x73, 0x22, 0x28, 0x0a, 0x0c, 0x48, 0x75, 0x6d, 0x61,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x64, 0x6d, 0x69,
	0x6e, 0x49, 0x44, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x41, 0x64, 0x6d, 0x69, 0x6e,
	0x49, 0x44, 0x22, 0x47, 0x0a, 0x0d, 0x48, 0x75, 0x6d, 0x61, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x49, 0x44, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x49, 0x44, 0x12, 0x1c, 0x0a,
	0x09, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x32, 0x3f, 0x0a, 0x0c, 0x48,
	0x75, 0x6d, 0x61, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2f, 0x0a, 0x04, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x11, 0x2e, 0x70, 0x62, 0x73, 0x2e, 0x48, 0x75, 0x6d, 0x61, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x70, 0x62, 0x73, 0x2e, 0x48, 0x75, 0x6d,
	0x61, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x07, 0x5a, 0x05,
	0x2e, 0x3b, 0x70, 0x62, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_human_proto_rawDescOnce sync.Once
	file_proto_human_proto_rawDescData = file_proto_human_proto_rawDesc
)

func file_proto_human_proto_rawDescGZIP() []byte {
	file_proto_human_proto_rawDescOnce.Do(func() {
		file_proto_human_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_human_proto_rawDescData)
	})
	return file_proto_human_proto_rawDescData
}

var file_proto_human_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_human_proto_goTypes = []interface{}{
	(*HumanRequest)(nil),  // 0: pbs.HumanRequest
	(*HumanResponse)(nil), // 1: pbs.HumanResponse
}
var file_proto_human_proto_depIdxs = []int32{
	0, // 0: pbs.HumanService.Info:input_type -> pbs.HumanRequest
	1, // 1: pbs.HumanService.Info:output_type -> pbs.HumanResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_human_proto_init() }
func file_proto_human_proto_init() {
	if File_proto_human_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_human_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HumanRequest); i {
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
		file_proto_human_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HumanResponse); i {
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
			RawDescriptor: file_proto_human_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_human_proto_goTypes,
		DependencyIndexes: file_proto_human_proto_depIdxs,
		MessageInfos:      file_proto_human_proto_msgTypes,
	}.Build()
	File_proto_human_proto = out.File
	file_proto_human_proto_rawDesc = nil
	file_proto_human_proto_goTypes = nil
	file_proto_human_proto_depIdxs = nil
}
