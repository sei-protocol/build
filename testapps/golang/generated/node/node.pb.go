// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v4.25.0
// source: node.proto

package node

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

// TransactionHeader message definition
type TransactionHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Properties map[string]string `protobuf:"bytes,1,rep,name=properties,proto3" json:"properties,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	EdgeNode   string            `protobuf:"bytes,2,opt,name=edge_node,json=edgeNode,proto3" json:"edge_node,omitempty"`
	Signature  *Signature        `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *TransactionHeader) Reset() {
	*x = TransactionHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransactionHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionHeader) ProtoMessage() {}

func (x *TransactionHeader) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionHeader.ProtoReflect.Descriptor instead.
func (*TransactionHeader) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{0}
}

func (x *TransactionHeader) GetProperties() map[string]string {
	if x != nil {
		return x.Properties
	}
	return nil
}

func (x *TransactionHeader) GetEdgeNode() string {
	if x != nil {
		return x.EdgeNode
	}
	return ""
}

func (x *TransactionHeader) GetSignature() *Signature {
	if x != nil {
		return x.Signature
	}
	return nil
}

// Signature message definition
type Signature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Algorithm string `protobuf:"bytes,1,opt,name=algorithm,proto3" json:"algorithm,omitempty"`
	Signature string `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *Signature) Reset() {
	*x = Signature{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Signature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Signature) ProtoMessage() {}

func (x *Signature) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Signature.ProtoReflect.Descriptor instead.
func (*Signature) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{1}
}

func (x *Signature) GetAlgorithm() string {
	if x != nil {
		return x.Algorithm
	}
	return ""
}

func (x *Signature) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

// Transaction message definition
type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash    string             `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Payload []byte             `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	GasUsed int64              `protobuf:"varint,3,opt,name=gasUsed,proto3" json:"gasUsed,omitempty"`
	Header  *TransactionHeader `protobuf:"bytes,4,opt,name=header,proto3" json:"header,omitempty"`
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{2}
}

func (x *Transaction) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *Transaction) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *Transaction) GetGasUsed() int64 {
	if x != nil {
		return x.GasUsed
	}
	return 0
}

func (x *Transaction) GetHeader() *TransactionHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

// TransactionResponse message definition
type TransactionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash    string `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Success bool   `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *TransactionResponse) Reset() {
	*x = TransactionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransactionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionResponse) ProtoMessage() {}

func (x *TransactionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionResponse.ProtoReflect.Descriptor instead.
func (*TransactionResponse) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{3}
}

func (x *TransactionResponse) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *TransactionResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *TransactionResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// Block header message
type BlockHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockNumber int64  `protobuf:"varint,1,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`
	Timestamp   int64  `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	TotalGas    int64  `protobuf:"varint,3,opt,name=total_gas,json=totalGas,proto3" json:"total_gas,omitempty"`
	StreamNode  string `protobuf:"bytes,4,opt,name=stream_node,json=streamNode,proto3" json:"stream_node,omitempty"`
	EdgeCount   int64  `protobuf:"varint,5,opt,name=edge_count,json=edgeCount,proto3" json:"edge_count,omitempty"`
}

func (x *BlockHeader) Reset() {
	*x = BlockHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockHeader) ProtoMessage() {}

func (x *BlockHeader) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockHeader.ProtoReflect.Descriptor instead.
func (*BlockHeader) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{4}
}

func (x *BlockHeader) GetBlockNumber() int64 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *BlockHeader) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *BlockHeader) GetTotalGas() int64 {
	if x != nil {
		return x.TotalGas
	}
	return 0
}

func (x *BlockHeader) GetStreamNode() string {
	if x != nil {
		return x.StreamNode
	}
	return ""
}

func (x *BlockHeader) GetEdgeCount() int64 {
	if x != nil {
		return x.EdgeCount
	}
	return 0
}

// Block message that can include headers or bodies
type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header   *BlockHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Hashes   []string     `protobuf:"bytes,2,rep,name=hashes,proto3" json:"hashes,omitempty"`
	Payloads [][]byte     `protobuf:"bytes,3,rep,name=payloads,proto3" json:"payloads,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{5}
}

func (x *Block) GetHeader() *BlockHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *Block) GetHashes() []string {
	if x != nil {
		return x.Hashes
	}
	return nil
}

