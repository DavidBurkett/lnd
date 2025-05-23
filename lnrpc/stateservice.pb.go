// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.12.4
// source: stateservice.proto

package lnrpc

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

type WalletState int32

const (
	// NON_EXISTING means that the wallet has not yet been initialized.
	WalletState_NON_EXISTING WalletState = 0
	// LOCKED means that the wallet is locked and requires a password to unlock.
	WalletState_LOCKED WalletState = 1
	// UNLOCKED means that the wallet was unlocked successfully, but RPC server
	// isn't ready.
	WalletState_UNLOCKED WalletState = 2
	// RPC_ACTIVE means that the lnd server is active but not fully ready for
	// calls.
	WalletState_RPC_ACTIVE WalletState = 3
	// SERVER_ACTIVE means that the lnd server is ready to accept calls.
	WalletState_SERVER_ACTIVE WalletState = 4
	// WAITING_TO_START means that node is waiting to become the leader in a
	// cluster and is not started yet.
	WalletState_WAITING_TO_START WalletState = 255
)

// Enum value maps for WalletState.
var (
	WalletState_name = map[int32]string{
		0:   "NON_EXISTING",
		1:   "LOCKED",
		2:   "UNLOCKED",
		3:   "RPC_ACTIVE",
		4:   "SERVER_ACTIVE",
		255: "WAITING_TO_START",
	}
	WalletState_value = map[string]int32{
		"NON_EXISTING":     0,
		"LOCKED":           1,
		"UNLOCKED":         2,
		"RPC_ACTIVE":       3,
		"SERVER_ACTIVE":    4,
		"WAITING_TO_START": 255,
	}
)

func (x WalletState) Enum() *WalletState {
	p := new(WalletState)
	*p = x
	return p
}

func (x WalletState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WalletState) Descriptor() protoreflect.EnumDescriptor {
	return file_stateservice_proto_enumTypes[0].Descriptor()
}

func (WalletState) Type() protoreflect.EnumType {
	return &file_stateservice_proto_enumTypes[0]
}

func (x WalletState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WalletState.Descriptor instead.
func (WalletState) EnumDescriptor() ([]byte, []int) {
	return file_stateservice_proto_rawDescGZIP(), []int{0}
}

type SubscribeStateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SubscribeStateRequest) Reset() {
	*x = SubscribeStateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stateservice_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeStateRequest) ProtoMessage() {}

func (x *SubscribeStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_stateservice_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeStateRequest.ProtoReflect.Descriptor instead.
func (*SubscribeStateRequest) Descriptor() ([]byte, []int) {
	return file_stateservice_proto_rawDescGZIP(), []int{0}
}

type SubscribeStateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State WalletState `protobuf:"varint,1,opt,name=state,proto3,enum=lnrpc.WalletState" json:"state,omitempty"`
}

func (x *SubscribeStateResponse) Reset() {
	*x = SubscribeStateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stateservice_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeStateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeStateResponse) ProtoMessage() {}

func (x *SubscribeStateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_stateservice_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeStateResponse.ProtoReflect.Descriptor instead.
func (*SubscribeStateResponse) Descriptor() ([]byte, []int) {
	return file_stateservice_proto_rawDescGZIP(), []int{1}
}

func (x *SubscribeStateResponse) GetState() WalletState {
	if x != nil {
		return x.State
	}
	return WalletState_NON_EXISTING
}

type GetStateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetStateRequest) Reset() {
	*x = GetStateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stateservice_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStateRequest) ProtoMessage() {}

func (x *GetStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_stateservice_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStateRequest.ProtoReflect.Descriptor instead.
func (*GetStateRequest) Descriptor() ([]byte, []int) {
	return file_stateservice_proto_rawDescGZIP(), []int{2}
}

type GetStateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State WalletState `protobuf:"varint,1,opt,name=state,proto3,enum=lnrpc.WalletState" json:"state,omitempty"`
}

func (x *GetStateResponse) Reset() {
	*x = GetStateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stateservice_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStateResponse) ProtoMessage() {}

