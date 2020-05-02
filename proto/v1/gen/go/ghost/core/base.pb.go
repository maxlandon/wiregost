// Wiregost - Golang Exploitation Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.8.0
// source: ghost/core/base.proto

package corepb

import (
	transport "../gen/go/transport"
	route "../gen/go/transport/route"
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Info - Send back multiple informations to the server at once.
type Info struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string                 `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Username    string                 `protobuf:"bytes,3,opt,name=Username,proto3" json:"Username,omitempty"`
	UID         string                 `protobuf:"bytes,4,opt,name=UID,proto3" json:"UID,omitempty"`
	GID         string                 `protobuf:"bytes,5,opt,name=GID,proto3" json:"GID,omitempty"`
	OS          string                 `protobuf:"bytes,6,opt,name=OS,proto3" json:"OS,omitempty"`
	Arch        string                 `protobuf:"bytes,7,opt,name=Arch,proto3" json:"Arch,omitempty"`
	PID         int32                  `protobuf:"varint,8,opt,name=PID,proto3" json:"PID,omitempty"`
	Filename    string                 `protobuf:"bytes,9,opt,name=Filename,proto3" json:"Filename,omitempty"`
	Version     string                 `protobuf:"bytes,10,opt,name=Version,proto3" json:"Version,omitempty"`
	WorkspaceID uint32                 `protobuf:"varint,12,opt,name=WorkspaceID,proto3" json:"WorkspaceID,omitempty"`
	Transports  []*transport.Transport `protobuf:"bytes,14,rep,name=Transports,proto3" json:"Transports,omitempty"` // Available C2 transports
	Interfaces  []*NetInterface        `protobuf:"bytes,15,rep,name=Interfaces,proto3" json:"Interfaces,omitempty"` // Network
	Netstat     []*SocketTabEntry      `protobuf:"bytes,18,rep,name=Netstat,proto3" json:"Netstat,omitempty"`       // Maybe we can use this to automatically devise some better ports to use for routing/etc
	Routes      []*route.Route         `protobuf:"bytes,19,rep,name=Routes,proto3" json:"Routes,omitempty"`         // We might send current route listeners, in case.
}

func (x *Info) Reset() {
	*x = Info{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ghost_core_base_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Info) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Info) ProtoMessage() {}

func (x *Info) ProtoReflect() protoreflect.Message {
	mi := &file_ghost_core_base_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Info.ProtoReflect.Descriptor instead.
func (*Info) Descriptor() ([]byte, []int) {
	return file_ghost_core_base_proto_rawDescGZIP(), []int{0}
}

func (x *Info) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Info) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Info) GetUID() string {
	if x != nil {
		return x.UID
	}
	return ""
}

func (x *Info) GetGID() string {
	if x != nil {
		return x.GID
	}
	return ""
}

func (x *Info) GetOS() string {
	if x != nil {
		return x.OS
	}
	return ""
}

func (x *Info) GetArch() string {
	if x != nil {
		return x.Arch
	}
	return ""
}

func (x *Info) GetPID() int32 {
	if x != nil {
		return x.PID
	}
	return 0
}

func (x *Info) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *Info) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Info) GetWorkspaceID() uint32 {
	if x != nil {
		return x.WorkspaceID
	}
	return 0
}

func (x *Info) GetTransports() []*transport.Transport {
	if x != nil {
		return x.Transports
	}
	return nil
}

func (x *Info) GetInterfaces() []*NetInterface {
	if x != nil {
		return x.Interfaces
	}
	return nil
}

func (x *Info) GetNetstat() []*SocketTabEntry {
	if x != nil {
		return x.Netstat
	}
	return nil
}

func (x *Info) GetRoutes() []*route.Route {
	if x != nil {
		return x.Routes
	}
	return nil
}

// Register - Implant calls back to C2 Server and sends its information.
type Register struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Information *Info `protobuf:"bytes,1,opt,name=Information,proto3" json:"Information,omitempty"`
}

func (x *Register) Reset() {
	*x = Register{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ghost_core_base_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Register) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Register) ProtoMessage() {}

func (x *Register) ProtoReflect() protoreflect.Message {
	mi := &file_ghost_core_base_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Register.ProtoReflect.Descriptor instead.
func (*Register) Descriptor() ([]byte, []int) {
	return file_ghost_core_base_proto_rawDescGZIP(), []int{1}
}

func (x *Register) GetInformation() *Info {
	if x != nil {
		return x.Information
	}
	return nil
}

// Ping - Test ghost connection
type Ping struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nonce int32 `protobuf:"varint,1,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
}

