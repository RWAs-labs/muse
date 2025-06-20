// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: musechain/musecore/observer/tss.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type TSS struct {
	TssPubkey           string   `protobuf:"bytes,3,opt,name=tss_pubkey,json=tssPubkey,proto3" json:"tss_pubkey,omitempty"`
	TssParticipantList  []string `protobuf:"bytes,4,rep,name=tss_participant_list,json=tssParticipantList,proto3" json:"tss_participant_list,omitempty"`
	OperatorAddressList []string `protobuf:"bytes,5,rep,name=operator_address_list,json=operatorAddressList,proto3" json:"operator_address_list,omitempty"`
	FinalizedMuseHeight int64    `protobuf:"varint,6,opt,name=finalizedMuseHeight,proto3" json:"finalizedMuseHeight,omitempty"`
	KeyGenMuseHeight    int64    `protobuf:"varint,7,opt,name=keyGenMuseHeight,proto3" json:"keyGenMuseHeight,omitempty"`
}

func (m *TSS) Reset()         { *m = TSS{} }
func (m *TSS) String() string { return proto.CompactTextString(m) }
func (*TSS) ProtoMessage()    {}
func (*TSS) Descriptor() ([]byte, []int) {
	return fileDescriptor_7d057a6fd30177de, []int{0}
}
func (m *TSS) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TSS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TSS.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TSS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TSS.Merge(m, src)
}
func (m *TSS) XXX_Size() int {
	return m.Size()
}
func (m *TSS) XXX_DiscardUnknown() {
	xxx_messageInfo_TSS.DiscardUnknown(m)
}

var xxx_messageInfo_TSS proto.InternalMessageInfo

func (m *TSS) GetTssPubkey() string {
	if m != nil {
		return m.TssPubkey
	}
	return ""
}

func (m *TSS) GetTssParticipantList() []string {
	if m != nil {
		return m.TssParticipantList
	}
	return nil
}

func (m *TSS) GetOperatorAddressList() []string {
	if m != nil {
		return m.OperatorAddressList
	}
	return nil
}

func (m *TSS) GetFinalizedMuseHeight() int64 {
	if m != nil {
		return m.FinalizedMuseHeight
	}
	return 0
}

func (m *TSS) GetKeyGenMuseHeight() int64 {
	if m != nil {
		return m.KeyGenMuseHeight
	}
	return 0
}

func init() {
	proto.RegisterType((*TSS)(nil), "musechain.musecore.observer.TSS")
}

func init() {
	proto.RegisterFile("musechain/musecore/observer/tss.proto", fileDescriptor_7d057a6fd30177de)
}

var fileDescriptor_7d057a6fd30177de = []byte{
	// 295 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xc1, 0x4a, 0xf3, 0x40,
	0x14, 0x46, 0x3b, 0xe4, 0xff, 0x2b, 0x9d, 0x95, 0x4c, 0x2b, 0x14, 0xc5, 0xa1, 0x08, 0x42, 0x29,
	0x98, 0x29, 0xfa, 0x04, 0x15, 0x41, 0x17, 0x0a, 0x92, 0x0a, 0x82, 0x9b, 0x32, 0x69, 0xaf, 0xc9,
	0xd0, 0x34, 0x13, 0xe6, 0x4e, 0xc4, 0xf8, 0x14, 0x3e, 0x96, 0xcb, 0x2e, 0x5d, 0x4a, 0xb2, 0xf1,
	0x31, 0xa4, 0x13, 0xa2, 0x05, 0xdd, 0x5d, 0xee, 0x39, 0x67, 0xf3, 0xd1, 0xe3, 0x55, 0x8e, 0x30,
	0x8f, 0xa5, 0x4a, 0x85, 0xbb, 0xb4, 0x01, 0xa1, 0x43, 0x04, 0xf3, 0x04, 0x46, 0x58, 0x44, 0x3f,
	0x33, 0xda, 0x6a, 0x76, 0xf0, 0xad, 0xf9, 0x8d, 0xe6, 0x37, 0xda, 0x7e, 0x2f, 0xd2, 0x91, 0x76,
	0x9e, 0xd8, 0x5c, 0x75, 0x72, 0xf4, 0x49, 0xa8, 0x77, 0x37, 0x9d, 0xb2, 0x43, 0x4a, 0x2d, 0xe2,
	0x2c, 0xcb, 0xc3, 0x25, 0x14, 0x7d, 0x6f, 0x40, 0x86, 0x9d, 0xa0, 0x63, 0x11, 0x6f, 0xdd, 0x83,
	0x8d, 0x69, 0xcf, 0x61, 0x69, 0xac, 0x9a, 0xab, 0x4c, 0xa6, 0x76, 0x96, 0x28, 0xb4, 0xfd, 0x7f,
	0x03, 0x6f, 0xd8, 0x09, 0xd8, 0x46, 0xfc, 0x41, 0xd7, 0x0a, 0x2d, 0x3b, 0xa5, 0x7b, 0x3a, 0x03,
	0x23, 0xad, 0x36, 0x33, 0xb9, 0x58, 0x18, 0x40, 0xac, 0x93, 0xff, 0x2e, 0xe9, 0x36, 0x70, 0x52,
	0x33, 0xd7, 0x8c, 0x69, 0xf7, 0x51, 0xa5, 0x32, 0x51, 0x2f, 0xb0, 0xb8, 0xc9, 0x11, 0xae, 0x40,
	0x45, 0xb1, 0xed, 0xb7, 0x07, 0x64, 0xe8, 0x05, 0x7f, 0x21, 0x36, 0xa2, 0xbb, 0x4b, 0x28, 0x2e,
	0x21, 0xdd, 0xd2, 0x77, 0x9c, 0xfe, 0xeb, 0x7f, 0x7e, 0xf1, 0x56, 0x72, 0xb2, 0x2e, 0x39, 0xf9,
	0x28, 0x39, 0x79, 0xad, 0x78, 0x6b, 0x5d, 0xf1, 0xd6, 0x7b, 0xc5, 0x5b, 0x0f, 0xa3, 0x48, 0xd9,
	0x38, 0x0f, 0xfd, 0xb9, 0x5e, 0x89, 0xe0, 0x7e, 0x82, 0x27, 0x89, 0x0c, 0xd1, 0x2d, 0x2d, 0x9e,
	0xb7, 0x76, 0x2e, 0x32, 0xc0, 0xb0, 0xed, 0x76, 0x3b, 0xfb, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x1a,
	0x8f, 0xb4, 0x22, 0x93, 0x01, 0x00, 0x00,
}

