// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: media/media.proto

package ssov6

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

type UploadMediaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data     []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	FileName string `protobuf:"bytes,2,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	MimeType string `protobuf:"bytes,3,opt,name=mime_type,json=mimeType,proto3" json:"mime_type,omitempty"`
}

func (x *UploadMediaRequest) Reset() {
	*x = UploadMediaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_media_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadMediaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadMediaRequest) ProtoMessage() {}

func (x *UploadMediaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_media_media_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadMediaRequest.ProtoReflect.Descriptor instead.
func (*UploadMediaRequest) Descriptor() ([]byte, []int) {
	return file_media_media_proto_rawDescGZIP(), []int{0}
}

func (x *UploadMediaRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *UploadMediaRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *UploadMediaRequest) GetMimeType() string {
	if x != nil {
		return x.MimeType
	}
	return ""
}

type UploadMediaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId int64 `protobuf:"varint,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
}

func (x *UploadMediaResponse) Reset() {
	*x = UploadMediaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_media_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadMediaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadMediaResponse) ProtoMessage() {}

func (x *UploadMediaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_media_media_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadMediaResponse.ProtoReflect.Descriptor instead.
func (*UploadMediaResponse) Descriptor() ([]byte, []int) {
	return file_media_media_proto_rawDescGZIP(), []int{1}
}

func (x *UploadMediaResponse) GetFileId() int64 {
	if x != nil {
		return x.FileId
	}
	return 0
}

type DownloadMediaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId int64 `protobuf:"varint,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
}

func (x *DownloadMediaRequest) Reset() {
	*x = DownloadMediaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_media_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadMediaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadMediaRequest) ProtoMessage() {}

func (x *DownloadMediaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_media_media_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadMediaRequest.ProtoReflect.Descriptor instead.
func (*DownloadMediaRequest) Descriptor() ([]byte, []int) {
	return file_media_media_proto_rawDescGZIP(), []int{2}
}

func (x *DownloadMediaRequest) GetFileId() int64 {
	if x != nil {
		return x.FileId
	}
	return 0
}

type DownloadMediaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data     []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	FileName string `protobuf:"bytes,2,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	MimeType string `protobuf:"bytes,3,opt,name=mime_type,json=mimeType,proto3" json:"mime_type,omitempty"`
}

func (x *DownloadMediaResponse) Reset() {
	*x = DownloadMediaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_media_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadMediaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadMediaResponse) ProtoMessage() {}

func (x *DownloadMediaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_media_media_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadMediaResponse.ProtoReflect.Descriptor instead.
func (*DownloadMediaResponse) Descriptor() ([]byte, []int) {
	return file_media_media_proto_rawDescGZIP(), []int{3}
}

func (x *DownloadMediaResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *DownloadMediaResponse) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *DownloadMediaResponse) GetMimeType() string {
	if x != nil {
		return x.MimeType
	}
	return ""
}

var File_media_media_proto protoreflect.FileDescriptor

var file_media_media_proto_rawDesc = []byte{
	0x0a, 0x11, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x2f, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x22, 0x62, 0x0a, 0x12, 0x55, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x69, 0x6d, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x69, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x22, 0x2e,
	0x0a, 0x13, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x22, 0x2f,
	0x0a, 0x14, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x22,
	0x65, 0x0a, 0x15, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09,
	0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x69, 0x6d,
	0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x69,
	0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x32, 0x99, 0x01, 0x0a, 0x05, 0x4d, 0x65, 0x64, 0x69, 0x61,
	0x12, 0x44, 0x0a, 0x0b, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x12,
	0x19, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x65,
	0x64, 0x69, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6d, 0x65, 0x64,
	0x69, 0x61, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x0d, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f,
	0x61, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x12, 0x1b, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x2e,
	0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x2e, 0x44, 0x6f, 0x77,
	0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x1d, 0x5a, 0x1b, 0x6d, 0x67, 0x72, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x2d, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x73, 0x73, 0x6f, 0x2e, 0x76, 0x31, 0x3b, 0x73, 0x73, 0x6f, 0x76,
	0x36, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_media_media_proto_rawDescOnce sync.Once
	file_media_media_proto_rawDescData = file_media_media_proto_rawDesc
)

func file_media_media_proto_rawDescGZIP() []byte {
	file_media_media_proto_rawDescOnce.Do(func() {
		file_media_media_proto_rawDescData = protoimpl.X.CompressGZIP(file_media_media_proto_rawDescData)
	})
	return file_media_media_proto_rawDescData
}

var file_media_media_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_media_media_proto_goTypes = []any{
	(*UploadMediaRequest)(nil),    // 0: media.UploadMediaRequest
	(*UploadMediaResponse)(nil),   // 1: media.UploadMediaResponse
	(*DownloadMediaRequest)(nil),  // 2: media.DownloadMediaRequest
	(*DownloadMediaResponse)(nil), // 3: media.DownloadMediaResponse
}
var file_media_media_proto_depIdxs = []int32{
	0, // 0: media.Media.UploadMedia:input_type -> media.UploadMediaRequest
	2, // 1: media.Media.DownloadMedia:input_type -> media.DownloadMediaRequest
	1, // 2: media.Media.UploadMedia:output_type -> media.UploadMediaResponse
	3, // 3: media.Media.DownloadMedia:output_type -> media.DownloadMediaResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_media_media_proto_init() }
func file_media_media_proto_init() {
	if File_media_media_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_media_media_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*UploadMediaRequest); i {
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
		file_media_media_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*UploadMediaResponse); i {
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
		file_media_media_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*DownloadMediaRequest); i {
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
		file_media_media_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*DownloadMediaResponse); i {
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
			RawDescriptor: file_media_media_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_media_media_proto_goTypes,
		DependencyIndexes: file_media_media_proto_depIdxs,
		MessageInfos:      file_media_media_proto_msgTypes,
	}.Build()
	File_media_media_proto = out.File
	file_media_media_proto_rawDesc = nil
	file_media_media_proto_goTypes = nil
	file_media_media_proto_depIdxs = nil
}
