// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.1
// 	protoc        v3.21.12
// source: detail.proto

package cotproto

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

type Detail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	XmlDetail string `protobuf:"bytes,1,opt,name=xmlDetail,proto3" json:"xmlDetail,omitempty"`
	// <contact>
	Contact *Contact `protobuf:"bytes,2,opt,name=contact,proto3" json:"contact,omitempty"`
	// <__group>
	Group *Group `protobuf:"bytes,3,opt,name=group,proto3" json:"group,omitempty"`
	// <precisionlocation>
	PrecisionLocation *PrecisionLocation `protobuf:"bytes,4,opt,name=precisionLocation,proto3" json:"precisionLocation,omitempty"`
	// <status>
	Status *Status `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
	// <takv>
	Takv *Takv `protobuf:"bytes,6,opt,name=takv,proto3" json:"takv,omitempty"`
	// <track>
	Track *Track `protobuf:"bytes,7,opt,name=track,proto3" json:"track,omitempty"`
}

func (x *Detail) Reset() {
	*x = Detail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_detail_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Detail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Detail) ProtoMessage() {}

func (x *Detail) ProtoReflect() protoreflect.Message {
	mi := &file_detail_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Detail.ProtoReflect.Descriptor instead.
func (*Detail) Descriptor() ([]byte, []int) {
	return file_detail_proto_rawDescGZIP(), []int{0}
}

func (x *Detail) GetXmlDetail() string {
	if x != nil {
		return x.XmlDetail
	}
	return ""
}

func (x *Detail) GetContact() *Contact {
	if x != nil {
		return x.Contact
	}
	return nil
}

func (x *Detail) GetGroup() *Group {
	if x != nil {
		return x.Group
	}
	return nil
}

func (x *Detail) GetPrecisionLocation() *PrecisionLocation {
	if x != nil {
		return x.PrecisionLocation
	}
	return nil
}

func (x *Detail) GetStatus() *Status {
	if x != nil {
		return x.Status
	}
	return nil
}

func (x *Detail) GetTakv() *Takv {
	if x != nil {
		return x.Takv
	}
	return nil
}

func (x *Detail) GetTrack() *Track {
	if x != nil {
		return x.Track
	}
	return nil
}

var File_detail_proto protoreflect.FileDescriptor

var file_detail_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d,
	0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x70, 0x72, 0x65, 0x63,
	0x69, 0x73, 0x69, 0x6f, 0x6e, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x0a, 0x74, 0x61, 0x6b, 0x76, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x74,
	0x72, 0x61, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x84, 0x02, 0x0a, 0x06, 0x44,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x78, 0x6d, 0x6c, 0x44, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x78, 0x6d, 0x6c, 0x44, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x12, 0x22, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x07,
	0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x12, 0x1c, 0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x05,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x40, 0x0a, 0x11, 0x70, 0x72, 0x65, 0x63, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x50, 0x72, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x4c, 0x6f, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x11, 0x70, 0x72, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x4c,
	0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x19, 0x0a, 0x04, 0x74, 0x61, 0x6b, 0x76,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x54, 0x61, 0x6b, 0x76, 0x52, 0x04, 0x74,
	0x61, 0x6b, 0x76, 0x12, 0x1c, 0x0a, 0x05, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x06, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x52, 0x05, 0x74, 0x72, 0x61, 0x63,
	0x6b, 0x42, 0x26, 0x48, 0x03, 0x5a, 0x22, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6b, 0x64, 0x75, 0x64, 0x6b, 0x6f, 0x76, 0x2f, 0x67, 0x6f, 0x61, 0x74, 0x61, 0x6b,
	0x2f, 0x63, 0x6f, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_detail_proto_rawDescOnce sync.Once
	file_detail_proto_rawDescData = file_detail_proto_rawDesc
)

func file_detail_proto_rawDescGZIP() []byte {
	file_detail_proto_rawDescOnce.Do(func() {
		file_detail_proto_rawDescData = protoimpl.X.CompressGZIP(file_detail_proto_rawDescData)
	})
	return file_detail_proto_rawDescData
}

var file_detail_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_detail_proto_goTypes = []interface{}{
	(*Detail)(nil),            // 0: Detail
	(*Contact)(nil),           // 1: Contact
	(*Group)(nil),             // 2: Group
	(*PrecisionLocation)(nil), // 3: PrecisionLocation
	(*Status)(nil),            // 4: Status
	(*Takv)(nil),              // 5: Takv
	(*Track)(nil),             // 6: Track
}
var file_detail_proto_depIdxs = []int32{
	1, // 0: Detail.contact:type_name -> Contact
	2, // 1: Detail.group:type_name -> Group
	3, // 2: Detail.precisionLocation:type_name -> PrecisionLocation
	4, // 3: Detail.status:type_name -> Status
	5, // 4: Detail.takv:type_name -> Takv
	6, // 5: Detail.track:type_name -> Track
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_detail_proto_init() }
func file_detail_proto_init() {
	if File_detail_proto != nil {
		return
	}
	file_contact_proto_init()
	file_group_proto_init()
	file_precisionlocation_proto_init()
	file_status_proto_init()
	file_takv_proto_init()
	file_track_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_detail_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Detail); i {
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
			RawDescriptor: file_detail_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_detail_proto_goTypes,
		DependencyIndexes: file_detail_proto_depIdxs,
		MessageInfos:      file_detail_proto_msgTypes,
	}.Build()
	File_detail_proto = out.File
	file_detail_proto_rawDesc = nil
	file_detail_proto_goTypes = nil
	file_detail_proto_depIdxs = nil
}