func (x *Ping) Reset() {
	*x = Ping{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ghost_core_base_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ping) ProtoMessage() {}

func (x *Ping) ProtoReflect() protoreflect.Message {
	mi := &file_ghost_core_base_proto_msgTypes[2]
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
	return file_ghost_core_base_proto_rawDescGZIP(), []int{2}
}

func (x *Ping) GetNonce() int32 {
	if x != nil {
		return x.Nonce
	}
	return 0
}

// KillRequest - Kills the ghost executable and connection
type KillRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GhostID uint32 `protobuf:"varint,1,opt,name=GhostID,proto3" json:"GhostID,omitempty"`
	Force   bool   `protobuf:"varint,2,opt,name=Force,proto3" json:"Force,omitempty"`
}

func (x *KillRequest) Reset() {
	*x = KillRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ghost_core_base_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KillRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KillRequest) ProtoMessage() {}

func (x *KillRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ghost_core_base_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KillRequest.ProtoReflect.Descriptor instead.
func (*KillRequest) Descriptor() ([]byte, []int) {
	return file_ghost_core_base_proto_rawDescGZIP(), []int{3}
}

func (x *KillRequest) GetGhostID() uint32 {
	if x != nil {
		return x.GhostID
	}
	return 0
}

func (x *KillRequest) GetForce() bool {
	if x != nil {
		return x.Force
	}
	return false
}

// Kill - Sends back status on implant kill
type Kill struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=Success,proto3" json:"Success,omitempty"`
	Warning string `protobuf:"bytes,2,opt,name=Warning,proto3" json:"Warning,omitempty"` // If force kill, will kill anyway but send a warning instead of an error
	Err     string `protobuf:"bytes,3,opt,name=Err,proto3" json:"Err,omitempty"`         // Sends an error if kill was not force and there things still running
}

func (x *Kill) Reset() {
	*x = Kill{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ghost_core_base_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Kill) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Kill) ProtoMessage() {}

func (x *Kill) ProtoReflect() protoreflect.Message {
	mi := &file_ghost_core_base_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Kill.ProtoReflect.Descriptor instead.
func (*Kill) Descriptor() ([]byte, []int) {
	return file_ghost_core_base_proto_rawDescGZIP(), []int{4}
}

func (x *Kill) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *Kill) GetWarning() string {
	if x != nil {
		return x.Warning
	}
	return ""
}

func (x *Kill) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

var File_ghost_core_base_proto protoreflect.FileDescriptor

