#pragma once

#include <cstdint>
#include <string>
#include <vector>

namespace aireya {
namespace wire {

// Core Serialization engine for Aireya IDL.
// Uses a memory-aligned or variable length encoding depending on the type.

enum class WireType : uint8_t {
    VARINT = 0,    // int32, int64, uint32, uint64, bool, enum
    FIXED64 = 1,   // double, fixed64
    LENGTH_DELIMITED = 2, // string, bytes, embedded messages, repeated fields
    FIXED32 = 5,   // float, fixed32
};

class Encoder {
public:
    Encoder(std::vector<uint8_t>& buffer) : buffer_(buffer) {}

    // Encode a tag (Field number + Wire type)
    void WriteTag(uint32_t field_number, WireType type);

    // Write primitives
    void WriteVarint(uint64_t value);
    void WriteInt64(uint32_t field_number, int64_t value);
    void WriteString(uint32_t field_number, const std::string& value);

private:
    std::vector<uint8_t>& buffer_;
};

class Decoder {
public:
    Decoder(const uint8_t* data, size_t size) : data_(data), size_(size), offset_(0) {}

    bool IsAtEnd() const { return offset_ >= size_; }
    
    // Read the next tag, returns false if EOF
    bool ReadTag(uint32_t& field_number, WireType& type);

    // Read primitives
    bool ReadVarint(uint64_t& value);
    bool ReadInt64(int64_t& value);
    bool ReadString(std::string& value);

    // Skip a field if unknown
    bool SkipField(WireType type);

private:
    const uint8_t* data_;
    size_t size_;
    size_t offset_;
};

} // namespace wire
} // namespace aireya
