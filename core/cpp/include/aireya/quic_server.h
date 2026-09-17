#pragma once

#include "aireya/server.h"

namespace aireya {
namespace quic {

// QuicServer represents the next-generation transport for Aireya-RPC.
// It runs over UDP and uses a robust QUIC implementation (like ngtcp2 or quiche)
// to multiplex streams without Head-of-Line blocking.

class QuicServer {
public:
    QuicServer(uint16_t port);
    ~QuicServer();

    void SetHandler(FrameHandler handler);
    void Run();
    void Stop();

private:
    uint16_t port_;
    int udp_fd_;
    bool running_{false};
    FrameHandler handler_;

    // In a full implementation, we would maintain a map of:
    // QuicConnectionID -> QuicSession
    // And each QuicSession has multiple QuicStreams (which map perfectly to our Frame.stream_id).

    void HandleUdpRead();
};

} // namespace quic
} // namespace aireya
