// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: moderator.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ModeratorUpdated_ModeratorOperation int32

const (
	ModeratorUpdated_ADD    ModeratorUpdated_ModeratorOperation = 0
	ModeratorUpdated_DELETE ModeratorUpdated_ModeratorOperation = 1
)

// Enum value maps for ModeratorUpdated_ModeratorOperation.
var (
	ModeratorUpdated_ModeratorOperation_name = map[int32]string{
		0: "ADD",
		1: "DELETE",
	}
	ModeratorUpdated_ModeratorOperation_value = map[string]int32{
		"ADD":    0,
		"DELETE": 1,
	}
)

func (x ModeratorUpdated_ModeratorOperation) Enum() *ModeratorUpdated_ModeratorOperation {
	p := new(ModeratorUpdated_ModeratorOperation)
	*p = x
	return p
}

func (x ModeratorUpdated_ModeratorOperation) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ModeratorUpdated_ModeratorOperation) Descriptor() protoreflect.EnumDescriptor {
	return file_moderator_proto_enumTypes[0].Descriptor()
}

func (ModeratorUpdated_ModeratorOperation) Type() protoreflect.EnumType {
	return &file_moderator_proto_enumTypes[0]
}

func (x ModeratorUpdated_ModeratorOperation) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ModeratorUpdated_ModeratorOperation.Descriptor instead.
func (ModeratorUpdated_ModeratorOperation) EnumDescriptor() ([]byte, []int) {
	return file_moderator_proto_rawDescGZIP(), []int{1, 0}
}

type SSNGUpdated_SSNGOperation int32

const (
	SSNGUpdated_ADD    SSNGUpdated_SSNGOperation = 0
	SSNGUpdated_DELETE SSNGUpdated_SSNGOperation = 1
)

// Enum value maps for SSNGUpdated_SSNGOperation.
var (
	SSNGUpdated_SSNGOperation_name = map[int32]string{
		0: "ADD",
		1: "DELETE",
	}
	SSNGUpdated_SSNGOperation_value = map[string]int32{
		"ADD":    0,
		"DELETE": 1,
	}
)

func (x SSNGUpdated_SSNGOperation) Enum() *SSNGUpdated_SSNGOperation {
	p := new(SSNGUpdated_SSNGOperation)
	*p = x
	return p
}

func (x SSNGUpdated_SSNGOperation) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SSNGUpdated_SSNGOperation) Descriptor() protoreflect.EnumDescriptor {
	return file_moderator_proto_enumTypes[1].Descriptor()
}

func (SSNGUpdated_SSNGOperation) Type() protoreflect.EnumType {
	return &file_moderator_proto_enumTypes[1]
}

func (x SSNGUpdated_SSNGOperation) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SSNGUpdated_SSNGOperation.Descriptor instead.
func (SSNGUpdated_SSNGOperation) EnumDescriptor() ([]byte, []int) {
	return file_moderator_proto_rawDescGZIP(), []int{2, 0}
}

type SSNGUpdated_SSNGType int32

const (
	SSNGUpdated_USER    SSNGUpdated_SSNGType = 0
	SSNGUpdated_WORD    SSNGUpdated_SSNGType = 1
	SSNGUpdated_COMMAND SSNGUpdated_SSNGType = 2
)

// Enum value maps for SSNGUpdated_SSNGType.
var (
	SSNGUpdated_SSNGType_name = map[int32]string{
		0: "USER",
		1: "WORD",
		2: "COMMAND",
	}
	SSNGUpdated_SSNGType_value = map[string]int32{
		"USER":    0,
		"WORD":    1,
		"COMMAND": 2,
	}
)

func (x SSNGUpdated_SSNGType) Enum() *SSNGUpdated_SSNGType {
	p := new(SSNGUpdated_SSNGType)
	*p = x
	return p
}

func (x SSNGUpdated_SSNGType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SSNGUpdated_SSNGType) Descriptor() protoreflect.EnumDescriptor {
	return file_moderator_proto_enumTypes[2].Descriptor()
}