func (x *Block) GetPayloads() [][]byte {
	if x != nil {
		return x.Payloads
	}
	return nil
}

// Request message for GetBlocks
type SubscribeBlocksRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartBlock      int64 `protobuf:"varint,1,opt,name=start_block,json=startBlock,proto3" json:"start_block,omitempty"` // If 0, start from the latest block
	IncludePayloads bool  `protobuf:"varint,2,opt,name=include_payloads,json=includePayloads,proto3" json:"include_payloads,omitempty"`
	Latest          bool  `protobuf:"varint,3,opt,name=latest,proto3" json:"latest,omitempty"` // start from latest (ignore start_block)
}

func (x *SubscribeBlocksRequest) Reset() {
	*x = SubscribeBlocksRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeBlocksRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeBlocksRequest) ProtoMessage() {}

func (x *SubscribeBlocksRequest) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeBlocksRequest.ProtoReflect.Descriptor instead.
func (*SubscribeBlocksRequest) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{6}
}

func (x *SubscribeBlocksRequest) GetStartBlock() int64 {
	if x != nil {
		return x.StartBlock
	}
	return 0
}

func (x *SubscribeBlocksRequest) GetIncludePayloads() bool {
	if x != nil {
		return x.IncludePayloads
	}
	return false
}

func (x *SubscribeBlocksRequest) GetLatest() bool {
	if x != nil {
		return x.Latest
	}
	return false
}

// Response message for GetBlocks
type GetBlocksResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Blocks []*Block `protobuf:"bytes,1,rep,name=blocks,proto3" json:"blocks,omitempty"`
}

func (x *GetBlocksResponse) Reset() {
	*x = GetBlocksResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBlocksResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBlocksResponse) ProtoMessage() {}

func (x *GetBlocksResponse) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBlocksResponse.ProtoReflect.Descriptor instead.
func (*GetBlocksResponse) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{7}
}

func (x *GetBlocksResponse) GetBlocks() []*Block {
	if x != nil {
		return x.Blocks
	}
	return nil
}

var File_node_proto protoreflect.FileDescriptor

var file_node_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x73, 0x65,
	0x69, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x73, 0x65, 0x69, 0x73, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x22, 0x93, 0x02, 0x0a, 0x11, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x5d,
	0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x3d, 0x2e, 0x73, 0x65, 0x69, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c,
	0x2e, 0x73, 0x65, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x0a, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x12, 0x1b, 0x0a,
	0x09, 0x65, 0x64, 0x67, 0x65, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x65, 0x64, 0x67, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x43, 0x0a, 0x09, 0x73, 0x69,
	0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e,
	0x73, 0x65, 0x69, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x73, 0x65, 0x69, 0x73,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x1a,
	0x3d, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x47,
	0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x61,
	0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x69,
	0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x9c, 0x01, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x70,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x67, 0x61, 0x73, 0x55, 0x73, 0x65, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x67, 0x61, 0x73, 0x55, 0x73, 0x65, 0x64, 0x12,
	0x45, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2d, 0x2e, 0x73, 0x65, 0x69, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x73, 0x65,
	0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x22, 0x5d, 0x0a, 0x13, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73,
	0x68, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xab, 0x01, 0x0a, 0x0b, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f,
	0x67, 0x61, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x74, 0x6f, 0x74, 0x61, 0x6c,
	0x47, 0x61, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x6e, 0x6f,
	0x64, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x4e, 0x6f, 0x64, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x64, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x65, 0x64, 0x67, 0x65, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x22, 0x7c, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x3f, 0x0a, 0x06,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x73,
	0x65, 0x69, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x73, 0x65, 0x69, 0x73, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x16, 0x0a,
	0x06, 0x68, 0x61, 0x73, 0x68, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x68,
	0x61, 0x73, 0x68, 0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x08, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x73, 0x22, 0x7c, 0x0a, 0x16, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x29, 0x0a, 0x10,
	0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x5f, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0f, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x50,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x61, 0x74, 0x65, 0x73,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x22,
	0x4e, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x73, 0x65, 0x69, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6c, 0x2e, 0x73, 0x65, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x6e, 0x6f, 0x64,
	0x65, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x32,
	0x80, 0x01, 0x0a, 0x0c, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x70, 0x0a, 0x10, 0x53, 0x65, 0x6e, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x27, 0x2e, 0x73, 0x65, 0x69, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6c, 0x2e, 0x73, 0x65, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x6e, 0x6f, 0x64,
	0x65, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x2f, 0x2e,
	0x73, 0x65, 0x69, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x73, 0x65, 0x69, 0x73,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01,
	0x30, 0x01, 0x32, 0x74, 0x0a, 0x06, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x12, 0x6a, 0x0a, 0x0f,
	0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x12,
	0x32, 0x2e, 0x73, 0x65, 0x69, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x73, 0x65,
	0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x73, 0x65, 0x69, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6c, 0x2e, 0x73, 0x65, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x30, 0x01, 0x42, 0x08, 0x5a, 0x06, 0x2e, 0x2f, 0x6e, 0x6f,
	0x64, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_node_proto_rawDescOnce sync.Once
	file_node_proto_rawDescData = file_node_proto_rawDesc
)

