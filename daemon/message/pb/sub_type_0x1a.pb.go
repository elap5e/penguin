// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: daemon/message/pb/sub_type_0x1a.proto

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

// Message generated by proto-message-gen. DO NOT EDIT.
// source: tencent.p1298im.s2c.msgtype0x210.submsgtype0x1a.MsgBody
//
// SubMsgType0x1a is the message type for the SubMsgType0x1a.
type SubMsgType0X1A struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SubMsgType0X1A) Reset() {
	*x = SubMsgType0X1A{}
	if protoimpl.UnsafeEnabled {
		mi := &file_daemon_message_pb_sub_type_0x1a_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubMsgType0X1A) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubMsgType0X1A) ProtoMessage() {}

func (x *SubMsgType0X1A) ProtoReflect() protoreflect.Message {
	mi := &file_daemon_message_pb_sub_type_0x1a_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubMsgType0X1A.ProtoReflect.Descriptor instead.
func (*SubMsgType0X1A) Descriptor() ([]byte, []int) {
	return file_daemon_message_pb_sub_type_0x1a_proto_rawDescGZIP(), []int{0}
}

// MsgBody is the message type for the MsgBody.
type SubMsgType0X1A_MsgBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileKey      []byte `protobuf:"bytes,1,opt,name=file_key,json=fileKey,proto3" json:"file_key,omitempty"`                 // bytes_file_key
	FromUin      uint32 `protobuf:"varint,2,opt,name=from_uin,json=fromUin,proto3" json:"from_uin,omitempty"`                // uint32_from_uin
	ToUin        uint32 `protobuf:"varint,3,opt,name=to_uin,json=toUin,proto3" json:"to_uin,omitempty"`                      // uint32_to_uin
	Status       uint32 `protobuf:"varint,4,opt,name=status,proto3" json:"status,omitempty"`                                 // uint32_status
	Ttl          uint32 `protobuf:"varint,5,opt,name=ttl,proto3" json:"ttl,omitempty"`                                       // uint32_ttl
	Desc         string `protobuf:"bytes,6,opt,name=desc,proto3" json:"desc,omitempty"`                                      // string_desc
	Type         uint32 `protobuf:"varint,7,opt,name=type,proto3" json:"type,omitempty"`                                     // uint32_type
	CaptureTimes uint32 `protobuf:"varint,8,opt,name=capture_times,json=captureTimes,proto3" json:"capture_times,omitempty"` // uint32_capture_times
	FromUin_64   uint64 `protobuf:"varint,9,opt,name=from_uin_64,json=fromUin64,proto3" json:"from_uin_64,omitempty"`        // uint64_from_uin
	ToUin_64     uint64 `protobuf:"varint,10,opt,name=to_uin_64,json=toUin64,proto3" json:"to_uin_64,omitempty"`             // uint64_to_uin
}

func (x *SubMsgType0X1A_MsgBody) Reset() {
	*x = SubMsgType0X1A_MsgBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_daemon_message_pb_sub_type_0x1a_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubMsgType0X1A_MsgBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubMsgType0X1A_MsgBody) ProtoMessage() {}

func (x *SubMsgType0X1A_MsgBody) ProtoReflect() protoreflect.Message {
	mi := &file_daemon_message_pb_sub_type_0x1a_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubMsgType0X1A_MsgBody.ProtoReflect.Descriptor instead.
func (*SubMsgType0X1A_MsgBody) Descriptor() ([]byte, []int) {
	return file_daemon_message_pb_sub_type_0x1a_proto_rawDescGZIP(), []int{0, 0}
}

func (x *SubMsgType0X1A_MsgBody) GetFileKey() []byte {
	if x != nil {
		return x.FileKey
	}
	return nil
}

func (x *SubMsgType0X1A_MsgBody) GetFromUin() uint32 {
	if x != nil {
		return x.FromUin
	}
	return 0
}

func (x *SubMsgType0X1A_MsgBody) GetToUin() uint32 {
	if x != nil {
		return x.ToUin
	}
	return 0
}

func (x *SubMsgType0X1A_MsgBody) GetStatus() uint32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *SubMsgType0X1A_MsgBody) GetTtl() uint32 {
	if x != nil {
		return x.Ttl
	}
	return 0
}

func (x *SubMsgType0X1A_MsgBody) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

func (x *SubMsgType0X1A_MsgBody) GetType() uint32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *SubMsgType0X1A_MsgBody) GetCaptureTimes() uint32 {
	if x != nil {
		return x.CaptureTimes
	}
	return 0
}

