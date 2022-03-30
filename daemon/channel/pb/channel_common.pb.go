// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: daemon/channel/pb/channel_common.proto

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
// source: guild.GuildCommon
//
// GuildCommon is the message type for the GuildCommon.
type GuildCommon struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GuildCommon) Reset() {
	*x = GuildCommon{}
	if protoimpl.UnsafeEnabled {
		mi := &file_daemon_channel_pb_channel_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuildCommon) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuildCommon) ProtoMessage() {}

func (x *GuildCommon) ProtoReflect() protoreflect.Message {
	mi := &file_daemon_channel_pb_channel_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuildCommon.ProtoReflect.Descriptor instead.
func (*GuildCommon) Descriptor() ([]byte, []int) {
	return file_daemon_channel_pb_channel_common_proto_rawDescGZIP(), []int{0}
}

// BytesEntry is the message type for the BytesEntry.
type GuildCommon_BytesEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`     // key
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"` // value
}

func (x *GuildCommon_BytesEntry) Reset() {
	*x = GuildCommon_BytesEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_daemon_channel_pb_channel_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuildCommon_BytesEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuildCommon_BytesEntry) ProtoMessage() {}

func (x *GuildCommon_BytesEntry) ProtoReflect() protoreflect.Message {
	mi := &file_daemon_channel_pb_channel_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuildCommon_BytesEntry.ProtoReflect.Descriptor instead.
func (*GuildCommon_BytesEntry) Descriptor() ([]byte, []int) {
	return file_daemon_channel_pb_channel_common_proto_rawDescGZIP(), []int{0, 0}
}

func (x *GuildCommon_BytesEntry) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *GuildCommon_BytesEntry) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

// Entry is the message type for the Entry.
type GuildCommon_Entry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`     // key
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"` // value
}

func (x *GuildCommon_Entry) Reset() {
	*x = GuildCommon_Entry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_daemon_channel_pb_channel_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuildCommon_Entry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuildCommon_Entry) ProtoMessage() {}

func (x *GuildCommon_Entry) ProtoReflect() protoreflect.Message {
	mi := &file_daemon_channel_pb_channel_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuildCommon_Entry.ProtoReflect.Descriptor instead.
func (*GuildCommon_Entry) Descriptor() ([]byte, []int) {
	return file_daemon_channel_pb_channel_common_proto_rawDescGZIP(), []int{0, 1}
}

func (x *GuildCommon_Entry) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *GuildCommon_Entry) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

// Result is the message type for the Result.
type GuildCommon_Result struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RetCode int32 `protobuf:"varint,1,opt,name=retCode,proto3" json:"retCode,omitempty"` // retCode
}

