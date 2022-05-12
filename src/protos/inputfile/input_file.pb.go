// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: input_file.proto

package inputfile

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

type InputFileType int32

const (
	InputFileType_CSV InputFileType = 0
)

// Enum value maps for InputFileType.
var (
	InputFileType_name = map[int32]string{
		0: "CSV",
	}
	InputFileType_value = map[string]int32{
		"CSV": 0,
	}
)

func (x InputFileType) Enum() *InputFileType {
	p := new(InputFileType)
	*p = x
	return p
}

func (x InputFileType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (InputFileType) Descriptor() protoreflect.EnumDescriptor {
	return file_input_file_proto_enumTypes[0].Descriptor()
}

func (InputFileType) Type() protoreflect.EnumType {
	return &file_input_file_proto_enumTypes[0]
}

func (x InputFileType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use InputFileType.Descriptor instead.
func (InputFileType) EnumDescriptor() ([]byte, []int) {
	return file_input_file_proto_rawDescGZIP(), []int{0}
}

type Type int32

const (
	Type_ORGANIZATION Type = 0
	Type_SCHOOL       Type = 1
	Type_CLASS        Type = 2
	Type_USER         Type = 3
	Type_ROLE         Type = 4
	Type_PROGRAM      Type = 5
	Type_UNKNOWN      Type = 6
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "ORGANIZATION",
		1: "SCHOOL",
		2: "CLASS",
		3: "USER",
		4: "ROLE",
		5: "PROGRAM",
		6: "UNKNOWN",
	}
	Type_value = map[string]int32{
		"ORGANIZATION": 0,
		"SCHOOL":       1,
		"CLASS":        2,
		"USER":         3,
		"ROLE":         4,
		"PROGRAM":      5,
		"UNKNOWN":      6,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_input_file_proto_enumTypes[1].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_input_file_proto_enumTypes[1]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_input_file_proto_rawDescGZIP(), []int{1}
}

type InputFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId        string        `protobuf:"bytes,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	InputFileType InputFileType `protobuf:"varint,2,opt,name=input_file_type,json=inputFileType,proto3,enum=protos.inputfile.InputFileType" json:"input_file_type,omitempty"`
	Path          string        `protobuf:"bytes,3,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *InputFile) Reset() {
	*x = InputFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_input_file_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InputFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputFile) ProtoMessage() {}

func (x *InputFile) ProtoReflect() protoreflect.Message {
	mi := &file_input_file_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputFile.ProtoReflect.Descriptor instead.
func (*InputFile) Descriptor() ([]byte, []int) {
	return file_input_file_proto_rawDescGZIP(), []int{0}
}

func (x *InputFile) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

func (x *InputFile) GetInputFileType() InputFileType {
	if x != nil {
		return x.InputFileType
	}
	return InputFileType_CSV
}

func (x *InputFile) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

type InputFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      Type       `protobuf:"varint,1,opt,name=type,proto3,enum=protos.inputfile.Type" json:"type,omitempty"`
	InputFile *InputFile `protobuf:"bytes,2,opt,name=input_file,json=inputFile,proto3" json:"input_file,omitempty"`
}

func (x *InputFileRequest) Reset() {
	*x = InputFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_input_file_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InputFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputFileRequest) ProtoMessage() {}

func (x *InputFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_input_file_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputFileRequest.ProtoReflect.Descriptor instead.
func (*InputFileRequest) Descriptor() ([]byte, []int) {
	return file_input_file_proto_rawDescGZIP(), []int{1}
}

func (x *InputFileRequest) GetType() Type {
	if x != nil {
		return x.Type
	}
	return Type_ORGANIZATION
}

func (x *InputFileRequest) GetInputFile() *InputFile {
	if x != nil {
		return x.InputFile
	}
	return nil
}

type InputFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool              `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	Errors  []*InputFileError `protobuf:"bytes,3,rep,name=errors,proto3" json:"errors,omitempty"`
}

func (x *InputFileResponse) Reset() {
	*x = InputFileResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_input_file_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InputFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputFileResponse) ProtoMessage() {}

func (x *InputFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_input_file_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputFileResponse.ProtoReflect.Descriptor instead.
func (*InputFileResponse) Descriptor() ([]byte, []int) {
	return file_input_file_proto_rawDescGZIP(), []int{2}
}

func (x *InputFileResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *InputFileResponse) GetErrors() []*InputFileError {
	if x != nil {
		return x.Errors
	}
	return nil
}

type InputFileError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId  string   `protobuf:"bytes,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	Message []string `protobuf:"bytes,2,rep,name=message,proto3" json:"message,omitempty"`
}

func (x *InputFileError) Reset() {
	*x = InputFileError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_input_file_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InputFileError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputFileError) ProtoMessage() {}

func (x *InputFileError) ProtoReflect() protoreflect.Message {
	mi := &file_input_file_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputFileError.ProtoReflect.Descriptor instead.
func (*InputFileError) Descriptor() ([]byte, []int) {
	return file_input_file_proto_rawDescGZIP(), []int{3}
}

func (x *InputFileError) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

func (x *InputFileError) GetMessage() []string {
	if x != nil {
		return x.Message
	}
	return nil
}

var File_input_file_proto protoreflect.FileDescriptor

var file_input_file_proto_rawDesc = []byte{
	0x0a, 0x10, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x69, 0x6e, 0x70, 0x75, 0x74,
	0x66, 0x69, 0x6c, 0x65, 0x22, 0x81, 0x01, 0x0a, 0x09, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x46, 0x69,
	0x6c, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x47, 0x0a, 0x0f, 0x69,
	0x6e, 0x70, 0x75, 0x74, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x69, 0x6e,
	0x70, 0x75, 0x74, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x46, 0x69, 0x6c,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0d, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x22, 0x7a, 0x0a, 0x10, 0x49, 0x6e, 0x70, 0x75,
	0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x2e, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x3a, 0x0a, 0x0a, 0x69, 0x6e, 0x70, 0x75,
	0x74, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x66, 0x69, 0x6c, 0x65, 0x2e,
	0x49, 0x6e, 0x70, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x09, 0x69, 0x6e, 0x70, 0x75, 0x74,
	0x46, 0x69, 0x6c, 0x65, 0x22, 0x67, 0x0a, 0x11, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x46, 0x69, 0x6c,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x12, 0x38, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x69, 0x6e, 0x70,
	0x75, 0x74, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x06, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x22, 0x43, 0x0a,
	0x0e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12,
	0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2a, 0x18, 0x0a, 0x0d, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x07, 0x0a, 0x03, 0x43, 0x53, 0x56, 0x10, 0x00, 0x2a, 0x5d, 0x0a, 0x04,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x0c, 0x4f, 0x52, 0x47, 0x41, 0x4e, 0x49, 0x5a, 0x41,
	0x54, 0x49, 0x4f, 0x4e, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x43, 0x48, 0x4f, 0x4f, 0x4c,
	0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x43, 0x4c, 0x41, 0x53, 0x53, 0x10, 0x02, 0x12, 0x08, 0x0a,
	0x04, 0x55, 0x53, 0x45, 0x52, 0x10, 0x03, 0x12, 0x08, 0x0a, 0x04, 0x52, 0x4f, 0x4c, 0x45, 0x10,
	0x04, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x52, 0x4f, 0x47, 0x52, 0x41, 0x4d, 0x10, 0x05, 0x12, 0x0b,
	0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x06, 0x32, 0xd3, 0x01, 0x0a, 0x11,
	0x49, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x5e, 0x0a, 0x0f, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x50,
	0x52, 0x4f, 0x54, 0x4f, 0x12, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x69, 0x6e,
	0x70, 0x75, 0x74, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x46, 0x69, 0x6c,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x73, 0x2e, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x49, 0x6e, 0x70, 0x75,
	0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28,
	0x01, 0x12, 0x5e, 0x0a, 0x0f, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x41,
	0x56, 0x52, 0x4f, 0x53, 0x12, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x69, 0x6e,
	0x70, 0x75, 0x74, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x46, 0x69, 0x6c,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x73, 0x2e, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x49, 0x6e, 0x70, 0x75,
	0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28,
	0x01, 0x42, 0x16, 0x5a, 0x14, 0x73, 0x72, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f,
	0x69, 0x6e, 0x70, 0x75, 0x74, 0x66, 0x69, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_input_file_proto_rawDescOnce sync.Once
	file_input_file_proto_rawDescData = file_input_file_proto_rawDesc
)

func file_input_file_proto_rawDescGZIP() []byte {
	file_input_file_proto_rawDescOnce.Do(func() {
		file_input_file_proto_rawDescData = protoimpl.X.CompressGZIP(file_input_file_proto_rawDescData)
	})
	return file_input_file_proto_rawDescData
}

var file_input_file_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_input_file_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_input_file_proto_goTypes = []interface{}{
	(InputFileType)(0),        // 0: protos.inputfile.InputFileType
	(Type)(0),                 // 1: protos.inputfile.Type
	(*InputFile)(nil),         // 2: protos.inputfile.InputFile
	(*InputFileRequest)(nil),  // 3: protos.inputfile.InputFileRequest
	(*InputFileResponse)(nil), // 4: protos.inputfile.InputFileResponse
	(*InputFileError)(nil),    // 5: protos.inputfile.InputFileError
}
var file_input_file_proto_depIdxs = []int32{
	0, // 0: protos.inputfile.InputFile.input_file_type:type_name -> protos.inputfile.InputFileType
	1, // 1: protos.inputfile.InputFileRequest.type:type_name -> protos.inputfile.Type
	2, // 2: protos.inputfile.InputFileRequest.input_file:type_name -> protos.inputfile.InputFile
	5, // 3: protos.inputfile.InputFileResponse.errors:type_name -> protos.inputfile.InputFileError
	3, // 4: protos.inputfile.IngestFileService.IngestFilePROTO:input_type -> protos.inputfile.InputFileRequest
	3, // 5: protos.inputfile.IngestFileService.IngestFileAVROS:input_type -> protos.inputfile.InputFileRequest
	4, // 6: protos.inputfile.IngestFileService.IngestFilePROTO:output_type -> protos.inputfile.InputFileResponse
	4, // 7: protos.inputfile.IngestFileService.IngestFileAVROS:output_type -> protos.inputfile.InputFileResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_input_file_proto_init() }
func file_input_file_proto_init() {
	if File_input_file_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_input_file_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InputFile); i {
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
		file_input_file_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InputFileRequest); i {
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
		file_input_file_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InputFileResponse); i {
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
		file_input_file_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InputFileError); i {
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
			RawDescriptor: file_input_file_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_input_file_proto_goTypes,
		DependencyIndexes: file_input_file_proto_depIdxs,
		EnumInfos:         file_input_file_proto_enumTypes,
		MessageInfos:      file_input_file_proto_msgTypes,
	}.Build()
	File_input_file_proto = out.File
	file_input_file_proto_rawDesc = nil
	file_input_file_proto_goTypes = nil
	file_input_file_proto_depIdxs = nil
}
