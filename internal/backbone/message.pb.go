// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.3
// source: internal/backbone/message.proto

package backbone

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

type Handshake struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	Token   []byte `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *Handshake) Reset() {
	*x = Handshake{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_backbone_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Handshake) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Handshake) ProtoMessage() {}

func (x *Handshake) ProtoReflect() protoreflect.Message {
	mi := &file_internal_backbone_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Handshake.ProtoReflect.Descriptor instead.
func (*Handshake) Descriptor() ([]byte, []int) {
	return file_internal_backbone_message_proto_rawDescGZIP(), []int{0}
}

func (x *Handshake) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Handshake) GetToken() []byte {
	if x != nil {
		return x.Token
	}
	return nil
}

type HandshakeRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error   string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *HandshakeRes) Reset() {
	*x = HandshakeRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_backbone_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandshakeRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandshakeRes) ProtoMessage() {}

func (x *HandshakeRes) ProtoReflect() protoreflect.Message {
	mi := &file_internal_backbone_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandshakeRes.ProtoReflect.Descriptor instead.
func (*HandshakeRes) Descriptor() ([]byte, []int) {
	return file_internal_backbone_message_proto_rawDescGZIP(), []int{1}
}

func (x *HandshakeRes) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *HandshakeRes) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type Ping struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Time *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *Ping) Reset() {
	*x = Ping{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_backbone_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ping) ProtoMessage() {}

func (x *Ping) ProtoReflect() protoreflect.Message {
	mi := &file_internal_backbone_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ping.ProtoReflect.Descriptor instead.
func (*Ping) Descriptor() ([]byte, []int) {
	return file_internal_backbone_message_proto_rawDescGZIP(), []int{2}
}

func (x *Ping) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

type Pong struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string                 `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Time  *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *Pong) Reset() {
	*x = Pong{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_backbone_message_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pong) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pong) ProtoMessage() {}

func (x *Pong) ProtoReflect() protoreflect.Message {
	mi := &file_internal_backbone_message_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pong.ProtoReflect.Descriptor instead.
func (*Pong) Descriptor() ([]byte, []int) {
	return file_internal_backbone_message_proto_rawDescGZIP(), []int{3}
}

func (x *Pong) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *Pong) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

type RegisterEdge struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EdgeName        string `protobuf:"bytes,1,opt,name=edge_name,json=edgeName,proto3" json:"edge_name,omitempty"`
	EdgeType        string `protobuf:"bytes,2,opt,name=edge_type,json=edgeType,proto3" json:"edge_type,omitempty"`
	LoadBalancerKey string `protobuf:"bytes,3,opt,name=load_balancer_key,json=loadBalancerKey,proto3" json:"load_balancer_key,omitempty"`
	EdgeServerAddr  string `protobuf:"bytes,4,opt,name=edge_server_addr,json=edgeServerAddr,proto3" json:"edge_server_addr,omitempty"`
}

func (x *RegisterEdge) Reset() {
	*x = RegisterEdge{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_backbone_message_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterEdge) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterEdge) ProtoMessage() {}

func (x *RegisterEdge) ProtoReflect() protoreflect.Message {
	mi := &file_internal_backbone_message_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterEdge.ProtoReflect.Descriptor instead.
func (*RegisterEdge) Descriptor() ([]byte, []int) {
	return file_internal_backbone_message_proto_rawDescGZIP(), []int{4}
}

func (x *RegisterEdge) GetEdgeName() string {
	if x != nil {
		return x.EdgeName
	}
	return ""
}

func (x *RegisterEdge) GetEdgeType() string {
	if x != nil {
		return x.EdgeType
	}
	return ""
}

func (x *RegisterEdge) GetLoadBalancerKey() string {
	if x != nil {
		return x.LoadBalancerKey
	}
	return ""
}

func (x *RegisterEdge) GetEdgeServerAddr() string {
	if x != nil {
		return x.EdgeServerAddr
	}
	return ""
}

type RegisterEdgeRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error        string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	EdgeRunnerId string `protobuf:"bytes,2,opt,name=edge_runner_id,json=edgeRunnerId,proto3" json:"edge_runner_id,omitempty"`
	EdgeName     string `protobuf:"bytes,3,opt,name=edge_name,json=edgeName,proto3" json:"edge_name,omitempty"`
}

