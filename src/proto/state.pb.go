// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: state.proto

package proto

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

type NicoliveState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Statistics             *Statistics             `protobuf:"bytes,1,opt,name=statistics,proto3,oneof" json:"statistics,omitempty"`
	Enquete                *Enquete                `protobuf:"bytes,2,opt,name=enquete,proto3,oneof" json:"enquete,omitempty"`
	MoveOrder              *MoveOrder              `protobuf:"bytes,3,opt,name=move_order,json=moveOrder,proto3,oneof" json:"move_order,omitempty"`
	Marquee                *Marquee                `protobuf:"bytes,4,opt,name=marquee,proto3,oneof" json:"marquee,omitempty"`
	CommentLock            *CommentLock            `protobuf:"bytes,5,opt,name=comment_lock,json=commentLock,proto3,oneof" json:"comment_lock,omitempty"`
	CommentMode            *CommentMode            `protobuf:"bytes,6,opt,name=comment_mode,json=commentMode,proto3,oneof" json:"comment_mode,omitempty"`
	TrialPanel             *TrialPanel             `protobuf:"bytes,7,opt,name=trial_panel,json=trialPanel,proto3,oneof" json:"trial_panel,omitempty"`
	ProgramStatus          *ProgramStatus          `protobuf:"bytes,9,opt,name=program_status,json=programStatus,proto3,oneof" json:"program_status,omitempty"`
	ModerationAnnouncement *ModerationAnnouncement `protobuf:"bytes,10,opt,name=moderation_announcement,json=moderationAnnouncement,proto3,oneof" json:"moderation_announcement,omitempty"`
}

func (x *NicoliveState) Reset() {
	*x = NicoliveState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_state_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NicoliveState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NicoliveState) ProtoMessage() {}

