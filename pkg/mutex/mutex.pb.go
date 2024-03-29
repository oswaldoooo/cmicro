// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.20.3
// source: mutex.proto

package mutex

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type MutexInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MutexId int32 `protobuf:"varint,1,opt,name=mutex_id,json=mutexId,proto3" json:"mutex_id,omitempty"`
}

func (x *MutexInfo) Reset() {
	*x = MutexInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mutex_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MutexInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MutexInfo) ProtoMessage() {}

func (x *MutexInfo) ProtoReflect() protoreflect.Message {
	mi := &file_mutex_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MutexInfo.ProtoReflect.Descriptor instead.
func (*MutexInfo) Descriptor() ([]byte, []int) {
	return file_mutex_proto_rawDescGZIP(), []int{0}
}

func (x *MutexInfo) GetMutexId() int32 {
	if x != nil {
		return x.MutexId
	}
	return 0
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mutex_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_mutex_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_mutex_proto_rawDescGZIP(), []int{1}
}

func (x *Response) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

type ResponseTwo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok  bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Rok bool `protobuf:"varint,2,opt,name=rok,proto3" json:"rok,omitempty"`
}

func (x *ResponseTwo) Reset() {
	*x = ResponseTwo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mutex_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseTwo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseTwo) ProtoMessage() {}

func (x *ResponseTwo) ProtoReflect() protoreflect.Message {
	mi := &file_mutex_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseTwo.ProtoReflect.Descriptor instead.
func (*ResponseTwo) Descriptor() ([]byte, []int) {
	return file_mutex_proto_rawDescGZIP(), []int{2}
}

func (x *ResponseTwo) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *ResponseTwo) GetRok() bool {
	if x != nil {
		return x.Rok
	}
	return false
}

var File_mutex_proto protoreflect.FileDescriptor

var file_mutex_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6d, 0x75, 0x74, 0x65, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6d,
	0x61, 0x69, 0x6e, 0x22, 0x27, 0x0a, 0x0a, 0x6d, 0x75, 0x74, 0x65, 0x78, 0x5f, 0x69, 0x6e, 0x66,
	0x6f, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x75, 0x74, 0x65, 0x78, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x07, 0x6d, 0x75, 0x74, 0x65, 0x78, 0x49, 0x64, 0x22, 0x1a, 0x0a, 0x08,
	0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x22, 0x30, 0x0a, 0x0c, 0x72, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x74, 0x77, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x6f, 0x6b, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x72, 0x6f, 0x6b, 0x32, 0xe6, 0x01, 0x0a, 0x05, 0x6d,
	0x75, 0x74, 0x65, 0x78, 0x12, 0x28, 0x0a, 0x04, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x10, 0x2e, 0x6d,
	0x61, 0x69, 0x6e, 0x2e, 0x6d, 0x75, 0x74, 0x65, 0x78, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x1a, 0x0e,
	0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a,
	0x0a, 0x06, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x10, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e,
	0x6d, 0x75, 0x74, 0x65, 0x78, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x1a, 0x0e, 0x2e, 0x6d, 0x61, 0x69,
	0x6e, 0x2e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a, 0x07, 0x67, 0x65,
	0x74, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x10, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x6d, 0x75, 0x74,
	0x65, 0x78, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x1a, 0x10, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x6d,
	0x75, 0x74, 0x65, 0x78, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x12, 0x2c, 0x0a, 0x08, 0x66, 0x72, 0x65,
	0x65, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x10, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x6d, 0x75, 0x74,
	0x65, 0x78, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x1a, 0x0e, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x72,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x06, 0x69, 0x73, 0x6c, 0x6f, 0x63,
	0x6b, 0x12, 0x10, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x6d, 0x75, 0x74, 0x65, 0x78, 0x5f, 0x69,
	0x6e, 0x66, 0x6f, 0x1a, 0x0e, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x72, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x32, 0xc4, 0x02, 0x0a, 0x07, 0x72, 0x77, 0x6d, 0x75, 0x74, 0x65, 0x78, 0x12,
	0x28, 0x0a, 0x04, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x10, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x6d,
	0x75, 0x74, 0x65, 0x78, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x1a, 0x0e, 0x2e, 0x6d, 0x61, 0x69, 0x6e,
	0x2e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x06, 0x75, 0x6e, 0x6c,
	0x6f, 0x63, 0x6b, 0x12, 0x10, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x6d, 0x75, 0x74, 0x65, 0x78,
	0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x1a, 0x0e, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x05, 0x72, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x10,
	0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x6d, 0x75, 0x74, 0x65, 0x78, 0x5f, 0x69, 0x6e, 0x66, 0x6f,
	0x1a, 0x0e, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x2b, 0x0a, 0x07, 0x72, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x10, 0x2e, 0x6d, 0x61,
	0x69, 0x6e, 0x2e, 0x6d, 0x75, 0x74, 0x65, 0x78, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x1a, 0x0e, 0x2e,
	0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a,
	0x07, 0x67, 0x65, 0x74, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x10, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e,
	0x6d, 0x75, 0x74, 0x65, 0x78, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x1a, 0x10, 0x2e, 0x6d, 0x61, 0x69,
	0x6e, 0x2e, 0x6d, 0x75, 0x74, 0x65, 0x78, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x12, 0x2c, 0x0a, 0x08,
	0x66, 0x72, 0x65, 0x65, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x10, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e,
	0x6d, 0x75, 0x74, 0x65, 0x78, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x1a, 0x0e, 0x2e, 0x6d, 0x61, 0x69,
	0x6e, 0x2e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x06, 0x69, 0x73,
	0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x10, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x6d, 0x75, 0x74, 0x65,
	0x78, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x1a, 0x12, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x74, 0x77, 0x6f, 0x42, 0x09, 0x5a, 0x07, 0x2f, 0x3b,
	0x6d, 0x75, 0x74, 0x65, 0x78, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mutex_proto_rawDescOnce sync.Once
	file_mutex_proto_rawDescData = file_mutex_proto_rawDesc
)