func file_node_proto_rawDescGZIP() []byte {
	file_node_proto_rawDescOnce.Do(func() {
		file_node_proto_rawDescData = protoimpl.X.CompressGZIP(file_node_proto_rawDescData)
	})
	return file_node_proto_rawDescData
}

var file_node_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_node_proto_goTypes = []any{
	(*TransactionHeader)(nil),      // 0: seiprotocol.seistream.node.TransactionHeader
	(*Signature)(nil),              // 1: seiprotocol.seistream.node.Signature
	(*Transaction)(nil),            // 2: seiprotocol.seistream.node.Transaction
	(*TransactionResponse)(nil),    // 3: seiprotocol.seistream.node.TransactionResponse
	(*BlockHeader)(nil),            // 4: seiprotocol.seistream.node.BlockHeader
	(*Block)(nil),                  // 5: seiprotocol.seistream.node.Block
	(*SubscribeBlocksRequest)(nil), // 6: seiprotocol.seistream.node.SubscribeBlocksRequest
	(*GetBlocksResponse)(nil),      // 7: seiprotocol.seistream.node.GetBlocksResponse
	nil,                            // 8: seiprotocol.seistream.node.TransactionHeader.PropertiesEntry
}
var file_node_proto_depIdxs = []int32{
	8, // 0: seiprotocol.seistream.node.TransactionHeader.properties:type_name -> seiprotocol.seistream.node.TransactionHeader.PropertiesEntry
	1, // 1: seiprotocol.seistream.node.TransactionHeader.signature:type_name -> seiprotocol.seistream.node.Signature
	0, // 2: seiprotocol.seistream.node.Transaction.header:type_name -> seiprotocol.seistream.node.TransactionHeader
	4, // 3: seiprotocol.seistream.node.Block.header:type_name -> seiprotocol.seistream.node.BlockHeader
	5, // 4: seiprotocol.seistream.node.GetBlocksResponse.blocks:type_name -> seiprotocol.seistream.node.Block
	2, // 5: seiprotocol.seistream.node.Transactions.SendTransactions:input_type -> seiprotocol.seistream.node.Transaction
	6, // 6: seiprotocol.seistream.node.Blocks.SubscribeBlocks:input_type -> seiprotocol.seistream.node.SubscribeBlocksRequest
	3, // 7: seiprotocol.seistream.node.Transactions.SendTransactions:output_type -> seiprotocol.seistream.node.TransactionResponse
	5, // 8: seiprotocol.seistream.node.Blocks.SubscribeBlocks:output_type -> seiprotocol.seistream.node.Block
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_node_proto_init() }
func file_node_proto_init() {
	if File_node_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_node_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*TransactionHeader); i {
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
		file_node_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Signature); i {
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
		file_node_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Transaction); i {
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
		file_node_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*TransactionResponse); i {
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
		file_node_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*BlockHeader); i {
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
		file_node_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*Block); i {
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
		file_node_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*SubscribeBlocksRequest); i {
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
		file_node_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*GetBlocksResponse); i {
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
			RawDescriptor: file_node_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_node_proto_goTypes,
		DependencyIndexes: file_node_proto_depIdxs,
		MessageInfos:      file_node_proto_msgTypes,
	}.Build()
	File_node_proto = out.File
	file_node_proto_rawDesc = nil
	file_node_proto_goTypes = nil
	file_node_proto_depIdxs = nil
}
