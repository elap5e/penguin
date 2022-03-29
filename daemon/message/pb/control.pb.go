// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: daemon/message/pb/control.proto

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
// source: msf.msgsvc.msg_ctrl
//
// msg_ctrl is the message type for the msg_ctrl.
type MsgControl struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MsgControl) Reset() {
	*x = MsgControl{}
	if protoimpl.UnsafeEnabled {
		mi := &file_daemon_message_pb_control_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgControl) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgControl) ProtoMessage() {}

func (x *MsgControl) ProtoReflect() protoreflect.Message {
	mi := &file_daemon_message_pb_control_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgControl.ProtoReflect.Descriptor instead.
func (*MsgControl) Descriptor() ([]byte, []int) {
	return file_daemon_message_pb_control_proto_rawDescGZIP(), []int{0}
}

// MsgCtrl is the message type for the MsgCtrl.
type MsgControl_MsgCtrl struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MsgFlag      uint32                   `protobuf:"varint,1,opt,name=msg_flag,json=msgFlag,proto3" json:"msg_flag,omitempty"`                 // msg_flag
	ResvResvInfo *MsgControl_ResvResvInfo `protobuf:"bytes,2,opt,name=resv_resv_info,json=resvResvInfo,proto3" json:"resv_resv_info,omitempty"` // resv_resv_info
}

func (x *MsgControl_MsgCtrl) Reset() {
	*x = MsgControl_MsgCtrl{}
	if protoimpl.UnsafeEnabled {
		mi := &file_daemon_message_pb_control_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgControl_MsgCtrl) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgControl_MsgCtrl) ProtoMessage() {}

func (x *MsgControl_MsgCtrl) ProtoReflect() protoreflect.Message {
	mi := &file_daemon_message_pb_control_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgControl_MsgCtrl.ProtoReflect.Descriptor instead.
func (*MsgControl_MsgCtrl) Descriptor() ([]byte, []int) {
	return file_daemon_message_pb_control_proto_rawDescGZIP(), []int{0, 0}
}

func (x *MsgControl_MsgCtrl) GetMsgFlag() uint32 {
	if x != nil {
		return x.MsgFlag
	}
	return 0
}

func (x *MsgControl_MsgCtrl) GetResvResvInfo() *MsgControl_ResvResvInfo {
	if x != nil {
		return x.ResvResvInfo
	}
	return nil
}

// ResvResvInfo is the message type for the ResvResvInfo.
type MsgControl_ResvResvInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Flag       uint32 `protobuf:"varint,1,opt,name=flag,proto3" json:"flag,omitempty"`                               // uint32_flag
	Reserv1    []byte `protobuf:"bytes,2,opt,name=reserv1,proto3" json:"reserv1,omitempty"`                          // bytes_reserv1
	Reserv2    uint64 `protobuf:"varint,3,opt,name=reserv2,proto3" json:"reserv2,omitempty"`                         // uint64_reserv2
	Reserv3    uint64 `protobuf:"varint,4,opt,name=reserv3,proto3" json:"reserv3,omitempty"`                         // uint64_reserv3
	CreateTime uint32 `protobuf:"varint,5,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"` // uint32_create_time
	PicHeight  uint32 `protobuf:"varint,6,opt,name=pic_height,json=picHeight,proto3" json:"pic_height,omitempty"`    // uint32_pic_height
	PicWidth   uint32 `protobuf:"varint,7,opt,name=pic_width,json=picWidth,proto3" json:"pic_width,omitempty"`       // uint32_pic_width
	ResvFlag   uint32 `protobuf:"varint,8,opt,name=resv_flag,json=resvFlag,proto3" json:"resv_flag,omitempty"`       // uint32_resv_flag
}

func (x *MsgControl_ResvResvInfo) Reset() {
	*x = MsgControl_ResvResvInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_daemon_message_pb_control_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgControl_ResvResvInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgControl_ResvResvInfo) ProtoMessage() {}