func file_mutex_proto_rawDescGZIP() []byte {
	file_mutex_proto_rawDescOnce.Do(func() {
		file_mutex_proto_rawDescData = protoimpl.X.CompressGZIP(file_mutex_proto_rawDescData)
	})
	return file_mutex_proto_rawDescData
}

var file_mutex_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_mutex_proto_goTypes = []interface{}{
	(*MutexInfo)(nil),   // 0: main.mutex_info
	(*Response)(nil),    // 1: main.response
	(*ResponseTwo)(nil), // 2: main.response_two
}
var file_mutex_proto_depIdxs = []int32{
	0,  // 0: main.mutex.lock:input_type -> main.mutex_info
	0,  // 1: main.mutex.unlock:input_type -> main.mutex_info
	0,  // 2: main.mutex.getlock:input_type -> main.mutex_info
	0,  // 3: main.mutex.freelock:input_type -> main.mutex_info
	0,  // 4: main.mutex.islock:input_type -> main.mutex_info
	0,  // 5: main.rwmutex.lock:input_type -> main.mutex_info
	0,  // 6: main.rwmutex.unlock:input_type -> main.mutex_info
	0,  // 7: main.rwmutex.rlock:input_type -> main.mutex_info
	0,  // 8: main.rwmutex.runlock:input_type -> main.mutex_info
	0,  // 9: main.rwmutex.getlock:input_type -> main.mutex_info
	0,  // 10: main.rwmutex.freelock:input_type -> main.mutex_info
	0,  // 11: main.rwmutex.islock:input_type -> main.mutex_info
	1,  // 12: main.mutex.lock:output_type -> main.response
	1,  // 13: main.mutex.unlock:output_type -> main.response
	0,  // 14: main.mutex.getlock:output_type -> main.mutex_info
	1,  // 15: main.mutex.freelock:output_type -> main.response
	1,  // 16: main.mutex.islock:output_type -> main.response
	1,  // 17: main.rwmutex.lock:output_type -> main.response
	1,  // 18: main.rwmutex.unlock:output_type -> main.response
	1,  // 19: main.rwmutex.rlock:output_type -> main.response
	1,  // 20: main.rwmutex.runlock:output_type -> main.response
	0,  // 21: main.rwmutex.getlock:output_type -> main.mutex_info
	1,  // 22: main.rwmutex.freelock:output_type -> main.response
	2,  // 23: main.rwmutex.islock:output_type -> main.response_two
	12, // [12:24] is the sub-list for method output_type
	0,  // [0:12] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_mutex_proto_init() }
func file_mutex_proto_init() {
	if File_mutex_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mutex_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MutexInfo); i {
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
		file_mutex_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_mutex_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseTwo); i {
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
			RawDescriptor: file_mutex_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_mutex_proto_goTypes,
		DependencyIndexes: file_mutex_proto_depIdxs,
		MessageInfos:      file_mutex_proto_msgTypes,
	}.Build()
	File_mutex_proto = out.File
	file_mutex_proto_rawDesc = nil
	file_mutex_proto_goTypes = nil
	file_mutex_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// MutexClient is the client API for Mutex service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MutexClient interface {
	Lock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error)
	Unlock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error)
	Getlock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*MutexInfo, error)
	Freelock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error)
	Islock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error)
}