func (x *GetStateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_stateservice_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStateResponse.ProtoReflect.Descriptor instead.
func (*GetStateResponse) Descriptor() ([]byte, []int) {
	return file_stateservice_proto_rawDescGZIP(), []int{3}
}

func (x *GetStateResponse) GetState() WalletState {
	if x != nil {
		return x.State
	}
	return WalletState_NON_EXISTING
}

var File_stateservice_proto protoreflect.FileDescriptor

var file_stateservice_proto_rawDesc = []byte{
	0x0a, 0x12, 0x73, 0x74, 0x61, 0x74, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6c, 0x6e, 0x72, 0x70, 0x63, 0x22, 0x17, 0x0a, 0x15, 0x53,
	0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x42, 0x0a, 0x16, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28,
	0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e,
	0x6c, 0x6e, 0x72, 0x70, 0x63, 0x2e, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x11, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3c, 0x0a, 0x10, 0x47,
	0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x28, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12,
	0x2e, 0x6c, 0x6e, 0x72, 0x70, 0x63, 0x2e, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2a, 0x73, 0x0a, 0x0b, 0x57, 0x61, 0x6c,
	0x6c, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x10, 0x0a, 0x0c, 0x4e, 0x4f, 0x4e, 0x5f,
	0x45, 0x58, 0x49, 0x53, 0x54, 0x49, 0x4e, 0x47, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x4c, 0x4f,
	0x43, 0x4b, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x55, 0x4e, 0x4c, 0x4f, 0x43, 0x4b,
	0x45, 0x44, 0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x52, 0x50, 0x43, 0x5f, 0x41, 0x43, 0x54, 0x49,
	0x56, 0x45, 0x10, 0x03, 0x12, 0x11, 0x0a, 0x0d, 0x53, 0x45, 0x52, 0x56, 0x45, 0x52, 0x5f, 0x41,
	0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x04, 0x12, 0x15, 0x0a, 0x10, 0x57, 0x41, 0x49, 0x54, 0x49,
	0x4e, 0x47, 0x5f, 0x54, 0x4f, 0x5f, 0x53, 0x54, 0x41, 0x52, 0x54, 0x10, 0xff, 0x01, 0x32, 0x95,
	0x01, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x4f, 0x0a, 0x0e, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x62, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x1c, 0x2e, 0x6c, 0x6e, 0x72,
	0x70, 0x63, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x6c, 0x6e, 0x72, 0x70, 0x63,
	0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x3b, 0x0a, 0x08, 0x47, 0x65, 0x74,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x16, 0x2e, 0x6c, 0x6e, 0x72, 0x70, 0x63, 0x2e, 0x47, 0x65,
	0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e,
	0x6c, 0x6e, 0x72, 0x70, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1f, 0x5a, 0x1d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x74, 0x63, 0x73, 0x75, 0x69, 0x74, 0x65, 0x2f, 0x6c, 0x6e,
	0x64, 0x2f, 0x6c, 0x6e, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_stateservice_proto_rawDescOnce sync.Once
	file_stateservice_proto_rawDescData = file_stateservice_proto_rawDesc
)

func file_stateservice_proto_rawDescGZIP() []byte {
	file_stateservice_proto_rawDescOnce.Do(func() {
		file_stateservice_proto_rawDescData = protoimpl.X.CompressGZIP(file_stateservice_proto_rawDescData)
	})
	return file_stateservice_proto_rawDescData
}

var file_stateservice_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_stateservice_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_stateservice_proto_goTypes = []interface{}{
	(WalletState)(0),               // 0: lnrpc.WalletState
	(*SubscribeStateRequest)(nil),  // 1: lnrpc.SubscribeStateRequest
	(*SubscribeStateResponse)(nil), // 2: lnrpc.SubscribeStateResponse
	(*GetStateRequest)(nil),        // 3: lnrpc.GetStateRequest
	(*GetStateResponse)(nil),       // 4: lnrpc.GetStateResponse
}
var file_stateservice_proto_depIdxs = []int32{
	0, // 0: lnrpc.SubscribeStateResponse.state:type_name -> lnrpc.WalletState
	0, // 1: lnrpc.GetStateResponse.state:type_name -> lnrpc.WalletState
	1, // 2: lnrpc.State.SubscribeState:input_type -> lnrpc.SubscribeStateRequest
	3, // 3: lnrpc.State.GetState:input_type -> lnrpc.GetStateRequest
	2, // 4: lnrpc.State.SubscribeState:output_type -> lnrpc.SubscribeStateResponse
	4, // 5: lnrpc.State.GetState:output_type -> lnrpc.GetStateResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_stateservice_proto_init() }
func file_stateservice_proto_init() {
	if File_stateservice_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_stateservice_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscribeStateRequest); i {
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
		file_stateservice_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscribeStateResponse); i {
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
		file_stateservice_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStateRequest); i {
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
		file_stateservice_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStateResponse); i {
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
			RawDescriptor: file_stateservice_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_stateservice_proto_goTypes,
		DependencyIndexes: file_stateservice_proto_depIdxs,
		EnumInfos:         file_stateservice_proto_enumTypes,
		MessageInfos:      file_stateservice_proto_msgTypes,
	}.Build()
	File_stateservice_proto = out.File
	file_stateservice_proto_rawDesc = nil
	file_stateservice_proto_goTypes = nil
	file_stateservice_proto_depIdxs = nil
}
