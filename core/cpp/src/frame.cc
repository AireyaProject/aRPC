#include "aireya/frame.h"
#include <cstring>
#include <arpa/inet.h> // For htonl/ntohl (Note: we use big-endian for network format)

namespace aireya {

std::vector<uint8_t> Frame::Encode() const {
    // 1. Calculate Metadata size
    uint16_t meta_len = 0;
    for (const auto& kv : metadata) {
        meta_len += 1 + kv.first.size() + 2 + kv.second.size();
    }

    // 2. Calculate total frame length (excluding the first 4 bytes of frame_length itself)
    uint32_t total_len = (HEADER_BASE_SIZE - 4) + meta_len + payload.size();

    // 3. Allocate buffer
    std::vector<uint8_t> buf(4 + total_len);
    size_t offset = 0;

    // Frame Length (Big Endian)
    uint32_t net_len = htonl(total_len);
    std::memcpy(buf.data() + offset, &net_len, 4);
    offset += 4;

    // Magic Version (8 bytes)
    std::memcpy(buf.data() + offset, MAGIC_VERSION, 8);
    offset += 8;

    // Type and Flags
    buf[offset++] = static_cast<uint8_t>(type);
    buf[offset++] = flags;

    // Stream ID
    uint32_t net_stream_id = htonl(stream_id);
    std::memcpy(buf.data() + offset, &net_stream_id, 4);
    offset += 4;

    // Metadata length
    uint16_t net_meta_len = htons(meta_len);
    std::memcpy(buf.data() + offset, &net_meta_len, 2);
    offset += 2;

    // Metadata Key-Values
    for (const auto& kv : metadata) {
        buf[offset++] = static_cast<uint8_t>(kv.first.size());
        std::memcpy(buf.data() + offset, kv.first.data(), kv.first.size());
        offset += kv.first.size();

        uint16_t val_len = htons(static_cast<uint16_t>(kv.second.size()));
        std::memcpy(buf.data() + offset, &val_len, 2);
        offset += 2;

        std::memcpy(buf.data() + offset, kv.second.data(), kv.second.size());
        offset += kv.second.size();
    }

    // Payload
    if (!payload.empty()) {
        std::memcpy(buf.data() + offset, payload.data(), payload.size());
    }

    return buf;
}

int32_t Frame::Decode(const uint8_t* data, size_t size, Frame& out_frame) {
    if (size < 4) return 0; // Not enough for length

    uint32_t net_len;
    std::memcpy(&net_len, data, 4);
    uint32_t frame_len = ntohl(net_len);

    if (size < 4 + frame_len) return 0; // Incomplete frame

    if (frame_len < (HEADER_BASE_SIZE - 4)) return -1; // Invalid frame size

    size_t offset = 4;

    // Check Magic
    if (std::memcmp(data + offset, MAGIC_VERSION, 8) != 0) {
        return -1; // Invalid Magic
    }
    offset += 8;

    out_frame.frame_length = frame_len;
    out_frame.type = static_cast<FrameType>(data[offset++]);
    out_frame.flags = data[offset++];

    uint32_t net_stream_id;
    std::memcpy(&net_stream_id, data + offset, 4);
    out_frame.stream_id = ntohl(net_stream_id);
    offset += 4;

    uint16_t net_meta_len;
    std::memcpy(&net_meta_len, data + offset, 2);
    uint16_t meta_len = ntohs(net_meta_len);
    offset += 2;

    if (frame_len < (HEADER_BASE_SIZE - 4) + meta_len) return -1;

    // Parse Metadata
    size_t meta_end = offset + meta_len;
    out_frame.metadata.clear();
    while (offset < meta_end) {
        uint8_t key_len = data[offset++];
        if (offset + key_len > meta_end) return -1;
        std::string key(reinterpret_cast<const char*>(data + offset), key_len);
        offset += key_len;

        if (offset + 2 > meta_end) return -1;
        uint16_t val_len_net;
        std::memcpy(&val_len_net, data + offset, 2);
        uint16_t val_len = ntohs(val_len_net);
        offset += 2;

        if (offset + val_len > meta_end) return -1;
        std::string val(reinterpret_cast<const char*>(data + offset), val_len);
        offset += val_len;

        out_frame.metadata[key] = val;
    }

    // Parse Payload
    size_t payload_size = (4 + frame_len) - offset;
    out_frame.payload.assign(data + offset, data + offset + payload_size);

    return 4 + frame_len;
}

} // namespace aireya