func (x *NicoliveState) ProtoReflect() protoreflect.Message {
	mi := &file_state_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NicoliveState.ProtoReflect.Descriptor instead.
func (*NicoliveState) Descriptor() ([]byte, []int) {
	return file_state_proto_rawDescGZIP(), []int{0}
}

func (x *NicoliveState) GetStatistics() *Statistics {
	if x != nil {
		return x.Statistics
	}
	return nil
}

func (x *NicoliveState) GetEnquete() *Enquete {
	if x != nil {
		return x.Enquete
	}
	return nil
}

func (x *NicoliveState) GetMoveOrder() *MoveOrder {
	if x != nil {
		return x.MoveOrder
	}
	return nil
}

func (x *NicoliveState) GetMarquee() *Marquee {
	if x != nil {
		return x.Marquee
	}
	return nil
}

func (x *NicoliveState) GetCommentLock() *CommentLock {
	if x != nil {
		return x.CommentLock
	}
	return nil
}

func (x *NicoliveState) GetCommentMode() *CommentMode {
	if x != nil {
		return x.CommentMode
	}
	return nil
}

func (x *NicoliveState) GetTrialPanel() *TrialPanel {
	if x != nil {
		return x.TrialPanel
	}
	return nil
}

func (x *NicoliveState) GetProgramStatus() *ProgramStatus {
	if x != nil {
		return x.ProgramStatus
	}
	return nil
}

func (x *NicoliveState) GetModerationAnnouncement() *ModerationAnnouncement {
	if x != nil {
		return x.ModerationAnnouncement
	}
	return nil
}

var File_state_proto protoreflect.FileDescriptor

var file_state_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x64,
	0x77, 0x61, 0x6e, 0x67, 0x6f, 0x2e, 0x6e, 0x69, 0x63, 0x6f, 0x6c, 0x69, 0x76, 0x65, 0x2e, 0x63,
	0x68, 0x61, 0x74, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x0b, 0x61, 0x74, 0x6f, 0x6d, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0f, 0x6d, 0x6f, 0x64, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfc, 0x06, 0x0a, 0x0d, 0x4e, 0x69, 0x63, 0x6f, 0x6c,
	0x69, 0x76, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x4a, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74,
	0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x64,
	0x77, 0x61, 0x6e, 0x67, 0x6f, 0x2e, 0x6e, 0x69, 0x63, 0x6f, 0x6c, 0x69, 0x76, 0x65, 0x2e, 0x63,
	0x68, 0x61, 0x74, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74,
	0x69, 0x63, 0x73, 0x48, 0x00, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63,
	0x73, 0x88, 0x01, 0x01, 0x12, 0x41, 0x0a, 0x07, 0x65, 0x6e, 0x71, 0x75, 0x65, 0x74, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x64, 0x77, 0x61, 0x6e, 0x67, 0x6f, 0x2e, 0x6e,
	0x69, 0x63, 0x6f, 0x6c, 0x69, 0x76, 0x65, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x64, 0x61, 0x74,
	0x61, 0x2e, 0x45, 0x6e, 0x71, 0x75, 0x65, 0x74, 0x65, 0x48, 0x01, 0x52, 0x07, 0x65, 0x6e, 0x71,
	0x75, 0x65, 0x74, 0x65, 0x88, 0x01, 0x01, 0x12, 0x48, 0x0a, 0x0a, 0x6d, 0x6f, 0x76, 0x65, 0x5f,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x64, 0x77,
	0x61, 0x6e, 0x67, 0x6f, 0x2e, 0x6e, 0x69, 0x63, 0x6f, 0x6c, 0x69, 0x76, 0x65, 0x2e, 0x63, 0x68,
	0x61, 0x74, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x4d, 0x6f, 0x76, 0x65, 0x4f, 0x72, 0x64, 0x65,
	0x72, 0x48, 0x02, 0x52, 0x09, 0x6d, 0x6f, 0x76, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x88, 0x01,
	0x01, 0x12, 0x41, 0x0a, 0x07, 0x6d, 0x61, 0x72, 0x71, 0x75, 0x65, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x22, 0x2e, 0x64, 0x77, 0x61, 0x6e, 0x67, 0x6f, 0x2e, 0x6e, 0x69, 0x63, 0x6f,
	0x6c, 0x69, 0x76, 0x65, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x4d,
	0x61, 0x72, 0x71, 0x75, 0x65, 0x65, 0x48, 0x03, 0x52, 0x07, 0x6d, 0x61, 0x72, 0x71, 0x75, 0x65,
	0x65, 0x88, 0x01, 0x01, 0x12, 0x4e, 0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x5f,
	0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x64, 0x77, 0x61,
	0x6e, 0x67, 0x6f, 0x2e, 0x6e, 0x69, 0x63, 0x6f, 0x6c, 0x69, 0x76, 0x65, 0x2e, 0x63, 0x68, 0x61,
	0x74, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x4c, 0x6f,
	0x63, 0x6b, 0x48, 0x04, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x63,
	0x6b, 0x88, 0x01, 0x01, 0x12, 0x4e, 0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x5f,
	0x6d, 0x6f, 0x64, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x64, 0x77, 0x61,
	0x6e, 0x67, 0x6f, 0x2e, 0x6e, 0x69, 0x63, 0x6f, 0x6c, 0x69, 0x76, 0x65, 0x2e, 0x63, 0x68, 0x61,
	0x74, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x6f,
	0x64, 0x65, 0x48, 0x05, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x6f, 0x64,
	0x65, 0x88, 0x01, 0x01, 0x12, 0x4b, 0x0a, 0x0b, 0x74, 0x72, 0x69, 0x61, 0x6c, 0x5f, 0x70, 0x61,
	0x6e, 0x65, 0x6c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x64, 0x77, 0x61, 0x6e,
	0x67, 0x6f, 0x2e, 0x6e, 0x69, 0x63, 0x6f, 0x6c, 0x69, 0x76, 0x65, 0x2e, 0x63, 0x68, 0x61, 0x74,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x54, 0x72, 0x69, 0x61, 0x6c, 0x50, 0x61, 0x6e, 0x65, 0x6c,
	0x48, 0x06, 0x52, 0x0a, 0x74, 0x72, 0x69, 0x61, 0x6c, 0x50, 0x61, 0x6e, 0x65, 0x6c, 0x88, 0x01,
	0x01, 0x12, 0x54, 0x0a, 0x0e, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x64, 0x77, 0x61, 0x6e,
	0x67, 0x6f, 0x2e, 0x6e, 0x69, 0x63, 0x6f, 0x6c, 0x69, 0x76, 0x65, 0x2e, 0x63, 0x68, 0x61, 0x74,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x48, 0x07, 0x52, 0x0d, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x88, 0x01, 0x01, 0x12, 0x75, 0x0a, 0x17, 0x6d, 0x6f, 0x64, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x61, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x37, 0x2e, 0x64, 0x77, 0x61, 0x6e, 0x67,
	0x6f, 0x2e, 0x6e, 0x69, 0x63, 0x6f, 0x6c, 0x69, 0x76, 0x65, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x61, 0x74, 0x6f, 0x6d, 0x73, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x48, 0x08, 0x52, 0x16, 0x6d, 0x6f, 0x64, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41,
	0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x88, 0x01, 0x01, 0x42, 0x0d,
	0x0a, 0x0b, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x42, 0x0a, 0x0a,
	0x08, 0x5f, 0x65, 0x6e, 0x71, 0x75, 0x65, 0x74, 0x65, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x6d, 0x6f,
	0x76, 0x65, 0x5f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x6d, 0x61, 0x72,
	0x71, 0x75, 0x65, 0x65, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x5f, 0x6c, 0x6f, 0x63, 0x6b, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x74, 0x72, 0x69, 0x61, 0x6c,
	0x5f, 0x70, 0x61, 0x6e, 0x65, 0x6c, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x70, 0x72, 0x6f, 0x67, 0x72,
	0x61, 0x6d, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x1a, 0x0a, 0x18, 0x5f, 0x6d, 0x6f,
	0x64, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x61, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_state_proto_rawDescOnce sync.Once
	file_state_proto_rawDescData = file_state_proto_rawDesc
)

