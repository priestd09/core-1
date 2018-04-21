// Code generated by protoc-gen-go. DO NOT EDIT.
// source: types/api/task.proto

package types

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ExecuteTaskRequest struct {
	Service *ProtoService `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
	Task    string        `protobuf:"bytes,2,opt,name=task" json:"task,omitempty"`
	Data    string        `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *ExecuteTaskRequest) Reset()                    { *m = ExecuteTaskRequest{} }
func (m *ExecuteTaskRequest) String() string            { return proto.CompactTextString(m) }
func (*ExecuteTaskRequest) ProtoMessage()               {}
func (*ExecuteTaskRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *ExecuteTaskRequest) GetService() *ProtoService {
	if m != nil {
		return m.Service
	}
	return nil
}

func (m *ExecuteTaskRequest) GetTask() string {
	if m != nil {
		return m.Task
	}
	return ""
}

func (m *ExecuteTaskRequest) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

type ListenTaskRequest struct {
	Service *ProtoService `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
}

func (m *ListenTaskRequest) Reset()                    { *m = ListenTaskRequest{} }
func (m *ListenTaskRequest) String() string            { return proto.CompactTextString(m) }
func (*ListenTaskRequest) ProtoMessage()               {}
func (*ListenTaskRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *ListenTaskRequest) GetService() *ProtoService {
	if m != nil {
		return m.Service
	}
	return nil
}

type TaskReply struct {
	Task string `protobuf:"bytes,1,opt,name=task" json:"task,omitempty"`
	Data string `protobuf:"bytes,2,opt,name=data" json:"data,omitempty"`
}

func (m *TaskReply) Reset()                    { *m = TaskReply{} }
func (m *TaskReply) String() string            { return proto.CompactTextString(m) }
func (*TaskReply) ProtoMessage()               {}
func (*TaskReply) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

func (m *TaskReply) GetTask() string {
	if m != nil {
		return m.Task
	}
	return ""
}

func (m *TaskReply) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterType((*ExecuteTaskRequest)(nil), "types.ExecuteTaskRequest")
	proto.RegisterType((*ListenTaskRequest)(nil), "types.ListenTaskRequest")
	proto.RegisterType((*TaskReply)(nil), "types.TaskReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Task service

type TaskClient interface {
	Execute(ctx context.Context, in *ExecuteTaskRequest, opts ...grpc.CallOption) (*TaskReply, error)
	Listen(ctx context.Context, in *ListenTaskRequest, opts ...grpc.CallOption) (Task_ListenClient, error)
}

type taskClient struct {
	cc *grpc.ClientConn
}

func NewTaskClient(cc *grpc.ClientConn) TaskClient {
	return &taskClient{cc}
}

func (c *taskClient) Execute(ctx context.Context, in *ExecuteTaskRequest, opts ...grpc.CallOption) (*TaskReply, error) {
	out := new(TaskReply)
	err := grpc.Invoke(ctx, "/types.Task/Execute", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskClient) Listen(ctx context.Context, in *ListenTaskRequest, opts ...grpc.CallOption) (Task_ListenClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Task_serviceDesc.Streams[0], c.cc, "/types.Task/Listen", opts...)
	if err != nil {
		return nil, err
	}
	x := &taskListenClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Task_ListenClient interface {
	Recv() (*TaskReply, error)
	grpc.ClientStream
}

type taskListenClient struct {
	grpc.ClientStream
}

func (x *taskListenClient) Recv() (*TaskReply, error) {
	m := new(TaskReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Task service

type TaskServer interface {
	Execute(context.Context, *ExecuteTaskRequest) (*TaskReply, error)
	Listen(*ListenTaskRequest, Task_ListenServer) error
}

func RegisterTaskServer(s *grpc.Server, srv TaskServer) {
	s.RegisterService(&_Task_serviceDesc, srv)
}

func _Task_Execute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServer).Execute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/types.Task/Execute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServer).Execute(ctx, req.(*ExecuteTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Task_Listen_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListenTaskRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TaskServer).Listen(m, &taskListenServer{stream})
}

type Task_ListenServer interface {
	Send(*TaskReply) error
	grpc.ServerStream
}

type taskListenServer struct {
	grpc.ServerStream
}

func (x *taskListenServer) Send(m *TaskReply) error {
	return x.ServerStream.SendMsg(m)
}

var _Task_serviceDesc = grpc.ServiceDesc{
	ServiceName: "types.Task",
	HandlerType: (*TaskServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Execute",
			Handler:    _Task_Execute_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Listen",
			Handler:       _Task_Listen_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "types/api/task.proto",
}

func init() { proto.RegisterFile("types/api/task.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 217 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x29, 0xa9, 0x2c, 0x48,
	0x2d, 0xd6, 0x4f, 0x2c, 0xc8, 0xd4, 0x2f, 0x49, 0x2c, 0xce, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x62, 0x05, 0x8b, 0x4a, 0x09, 0x43, 0x24, 0x8b, 0x53, 0x8b, 0xca, 0x32, 0x93, 0x53, 0x21,
	0x72, 0x4a, 0xd9, 0x5c, 0x42, 0xae, 0x15, 0xa9, 0xc9, 0xa5, 0x25, 0xa9, 0x21, 0x89, 0xc5, 0xd9,
	0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0xba, 0x5c, 0xec, 0x50, 0x65, 0x12, 0x8c, 0x0a,
	0x8c, 0x1a, 0xdc, 0x46, 0xc2, 0x7a, 0x60, 0xcd, 0x7a, 0x01, 0x20, 0x4d, 0xc1, 0x10, 0xa9, 0x20,
	0x98, 0x1a, 0x21, 0x21, 0x2e, 0x16, 0x90, 0x75, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x60,
	0x36, 0x48, 0x2c, 0x25, 0xb1, 0x24, 0x51, 0x82, 0x19, 0x22, 0x06, 0x62, 0x2b, 0x39, 0x71, 0x09,
	0xfa, 0x64, 0x16, 0x97, 0xa4, 0xe6, 0x91, 0x6f, 0x97, 0x92, 0x31, 0x17, 0x27, 0x44, 0x77, 0x41,
	0x4e, 0x25, 0xdc, 0x62, 0x46, 0x2c, 0x16, 0x33, 0x21, 0x2c, 0x36, 0xaa, 0xe2, 0x62, 0x01, 0x69,
	0x12, 0xb2, 0xe0, 0x62, 0x87, 0xfa, 0x56, 0x48, 0x12, 0x6a, 0x0b, 0xa6, 0xef, 0xa5, 0x04, 0xa0,
	0x52, 0x70, 0x7b, 0x94, 0x18, 0x84, 0x2c, 0xb8, 0xd8, 0x20, 0x4e, 0x17, 0x92, 0x80, 0xca, 0x62,
	0xf8, 0x04, 0x9b, 0x3e, 0x03, 0xc6, 0x24, 0x36, 0x70, 0x40, 0x1b, 0x03, 0x02, 0x00, 0x00, 0xff,
	0xff, 0x24, 0xe8, 0x92, 0xd7, 0x9c, 0x01, 0x00, 0x00,
}
