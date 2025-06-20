// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: musechain/musecore/observer/observer.proto

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

type ObservationType int32

const (
	ObservationType_EmptyObserverType ObservationType = 0
	ObservationType_InboundTx         ObservationType = 1
	ObservationType_OutboundTx        ObservationType = 2
	ObservationType_TSSKeyGen         ObservationType = 3
	ObservationType_TSSKeySign        ObservationType = 4
)

var ObservationType_name = map[int32]string{
	0: "EmptyObserverType",
	1: "InboundTx",
	2: "OutboundTx",
	3: "TSSKeyGen",
	4: "TSSKeySign",
}

var ObservationType_value = map[string]int32{
	"EmptyObserverType": 0,
	"InboundTx":         1,
	"OutboundTx":        2,
	"TSSKeyGen":         3,
	"TSSKeySign":        4,
}

func (x ObservationType) String() string {
	return proto.EnumName(ObservationType_name, int32(x))
}

func (ObservationType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c1bf442c6827e5bb, []int{0}
}

type ObserverUpdateReason int32

const (
	ObserverUpdateReason_Undefined   ObserverUpdateReason = 0
	ObserverUpdateReason_Tombstoned  ObserverUpdateReason = 1
	ObserverUpdateReason_AdminUpdate ObserverUpdateReason = 2
)

var ObserverUpdateReason_name = map[int32]string{
	0: "Undefined",
	1: "Tombstoned",
	2: "AdminUpdate",
}

var ObserverUpdateReason_value = map[string]int32{
	"Undefined":   0,
	"Tombstoned":  1,
	"AdminUpdate": 2,
}

func (x ObserverUpdateReason) String() string {
	return proto.EnumName(ObserverUpdateReason_name, int32(x))
}

func (ObserverUpdateReason) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c1bf442c6827e5bb, []int{1}
}

type ObserverSet struct {
	ObserverList []string `protobuf:"bytes,1,rep,name=observer_list,json=observerList,proto3" json:"observer_list,omitempty"`
}

func (m *ObserverSet) Reset()         { *m = ObserverSet{} }
func (m *ObserverSet) String() string { return proto.CompactTextString(m) }
func (*ObserverSet) ProtoMessage()    {}
func (*ObserverSet) Descriptor() ([]byte, []int) {
	return fileDescriptor_c1bf442c6827e5bb, []int{0}
}
func (m *ObserverSet) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ObserverSet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ObserverSet.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ObserverSet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ObserverSet.Merge(m, src)
}
func (m *ObserverSet) XXX_Size() int {
	return m.Size()
}
func (m *ObserverSet) XXX_DiscardUnknown() {
	xxx_messageInfo_ObserverSet.DiscardUnknown(m)
}

var xxx_messageInfo_ObserverSet proto.InternalMessageInfo

func (m *ObserverSet) GetObserverList() []string {
	if m != nil {
		return m.ObserverList
	}
	return nil
}

type LastObserverCount struct {
	Count            uint64 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	LastChangeHeight int64  `protobuf:"varint,2,opt,name=last_change_height,json=lastChangeHeight,proto3" json:"last_change_height,omitempty"`
}

func (m *LastObserverCount) Reset()         { *m = LastObserverCount{} }
func (m *LastObserverCount) String() string { return proto.CompactTextString(m) }
func (*LastObserverCount) ProtoMessage()    {}
func (*LastObserverCount) Descriptor() ([]byte, []int) {
	return fileDescriptor_c1bf442c6827e5bb, []int{1}
}
func (m *LastObserverCount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LastObserverCount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LastObserverCount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LastObserverCount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LastObserverCount.Merge(m, src)
}
func (m *LastObserverCount) XXX_Size() int {
	return m.Size()
}
func (m *LastObserverCount) XXX_DiscardUnknown() {
	xxx_messageInfo_LastObserverCount.DiscardUnknown(m)
}

var xxx_messageInfo_LastObserverCount proto.InternalMessageInfo

func (m *LastObserverCount) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *LastObserverCount) GetLastChangeHeight() int64 {
	if m != nil {
		return m.LastChangeHeight
	}
	return 0
}