func (SSNGUpdated_SSNGType) Type() protoreflect.EnumType {
	return &file_moderator_proto_enumTypes[2]
}

func (x SSNGUpdated_SSNGType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SSNGUpdated_SSNGType.Descriptor instead.
func (SSNGUpdated_SSNGType) EnumDescriptor() ([]byte, []int) {
	return file_moderator_proto_rawDescGZIP(), []int{2, 1}
}

type ModerationAnnouncement_GuidelineItem int32

const (
	ModerationAnnouncement_UNKNOWN              ModerationAnnouncement_GuidelineItem = 0
	ModerationAnnouncement_SEXUAL               ModerationAnnouncement_GuidelineItem = 1
	ModerationAnnouncement_SPAM                 ModerationAnnouncement_GuidelineItem = 2
	ModerationAnnouncement_SLANDER              ModerationAnnouncement_GuidelineItem = 3
	ModerationAnnouncement_PERSONAL_INFORMATION ModerationAnnouncement_GuidelineItem = 4
)

// Enum value maps for ModerationAnnouncement_GuidelineItem.
var (
	ModerationAnnouncement_GuidelineItem_name = map[int32]string{
		0: "UNKNOWN",
		1: "SEXUAL",
		2: "SPAM",
		3: "SLANDER",
		4: "PERSONAL_INFORMATION",
	}
	ModerationAnnouncement_GuidelineItem_value = map[string]int32{
		"UNKNOWN":              0,
		"SEXUAL":               1,
		"SPAM":                 2,
		"SLANDER":              3,
		"PERSONAL_INFORMATION": 4,
	}
)

func (x ModerationAnnouncement_GuidelineItem) Enum() *ModerationAnnouncement_GuidelineItem {
	p := new(ModerationAnnouncement_GuidelineItem)
	*p = x
	return p
}

func (x ModerationAnnouncement_GuidelineItem) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ModerationAnnouncement_GuidelineItem) Descriptor() protoreflect.EnumDescriptor {
	return file_moderator_proto_enumTypes[3].Descriptor()
}

func (ModerationAnnouncement_GuidelineItem) Type() protoreflect.EnumType {
	return &file_moderator_proto_enumTypes[3]
}

func (x ModerationAnnouncement_GuidelineItem) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ModerationAnnouncement_GuidelineItem.Descriptor instead.
func (ModerationAnnouncement_GuidelineItem) EnumDescriptor() ([]byte, []int) {
	return file_moderator_proto_rawDescGZIP(), []int{3, 0}
}

type ModeratorUserInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   int64   `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Nickname *string `protobuf:"bytes,2,opt,name=nickname,proto3,oneof" json:"nickname,omitempty"`
	IconUrl  *string `protobuf:"bytes,3,opt,name=iconUrl,proto3,oneof" json:"iconUrl,omitempty"`
}

func (x *ModeratorUserInfo) Reset() {
	*x = ModeratorUserInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_moderator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ModeratorUserInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ModeratorUserInfo) ProtoMessage() {}

func (x *ModeratorUserInfo) ProtoReflect() protoreflect.Message {
	mi := &file_moderator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ModeratorUserInfo.ProtoReflect.Descriptor instead.
func (*ModeratorUserInfo) Descriptor() ([]byte, []int) {
	return file_moderator_proto_rawDescGZIP(), []int{0}
}

func (x *ModeratorUserInfo) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ModeratorUserInfo) GetNickname() string {
	if x != nil && x.Nickname != nil {
		return *x.Nickname
	}
	return ""
}

func (x *ModeratorUserInfo) GetIconUrl() string {
	if x != nil && x.IconUrl != nil {
		return *x.IconUrl
	}
	return ""
}

type ModeratorUpdated struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Operation ModeratorUpdated_ModeratorOperation `protobuf:"varint,1,opt,name=operation,proto3,enum=dwango.nicolive.chat.data.atoms.ModeratorUpdated_ModeratorOperation" json:"operation,omitempty"`
	Operator  *ModeratorUserInfo                  `protobuf:"bytes,2,opt,name=operator,proto3" json:"operator,omitempty"`
	UpdatedAt *timestamppb.Timestamp              `protobuf:"bytes,3,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
}