var file_ghost_core_base_proto_rawDesc = []byte{
	0x0a, 0x15, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x62, 0x61, 0x73,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x1a, 0x19, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14,
	0x67, 0x68, 0x6f, 0x73, 0x74, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x6e, 0x65, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2f,
	0x72, 0x6f, 0x75, 0x74, 0x65, 0x2f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xbe, 0x03, 0x0a, 0x04, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x55, 0x49,
	0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x55, 0x49, 0x44, 0x12, 0x10, 0x0a, 0x03,
	0x47, 0x49, 0x44, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x47, 0x49, 0x44, 0x12, 0x0e,
	0x0a, 0x02, 0x4f, 0x53, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x4f, 0x53, 0x12, 0x12,
	0x0a, 0x04, 0x41, 0x72, 0x63, 0x68, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x41, 0x72,
	0x63, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x50, 0x49, 0x44, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x03, 0x50, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x57, 0x6f,
	0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x44, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x0b, 0x57, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x44, 0x12, 0x34, 0x0a, 0x0a,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x18, 0x0e, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x0a, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72,
	0x74, 0x73, 0x12, 0x38, 0x0a, 0x0a, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73,
	0x18, 0x0f, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x4e, 0x65, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65,
	0x52, 0x0a, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x12, 0x34, 0x0a, 0x07,
	0x4e, 0x65, 0x74, 0x73, 0x74, 0x61, 0x74, 0x18, 0x12, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x68, 0x6f, 0x73, 0x74, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x53, 0x6f, 0x63, 0x6b, 0x65,
	0x74, 0x54, 0x61, 0x62, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x4e, 0x65, 0x74, 0x73, 0x74,
	0x61, 0x74, 0x12, 0x2e, 0x0a, 0x06, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x18, 0x13, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x72,
	0x6f, 0x75, 0x74, 0x65, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x06, 0x52, 0x6f, 0x75, 0x74,
	0x65, 0x73, 0x22, 0x3e, 0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x32,
	0x0a, 0x0b, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x2e, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0b, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x22, 0x1c, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x4e, 0x6f,
	0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x4e, 0x6f, 0x6e, 0x63, 0x65,
	0x22, 0x3d, 0x0a, 0x0b, 0x4b, 0x69, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x18, 0x0a, 0x07, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x07, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x46, 0x6f, 0x72,
	0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x46, 0x6f, 0x72, 0x63, 0x65, 0x22,
	0x4c, 0x0a, 0x04, 0x4b, 0x69, 0x6c, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x12, 0x18, 0x0a, 0x07, 0x57, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x57, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x45,
	0x72, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x45, 0x72, 0x72, 0x42, 0x08, 0x5a,
	0x06, 0x63, 0x6f, 0x72, 0x65, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ghost_core_base_proto_rawDescOnce sync.Once
	file_ghost_core_base_proto_rawDescData = file_ghost_core_base_proto_rawDesc
)

func file_ghost_core_base_proto_rawDescGZIP() []byte {
	file_ghost_core_base_proto_rawDescOnce.Do(func() {
		file_ghost_core_base_proto_rawDescData = protoimpl.X.CompressGZIP(file_ghost_core_base_proto_rawDescData)
	})
	return file_ghost_core_base_proto_rawDescData
}

var file_ghost_core_base_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_ghost_core_base_proto_goTypes = []interface{}{
	(*Info)(nil),                // 0: ghost.core.Info
	(*Register)(nil),            // 1: ghost.core.Register
	(*Ping)(nil),                // 2: ghost.core.Ping
	(*KillRequest)(nil),         // 3: ghost.core.KillRequest
	(*Kill)(nil),                // 4: ghost.core.Kill
	(*transport.Transport)(nil), // 5: transport.Transport
	(*NetInterface)(nil),        // 6: ghost.core.NetInterface
	(*SocketTabEntry)(nil),      // 7: ghost.core.SocketTabEntry
	(*route.Route)(nil),         // 8: transport.route.Route
}
var file_ghost_core_base_proto_depIdxs = []int32{
	5, // 0: ghost.core.Info.Transports:type_name -> transport.Transport
	6, // 1: ghost.core.Info.Interfaces:type_name -> ghost.core.NetInterface
	7, // 2: ghost.core.Info.Netstat:type_name -> ghost.core.SocketTabEntry
	8, // 3: ghost.core.Info.Routes:type_name -> transport.route.Route
	0, // 4: ghost.core.Register.Information:type_name -> ghost.core.Info
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_ghost_core_base_proto_init() }
func file_ghost_core_base_proto_init() {
	if File_ghost_core_base_proto != nil {
		return
	}
	file_ghost_core_net_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_ghost_core_base_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Info); i {
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
		file_ghost_core_base_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Register); i {
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
		file_ghost_core_base_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_ghost_core_base_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KillRequest); i {
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
		file_ghost_core_base_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Kill); i {
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
			RawDescriptor: file_ghost_core_base_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ghost_core_base_proto_goTypes,
		DependencyIndexes: file_ghost_core_base_proto_depIdxs,
		MessageInfos:      file_ghost_core_base_proto_msgTypes,
	}.Build()
	File_ghost_core_base_proto = out.File
	file_ghost_core_base_proto_rawDesc = nil
	file_ghost_core_base_proto_goTypes = nil
	file_ghost_core_base_proto_depIdxs = nil
}