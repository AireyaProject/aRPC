#include "aireya/wire_format.h"

namespace aireya {
namespace wire {

void Encoder::WriteTag(uint32_t field_number, WireType type) {
    uint32_t tag = (field_number << 3) | static_cast<uint32_t>(type);
    WriteVarint(tag);
}

void Encoder::WriteVarint(uint64_t value) {
    while (value >= 0x80) {
        buffer_.push_back(static_cast<uint8_t>(value | 0x80));
        value >>= 7;
    }
    buffer_.push_back(static_cast<uint8_t>(value));
}

void Encoder::WriteInt64(uint32_t field_number, int64_t value) {
    WriteTag(field_number, WireType::VARINT);
    WriteVarint(static_cast<uint64_t>(value));
}

void Encoder::WriteString(uint32_t field_number, const std::string& value) {
    WriteTag(field_number, WireType::LENGTH_DELIMITED);
    WriteVarint(value.size());
    buffer_.insert(buffer_.end(), value.begin(), value.end());
}

bool Decoder::ReadTag(uint32_t& field_number, WireType& type) {
    if (IsAtEnd()) return false;
    uint64_t tag;
    if (!ReadVarint(tag)) return false;
    field_number = tag >> 3;
    type = static_cast<WireType>(tag & 0x07);
    return true;
}

bool Decoder::ReadVarint(uint64_t& value) {
    value = 0;
    uint32_t shift = 0;
    while (offset_ < size_) {
        uint8_t b = data_[offset_++];
        value |= static_cast<uint64_t>(b & 0x7F) << shift;
        if (!(b & 0x80)) return true;
        shift += 7;
        if (shift >= 64) return false; // Invalid varint
    }
    return false;
}

bool Decoder::ReadInt64(int64_t& value) {
    uint64_t v;
    if (!ReadVarint(v)) return false;
    value = static_cast<int64_t>(v);
    return true;
}

bool Decoder::ReadString(std::string& value) {
    uint64_t length;
    if (!ReadVarint(length)) return false;
    if (offset_ + length > size_) return false;
    value.assign(reinterpret_cast<const char*>(data_ + offset_), length);
    offset_ += length;
    return true;
}

bool Decoder::SkipField(WireType type) {
    switch (type) {
        case WireType::VARINT: {
            uint64_t v;
            return ReadVarint(v);
        }
        case WireType::FIXED64: {
            if (offset_ + 8 > size_) return false;
            offset_ += 8;
            return true;
        }
        case WireType::LENGTH_DELIMITED: {
            uint64_t length;
            if (!ReadVarint(length)) return false;
            if (offset_ + length > size_) return false;
            offset_ += length;
            return true;
        }
        case WireType::FIXED32: {
            if (offset_ + 4 > size_) return false;
            offset_ += 4;
            return true;
        }
        default:
            return false;
    }
}

} // namespace wire
} // namespace aireya