func (x *GuildCommon_Result) Reset() {
	*x = GuildCommon_Result{}
	if protoimpl.UnsafeEnabled {
		mi := &file_daemon_channel_pb_channel_common_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuildCommon_Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuildCommon_Result) ProtoMessage() {}

func (x *GuildCommon_Result) ProtoReflect() protoreflect.Message {
	mi := &file_daemon_channel_pb_channel_common_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuildCommon_Result.ProtoReflect.Descriptor instead.
func (*GuildCommon_Result) Descriptor() ([]byte, []int) {
	return file_daemon_channel_pb_channel_common_proto_rawDescGZIP(), []int{0, 2}
}

func (x *GuildCommon_Result) GetRetCode() int32 {
	if x != nil {
		return x.RetCode
	}
	return 0
}

// StCommonExt is the message type for the StCommonExt.
type GuildCommon_StCommonExt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MapInfo      []*GuildCommon_Entry      `protobuf:"bytes,1,rep,name=mapInfo,proto3" json:"mapInfo,omitempty"`           // mapInfo
	AttachInfo   string                    `protobuf:"bytes,2,opt,name=attachInfo,proto3" json:"attachInfo,omitempty"`     // attachInfo
	MapBytesInfo []*GuildCommon_BytesEntry `protobuf:"bytes,3,rep,name=mapBytesInfo,proto3" json:"mapBytesInfo,omitempty"` // mapBytesInfo
}

func (x *GuildCommon_StCommonExt) Reset() {
	*x = GuildCommon_StCommonExt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_daemon_channel_pb_channel_common_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuildCommon_StCommonExt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuildCommon_StCommonExt) ProtoMessage() {}

func (x *GuildCommon_StCommonExt) ProtoReflect() protoreflect.Message {
	mi := &file_daemon_channel_pb_channel_common_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuildCommon_StCommonExt.ProtoReflect.Descriptor instead.
func (*GuildCommon_StCommonExt) Descriptor() ([]byte, []int) {
	return file_daemon_channel_pb_channel_common_proto_rawDescGZIP(), []int{0, 3}
}

func (x *GuildCommon_StCommonExt) GetMapInfo() []*GuildCommon_Entry {
	if x != nil {
		return x.MapInfo
	}
	return nil
}

func (x *GuildCommon_StCommonExt) GetAttachInfo() string {
	if x != nil {
		return x.AttachInfo
	}
	return ""
}

func (x *GuildCommon_StCommonExt) GetMapBytesInfo() []*GuildCommon_BytesEntry {
	if x != nil {
		return x.MapBytesInfo
	}
	return nil
}

// SyncMessage is the message type for the SyncMessage.
type GuildCommon_SyncMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data        []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`                // data
	TimestampMS int64  `protobuf:"varint,2,opt,name=timestampMS,proto3" json:"timestampMS,omitempty"` // timestampMS
}

func (x *GuildCommon_SyncMessage) Reset() {
	*x = GuildCommon_SyncMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_daemon_channel_pb_channel_common_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuildCommon_SyncMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuildCommon_SyncMessage) ProtoMessage() {}

func (x *GuildCommon_SyncMessage) ProtoReflect() protoreflect.Message {
	mi := &file_daemon_channel_pb_channel_common_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuildCommon_SyncMessage.ProtoReflect.Descriptor instead.
func (*GuildCommon_SyncMessage) Descriptor() ([]byte, []int) {
	return file_daemon_channel_pb_channel_common_proto_rawDescGZIP(), []int{0, 4}
}

func (x *GuildCommon_SyncMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *GuildCommon_SyncMessage) GetTimestampMS() int64 {
	if x != nil {
		return x.TimestampMS
	}
	return 0
}

var File_daemon_channel_pb_channel_common_proto protoreflect.FileDescriptor

var file_daemon_channel_pb_channel_common_proto_rawDesc = []byte{
	0x0a, 0x26, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x2f, 0x70, 0x62, 0x2f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf8, 0x02, 0x0a, 0x0b, 0x47, 0x75, 0x69,
	0x6c, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x1a, 0x34, 0x0a, 0x0a, 0x42, 0x79, 0x74, 0x65,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x2f,
	0x0a, 0x05, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x1a,
	0x22, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x74,
	0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x72, 0x65, 0x74, 0x43,
	0x6f, 0x64, 0x65, 0x1a, 0x98, 0x01, 0x0a, 0x0b, 0x53, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x45, 0x78, 0x74, 0x12, 0x2c, 0x0a, 0x07, 0x6d, 0x61, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x43, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x6d, 0x61, 0x70, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x1e, 0x0a, 0x0a, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x49, 0x6e, 0x66, 0x6f, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x3b, 0x0a, 0x0c, 0x6d, 0x61, 0x70, 0x42, 0x79, 0x74, 0x65, 0x73, 0x49, 0x6e, 0x66,
	0x6f, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x43,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x42, 0x79, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x0c, 0x6d, 0x61, 0x70, 0x42, 0x79, 0x74, 0x65, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x43,
	0x0a, 0x0b, 0x53, 0x79, 0x6e, 0x63, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x20, 0x0a, 0x0b, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x4d, 0x53,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x4d, 0x53, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x65, 0x6c, 0x61, 0x70, 0x35, 0x65, 0x2f, 0x70, 0x65, 0x6e, 0x67, 0x75, 0x69, 0x6e,
	0x2f, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x2f,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_daemon_channel_pb_channel_common_proto_rawDescOnce sync.Once
	file_daemon_channel_pb_channel_common_proto_rawDescData = file_daemon_channel_pb_channel_common_proto_rawDesc
)

func file_daemon_channel_pb_channel_common_proto_rawDescGZIP() []byte {
	file_daemon_channel_pb_channel_common_proto_rawDescOnce.Do(func() {
		file_daemon_channel_pb_channel_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_daemon_channel_pb_channel_common_proto_rawDescData)
	})
	return file_daemon_channel_pb_channel_common_proto_rawDescData
}

var file_daemon_channel_pb_channel_common_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_daemon_channel_pb_channel_common_proto_goTypes = []interface{}{
	(*GuildCommon)(nil),             // 0: GuildCommon
	(*GuildCommon_BytesEntry)(nil),  // 1: GuildCommon.BytesEntry
	(*GuildCommon_Entry)(nil),       // 2: GuildCommon.Entry
	(*GuildCommon_Result)(nil),      // 3: GuildCommon.Result
	(*GuildCommon_StCommonExt)(nil), // 4: GuildCommon.StCommonExt
	(*GuildCommon_SyncMessage)(nil), // 5: GuildCommon.SyncMessage
}
var file_daemon_channel_pb_channel_common_proto_depIdxs = []int32{
	2, // 0: GuildCommon.StCommonExt.mapInfo:type_name -> GuildCommon.Entry
	1, // 1: GuildCommon.StCommonExt.mapBytesInfo:type_name -> GuildCommon.BytesEntry
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_daemon_channel_pb_channel_common_proto_init() }
func file_daemon_channel_pb_channel_common_proto_init() {
	if File_daemon_channel_pb_channel_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_daemon_channel_pb_channel_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GuildCommon); i {
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
		file_daemon_channel_pb_channel_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GuildCommon_BytesEntry); i {
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
		file_daemon_channel_pb_channel_common_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GuildCommon_Entry); i {
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
		file_daemon_channel_pb_channel_common_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GuildCommon_Result); i {
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
		file_daemon_channel_pb_channel_common_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GuildCommon_StCommonExt); i {
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
		file_daemon_channel_pb_channel_common_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GuildCommon_SyncMessage); i {
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
			RawDescriptor: file_daemon_channel_pb_channel_common_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_daemon_channel_pb_channel_common_proto_goTypes,
		DependencyIndexes: file_daemon_channel_pb_channel_common_proto_depIdxs,
		MessageInfos:      file_daemon_channel_pb_channel_common_proto_msgTypes,
	}.Build()
	File_daemon_channel_pb_channel_common_proto = out.File
	file_daemon_channel_pb_channel_common_proto_rawDesc = nil
	file_daemon_channel_pb_channel_common_proto_goTypes = nil
	file_daemon_channel_pb_channel_common_proto_depIdxs = nil
}