func (x *RegisterEdgeRes) Reset() {
	*x = RegisterEdgeRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_backbone_message_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterEdgeRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterEdgeRes) ProtoMessage() {}

func (x *RegisterEdgeRes) ProtoReflect() protoreflect.Message {
	mi := &file_internal_backbone_message_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterEdgeRes.ProtoReflect.Descriptor instead.
func (*RegisterEdgeRes) Descriptor() ([]byte, []int) {
	return file_internal_backbone_message_proto_rawDescGZIP(), []int{5}
}

func (x *RegisterEdgeRes) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *RegisterEdgeRes) GetEdgeRunnerId() string {
	if x != nil {
		return x.EdgeRunnerId
	}
	return ""
}

func (x *RegisterEdgeRes) GetEdgeName() string {
	if x != nil {
		return x.EdgeName
	}
	return ""
}

type UnregisterEdge struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EdgeRunnerId string `protobuf:"bytes,1,opt,name=edge_runner_id,json=edgeRunnerId,proto3" json:"edge_runner_id,omitempty"`
}

func (x *UnregisterEdge) Reset() {
	*x = UnregisterEdge{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_backbone_message_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnregisterEdge) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnregisterEdge) ProtoMessage() {}

func (x *UnregisterEdge) ProtoReflect() protoreflect.Message {
	mi := &file_internal_backbone_message_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnregisterEdge.ProtoReflect.Descriptor instead.
func (*UnregisterEdge) Descriptor() ([]byte, []int) {
	return file_internal_backbone_message_proto_rawDescGZIP(), []int{6}
}

func (x *UnregisterEdge) GetEdgeRunnerId() string {
	if x != nil {
		return x.EdgeRunnerId
	}
	return ""
}

type UnregisterEdgeRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error        string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	EdgeRunnerId string `protobuf:"bytes,2,opt,name=edge_runner_id,json=edgeRunnerId,proto3" json:"edge_runner_id,omitempty"`
	EdgeName     string `protobuf:"bytes,3,opt,name=edge_name,json=edgeName,proto3" json:"edge_name,omitempty"`
}

func (x *UnregisterEdgeRes) Reset() {
	*x = UnregisterEdgeRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_backbone_message_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnregisterEdgeRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnregisterEdgeRes) ProtoMessage() {}

func (x *UnregisterEdgeRes) ProtoReflect() protoreflect.Message {
	mi := &file_internal_backbone_message_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnregisterEdgeRes.ProtoReflect.Descriptor instead.
func (*UnregisterEdgeRes) Descriptor() ([]byte, []int) {
	return file_internal_backbone_message_proto_rawDescGZIP(), []int{7}
}

func (x *UnregisterEdgeRes) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *UnregisterEdgeRes) GetEdgeRunnerId() string {
	if x != nil {
		return x.EdgeRunnerId
	}
	return ""
}

func (x *UnregisterEdgeRes) GetEdgeName() string {
	if x != nil {
		return x.EdgeName
	}
	return ""
}

var File_internal_backbone_message_proto protoreflect.FileDescriptor