func (m *TSS) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TSS) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TSS) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.KeyGenMuseHeight != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.KeyGenMuseHeight))
		i--
		dAtA[i] = 0x38
	}
	if m.FinalizedMuseHeight != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.FinalizedMuseHeight))
		i--
		dAtA[i] = 0x30
	}
	if len(m.OperatorAddressList) > 0 {
		for iNdEx := len(m.OperatorAddressList) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.OperatorAddressList[iNdEx])
			copy(dAtA[i:], m.OperatorAddressList[iNdEx])
			i = encodeVarintTss(dAtA, i, uint64(len(m.OperatorAddressList[iNdEx])))
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.TssParticipantList) > 0 {
		for iNdEx := len(m.TssParticipantList) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.TssParticipantList[iNdEx])
			copy(dAtA[i:], m.TssParticipantList[iNdEx])
			i = encodeVarintTss(dAtA, i, uint64(len(m.TssParticipantList[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.TssPubkey) > 0 {
		i -= len(m.TssPubkey)
		copy(dAtA[i:], m.TssPubkey)
		i = encodeVarintTss(dAtA, i, uint64(len(m.TssPubkey)))
		i--
		dAtA[i] = 0x1a
	}
	return len(dAtA) - i, nil
}

func encodeVarintTss(dAtA []byte, offset int, v uint64) int {
	offset -= sovTss(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TSS) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.TssPubkey)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	if len(m.TssParticipantList) > 0 {
		for _, s := range m.TssParticipantList {
			l = len(s)
			n += 1 + l + sovTss(uint64(l))
		}
	}
	if len(m.OperatorAddressList) > 0 {
		for _, s := range m.OperatorAddressList {
			l = len(s)
			n += 1 + l + sovTss(uint64(l))
		}
	}
	if m.FinalizedMuseHeight != 0 {
		n += 1 + sovTss(uint64(m.FinalizedMuseHeight))
	}
	if m.KeyGenMuseHeight != 0 {
		n += 1 + sovTss(uint64(m.KeyGenMuseHeight))
	}
	return n
}

func sovTss(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTss(x uint64) (n int) {
	return sovTss(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TSS) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TSS: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TSS: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TssPubkey", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TssPubkey = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TssParticipantList", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TssParticipantList = append(m.TssParticipantList, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OperatorAddressList", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OperatorAddressList = append(m.OperatorAddressList, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FinalizedMuseHeight", wireType)
			}
			m.FinalizedMuseHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FinalizedMuseHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field KeyGenMuseHeight", wireType)
			}
			m.KeyGenMuseHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.KeyGenMuseHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTss(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTss
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTss
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTss
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTss
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTss
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTss
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTss        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTss          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTss = fmt.Errorf("proto: unexpected end of group")
)
