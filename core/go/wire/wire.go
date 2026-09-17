package wire

type WireType uint8

const (
	Varint          WireType = 0
	Fixed64         WireType = 1
	LengthDelimited WireType = 2
	Fixed32         WireType = 5
)

type Encoder struct {
	Buf []byte
}

func (e *Encoder) WriteTag(fieldNumber uint32, t WireType) {
	tag := (fieldNumber << 3) | uint32(t)
	e.WriteVarint(uint64(tag))
}

func (e *Encoder) WriteVarint(v uint64) {
	for v >= 0x80 {
		e.Buf = append(e.Buf, byte(v|0x80))
		v >>= 7
	}
	e.Buf = append(e.Buf, byte(v))
}

func (e *Encoder) WriteInt64(fieldNumber uint32, v int64) {
	e.WriteTag(fieldNumber, Varint)
	e.WriteVarint(uint64(v))
}

func (e *Encoder) WriteString(fieldNumber uint32, v string) {
	e.WriteTag(fieldNumber, LengthDelimited)
	e.WriteVarint(uint64(len(v)))
	e.Buf = append(e.Buf, []byte(v)...)
}

type Decoder struct {
	Buf    []byte
	Offset int
}

func (d *Decoder) IsAtEnd() bool {
	return d.Offset >= len(d.Buf)
}

func (d *Decoder) ReadTag() (uint32, WireType, bool) {
	if d.IsAtEnd() {
		return 0, 0, false
	}
	tag, ok := d.ReadVarint()
	if !ok {
		return 0, 0, false
	}
	return uint32(tag >> 3), WireType(tag & 0x07), true
}

func (d *Decoder) ReadVarint() (uint64, bool) {
	var v uint64
	var shift uint
	for d.Offset < len(d.Buf) {
		b := d.Buf[d.Offset]
		d.Offset++
		v |= uint64(b&0x7F) << shift
		if b&0x80 == 0 {
			return v, true
		}
		shift += 7
		if shift >= 64 {
			return 0, false
		}
	}
	return 0, false
}

func (d *Decoder) ReadInt64() (int64, bool) {
	v, ok := d.ReadVarint()
	return int64(v), ok
}

func (d *Decoder) ReadString() (string, bool) {
	length, ok := d.ReadVarint()
	if !ok || d.Offset+int(length) > len(d.Buf) {
		return "", false
	}
	str := string(d.Buf[d.Offset : d.Offset+int(length)])
	d.Offset += int(length)
	return str, true
}
