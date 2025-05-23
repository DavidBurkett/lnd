// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.12.4
// source: lnclipb/lncli.proto

package lnclipb

import (
	verrpc "github.com/ltcsuite/lnd/lnrpc/verrpc"
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

type VersionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The version information for lncli.
	Lncli *verrpc.Version `protobuf:"bytes,1,opt,name=lncli,proto3" json:"lncli,omitempty"`
	// The version information for lnd.
	Lnd *verrpc.Version `protobuf:"bytes,2,opt,name=lnd,proto3" json:"lnd,omitempty"`
}

func (x *VersionResponse) Reset() {
	*x = VersionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lnclipb_lncli_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionResponse) ProtoMessage() {}

func (x *VersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_lnclipb_lncli_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionResponse.ProtoReflect.Descriptor instead.
func (*VersionResponse) Descriptor() ([]byte, []int) {
	return file_lnclipb_lncli_proto_rawDescGZIP(), []int{0}
}

func (x *VersionResponse) GetLncli() *verrpc.Version {
	if x != nil {
		return x.Lncli
	}
	return nil
}

func (x *VersionResponse) GetLnd() *verrpc.Version {
	if x != nil {
		return x.Lnd
	}
	return nil
}

var File_lnclipb_lncli_proto protoreflect.FileDescriptor

var file_lnclipb_lncli_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6c, 0x6e, 0x63, 0x6c, 0x69, 0x70, 0x62, 0x2f, 0x6c, 0x6e, 0x63, 0x6c, 0x69, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6c, 0x6e, 0x63, 0x6c, 0x69, 0x70, 0x62, 0x1a, 0x13,
	0x76, 0x65, 0x72, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x65, 0x72, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x5b, 0x0a, 0x0f, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x05, 0x6c, 0x6e, 0x63, 0x6c, 0x69, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x76, 0x65, 0x72, 0x72, 0x70, 0x63, 0x2e, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x6c, 0x6e, 0x63, 0x6c, 0x69, 0x12, 0x21, 0x0a,
	0x03, 0x6c, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x76, 0x65, 0x72,
	0x72, 0x70, 0x63, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x03, 0x6c, 0x6e, 0x64,
	0x42, 0x27, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c,
	0x74, 0x63, 0x73, 0x75, 0x69, 0x74, 0x65, 0x2f, 0x6c, 0x6e, 0x64, 0x2f, 0x6c, 0x6e, 0x72, 0x70,
	0x63, 0x2f, 0x6c, 0x6e, 0x63, 0x6c, 0x69, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_lnclipb_lncli_proto_rawDescOnce sync.Once
	file_lnclipb_lncli_proto_rawDescData = file_lnclipb_lncli_proto_rawDesc
)

func file_lnclipb_lncli_proto_rawDescGZIP() []byte {
	file_lnclipb_lncli_proto_rawDescOnce.Do(func() {
		file_lnclipb_lncli_proto_rawDescData = protoimpl.X.CompressGZIP(file_lnclipb_lncli_proto_rawDescData)
	})
	return file_lnclipb_lncli_proto_rawDescData
}

var file_lnclipb_lncli_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_lnclipb_lncli_proto_goTypes = []interface{}{
	(*VersionResponse)(nil), // 0: lnclipb.VersionResponse
	(*verrpc.Version)(nil),  // 1: verrpc.Version
}
var file_lnclipb_lncli_proto_depIdxs = []int32{
	1, // 0: lnclipb.VersionResponse.lncli:type_name -> verrpc.Version
	1, // 1: lnclipb.VersionResponse.lnd:type_name -> verrpc.Version
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_lnclipb_lncli_proto_init() }
func file_lnclipb_lncli_proto_init() {
	if File_lnclipb_lncli_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_lnclipb_lncli_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VersionResponse); i {
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
			RawDescriptor: file_lnclipb_lncli_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_lnclipb_lncli_proto_goTypes,
		DependencyIndexes: file_lnclipb_lncli_proto_depIdxs,
		MessageInfos:      file_lnclipb_lncli_proto_msgTypes,
	}.Build()
	File_lnclipb_lncli_proto = out.File
	file_lnclipb_lncli_proto_rawDesc = nil
	file_lnclipb_lncli_proto_goTypes = nil
	file_lnclipb_lncli_proto_depIdxs = nil
}