func (x *MsgControl_ResvResvInfo) ProtoReflect() protoreflect.Message {
	mi := &file_daemon_message_pb_control_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgControl_ResvResvInfo.ProtoReflect.Descriptor instead.
func (*MsgControl_ResvResvInfo) Descriptor() ([]byte, []int) {
	return file_daemon_message_pb_control_proto_rawDescGZIP(), []int{0, 1}
}

func (x *MsgControl_ResvResvInfo) GetFlag() uint32 {
	if x != nil {
		return x.Flag
	}
	return 0
}

func (x *MsgControl_ResvResvInfo) GetReserv1() []byte {
	if x != nil {
		return x.Reserv1
	}
	return nil
}

func (x *MsgControl_ResvResvInfo) GetReserv2() uint64 {
	if x != nil {
		return x.Reserv2
	}
	return 0
}

func (x *MsgControl_ResvResvInfo) GetReserv3() uint64 {
	if x != nil {
		return x.Reserv3
	}
	return 0
}

func (x *MsgControl_ResvResvInfo) GetCreateTime() uint32 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

func (x *MsgControl_ResvResvInfo) GetPicHeight() uint32 {
	if x != nil {
		return x.PicHeight
	}
	return 0
}

func (x *MsgControl_ResvResvInfo) GetPicWidth() uint32 {
	if x != nil {
		return x.PicWidth
	}
	return 0
}

func (x *MsgControl_ResvResvInfo) GetResvFlag() uint32 {
	if x != nil {
		return x.ResvFlag
	}
	return 0
}

var File_daemon_message_pb_control_proto protoreflect.FileDescriptor

var file_daemon_message_pb_control_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2f, 0x70, 0x62, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xdf, 0x02, 0x0a, 0x0a, 0x4d, 0x73, 0x67, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x1a, 0x64, 0x0a, 0x07, 0x4d, 0x73, 0x67, 0x43, 0x74, 0x72, 0x6c, 0x12, 0x19, 0x0a, 0x08, 0x6d,
	0x73, 0x67, 0x5f, 0x66, 0x6c, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6d,
	0x73, 0x67, 0x46, 0x6c, 0x61, 0x67, 0x12, 0x3e, 0x0a, 0x0e, 0x72, 0x65, 0x73, 0x76, 0x5f, 0x72,
	0x65, 0x73, 0x76, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18,
	0x2e, 0x4d, 0x73, 0x67, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x52, 0x65, 0x73, 0x76,
	0x52, 0x65, 0x73, 0x76, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x76, 0x52, 0x65,
	0x73, 0x76, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0xea, 0x01, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x76, 0x52,
	0x65, 0x73, 0x76, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x6c, 0x61, 0x67, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x66, 0x6c, 0x61, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x72,
	0x65, 0x73, 0x65, 0x72, 0x76, 0x31, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x72, 0x65,
	0x73, 0x65, 0x72, 0x76, 0x31, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x32,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x32, 0x12,
	0x18, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x33, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x07, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x33, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x69,
	0x63, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09,
	0x70, 0x69, 0x63, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x69, 0x63,
	0x5f, 0x77, 0x69, 0x64, 0x74, 0x68, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x70, 0x69,
	0x63, 0x57, 0x69, 0x64, 0x74, 0x68, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x73, 0x76, 0x5f, 0x66,
	0x6c, 0x61, 0x67, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x72, 0x65, 0x73, 0x76, 0x46,
	0x6c, 0x61, 0x67, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x65, 0x6c, 0x61, 0x70, 0x35, 0x65, 0x2f, 0x70, 0x65, 0x6e, 0x67, 0x75, 0x69, 0x6e,
	0x2f, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2f,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_daemon_message_pb_control_proto_rawDescOnce sync.Once
	file_daemon_message_pb_control_proto_rawDescData = file_daemon_message_pb_control_proto_rawDesc
)

func file_daemon_message_pb_control_proto_rawDescGZIP() []byte {
	file_daemon_message_pb_control_proto_rawDescOnce.Do(func() {
		file_daemon_message_pb_control_proto_rawDescData = protoimpl.X.CompressGZIP(file_daemon_message_pb_control_proto_rawDescData)
	})
	return file_daemon_message_pb_control_proto_rawDescData
}

var file_daemon_message_pb_control_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_daemon_message_pb_control_proto_goTypes = []interface{}{
	(*MsgControl)(nil),              // 0: MsgControl
	(*MsgControl_MsgCtrl)(nil),      // 1: MsgControl.MsgCtrl
	(*MsgControl_ResvResvInfo)(nil), // 2: MsgControl.ResvResvInfo
}
var file_daemon_message_pb_control_proto_depIdxs = []int32{
	2, // 0: MsgControl.MsgCtrl.resv_resv_info:type_name -> MsgControl.ResvResvInfo
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_daemon_message_pb_control_proto_init() }
func file_daemon_message_pb_control_proto_init() {
	if File_daemon_message_pb_control_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_daemon_message_pb_control_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgControl); i {
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
		file_daemon_message_pb_control_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgControl_MsgCtrl); i {
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
		file_daemon_message_pb_control_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgControl_ResvResvInfo); i {
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
			RawDescriptor: file_daemon_message_pb_control_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_daemon_message_pb_control_proto_goTypes,
		DependencyIndexes: file_daemon_message_pb_control_proto_depIdxs,
		MessageInfos:      file_daemon_message_pb_control_proto_msgTypes,
	}.Build()
	File_daemon_message_pb_control_proto = out.File
	file_daemon_message_pb_control_proto_rawDesc = nil
	file_daemon_message_pb_control_proto_goTypes = nil
	file_daemon_message_pb_control_proto_depIdxs = nil
}