func (x *ModeratorUpdated) Reset() {
	*x = ModeratorUpdated{}
	if protoimpl.UnsafeEnabled {
		mi := &file_moderator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ModeratorUpdated) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ModeratorUpdated) ProtoMessage() {}

func (x *ModeratorUpdated) ProtoReflect() protoreflect.Message {
	mi := &file_moderator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ModeratorUpdated.ProtoReflect.Descriptor instead.
func (*ModeratorUpdated) Descriptor() ([]byte, []int) {
	return file_moderator_proto_rawDescGZIP(), []int{1}
}

func (x *ModeratorUpdated) GetOperation() ModeratorUpdated_ModeratorOperation {
	if x != nil {
		return x.Operation
	}
	return ModeratorUpdated_ADD
}

func (x *ModeratorUpdated) GetOperator() *ModeratorUserInfo {
	if x != nil {
		return x.Operator
	}
	return nil
}

func (x *ModeratorUpdated) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type SSNGUpdated struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Operation SSNGUpdated_SSNGOperation `protobuf:"varint,1,opt,name=operation,proto3,enum=dwango.nicolive.chat.data.atoms.SSNGUpdated_SSNGOperation" json:"operation,omitempty"`
	SsngId    int64                     `protobuf:"varint,2,opt,name=ssng_id,json=ssngId,proto3" json:"ssng_id,omitempty"`
	Operator  *ModeratorUserInfo        `protobuf:"bytes,3,opt,name=operator,proto3" json:"operator,omitempty"`
	Type      *SSNGUpdated_SSNGType     `protobuf:"varint,4,opt,name=type,proto3,enum=dwango.nicolive.chat.data.atoms.SSNGUpdated_SSNGType,oneof" json:"type,omitempty"`
	Source    *string                   `protobuf:"bytes,5,opt,name=source,proto3,oneof" json:"source,omitempty"`
	UpdatedAt *timestamppb.Timestamp    `protobuf:"bytes,6,opt,name=updatedAt,proto3,oneof" json:"updatedAt,omitempty"`
}

func (x *SSNGUpdated) Reset() {
	*x = SSNGUpdated{}
	if protoimpl.UnsafeEnabled {
		mi := &file_moderator_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SSNGUpdated) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SSNGUpdated) ProtoMessage() {}