var file_internal_backbone_message_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x62,
	0x6f, 0x6e, 0x65, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3b, 0x0a, 0x09, 0x48,
	0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x3e, 0x0a, 0x0c, 0x48, 0x61, 0x6e, 0x64,
	0x73, 0x68, 0x61, 0x6b, 0x65, 0x52, 0x65, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x18,
	0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x36, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67,
	0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65,
	0x22, 0x4c, 0x0a, 0x04, 0x50, 0x6f, 0x6e, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x2e,
	0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x9e,
	0x01, 0x0a, 0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x45, 0x64, 0x67, 0x65, 0x12,
	0x1b, 0x0a, 0x09, 0x65, 0x64, 0x67, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x65, 0x64, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x65, 0x64, 0x67, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x65, 0x64, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2a, 0x0a, 0x11, 0x6c, 0x6f, 0x61,
	0x64, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x72, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63,
	0x65, 0x72, 0x4b, 0x65, 0x79, 0x12, 0x28, 0x0a, 0x10, 0x65, 0x64, 0x67, 0x65, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0e, 0x65, 0x64, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x22,
	0x6a, 0x0a, 0x0f, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x45, 0x64, 0x67, 0x65, 0x52,
	0x65, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x24, 0x0a, 0x0e, 0x65, 0x64, 0x67, 0x65,
	0x5f, 0x72, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x65, 0x64, 0x67, 0x65, 0x52, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1b,
	0x0a, 0x09, 0x65, 0x64, 0x67, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x65, 0x64, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x36, 0x0a, 0x0e, 0x55,
	0x6e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x45, 0x64, 0x67, 0x65, 0x12, 0x24, 0x0a,
	0x0e, 0x65, 0x64, 0x67, 0x65, 0x5f, 0x72, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x64, 0x67, 0x65, 0x52, 0x75, 0x6e, 0x6e, 0x65,
	0x72, 0x49, 0x64, 0x22, 0x6c, 0x0a, 0x11, 0x55, 0x6e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x45, 0x64, 0x67, 0x65, 0x52, 0x65, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x24,
	0x0a, 0x0e, 0x65, 0x64, 0x67, 0x65, 0x5f, 0x72, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x64, 0x67, 0x65, 0x52, 0x75, 0x6e, 0x6e,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x64, 0x67, 0x65, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x64, 0x67, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x42, 0x15, 0x5a, 0x13, 0x2e, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x62, 0x61, 0x63, 0x6b, 0x62, 0x6f, 0x6e, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_backbone_message_proto_rawDescOnce sync.Once
	file_internal_backbone_message_proto_rawDescData = file_internal_backbone_message_proto_rawDesc
)

func file_internal_backbone_message_proto_rawDescGZIP() []byte {
	file_internal_backbone_message_proto_rawDescOnce.Do(func() {
		file_internal_backbone_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_backbone_message_proto_rawDescData)
	})
	return file_internal_backbone_message_proto_rawDescData
}

var file_internal_backbone_message_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_internal_backbone_message_proto_goTypes = []interface{}{
	(*Handshake)(nil),             // 0: message.Handshake
	(*HandshakeRes)(nil),          // 1: message.HandshakeRes
	(*Ping)(nil),                  // 2: message.Ping
	(*Pong)(nil),                  // 3: message.Pong
	(*RegisterEdge)(nil),          // 4: message.RegisterEdge
	(*RegisterEdgeRes)(nil),       // 5: message.RegisterEdgeRes
	(*UnregisterEdge)(nil),        // 6: message.UnregisterEdge
	(*UnregisterEdgeRes)(nil),     // 7: message.UnregisterEdgeRes
	(*timestamppb.Timestamp)(nil), // 8: google.protobuf.Timestamp
}
var file_internal_backbone_message_proto_depIdxs = []int32{
	8, // 0: message.Ping.time:type_name -> google.protobuf.Timestamp
	8, // 1: message.Pong.time:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_internal_backbone_message_proto_init() }
func file_internal_backbone_message_proto_init() {
	if File_internal_backbone_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_backbone_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Handshake); i {
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
		file_internal_backbone_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandshakeRes); i {
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
		file_internal_backbone_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ping); i {
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
		file_internal_backbone_message_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pong); i {
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
		file_internal_backbone_message_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterEdge); i {
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
		file_internal_backbone_message_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterEdgeRes); i {
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
		file_internal_backbone_message_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnregisterEdge); i {
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
		file_internal_backbone_message_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnregisterEdgeRes); i {
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
			RawDescriptor: file_internal_backbone_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_backbone_message_proto_goTypes,
		DependencyIndexes: file_internal_backbone_message_proto_depIdxs,
		MessageInfos:      file_internal_backbone_message_proto_msgTypes,
	}.Build()
	File_internal_backbone_message_proto = out.File
	file_internal_backbone_message_proto_rawDesc = nil
	file_internal_backbone_message_proto_goTypes = nil
	file_internal_backbone_message_proto_depIdxs = nil
}