func (x *SubMsgType0X1A_MsgBody) GetFromUin_64() uint64 {
	if x != nil {
		return x.FromUin_64
	}
	return 0
}

func (x *SubMsgType0X1A_MsgBody) GetToUin_64() uint64 {
	if x != nil {
		return x.ToUin_64
	}
	return 0
}

var File_daemon_message_pb_sub_type_0x1a_proto protoreflect.FileDescriptor

var file_daemon_message_pb_sub_type_0x1a_proto_rawDesc = []byte{
	0x0a, 0x25, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2f, 0x70, 0x62, 0x2f, 0x73, 0x75, 0x62, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x30, 0x78, 0x31,
	0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9c, 0x02, 0x0a, 0x0e, 0x53, 0x75, 0x62, 0x4d,
	0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x30, 0x78, 0x31, 0x61, 0x1a, 0x89, 0x02, 0x0a, 0x07, 0x4d,
	0x73, 0x67, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x4b, 0x65,
	0x79, 0x12, 0x19, 0x0a, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x75, 0x69, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x07, 0x66, 0x72, 0x6f, 0x6d, 0x55, 0x69, 0x6e, 0x12, 0x15, 0x0a, 0x06,
	0x74, 0x6f, 0x5f, 0x75, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x74, 0x6f,
	0x55, 0x69, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x74,
	0x74, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x74, 0x74, 0x6c, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x65, 0x73, 0x63, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x65, 0x73,
	0x63, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x63, 0x61,
	0x70, 0x74, 0x75, 0x72, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x12, 0x1e, 0x0a, 0x0b, 0x66, 0x72,
	0x6f, 0x6d, 0x5f, 0x75, 0x69, 0x6e, 0x5f, 0x36, 0x34, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x09, 0x66, 0x72, 0x6f, 0x6d, 0x55, 0x69, 0x6e, 0x36, 0x34, 0x12, 0x1a, 0x0a, 0x09, 0x74, 0x6f,
	0x5f, 0x75, 0x69, 0x6e, 0x5f, 0x36, 0x34, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x74,
	0x6f, 0x55, 0x69, 0x6e, 0x36, 0x34, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x6c, 0x61, 0x70, 0x35, 0x65, 0x2f, 0x70, 0x65, 0x6e, 0x67,
	0x75, 0x69, 0x6e, 0x2f, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_daemon_message_pb_sub_type_0x1a_proto_rawDescOnce sync.Once
	file_daemon_message_pb_sub_type_0x1a_proto_rawDescData = file_daemon_message_pb_sub_type_0x1a_proto_rawDesc
)

func file_daemon_message_pb_sub_type_0x1a_proto_rawDescGZIP() []byte {
	file_daemon_message_pb_sub_type_0x1a_proto_rawDescOnce.Do(func() {
		file_daemon_message_pb_sub_type_0x1a_proto_rawDescData = protoimpl.X.CompressGZIP(file_daemon_message_pb_sub_type_0x1a_proto_rawDescData)
	})
	return file_daemon_message_pb_sub_type_0x1a_proto_rawDescData
}

var file_daemon_message_pb_sub_type_0x1a_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_daemon_message_pb_sub_type_0x1a_proto_goTypes = []interface{}{
	(*SubMsgType0X1A)(nil),         // 0: SubMsgType0x1a
	(*SubMsgType0X1A_MsgBody)(nil), // 1: SubMsgType0x1a.MsgBody
}
var file_daemon_message_pb_sub_type_0x1a_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_daemon_message_pb_sub_type_0x1a_proto_init() }
func file_daemon_message_pb_sub_type_0x1a_proto_init() {
	if File_daemon_message_pb_sub_type_0x1a_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_daemon_message_pb_sub_type_0x1a_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubMsgType0X1A); i {
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
		file_daemon_message_pb_sub_type_0x1a_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubMsgType0X1A_MsgBody); i {
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
			RawDescriptor: file_daemon_message_pb_sub_type_0x1a_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_daemon_message_pb_sub_type_0x1a_proto_goTypes,
		DependencyIndexes: file_daemon_message_pb_sub_type_0x1a_proto_depIdxs,
		MessageInfos:      file_daemon_message_pb_sub_type_0x1a_proto_msgTypes,
	}.Build()
	File_daemon_message_pb_sub_type_0x1a_proto = out.File
	file_daemon_message_pb_sub_type_0x1a_proto_rawDesc = nil
	file_daemon_message_pb_sub_type_0x1a_proto_goTypes = nil
	file_daemon_message_pb_sub_type_0x1a_proto_depIdxs = nil
}