type mutexClient struct {
	cc grpc.ClientConnInterface
}

func NewMutexClient(cc grpc.ClientConnInterface) MutexClient {
	return &mutexClient{cc}
}

func (c *mutexClient) Lock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/main.mutex/lock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mutexClient) Unlock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/main.mutex/unlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mutexClient) Getlock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*MutexInfo, error) {
	out := new(MutexInfo)
	err := c.cc.Invoke(ctx, "/main.mutex/getlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mutexClient) Freelock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/main.mutex/freelock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mutexClient) Islock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/main.mutex/islock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MutexServer is the server API for Mutex service.
type MutexServer interface {
	Lock(context.Context, *MutexInfo) (*Response, error)
	Unlock(context.Context, *MutexInfo) (*Response, error)
	Getlock(context.Context, *MutexInfo) (*MutexInfo, error)
	Freelock(context.Context, *MutexInfo) (*Response, error)
	Islock(context.Context, *MutexInfo) (*Response, error)
}

// UnimplementedMutexServer can be embedded to have forward compatible implementations.
type UnimplementedMutexServer struct {
}

func (*UnimplementedMutexServer) Lock(context.Context, *MutexInfo) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Lock not implemented")
}
func (*UnimplementedMutexServer) Unlock(context.Context, *MutexInfo) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unlock not implemented")
}
func (*UnimplementedMutexServer) Getlock(context.Context, *MutexInfo) (*MutexInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Getlock not implemented")
}
func (*UnimplementedMutexServer) Freelock(context.Context, *MutexInfo) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Freelock not implemented")
}
func (*UnimplementedMutexServer) Islock(context.Context, *MutexInfo) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Islock not implemented")
}

func RegisterMutexServer(s *grpc.Server, srv MutexServer) {
	s.RegisterService(&_Mutex_serviceDesc, srv)
}

func _Mutex_Lock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutexInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MutexServer).Lock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.mutex/Lock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MutexServer).Lock(ctx, req.(*MutexInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mutex_Unlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutexInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MutexServer).Unlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.mutex/Unlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MutexServer).Unlock(ctx, req.(*MutexInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mutex_Getlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutexInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MutexServer).Getlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.mutex/Getlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MutexServer).Getlock(ctx, req.(*MutexInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mutex_Freelock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutexInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MutexServer).Freelock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.mutex/Freelock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MutexServer).Freelock(ctx, req.(*MutexInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mutex_Islock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutexInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MutexServer).Islock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.mutex/Islock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MutexServer).Islock(ctx, req.(*MutexInfo))
	}
	return interceptor(ctx, in, info, handler)
}

var _Mutex_serviceDesc = grpc.ServiceDesc{
	ServiceName: "main.mutex",
	HandlerType: (*MutexServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "lock",
			Handler:    _Mutex_Lock_Handler,
		},
		{
			MethodName: "unlock",
			Handler:    _Mutex_Unlock_Handler,
		},
		{
			MethodName: "getlock",
			Handler:    _Mutex_Getlock_Handler,
		},
		{
			MethodName: "freelock",
			Handler:    _Mutex_Freelock_Handler,
		},
		{
			MethodName: "islock",
			Handler:    _Mutex_Islock_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mutex.proto",
}

// RwmutexClient is the client API for Rwmutex service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RwmutexClient interface {
	Lock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error)
	Unlock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error)
	Rlock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error)
	Runlock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error)
	Getlock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*MutexInfo, error)
	Freelock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error)
	Islock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*ResponseTwo, error)
}

type rwmutexClient struct {
	cc grpc.ClientConnInterface
}

func NewRwmutexClient(cc grpc.ClientConnInterface) RwmutexClient {
	return &rwmutexClient{cc}
}

