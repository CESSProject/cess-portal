// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: file.proto

package module

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

type FileUploadInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId    string `protobuf:"bytes,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	FileHash  string `protobuf:"bytes,2,opt,name=file_hash,json=fileHash,proto3" json:"file_hash,omitempty"`
	Backups   string `protobuf:"bytes,3,opt,name=backups,proto3" json:"backups,omitempty"`
	Blocks    int32  `protobuf:"varint,4,opt,name=blocks,proto3" json:"blocks,omitempty"`
	BlockSize int32  `protobuf:"varint,5,opt,name=block_size,json=blockSize,proto3" json:"block_size,omitempty"`
	BlockNum  int32  `protobuf:"varint,6,opt,name=block_num,json=blockNum,proto3" json:"block_num,omitempty"`
	Data      []byte `protobuf:"bytes,7,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *FileUploadInfo) Reset() {
	*x = FileUploadInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileUploadInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileUploadInfo) ProtoMessage() {}

func (x *FileUploadInfo) ProtoReflect() protoreflect.Message {
	mi := &file_file_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileUploadInfo.ProtoReflect.Descriptor instead.
func (*FileUploadInfo) Descriptor() ([]byte, []int) {
	return file_file_proto_rawDescGZIP(), []int{0}
}

func (x *FileUploadInfo) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

func (x *FileUploadInfo) GetFileHash() string {
	if x != nil {
		return x.FileHash
	}
	return ""
}

func (x *FileUploadInfo) GetBackups() string {
	if x != nil {
		return x.Backups
	}
	return ""
}

func (x *FileUploadInfo) GetBlocks() int32 {
	if x != nil {
		return x.Blocks
	}
	return 0
}

func (x *FileUploadInfo) GetBlockSize() int32 {
	if x != nil {
		return x.BlockSize
	}
	return 0
}

func (x *FileUploadInfo) GetBlockNum() int32 {
	if x != nil {
		return x.BlockNum
	}
	return 0
}

func (x *FileUploadInfo) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type FileDownloadInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId    string `protobuf:"bytes,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	Blocks    int32  `protobuf:"varint,2,opt,name=blocks,proto3" json:"blocks,omitempty"`
	BlockSize int32  `protobuf:"varint,3,opt,name=block_size,json=blockSize,proto3" json:"block_size,omitempty"`
	BlockNum  int32  `protobuf:"varint,4,opt,name=block_num,json=blockNum,proto3" json:"block_num,omitempty"`
	Data      []byte `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *FileDownloadInfo) Reset() {
	*x = FileDownloadInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileDownloadInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileDownloadInfo) ProtoMessage() {}

func (x *FileDownloadInfo) ProtoReflect() protoreflect.Message {
	mi := &file_file_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileDownloadInfo.ProtoReflect.Descriptor instead.
func (*FileDownloadInfo) Descriptor() ([]byte, []int) {
	return file_file_proto_rawDescGZIP(), []int{1}
}

func (x *FileDownloadInfo) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

func (x *FileDownloadInfo) GetBlocks() int32 {
	if x != nil {
		return x.Blocks
	}
	return 0
}

func (x *FileDownloadInfo) GetBlockSize() int32 {
	if x != nil {
		return x.BlockSize
	}
	return 0
}

func (x *FileDownloadInfo) GetBlockNum() int32 {
	if x != nil {
		return x.BlockNum
	}
	return 0
}

func (x *FileDownloadInfo) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type FileDownloadReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId        string `protobuf:"bytes,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	WalletAddress string `protobuf:"bytes,2,opt,name=walletAddress,proto3" json:"walletAddress,omitempty"`
	Blocks        int32  `protobuf:"varint,3,opt,name=blocks,proto3" json:"blocks,omitempty"`
}

func (x *FileDownloadReq) Reset() {
	*x = FileDownloadReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileDownloadReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileDownloadReq) ProtoMessage() {}

func (x *FileDownloadReq) ProtoReflect() protoreflect.Message {
	mi := &file_file_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileDownloadReq.ProtoReflect.Descriptor instead.
func (*FileDownloadReq) Descriptor() ([]byte, []int) {
	return file_file_proto_rawDescGZIP(), []int{2}
}

func (x *FileDownloadReq) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

func (x *FileDownloadReq) GetWalletAddress() string {
	if x != nil {
		return x.WalletAddress
	}
	return ""
}

func (x *FileDownloadReq) GetBlocks() int32 {
	if x != nil {
		return x.Blocks
	}
	return 0
}

var File_file_proto protoreflect.FileDescriptor

var file_file_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x22, 0xca, 0x01, 0x0a, 0x10, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x75, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65,
	0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x48, 0x61, 0x73, 0x68, 0x12,
	0x18, 0x0a, 0x07, 0x62, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x62, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x69, 0x7a, 0x65,
	0x12, 0x1b, 0x0a, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x08, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x22, 0x95, 0x01, 0x0a, 0x12, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x64, 0x6f, 0x77, 0x6e, 0x6c,
	0x6f, 0x61, 0x64, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x6a, 0x0a, 0x11, 0x66, 0x69, 0x6c,
	0x65, 0x5f, 0x64, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x72, 0x65, 0x71, 0x12, 0x17,
	0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0d, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x16, 0x0a,
	0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x2f, 0x3b, 0x6d, 0x6f, 0x64, 0x75,
	0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_file_proto_rawDescOnce sync.Once
	file_file_proto_rawDescData = file_file_proto_rawDesc
)

func file_file_proto_rawDescGZIP() []byte {
	file_file_proto_rawDescOnce.Do(func() {
		file_file_proto_rawDescData = protoimpl.X.CompressGZIP(file_file_proto_rawDescData)
	})
	return file_file_proto_rawDescData
}

var file_file_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_file_proto_goTypes = []interface{}{
	(*FileUploadInfo)(nil),   // 0: client.file_upload_info
	(*FileDownloadInfo)(nil), // 1: client.file_download_info
	(*FileDownloadReq)(nil),  // 2: client.file_download_req
}
var file_file_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_file_proto_init() }
func file_file_proto_init() {
	if File_file_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_file_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileUploadInfo); i {
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
		file_file_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileDownloadInfo); i {
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
		file_file_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileDownloadReq); i {
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
			RawDescriptor: file_file_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_file_proto_goTypes,
		DependencyIndexes: file_file_proto_depIdxs,
		MessageInfos:      file_file_proto_msgTypes,
	}.Build()
	File_file_proto = out.File
	file_file_proto_rawDesc = nil
	file_file_proto_goTypes = nil
	file_file_proto_depIdxs = nil
}