// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.7
// source: pb/dindin.proto

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

type PingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *PingResponse) Reset() {
	*x = PingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_dindin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingResponse) ProtoMessage() {}

func (x *PingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_dindin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.
func (*PingResponse) Descriptor() ([]byte, []int) {
	return file_pb_dindin_proto_rawDescGZIP(), []int{0}
}

func (x *PingResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type GetMembersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FirstName   string `protobuf:"bytes,1,opt,name=firstName,proto3" json:"firstName,omitempty"`
	LastName    string `protobuf:"bytes,2,opt,name=lastName,proto3" json:"lastName,omitempty"`
	RealName    string `protobuf:"bytes,3,opt,name=realName,proto3" json:"realName,omitempty"`
	DisplayName string `protobuf:"bytes,4,opt,name=displayName,proto3" json:"displayName,omitempty"`
	SlackUID    string `protobuf:"bytes,5,opt,name=slackUID,proto3" json:"slackUID,omitempty"`
}

func (x *GetMembersResponse) Reset() {
	*x = GetMembersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_dindin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMembersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMembersResponse) ProtoMessage() {}

func (x *GetMembersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_dindin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMembersResponse.ProtoReflect.Descriptor instead.
func (*GetMembersResponse) Descriptor() ([]byte, []int) {
	return file_pb_dindin_proto_rawDescGZIP(), []int{1}
}

func (x *GetMembersResponse) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *GetMembersResponse) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *GetMembersResponse) GetRealName() string {
	if x != nil {
		return x.RealName
	}
	return ""
}

func (x *GetMembersResponse) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

func (x *GetMembersResponse) GetSlackUID() string {
	if x != nil {
		return x.SlackUID
	}
	return ""
}

type CookingDay struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Day      int32  `protobuf:"varint,1,opt,name=day,proto3" json:"day,omitempty"`
	Month    int32  `protobuf:"varint,2,opt,name=month,proto3" json:"month,omitempty"`
	Year     int32  `protobuf:"varint,3,opt,name=year,proto3" json:"year,omitempty"`
	SlackUID string `protobuf:"bytes,4,opt,name=slackUID,proto3" json:"slackUID,omitempty"`
}

func (x *CookingDay) Reset() {
	*x = CookingDay{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_dindin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CookingDay) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CookingDay) ProtoMessage() {}

func (x *CookingDay) ProtoReflect() protoreflect.Message {
	mi := &file_pb_dindin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CookingDay.ProtoReflect.Descriptor instead.
func (*CookingDay) Descriptor() ([]byte, []int) {
	return file_pb_dindin_proto_rawDescGZIP(), []int{2}
}

func (x *CookingDay) GetDay() int32 {
	if x != nil {
		return x.Day
	}
	return 0
}

func (x *CookingDay) GetMonth() int32 {
	if x != nil {
		return x.Month
	}
	return 0
}

func (x *CookingDay) GetYear() int32 {
	if x != nil {
		return x.Year
	}
	return 0
}

func (x *CookingDay) GetSlackUID() string {
	if x != nil {
		return x.SlackUID
	}
	return ""
}

type AssignCooksRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CookingDays []*CookingDay `protobuf:"bytes,1,rep,name=CookingDays,proto3" json:"CookingDays,omitempty"`
}

func (x *AssignCooksRequest) Reset() {
	*x = AssignCooksRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_dindin_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssignCooksRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignCooksRequest) ProtoMessage() {}

func (x *AssignCooksRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_dindin_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignCooksRequest.ProtoReflect.Descriptor instead.
func (*AssignCooksRequest) Descriptor() ([]byte, []int) {
	return file_pb_dindin_proto_rawDescGZIP(), []int{3}
}

func (x *AssignCooksRequest) GetCookingDays() []*CookingDay {
	if x != nil {
		return x.CookingDays
	}
	return nil
}

type EmptyMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmptyMessage) Reset() {
	*x = EmptyMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_dindin_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmptyMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyMessage) ProtoMessage() {}