func file_state_proto_rawDescGZIP() []byte {
	file_state_proto_rawDescOnce.Do(func() {
		file_state_proto_rawDescData = protoimpl.X.CompressGZIP(file_state_proto_rawDescData)
	})
	return file_state_proto_rawDescData
}

var file_state_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_state_proto_goTypes = []any{
	(*NicoliveState)(nil),          // 0: dwango.nicolive.chat.data.NicoliveState
	(*Statistics)(nil),             // 1: dwango.nicolive.chat.data.Statistics
	(*Enquete)(nil),                // 2: dwango.nicolive.chat.data.Enquete
	(*MoveOrder)(nil),              // 3: dwango.nicolive.chat.data.MoveOrder
	(*Marquee)(nil),                // 4: dwango.nicolive.chat.data.Marquee
	(*CommentLock)(nil),            // 5: dwango.nicolive.chat.data.CommentLock
	(*CommentMode)(nil),            // 6: dwango.nicolive.chat.data.CommentMode
	(*TrialPanel)(nil),             // 7: dwango.nicolive.chat.data.TrialPanel
	(*ProgramStatus)(nil),          // 8: dwango.nicolive.chat.data.ProgramStatus
	(*ModerationAnnouncement)(nil), // 9: dwango.nicolive.chat.data.atoms.ModerationAnnouncement
}
var file_state_proto_depIdxs = []int32{
	1, // 0: dwango.nicolive.chat.data.NicoliveState.statistics:type_name -> dwango.nicolive.chat.data.Statistics
	2, // 1: dwango.nicolive.chat.data.NicoliveState.enquete:type_name -> dwango.nicolive.chat.data.Enquete
	3, // 2: dwango.nicolive.chat.data.NicoliveState.move_order:type_name -> dwango.nicolive.chat.data.MoveOrder
	4, // 3: dwango.nicolive.chat.data.NicoliveState.marquee:type_name -> dwango.nicolive.chat.data.Marquee
	5, // 4: dwango.nicolive.chat.data.NicoliveState.comment_lock:type_name -> dwango.nicolive.chat.data.CommentLock
	6, // 5: dwango.nicolive.chat.data.NicoliveState.comment_mode:type_name -> dwango.nicolive.chat.data.CommentMode
	7, // 6: dwango.nicolive.chat.data.NicoliveState.trial_panel:type_name -> dwango.nicolive.chat.data.TrialPanel
	8, // 7: dwango.nicolive.chat.data.NicoliveState.program_status:type_name -> dwango.nicolive.chat.data.ProgramStatus
	9, // 8: dwango.nicolive.chat.data.NicoliveState.moderation_announcement:type_name -> dwango.nicolive.chat.data.atoms.ModerationAnnouncement
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_state_proto_init() }
func file_state_proto_init() {
	if File_state_proto != nil {
		return
	}
	file_atoms_proto_init()
	file_moderator_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_state_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*NicoliveState); i {
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
	file_state_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_state_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_state_proto_goTypes,
		DependencyIndexes: file_state_proto_depIdxs,
		MessageInfos:      file_state_proto_msgTypes,
	}.Build()
	File_state_proto = out.File
	file_state_proto_rawDesc = nil
	file_state_proto_goTypes = nil
	file_state_proto_depIdxs = nil
}