func init() {
	proto.RegisterEnum("musechain.musecore.observer.ObservationType", ObservationType_name, ObservationType_value)
	proto.RegisterEnum("musechain.musecore.observer.ObserverUpdateReason", ObserverUpdateReason_name, ObserverUpdateReason_value)
	proto.RegisterType((*ObserverSet)(nil), "musechain.musecore.observer.ObserverSet")
	proto.RegisterType((*LastObserverCount)(nil), "musechain.musecore.observer.LastObserverCount")
}

func init() {
	proto.RegisterFile("musechain/musecore/observer/observer.proto", fileDescriptor_c1bf442c6827e5bb)
}

var fileDescriptor_c1bf442c6827e5bb = []byte{
	// 374 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x91, 0xd1, 0x8a, 0xda, 0x40,
	0x14, 0x86, 0x33, 0x6a, 0x0b, 0x8e, 0xb5, 0xc6, 0xc1, 0x82, 0x58, 0x08, 0x62, 0x6f, 0x24, 0xb4,
	0x09, 0xb4, 0x4f, 0x60, 0x6d, 0x69, 0x97, 0x15, 0x84, 0x44, 0x11, 0xf6, 0x46, 0x26, 0xc9, 0x6c,
	0x32, 0x90, 0xcc, 0x84, 0xcc, 0x64, 0x31, 0x6f, 0xb1, 0x0f, 0xb1, 0x17, 0xfb, 0x28, 0x7b, 0xe9,
	0xe5, 0x5e, 0x2e, 0xfa, 0x22, 0xcb, 0x24, 0x3b, 0xee, 0xdd, 0xff, 0x9f, 0xff, 0xff, 0xce, 0xb9,
	0x38, 0xd0, 0xce, 0x4a, 0x41, 0xc2, 0x04, 0x53, 0xe6, 0xd6, 0x8a, 0x17, 0xc4, 0xe5, 0x81, 0x20,
	0xc5, 0x1d, 0x29, 0x2e, 0xc2, 0xc9, 0x0b, 0x2e, 0x39, 0xfa, 0x7a, 0xe9, 0x3a, 0xba, 0xeb, 0xe8,
	0xca, 0x64, 0x14, 0xf3, 0x98, 0xd7, 0x3d, 0x57, 0xa9, 0x06, 0x99, 0xfd, 0x84, 0xbd, 0xf5, 0x5b,
	0xc3, 0x27, 0x12, 0x7d, 0x83, 0x7d, 0x0d, 0xec, 0x53, 0x2a, 0xe4, 0x18, 0x4c, 0xdb, 0xf3, 0xae,
	0xf7, 0x49, 0x0f, 0x57, 0x54, 0xc8, 0xd9, 0x0e, 0x0e, 0x57, 0x58, 0x48, 0xcd, 0x2d, 0x79, 0xc9,
	0x24, 0x1a, 0xc1, 0x0f, 0xa1, 0x12, 0x63, 0x30, 0x05, 0xf3, 0x8e, 0xd7, 0x18, 0xf4, 0x1d, 0xa2,
	0x14, 0x0b, 0xb9, 0x0f, 0x13, 0xcc, 0x62, 0xb2, 0x4f, 0x08, 0x8d, 0x13, 0x39, 0x6e, 0x4d, 0xc1,
	0xbc, 0xed, 0x99, 0x2a, 0x59, 0xd6, 0xc1, 0xff, 0x7a, 0x6e, 0xa7, 0x70, 0xd0, 0x2c, 0xc5, 0x92,
	0x72, 0xb6, 0xa9, 0x72, 0x82, 0xbe, 0xc0, 0xe1, 0xdf, 0x2c, 0x97, 0x95, 0x3e, 0xa6, 0x86, 0xa6,
	0x81, 0xfa, 0xb0, 0x7b, 0xc5, 0x02, 0x5e, 0xb2, 0x68, 0x73, 0x30, 0x01, 0xfa, 0x0c, 0xe1, 0xba,
	0x94, 0xda, 0xb7, 0x54, 0xbc, 0xf1, 0xfd, 0x6b, 0x52, 0xfd, 0x23, 0xcc, 0x6c, 0xab, 0xb8, 0xb1,
	0x3e, 0x8d, 0x99, 0xd9, 0x99, 0x74, 0x1e, 0x1f, 0x2c, 0x60, 0xaf, 0xe0, 0x48, 0x6f, 0xdd, 0xe6,
	0x11, 0x96, 0xc4, 0x23, 0x58, 0x70, 0xa6, 0xe0, 0x2d, 0x8b, 0xc8, 0x2d, 0x65, 0x24, 0x32, 0x8d,
	0x1a, 0xe6, 0x59, 0x20, 0x24, 0x57, 0x1e, 0xa0, 0x01, 0xec, 0x2d, 0xa2, 0x8c, 0xb2, 0x86, 0x31,
	0x5b, 0xcd, 0xb6, 0xdf, 0x7f, 0x9e, 0x4e, 0x16, 0x38, 0x9e, 0x2c, 0xf0, 0x72, 0xb2, 0xc0, 0xfd,
	0xd9, 0x32, 0x8e, 0x67, 0xcb, 0x78, 0x3e, 0x5b, 0xc6, 0x8d, 0x1d, 0x53, 0x99, 0x94, 0x81, 0x13,
	0xf2, 0xcc, 0xf5, 0x76, 0x0b, 0xf1, 0x23, 0xc5, 0x81, 0xa8, 0x9f, 0xe9, 0x1e, 0xde, 0x5f, 0x29,
	0xab, 0x9c, 0x88, 0xe0, 0x63, 0xfd, 0x95, 0x5f, 0xaf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x76, 0x0c,
	0x41, 0x45, 0xf6, 0x01, 0x00, 0x00,
}