func (x *EmptyMessage) ProtoReflect() protoreflect.Message {
	mi := &file_pb_dindin_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyMessage.ProtoReflect.Descriptor instead.
func (*EmptyMessage) Descriptor() ([]byte, []int) {
	return file_pb_dindin_proto_rawDescGZIP(), []int{4}
}

var File_pb_dindin_proto protoreflect.FileDescriptor

var file_pb_dindin_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x70, 0x62, 0x2f, 0x64, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x02, 0x70, 0x62, 0x22, 0x28, 0x0a, 0x0c, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22,
	0xa8, 0x01, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x61, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x61, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b,
	0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x73, 0x6c, 0x61, 0x63, 0x6b, 0x55, 0x49, 0x44, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x73, 0x6c, 0x61, 0x63, 0x6b, 0x55, 0x49, 0x44, 0x22, 0x64, 0x0a, 0x0a, 0x43, 0x6f,
	0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x44, 0x61, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x61, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x64, 0x61, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x6f,
	0x6e, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6d, 0x6f, 0x6e, 0x74, 0x68,
	0x12, 0x12, 0x0a, 0x04, 0x79, 0x65, 0x61, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x79, 0x65, 0x61, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x6c, 0x61, 0x63, 0x6b, 0x55, 0x49, 0x44,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x6c, 0x61, 0x63, 0x6b, 0x55, 0x49, 0x44,
	0x22, 0x46, 0x0a, 0x12, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x43, 0x6f, 0x6f, 0x6b, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x30, 0x0a, 0x0b, 0x43, 0x6f, 0x6f, 0x6b, 0x69, 0x6e,
	0x67, 0x44, 0x61, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x62,
	0x2e, 0x43, 0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x67, 0x44, 0x61, 0x79, 0x52, 0x0b, 0x43, 0x6f, 0x6f,
	0x6b, 0x69, 0x6e, 0x67, 0x44, 0x61, 0x79, 0x73, 0x22, 0x0e, 0x0a, 0x0c, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0xa1, 0x02, 0x0a, 0x0c, 0x53, 0x6c, 0x61,
	0x63, 0x6b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x36, 0x0a, 0x0e, 0x45, 0x61, 0x74,
	0x69, 0x6e, 0x67, 0x54, 0x6f, 0x6d, 0x6f, 0x72, 0x72, 0x6f, 0x77, 0x12, 0x10, 0x2e, 0x70, 0x62,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x10, 0x2e,
	0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22,
	0x00, 0x12, 0x2c, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x10, 0x2e, 0x70, 0x62,
	0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x3a, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x12, 0x10, 0x2e,
	0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a,
	0x16, 0x2e, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x12, 0x34, 0x0a, 0x0c, 0x57,
	0x65, 0x65, 0x6b, 0x6c, 0x79, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x10, 0x2e, 0x70, 0x62,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x10, 0x2e,
	0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22,
	0x00, 0x12, 0x39, 0x0a, 0x0b, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x43, 0x6f, 0x6f, 0x6b, 0x73,
	0x12, 0x16, 0x2e, 0x70, 0x62, 0x2e, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x43, 0x6f, 0x6f, 0x6b,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x42, 0x35, 0x5a, 0x33,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x64, 0x72, 0x69, 0x74,
	0x7a, 0x65, 0x6e, 0x68, 0x6f, 0x66, 0x66, 0x2f, 0x64, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x2f, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x2f, 0x72, 0x70, 0x63,
	0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_dindin_proto_rawDescOnce sync.Once
	file_pb_dindin_proto_rawDescData = file_pb_dindin_proto_rawDesc
)

func file_pb_dindin_proto_rawDescGZIP() []byte {
	file_pb_dindin_proto_rawDescOnce.Do(func() {
		file_pb_dindin_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_dindin_proto_rawDescData)
	})
	return file_pb_dindin_proto_rawDescData
}

var file_pb_dindin_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pb_dindin_proto_goTypes = []interface{}{
	(*PingResponse)(nil),       // 0: pb.PingResponse
	(*GetMembersResponse)(nil), // 1: pb.GetMembersResponse
	(*CookingDay)(nil),         // 2: pb.CookingDay
	(*AssignCooksRequest)(nil), // 3: pb.AssignCooksRequest
	(*EmptyMessage)(nil),       // 4: pb.EmptyMessage
}
var file_pb_dindin_proto_depIdxs = []int32{
	2, // 0: pb.AssignCooksRequest.CookingDays:type_name -> pb.CookingDay
	4, // 1: pb.SlackActions.EatingTomorrow:input_type -> pb.EmptyMessage
	4, // 2: pb.SlackActions.Ping:input_type -> pb.EmptyMessage
	4, // 3: pb.SlackActions.GetMembers:input_type -> pb.EmptyMessage
	4, // 4: pb.SlackActions.WeeklyUpdate:input_type -> pb.EmptyMessage
	3, // 5: pb.SlackActions.AssignCooks:input_type -> pb.AssignCooksRequest
	4, // 6: pb.SlackActions.EatingTomorrow:output_type -> pb.EmptyMessage
	0, // 7: pb.SlackActions.Ping:output_type -> pb.PingResponse
	1, // 8: pb.SlackActions.GetMembers:output_type -> pb.GetMembersResponse
	4, // 9: pb.SlackActions.WeeklyUpdate:output_type -> pb.EmptyMessage
	4, // 10: pb.SlackActions.AssignCooks:output_type -> pb.EmptyMessage
	6, // [6:11] is the sub-list for method output_type
	1, // [1:6] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pb_dindin_proto_init() }
func file_pb_dindin_proto_init() {
	if File_pb_dindin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_dindin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingResponse); i {
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
		file_pb_dindin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMembersResponse); i {
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
		file_pb_dindin_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CookingDay); i {
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
		file_pb_dindin_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssignCooksRequest); i {
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
		file_pb_dindin_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmptyMessage); i {
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
			RawDescriptor: file_pb_dindin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_dindin_proto_goTypes,
		DependencyIndexes: file_pb_dindin_proto_depIdxs,
		MessageInfos:      file_pb_dindin_proto_msgTypes,
	}.Build()
	File_pb_dindin_proto = out.File
	file_pb_dindin_proto_rawDesc = nil
	file_pb_dindin_proto_goTypes = nil
	file_pb_dindin_proto_depIdxs = nil
}
