// Code generated by protoc-gen-go.
// source: p.proto
// DO NOT EDIT!

/*
Package pktComms is a generated protocol buffer package.

It is generated from these files:
	p.proto

It has these top-level messages:
	AppMsg
	Hello
	Bye
	KeepAlive
	Ack
	Error
*/
package pktComms

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

// Opaque contents get copied through to the consensus layer
type AppMsg struct {
	AppNdx           *uint64 `protobuf:"varint,1,opt" json:"AppNdx,omitempty"`
	MsgN             *uint64 `protobuf:"varint,2,opt" json:"MsgN,omitempty"`
	ID               []byte  `protobuf:"bytes,3,opt" json:"ID,omitempty"`
	Contents         []byte  `protobuf:"bytes,4,opt" json:"Contents,omitempty"`
	Salt             []byte  `protobuf:"bytes,12,opt" json:"Salt,omitempty"`
	DigSig           []byte  `protobuf:"bytes,13,opt" json:"DigSig,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *AppMsg) Reset()         { *m = AppMsg{} }
func (m *AppMsg) String() string { return proto.CompactTextString(m) }
func (*AppMsg) ProtoMessage()    {}

func (m *AppMsg) GetAppNdx() uint64 {
	if m != nil && m.AppNdx != nil {
		return *m.AppNdx
	}
	return 0
}

func (m *AppMsg) GetMsgN() uint64 {
	if m != nil && m.MsgN != nil {
		return *m.MsgN
	}
	return 0
}

func (m *AppMsg) GetID() []byte {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *AppMsg) GetContents() []byte {
	if m != nil {
		return m.Contents
	}
	return nil
}

func (m *AppMsg) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

func (m *AppMsg) GetDigSig() []byte {
	if m != nil {
		return m.DigSig
	}
	return nil
}

// First message sent; initiates a communications cycle.  The first MsgN
// must be 1.  ID is the nodeID of the sender, a 20- (SHA1) or 32-byte
// (SHA256 or SHA3) value.  sigPubKey is the sender's RSA sig public key;
// commsPubKey is the sender's RSA comms public key.  The TCP address
// is conventionally a dotted quad followed by a port number (A.B.C.D:P).
//
type Hello struct {
	MsgN             *uint64 `protobuf:"varint,2,opt" json:"MsgN,omitempty"`
	ID               []byte  `protobuf:"bytes,3,opt" json:"ID,omitempty"`
	SigPubKey        []byte  `protobuf:"bytes,4,opt" json:"SigPubKey,omitempty"`
	CommsPubKey      []byte  `protobuf:"bytes,5,opt" json:"CommsPubKey,omitempty"`
	TCPAddr          *string `protobuf:"bytes,6,opt" json:"TCPAddr,omitempty"`
	Salt             []byte  `protobuf:"bytes,12,opt" json:"Salt,omitempty"`
	DigSig           []byte  `protobuf:"bytes,13,opt" json:"DigSig,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Hello) Reset()         { *m = Hello{} }
func (m *Hello) String() string { return proto.CompactTextString(m) }
func (*Hello) ProtoMessage()    {}

func (m *Hello) GetMsgN() uint64 {
	if m != nil && m.MsgN != nil {
		return *m.MsgN
	}
	return 0
}

func (m *Hello) GetID() []byte {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *Hello) GetSigPubKey() []byte {
	if m != nil {
		return m.SigPubKey
	}
	return nil
}

func (m *Hello) GetCommsPubKey() []byte {
	if m != nil {
		return m.CommsPubKey
	}
	return nil
}

func (m *Hello) GetTCPAddr() string {
	if m != nil && m.TCPAddr != nil {
		return *m.TCPAddr
	}
	return ""
}

func (m *Hello) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

func (m *Hello) GetDigSig() []byte {
	if m != nil {
		return m.DigSig
	}
	return nil
}

// Ends a communications cycle.
type Bye struct {
	MsgN             *uint64 `protobuf:"varint,2,opt" json:"MsgN,omitempty"`
	ID               []byte  `protobuf:"bytes,3,opt" json:"ID,omitempty"`
	Salt             []byte  `protobuf:"bytes,12,opt" json:"Salt,omitempty"`
	DigSig           []byte  `protobuf:"bytes,13,opt" json:"DigSig,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Bye) Reset()         { *m = Bye{} }
func (m *Bye) String() string { return proto.CompactTextString(m) }
func (*Bye) ProtoMessage()    {}

func (m *Bye) GetMsgN() uint64 {
	if m != nil && m.MsgN != nil {
		return *m.MsgN
	}
	return 0
}

func (m *Bye) GetID() []byte {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *Bye) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

func (m *Bye) GetDigSig() []byte {
	if m != nil {
		return m.DigSig
	}
	return nil
}

// Sent at regular intervals.
type KeepAlive struct {
	MsgN             *uint64 `protobuf:"varint,2,opt" json:"MsgN,omitempty"`
	ID               []byte  `protobuf:"bytes,3,opt" json:"ID,omitempty"`
	Salt             []byte  `protobuf:"bytes,12,opt" json:"Salt,omitempty"`
	DigSig           []byte  `protobuf:"bytes,13,opt" json:"DigSig,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *KeepAlive) Reset()         { *m = KeepAlive{} }
func (m *KeepAlive) String() string { return proto.CompactTextString(m) }
func (*KeepAlive) ProtoMessage()    {}

func (m *KeepAlive) GetMsgN() uint64 {
	if m != nil && m.MsgN != nil {
		return *m.MsgN
	}
	return 0
}

func (m *KeepAlive) GetID() []byte {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *KeepAlive) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

func (m *KeepAlive) GetDigSig() []byte {
	if m != nil {
		return m.DigSig
	}
	return nil
}

type Ack struct {
	MsgN             *uint64 `protobuf:"varint,2,opt" json:"MsgN,omitempty"`
	ID               []byte  `protobuf:"bytes,3,opt" json:"ID,omitempty"`
	YourMsgN         *uint64 `protobuf:"varint,4,opt" json:"YourMsgN,omitempty"`
	YourID           []byte  `protobuf:"bytes,5,opt" json:"YourID,omitempty"`
	Salt             []byte  `protobuf:"bytes,12,opt" json:"Salt,omitempty"`
	DigSig           []byte  `protobuf:"bytes,13,opt" json:"DigSig,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Ack) Reset()         { *m = Ack{} }
func (m *Ack) String() string { return proto.CompactTextString(m) }
func (*Ack) ProtoMessage()    {}

func (m *Ack) GetMsgN() uint64 {
	if m != nil && m.MsgN != nil {
		return *m.MsgN
	}
	return 0
}

func (m *Ack) GetID() []byte {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *Ack) GetYourMsgN() uint64 {
	if m != nil && m.YourMsgN != nil {
		return *m.YourMsgN
	}
	return 0
}

func (m *Ack) GetYourID() []byte {
	if m != nil {
		return m.YourID
	}
	return nil
}

func (m *Ack) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

func (m *Ack) GetDigSig() []byte {
	if m != nil {
		return m.DigSig
	}
	return nil
}

type Error struct {
	MsgN             *uint64 `protobuf:"varint,2,opt" json:"MsgN,omitempty"`
	ID               []byte  `protobuf:"bytes,3,opt" json:"ID,omitempty"`
	YourMsgN         *uint64 `protobuf:"varint,4,opt" json:"YourMsgN,omitempty"`
	YourID           []byte  `protobuf:"bytes,5,opt" json:"YourID,omitempty"`
	ErrCode          *uint64 `protobuf:"varint,6,opt" json:"ErrCode,omitempty"`
	ErrDesc          *string `protobuf:"bytes,7,opt" json:"ErrDesc,omitempty"`
	Salt             []byte  `protobuf:"bytes,12,opt" json:"Salt,omitempty"`
	DigSig           []byte  `protobuf:"bytes,13,opt" json:"DigSig,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}

func (m *Error) GetMsgN() uint64 {
	if m != nil && m.MsgN != nil {
		return *m.MsgN
	}
	return 0
}

func (m *Error) GetID() []byte {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *Error) GetYourMsgN() uint64 {
	if m != nil && m.YourMsgN != nil {
		return *m.YourMsgN
	}
	return 0
}

func (m *Error) GetYourID() []byte {
	if m != nil {
		return m.YourID
	}
	return nil
}

func (m *Error) GetErrCode() uint64 {
	if m != nil && m.ErrCode != nil {
		return *m.ErrCode
	}
	return 0
}

func (m *Error) GetErrDesc() string {
	if m != nil && m.ErrDesc != nil {
		return *m.ErrDesc
	}
	return ""
}

func (m *Error) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

func (m *Error) GetDigSig() []byte {
	if m != nil {
		return m.DigSig
	}
	return nil
}

func init() {
}