func (m *ObserverSet) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ObserverSet) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ObserverSet) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ObserverList) > 0 {
		for iNdEx := len(m.ObserverList) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ObserverList[iNdEx])
			copy(dAtA[i:], m.ObserverList[iNdEx])
			i = encodeVarintObserver(dAtA, i, uint64(len(m.ObserverList[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *LastObserverCount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LastObserverCount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LastObserverCount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.LastChangeHeight != 0 {
		i = encodeVarintObserver(dAtA, i, uint64(m.LastChangeHeight))
		i--
		dAtA[i] = 0x10
	}
	if m.Count != 0 {
		i = encodeVarintObserver(dAtA, i, uint64(m.Count))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintObserver(dAtA []byte, offset int, v uint64) int {
	offset -= sovObserver(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ObserverSet) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ObserverList) > 0 {
		for _, s := range m.ObserverList {
			l = len(s)
			n += 1 + l + sovObserver(uint64(l))
		}
	}
	return n
}

func (m *LastObserverCount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Count != 0 {
		n += 1 + sovObserver(uint64(m.Count))
	}
	if m.LastChangeHeight != 0 {
		n += 1 + sovObserver(uint64(m.LastChangeHeight))
	}
	return n
}

func sovObserver(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozObserver(x uint64) (n int) {
	return sovObserver(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ObserverSet) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowObserver
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
			return fmt.Errorf("proto: ObserverSet: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ObserverSet: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ObserverList", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowObserver
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
				return ErrInvalidLengthObserver
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthObserver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ObserverList = append(m.ObserverList, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipObserver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthObserver
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
func (m *LastObserverCount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowObserver
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
			return fmt.Errorf("proto: LastObserverCount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LastObserverCount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Count", wireType)
			}
			m.Count = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowObserver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Count |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastChangeHeight", wireType)
			}
			m.LastChangeHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowObserver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastChangeHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipObserver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthObserver
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
func skipObserver(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowObserver
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
					return 0, ErrIntOverflowObserver
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
					return 0, ErrIntOverflowObserver
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
				return 0, ErrInvalidLengthObserver
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupObserver
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthObserver
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthObserver        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowObserver          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupObserver = fmt.Errorf("proto: unexpected end of group")
)