func (x *SSNGUpdated) ProtoReflect() protoreflect.Message {
	mi := &file_moderator_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SSNGUpdated.ProtoReflect.Descriptor instead.
func (*SSNGUpdated) Descriptor() ([]byte, []int) {
	return file_moderator_proto_rawDescGZIP(), []int{2}
}

func (x *SSNGUpdated) GetOperation() SSNGUpdated_SSNGOperation {
	if x != nil {
		return x.Operation
	}
	return SSNGUpdated_ADD
}

func (x *SSNGUpdated) GetSsngId() int64 {
	if x != nil {
		return x.SsngId
	}
	return 0
}

func (x *SSNGUpdated) GetOperator() *ModeratorUserInfo {
	if x != nil {
		return x.Operator
	}
	return nil
}

func (x *SSNGUpdated) GetType() SSNGUpdated_SSNGType {
	if x != nil && x.Type != nil {
		return *x.Type
	}
	return SSNGUpdated_USER
}

func (x *SSNGUpdated) GetSource() string {
	if x != nil && x.Source != nil {
		return *x.Source
	}
	return ""
}

func (x *SSNGUpdated) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type ModerationAnnouncement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message        *string                                `protobuf:"bytes,1,opt,name=message,proto3,oneof" json:"message,omitempty"`
	GuidelineItems []ModerationAnnouncement_GuidelineItem `protobuf:"varint,2,rep,packed,name=guidelineItems,proto3,enum=dwango.nicolive.chat.data.atoms.ModerationAnnouncement_GuidelineItem" json:"guidelineItems,omitempty"`
	UpdatedAt      *timestamppb.Timestamp                 `protobuf:"bytes,3,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
}

func (x *ModerationAnnouncement) Reset() {
	*x = ModerationAnnouncement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_moderator_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ModerationAnnouncement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ModerationAnnouncement) ProtoMessage() {}

func (x *ModerationAnnouncement) ProtoReflect() protoreflect.Message {
	mi := &file_moderator_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ModerationAnnouncement.ProtoReflect.Descriptor instead.
func (*ModerationAnnouncement) Descriptor() ([]byte, []int) {
	return file_moderator_proto_rawDescGZIP(), []int{3}
}

func (x *ModerationAnnouncement) GetMessage() string {
	if x != nil && x.Message != nil {
		return *x.Message
	}
	return ""
}

func (x *ModerationAnnouncement) GetGuidelineItems() []ModerationAnnouncement_GuidelineItem {
	if x != nil {
		return x.GuidelineItems
	}
	return nil
}

func (x *ModerationAnnouncement) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

var File_moderator_proto protoreflect.FileDescriptor

var file_moderator_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x6d, 0x6f, 0x64, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x1f, 0x64, 0x77, 0x61, 0x6e, 0x67, 0x6f, 0x2e, 0x6e, 0x69, 0x63, 0x6f, 0x6c, 0x69,
	0x76, 0x65, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x61, 0x74, 0x6f,
	0x6d, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x85, 0x01, 0x0a, 0x11, 0x4d, 0x6f, 0x64, 0x65, 0x72, 0x61, 0x74, 0x6f,
	0x72, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x1f, 0x0a, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65,
	0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x07, 0x69, 0x63, 0x6f, 0x6e, 0x55, 0x72, 0x6c, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x07, 0x69, 0x63, 0x6f, 0x6e, 0x55, 0x72, 0x6c, 0x88,
	0x01, 0x01, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x42,
	0x0a, 0x0a, 0x08, 0x5f, 0x69, 0x63, 0x6f, 0x6e, 0x55, 0x72, 0x6c, 0x22, 0xab, 0x02, 0x0a, 0x10,
	0x4d, 0x6f, 0x64, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x12, 0x62, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x44, 0x2e, 0x64, 0x77, 0x61, 0x6e, 0x67, 0x6f, 0x2e, 0x6e, 0x69, 0x63,
	0x6f, 0x6c, 0x69, 0x76, 0x65, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e,
	0x61, 0x74, 0x6f, 0x6d, 0x73, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72,
	0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x4e, 0x0a, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x64, 0x77, 0x61, 0x6e, 0x67, 0x6f, 0x2e,
	0x6e, 0x69, 0x63, 0x6f, 0x6c, 0x69, 0x76, 0x65, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x64, 0x61,
	0x74, 0x61, 0x2e, 0x61, 0x74, 0x6f, 0x6d, 0x73, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x72, 0x61, 0x74,
	0x6f, 0x72, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x12, 0x38, 0x0a, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x29,
	0x0a, 0x12, 0x4d, 0x6f, 0x64, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x4f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x44, 0x44, 0x10, 0x00, 0x12, 0x0a, 0x0a,
	0x06, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x01, 0x22, 0xf1, 0x03, 0x0a, 0x0b, 0x53, 0x53,
	0x4e, 0x47, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x12, 0x58, 0x0a, 0x09, 0x6f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x3a, 0x2e, 0x64,
	0x77, 0x61, 0x6e, 0x67, 0x6f, 0x2e, 0x6e, 0x69, 0x63, 0x6f, 0x6c, 0x69, 0x76, 0x65, 0x2e, 0x63,
	0x68, 0x61, 0x74, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x61, 0x74, 0x6f, 0x6d, 0x73, 0x2e, 0x53,
	0x53, 0x4e, 0x47, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x53, 0x53, 0x4e, 0x47, 0x4f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x17, 0x0a, 0x07, 0x73, 0x73, 0x6e, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x73, 0x73, 0x6e, 0x67, 0x49, 0x64, 0x12, 0x4e, 0x0a, 0x08,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x32,
	0x2e, 0x64, 0x77, 0x61, 0x6e, 0x67, 0x6f, 0x2e, 0x6e, 0x69, 0x63, 0x6f, 0x6c, 0x69, 0x76, 0x65,
	0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x61, 0x74, 0x6f, 0x6d, 0x73,
	0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x4e, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x35, 0x2e, 0x64, 0x77, 0x61,
	0x6e, 0x67, 0x6f, 0x2e, 0x6e, 0x69, 0x63, 0x6f, 0x6c, 0x69, 0x76, 0x65, 0x2e, 0x63, 0x68, 0x61,
	0x74, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x61, 0x74, 0x6f, 0x6d, 0x73, 0x2e, 0x53, 0x53, 0x4e,
	0x47, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x53, 0x53, 0x4e, 0x47, 0x54, 0x79, 0x70,
	0x65, 0x48, 0x00, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x06,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x88, 0x01, 0x01, 0x12, 0x3d, 0x0a, 0x09, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x02, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x22, 0x24, 0x0a, 0x0d, 0x53, 0x53, 0x4e, 0x47,
	0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x44, 0x44,
	0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x01, 0x22, 0x2b,
	0x0a, 0x08, 0x53, 0x53, 0x4e, 0x47, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x55, 0x53,
	0x45, 0x52, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x57, 0x4f, 0x52, 0x44, 0x10, 0x01, 0x12, 0x0b,
	0x0a, 0x07, 0x43, 0x4f, 0x4d, 0x4d, 0x41, 0x4e, 0x44, 0x10, 0x02, 0x42, 0x07, 0x0a, 0x05, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x42,
	0x0c, 0x0a, 0x0a, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0xc7, 0x02,
	0x0a, 0x16, 0x4d, 0x6f, 0x64, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x6e, 0x6e, 0x6f,
	0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x88, 0x01, 0x01, 0x12, 0x6d, 0x0a, 0x0e, 0x67, 0x75, 0x69, 0x64, 0x65,
	0x6c, 0x69, 0x6e, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0e, 0x32,
	0x45, 0x2e, 0x64, 0x77, 0x61, 0x6e, 0x67, 0x6f, 0x2e, 0x6e, 0x69, 0x63, 0x6f, 0x6c, 0x69, 0x76,
	0x65, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x61, 0x74, 0x6f, 0x6d,
	0x73, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x6e, 0x6e, 0x6f,
	0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x47, 0x75, 0x69, 0x64, 0x65, 0x6c, 0x69,
	0x6e, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x0e, 0x67, 0x75, 0x69, 0x64, 0x65, 0x6c, 0x69, 0x6e,
	0x65, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x38, 0x0a, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x22, 0x59, 0x0a, 0x0d, 0x47, 0x75, 0x69, 0x64, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x49, 0x74, 0x65,
	0x6d, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0a,
	0x0a, 0x06, 0x53, 0x45, 0x58, 0x55, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x53, 0x50,
	0x41, 0x4d, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x4c, 0x41, 0x4e, 0x44, 0x45, 0x52, 0x10,
	0x03, 0x12, 0x18, 0x0a, 0x14, 0x50, 0x45, 0x52, 0x53, 0x4f, 0x4e, 0x41, 0x4c, 0x5f, 0x49, 0x4e,
	0x46, 0x4f, 0x52, 0x4d, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x04, 0x42, 0x0a, 0x0a, 0x08, 0x5f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_moderator_proto_rawDescOnce sync.Once
	file_moderator_proto_rawDescData = file_moderator_proto_rawDesc
)

func file_moderator_proto_rawDescGZIP() []byte {
	file_moderator_proto_rawDescOnce.Do(func() {
		file_moderator_proto_rawDescData = protoimpl.X.CompressGZIP(file_moderator_proto_rawDescData)
	})
	return file_moderator_proto_rawDescData
}

var file_moderator_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_moderator_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_moderator_proto_goTypes = []any{
	(ModeratorUpdated_ModeratorOperation)(0),  // 0: dwango.nicolive.chat.data.atoms.ModeratorUpdated.ModeratorOperation
	(SSNGUpdated_SSNGOperation)(0),            // 1: dwango.nicolive.chat.data.atoms.SSNGUpdated.SSNGOperation
	(SSNGUpdated_SSNGType)(0),                 // 2: dwango.nicolive.chat.data.atoms.SSNGUpdated.SSNGType
	(ModerationAnnouncement_GuidelineItem)(0), // 3: dwango.nicolive.chat.data.atoms.ModerationAnnouncement.GuidelineItem
	(*ModeratorUserInfo)(nil),                 // 4: dwango.nicolive.chat.data.atoms.ModeratorUserInfo
	(*ModeratorUpdated)(nil),                  // 5: dwango.nicolive.chat.data.atoms.ModeratorUpdated
	(*SSNGUpdated)(nil),                       // 6: dwango.nicolive.chat.data.atoms.SSNGUpdated
	(*ModerationAnnouncement)(nil),            // 7: dwango.nicolive.chat.data.atoms.ModerationAnnouncement
	(*timestamppb.Timestamp)(nil),             // 8: google.protobuf.Timestamp
}
var file_moderator_proto_depIdxs = []int32{
	0, // 0: dwango.nicolive.chat.data.atoms.ModeratorUpdated.operation:type_name -> dwango.nicolive.chat.data.atoms.ModeratorUpdated.ModeratorOperation
	4, // 1: dwango.nicolive.chat.data.atoms.ModeratorUpdated.operator:type_name -> dwango.nicolive.chat.data.atoms.ModeratorUserInfo
	8, // 2: dwango.nicolive.chat.data.atoms.ModeratorUpdated.updatedAt:type_name -> google.protobuf.Timestamp
	1, // 3: dwango.nicolive.chat.data.atoms.SSNGUpdated.operation:type_name -> dwango.nicolive.chat.data.atoms.SSNGUpdated.SSNGOperation
	4, // 4: dwango.nicolive.chat.data.atoms.SSNGUpdated.operator:type_name -> dwango.nicolive.chat.data.atoms.ModeratorUserInfo
	2, // 5: dwango.nicolive.chat.data.atoms.SSNGUpdated.type:type_name -> dwango.nicolive.chat.data.atoms.SSNGUpdated.SSNGType
	8, // 6: dwango.nicolive.chat.data.atoms.SSNGUpdated.updatedAt:type_name -> google.protobuf.Timestamp
	3, // 7: dwango.nicolive.chat.data.atoms.ModerationAnnouncement.guidelineItems:type_name -> dwango.nicolive.chat.data.atoms.ModerationAnnouncement.GuidelineItem
	8, // 8: dwango.nicolive.chat.data.atoms.ModerationAnnouncement.updatedAt:type_name -> google.protobuf.Timestamp
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_moderator_proto_init() }
func file_moderator_proto_init() {
	if File_moderator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_moderator_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ModeratorUserInfo); i {
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
		file_moderator_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ModeratorUpdated); i {
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
		file_moderator_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*SSNGUpdated); i {
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
		file_moderator_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*ModerationAnnouncement); i {
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
	file_moderator_proto_msgTypes[0].OneofWrappers = []any{}
	file_moderator_proto_msgTypes[2].OneofWrappers = []any{}
	file_moderator_proto_msgTypes[3].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_moderator_proto_rawDesc,
			NumEnums:      4,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_moderator_proto_goTypes,
		DependencyIndexes: file_moderator_proto_depIdxs,
		EnumInfos:         file_moderator_proto_enumTypes,
		MessageInfos:      file_moderator_proto_msgTypes,
	}.Build()
	File_moderator_proto = out.File
	file_moderator_proto_rawDesc = nil
	file_moderator_proto_goTypes = nil
	file_moderator_proto_depIdxs = nil
}