func (c *rwmutexClient) Lock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/main.rwmutex/lock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rwmutexClient) Unlock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/main.rwmutex/unlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rwmutexClient) Rlock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/main.rwmutex/rlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rwmutexClient) Runlock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/main.rwmutex/runlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rwmutexClient) Getlock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*MutexInfo, error) {
	out := new(MutexInfo)
	err := c.cc.Invoke(ctx, "/main.rwmutex/getlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rwmutexClient) Freelock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/main.rwmutex/freelock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rwmutexClient) Islock(ctx context.Context, in *MutexInfo, opts ...grpc.CallOption) (*ResponseTwo, error) {
	out := new(ResponseTwo)
	err := c.cc.Invoke(ctx, "/main.rwmutex/islock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RwmutexServer is the server API for Rwmutex service.
type RwmutexServer interface {
	Lock(context.Context, *MutexInfo) (*Response, error)
	Unlock(context.Context, *MutexInfo) (*Response, error)
	Rlock(context.Context, *MutexInfo) (*Response, error)
	Runlock(context.Context, *MutexInfo) (*Response, error)
	Getlock(context.Context, *MutexInfo) (*MutexInfo, error)
	Freelock(context.Context, *MutexInfo) (*Response, error)
	Islock(context.Context, *MutexInfo) (*ResponseTwo, error)
}

// UnimplementedRwmutexServer can be embedded to have forward compatible implementations.
type UnimplementedRwmutexServer struct {
}

func (*UnimplementedRwmutexServer) Lock(context.Context, *MutexInfo) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Lock not implemented")
}
func (*UnimplementedRwmutexServer) Unlock(context.Context, *MutexInfo) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unlock not implemented")
}
func (*UnimplementedRwmutexServer) Rlock(context.Context, *MutexInfo) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rlock not implemented")
}
func (*UnimplementedRwmutexServer) Runlock(context.Context, *MutexInfo) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Runlock not implemented")
}
func (*UnimplementedRwmutexServer) Getlock(context.Context, *MutexInfo) (*MutexInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Getlock not implemented")
}
func (*UnimplementedRwmutexServer) Freelock(context.Context, *MutexInfo) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Freelock not implemented")
}
func (*UnimplementedRwmutexServer) Islock(context.Context, *MutexInfo) (*ResponseTwo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Islock not implemented")
}

func RegisterRwmutexServer(s *grpc.Server, srv RwmutexServer) {
	s.RegisterService(&_Rwmutex_serviceDesc, srv)
}

func _Rwmutex_Lock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutexInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RwmutexServer).Lock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.rwmutex/Lock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RwmutexServer).Lock(ctx, req.(*MutexInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rwmutex_Unlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutexInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RwmutexServer).Unlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.rwmutex/Unlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RwmutexServer).Unlock(ctx, req.(*MutexInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rwmutex_Rlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutexInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RwmutexServer).Rlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.rwmutex/Rlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RwmutexServer).Rlock(ctx, req.(*MutexInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rwmutex_Runlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutexInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RwmutexServer).Runlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.rwmutex/Runlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RwmutexServer).Runlock(ctx, req.(*MutexInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rwmutex_Getlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutexInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RwmutexServer).Getlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.rwmutex/Getlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RwmutexServer).Getlock(ctx, req.(*MutexInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rwmutex_Freelock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutexInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RwmutexServer).Freelock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.rwmutex/Freelock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RwmutexServer).Freelock(ctx, req.(*MutexInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rwmutex_Islock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutexInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RwmutexServer).Islock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.rwmutex/Islock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RwmutexServer).Islock(ctx, req.(*MutexInfo))
	}
	return interceptor(ctx, in, info, handler)
}

var _Rwmutex_serviceDesc = grpc.ServiceDesc{
	ServiceName: "main.rwmutex",
	HandlerType: (*RwmutexServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "lock",
			Handler:    _Rwmutex_Lock_Handler,
		},
		{
			MethodName: "unlock",
			Handler:    _Rwmutex_Unlock_Handler,
		},
		{
			MethodName: "rlock",
			Handler:    _Rwmutex_Rlock_Handler,
		},
		{
			MethodName: "runlock",
			Handler:    _Rwmutex_Runlock_Handler,
		},
		{
			MethodName: "getlock",
			Handler:    _Rwmutex_Getlock_Handler,
		},
		{
			MethodName: "freelock",
			Handler:    _Rwmutex_Freelock_Handler,
		},
		{
			MethodName: "islock",
			Handler:    _Rwmutex_Islock_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mutex.proto",
}
