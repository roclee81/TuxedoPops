// Code generated by protoc-gen-go.
// source: TuxedoPopsStore.proto
// DO NOT EDIT!

/*
Package TuxedoPopsStore is a generated protocol buffer package.

It is generated from these files:
	TuxedoPopsStore.proto

It has these top-level messages:
	TuxedoPops
	OTX
	Ingedient
	Recipe
*/
package TuxedoPopsStore

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type TuxedoPops struct {
	Address string `protobuf:"bytes,1,opt,name=Address" json:"Address,omitempty"`
	Counter []byte `protobuf:"bytes,3,opt,name=Counter,proto3" json:"Counter,omitempty"`
	Outputs []*OTX `protobuf:"bytes,4,rep,name=Outputs" json:"Outputs,omitempty"`
}

func (m *TuxedoPops) Reset()         { *m = TuxedoPops{} }
func (m *TuxedoPops) String() string { return proto.CompactTextString(m) }
func (*TuxedoPops) ProtoMessage()    {}

func (m *TuxedoPops) GetOutputs() []*OTX {
	if m != nil {
		return m.Outputs
	}
	return nil
}

type OTX struct {
	Owners      [][]byte `protobuf:"bytes,1,rep,name=Owners,proto3" json:"Owners,omitempty"`
	Threshold   int64    `protobuf:"varint,2,opt,name=Threshold" json:"Threshold,omitempty"`
	Amount      int64    `protobuf:"varint,3,opt,name=Amount" json:"Amount,omitempty"`
	Type        string   `protobuf:"bytes,4,opt,name=Type" json:"Type,omitempty"`
	Data        string   `protobuf:"bytes,5,opt,name=Data" json:"Data,omitempty"`
	Recipe      string   `protobuf:"bytes,6,opt,name=Recipe" json:"Recipe,omitempty"`
	Creator     []byte   `protobuf:"bytes,7,opt,name=Creator,proto3" json:"Creator,omitempty"`
	PrevCounter []byte   `protobuf:"bytes,8,opt,name=PrevCounter,proto3" json:"PrevCounter,omitempty"`
}

func (m *OTX) Reset()         { *m = OTX{} }
func (m *OTX) String() string { return proto.CompactTextString(m) }
func (*OTX) ProtoMessage()    {}

type Ingedient struct {
	Numerator   int64  `protobuf:"varint,1,opt,name=Numerator" json:"Numerator,omitempty"`
	Denominator int64  `protobuf:"varint,2,opt,name=Denominator" json:"Denominator,omitempty"`
	Type        string `protobuf:"bytes,3,opt,name=Type" json:"Type,omitempty"`
}

func (m *Ingedient) Reset()         { *m = Ingedient{} }
func (m *Ingedient) String() string { return proto.CompactTextString(m) }
func (*Ingedient) ProtoMessage()    {}

type Recipe struct {
	CreatedType string       `protobuf:"bytes,1,opt,name=CreatedType" json:"CreatedType,omitempty"`
	Ingrediants []*Ingedient `protobuf:"bytes,2,rep,name=Ingrediants" json:"Ingrediants,omitempty"`
	Creator     []byte       `protobuf:"bytes,3,opt,name=Creator,proto3" json:"Creator,omitempty"`
}

func (m *Recipe) Reset()         { *m = Recipe{} }
func (m *Recipe) String() string { return proto.CompactTextString(m) }
func (*Recipe) ProtoMessage()    {}

func (m *Recipe) GetIngrediants() []*Ingedient {
	if m != nil {
		return m.Ingrediants
	}
	return nil
}
