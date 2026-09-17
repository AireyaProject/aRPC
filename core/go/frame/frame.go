package frame

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type FrameType uint8

const (
	Request  FrameType = 0x00
	Response FrameType = 0x01
	Error    FrameType = 0x02
	Ping     FrameType = 0x03
	Pong     FrameType = 0x04
)

const MagicVersion = "Aireyav1"
const HeaderBaseSize = 20

type Frame struct {
	FrameLength uint32
	Type        FrameType
	Flags       uint8
	StreamID    uint32
	Metadata    map[string]string
	Payload     []byte
}

func (f *Frame) Encode() []byte {
	metaLen := uint16(0)
	for k, v := range f.Metadata {
		metaLen += 1 + uint16(len(k)) + 2 + uint16(len(v))
	}

	totalLen := (HeaderBaseSize - 4) + uint32(metaLen) + uint32(len(f.Payload))
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, totalLen)
	buf.WriteString(MagicVersion)
	buf.WriteByte(byte(f.Type))
	buf.WriteByte(f.Flags)
	binary.Write(buf, binary.BigEndian, f.StreamID)
	binary.Write(buf, binary.BigEndian, metaLen)

	for k, v := range f.Metadata {
		buf.WriteByte(byte(len(k)))
		buf.WriteString(k)
		binary.Write(buf, binary.BigEndian, uint16(len(v)))
		buf.WriteString(v)
	}
	buf.Write(f.Payload)
	return buf.Bytes()
}

func Decode(data []byte) (*Frame, int, error) {
	if len(data) < 4 {
		return nil, 0, nil // need more data
	}
	
	frameLen := binary.BigEndian.Uint32(data[0:4])
	if len(data) < int(4+frameLen) {
		return nil, 0, nil // incomplete
	}

	if frameLen < HeaderBaseSize-4 {
		return nil, -1, errors.New("invalid frame size")
	}

	offset := 4
	magic := string(data[offset : offset+8])
	if magic != MagicVersion {
		return nil, -1, errors.New("invalid magic")
	}
	offset += 8

	f := &Frame{
		FrameLength: frameLen,
		Type:        FrameType(data[offset]),
		Flags:       data[offset+1],
		Metadata:    make(map[string]string),
	}
	offset += 2

	f.StreamID = binary.BigEndian.Uint32(data[offset : offset+4])
	offset += 4

	metaLen := binary.BigEndian.Uint16(data[offset : offset+2])
	offset += 2

	metaEnd := offset + int(metaLen)
	for offset < metaEnd {
		keyLen := int(data[offset])
		offset++
		k := string(data[offset : offset+keyLen])
		offset += keyLen

		valLen := int(binary.BigEndian.Uint16(data[offset : offset+2]))
		offset += 2
		v := string(data[offset : offset+valLen])
		offset += valLen

		f.Metadata[k] = v
	}

	f.Payload = data[offset : 4+frameLen]
	return f, 4 + int(frameLen), nil
}
