// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: musechain/musecore/authority/chain_info.proto

package types

import (
	fmt "fmt"
	chains "github.com/RWAs-labs/muse/pkg/chains"
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

// ChainInfo contains static information about the chains
// This structure is used to dynamically update these info on a live network
// before hardcoding the values in a upgrade
type ChainInfo struct {
	Chains []chains.Chain `protobuf:"bytes,1,rep,name=chains,proto3" json:"chains"`
}

func (m *ChainInfo) Reset()         { *m = ChainInfo{} }
func (m *ChainInfo) String() string { return proto.CompactTextString(m) }
func (*ChainInfo) ProtoMessage()    {}
func (*ChainInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_725803a35e91141a, []int{0}
}
func (m *ChainInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ChainInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ChainInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ChainInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChainInfo.Merge(m, src)
}
func (m *ChainInfo) XXX_Size() int {
	return m.Size()
}
func (m *ChainInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ChainInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ChainInfo proto.InternalMessageInfo

func (m *ChainInfo) GetChains() []chains.Chain {
	if m != nil {
		return m.Chains
	}
	return nil
}

func init() {
	proto.RegisterType((*ChainInfo)(nil), "musechain.musecore.authority.ChainInfo")
}

func init() {
	proto.RegisterFile("musechain/musecore/authority/chain_info.proto", fileDescriptor_725803a35e91141a)
}

var fileDescriptor_725803a35e91141a = []byte{
	// 214 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0xcd, 0x2d, 0x2d, 0x4e,
	0x4d, 0xce, 0x48, 0xcc, 0xcc, 0xd3, 0x07, 0xb3, 0xf2, 0x8b, 0x52, 0xf5, 0x13, 0x4b, 0x4b, 0x32,
	0xf2, 0x8b, 0x32, 0x4b, 0x2a, 0xf5, 0xc1, 0x12, 0xf1, 0x99, 0x79, 0x69, 0xf9, 0x7a, 0x05, 0x45,
	0xf9, 0x25, 0xf9, 0x42, 0x32, 0x70, 0xe5, 0x7a, 0x30, 0xe5, 0x7a, 0x70, 0xe5, 0x52, 0x5a, 0x58,
	0x0c, 0x2b, 0xc8, 0x4e, 0x87, 0x18, 0x53, 0x0c, 0xa5, 0x20, 0x26, 0x49, 0x89, 0xa4, 0xe7, 0xa7,
	0xe7, 0x83, 0x99, 0xfa, 0x20, 0x16, 0x44, 0x54, 0xc9, 0x9f, 0x8b, 0xd3, 0x19, 0xa4, 0xca, 0x33,
	0x2f, 0x2d, 0x5f, 0xc8, 0x89, 0x8b, 0x0d, 0xa2, 0x45, 0x82, 0x51, 0x81, 0x59, 0x83, 0xdb, 0x48,
	0x45, 0x0f, 0x8b, 0xed, 0x05, 0xd9, 0xe9, 0x7a, 0x50, 0x83, 0xc1, 0x3a, 0x9d, 0x58, 0x4e, 0xdc,
	0x93, 0x67, 0x08, 0x82, 0xea, 0x74, 0x72, 0x3d, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6,
	0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39,
	0x86, 0x28, 0xed, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0xfd, 0xa0, 0x70,
	0xc7, 0x62, 0xdd, 0x9c, 0xc4, 0xa4, 0x62, 0xb0, 0xbb, 0xf5, 0x2b, 0x90, 0x82, 0xa0, 0xa4, 0xb2,
	0x20, 0xb5, 0x38, 0x89, 0x0d, 0xec, 0x3c, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf4, 0x30,
	0xe3, 0xf3, 0x2f, 0x01, 0x00, 0x00,
}

func (m *ChainInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ChainInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ChainInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Chains) > 0 {
		for iNdEx := len(m.Chains) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Chains[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintChainInfo(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintChainInfo(dAtA []byte, offset int, v uint64) int {
	offset -= sovChainInfo(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ChainInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Chains) > 0 {
		for _, e := range m.Chains {
			l = e.Size()
			n += 1 + l + sovChainInfo(uint64(l))
		}
	}
	return n
}

func sovChainInfo(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozChainInfo(x uint64) (n int) {
	return sovChainInfo(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ChainInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowChainInfo
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
			return fmt.Errorf("proto: ChainInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ChainInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chains", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChainInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthChainInfo
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthChainInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Chains = append(m.Chains, chains.Chain{})
			if err := m.Chains[len(m.Chains)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipChainInfo(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthChainInfo
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
func skipChainInfo(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowChainInfo
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
					return 0, ErrIntOverflowChainInfo
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
					return 0, ErrIntOverflowChainInfo
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
				return 0, ErrInvalidLengthChainInfo
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupChainInfo
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthChainInfo
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthChainInfo        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowChainInfo          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupChainInfo = fmt.Errorf("proto: unexpected end of group")
)
