// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: daemon/service/pb/online_push.proto

package pb

import (
	pb "github.com/elap5e/penguin/daemon/message/pb"
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

// OnlinePush is the message type for the msg_onlinepush.
type OnlinePush struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *OnlinePush) Reset() {
	*x = OnlinePush{}
	if protoimpl.UnsafeEnabled {
		mi := &file_daemon_service_pb_online_push_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnlinePush) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnlinePush) ProtoMessage() {}

func (x *OnlinePush) ProtoReflect() protoreflect.Message {
	mi := &file_daemon_service_pb_online_push_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnlinePush.ProtoReflect.Descriptor instead.
func (*OnlinePush) Descriptor() ([]byte, []int) {
	return file_daemon_service_pb_online_push_proto_rawDescGZIP(), []int{0}
}

// Message generated by proto-message-gen. DO NOT EDIT.
// source: msf.onlinepush.PbPushMsg
//
// PbPushMsg is the message type for the PbPushMsg.
type OnlinePush_PbPushMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg         *pb.MsgCommon_Msg `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`                                     // f425388msg
	Svrip       int32             `protobuf:"varint,2,opt,name=svrip,proto3" json:"svrip,omitempty"`                                // svrip
	PushToken   []byte            `protobuf:"bytes,3,opt,name=push_token,json=pushToken,proto3" json:"push_token,omitempty"`        // bytes_push_token
	PingFlag    uint32            `protobuf:"varint,4,opt,name=ping_flag,json=pingFlag,proto3" json:"ping_flag,omitempty"`          // ping_flag
	GeneralFlag uint32            `protobuf:"varint,9,opt,name=general_flag,json=generalFlag,proto3" json:"general_flag,omitempty"` // uint32_general_flag
	BindUin     uint64            `protobuf:"varint,10,opt,name=bind_uin,json=bindUin,proto3" json:"bind_uin,omitempty"`            // uint64_bind_uin
}

func (x *OnlinePush_PbPushMsg) Reset() {
	*x = OnlinePush_PbPushMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_daemon_service_pb_online_push_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnlinePush_PbPushMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnlinePush_PbPushMsg) ProtoMessage() {}

func (x *OnlinePush_PbPushMsg) ProtoReflect() protoreflect.Message {
	mi := &file_daemon_service_pb_online_push_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnlinePush_PbPushMsg.ProtoReflect.Descriptor instead.
func (*OnlinePush_PbPushMsg) Descriptor() ([]byte, []int) {
	return file_daemon_service_pb_online_push_proto_rawDescGZIP(), []int{0, 0}
}

func (x *OnlinePush_PbPushMsg) GetMsg() *pb.MsgCommon_Msg {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *OnlinePush_PbPushMsg) GetSvrip() int32 {
	if x != nil {
		return x.Svrip
	}
	return 0
}

func (x *OnlinePush_PbPushMsg) GetPushToken() []byte {
	if x != nil {
		return x.PushToken
	}
	return nil
}

func (x *OnlinePush_PbPushMsg) GetPingFlag() uint32 {
	if x != nil {
		return x.PingFlag
	}
	return 0
}

func (x *OnlinePush_PbPushMsg) GetGeneralFlag() uint32 {
	if x != nil {
		return x.GeneralFlag
	}
	return 0
}

func (x *OnlinePush_PbPushMsg) GetBindUin() uint64 {
	if x != nil {
		return x.BindUin
	}
	return 0
}

var File_daemon_service_pb_online_push_proto protoreflect.FileDescriptor

var file_daemon_service_pb_online_push_proto_rawDesc = []byte{
	0x0a, 0x23, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2f, 0x70, 0x62, 0x2f, 0x6f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x70, 0x75, 0x73, 0x68, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2f, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x70, 0x62, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcc, 0x01, 0x0a, 0x0a, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65,
	0x50, 0x75, 0x73, 0x68, 0x1a, 0xbd, 0x01, 0x0a, 0x09, 0x50, 0x62, 0x50, 0x75, 0x73, 0x68, 0x4d,
	0x73, 0x67, 0x12, 0x20, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0e, 0x2e, 0x4d, 0x73, 0x67, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4d, 0x73, 0x67, 0x52,
	0x03, 0x6d, 0x73, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x76, 0x72, 0x69, 0x70, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x73, 0x76, 0x72, 0x69, 0x70, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75,
	0x73, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09,
	0x70, 0x75, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x69, 0x6e,
	0x67, 0x5f, 0x66, 0x6c, 0x61, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x70, 0x69,
	0x6e, 0x67, 0x46, 0x6c, 0x61, 0x67, 0x12, 0x21, 0x0a, 0x0c, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61,
	0x6c, 0x5f, 0x66, 0x6c, 0x61, 0x67, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x67, 0x65,
	0x6e, 0x65, 0x72, 0x61, 0x6c, 0x46, 0x6c, 0x61, 0x67, 0x12, 0x19, 0x0a, 0x08, 0x62, 0x69, 0x6e,
	0x64, 0x5f, 0x75, 0x69, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x62, 0x69, 0x6e,
	0x64, 0x55, 0x69, 0x6e, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x65, 0x6c, 0x61, 0x70, 0x35, 0x65, 0x2f, 0x70, 0x65, 0x6e, 0x67, 0x75, 0x69,
	0x6e, 0x2f, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_daemon_service_pb_online_push_proto_rawDescOnce sync.Once
	file_daemon_service_pb_online_push_proto_rawDescData = file_daemon_service_pb_online_push_proto_rawDesc
)

func file_daemon_service_pb_online_push_proto_rawDescGZIP() []byte {
	file_daemon_service_pb_online_push_proto_rawDescOnce.Do(func() {
		file_daemon_service_pb_online_push_proto_rawDescData = protoimpl.X.CompressGZIP(file_daemon_service_pb_online_push_proto_rawDescData)
	})
	return file_daemon_service_pb_online_push_proto_rawDescData
}

var file_daemon_service_pb_online_push_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_daemon_service_pb_online_push_proto_goTypes = []interface{}{
	(*OnlinePush)(nil),           // 0: OnlinePush
	(*OnlinePush_PbPushMsg)(nil), // 1: OnlinePush.PbPushMsg
	(*pb.MsgCommon_Msg)(nil),     // 2: MsgCommon.Msg
}
var file_daemon_service_pb_online_push_proto_depIdxs = []int32{
	2, // 0: OnlinePush.PbPushMsg.msg:type_name -> MsgCommon.Msg
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_daemon_service_pb_online_push_proto_init() }
func file_daemon_service_pb_online_push_proto_init() {
	if File_daemon_service_pb_online_push_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_daemon_service_pb_online_push_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnlinePush); i {
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
		file_daemon_service_pb_online_push_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnlinePush_PbPushMsg); i {
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
			RawDescriptor: file_daemon_service_pb_online_push_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_daemon_service_pb_online_push_proto_goTypes,
		DependencyIndexes: file_daemon_service_pb_online_push_proto_depIdxs,
		MessageInfos:      file_daemon_service_pb_online_push_proto_msgTypes,
	}.Build()
	File_daemon_service_pb_online_push_proto = out.File
	file_daemon_service_pb_online_push_proto_rawDesc = nil
	file_daemon_service_pb_online_push_proto_goTypes = nil
	file_daemon_service_pb_online_push_proto_depIdxs = nil
}