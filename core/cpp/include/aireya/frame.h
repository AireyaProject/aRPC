#pragma once

#include <cstdint>
#include <string>
#include <vector>
#include <unordered_map>

namespace aireya {

// Frame Type enum
enum class FrameType : uint8_t {
    REQUEST = 0x00,
    RESPONSE = 0x01,
    ERROR = 0x02,
    PING = 0x03,
    PONG = 0x04
};

// Flags Bitmask
enum FrameFlags : uint8_t {
    NONE = 0x00,
    END_OF_STREAM = 0x01,
    COMPRESSED = 0x02
};

// Represents a parsed Aireya RPC Frame
class Frame {
public:
    static constexpr const char* MAGIC_VERSION = "Aireyav1";
    static constexpr size_t HEADER_BASE_SIZE = 4 + 8 + 1 + 1 + 4 + 2; // 20 bytes

    uint32_t frame_length = 0;
    FrameType type = FrameType::REQUEST;
    uint8_t flags = FrameFlags::NONE;
    uint32_t stream_id = 0;
    
    std::unordered_map<std::string, std::string> metadata;
    std::vector<uint8_t> payload;

    // Serialize the frame into a byte buffer
    std::vector<uint8_t> Encode() const;

    // Decode a frame from a byte buffer. 
    // Returns the number of bytes consumed, or 0 if incomplete, or -1 if invalid magic.
    static int32_t Decode(const uint8_t* data, size_t size, Frame& out_frame);
};

} // namespace aireya
