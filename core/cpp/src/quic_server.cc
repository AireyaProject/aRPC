#include "aireya/quic_server.h"
#include <sys/socket.h>
#include <netinet/in.h>
#include <unistd.h>
#include <iostream>
#include <stdexcept>

namespace aireya {
namespace quic {

QuicServer::QuicServer(uint16_t port) : port_(port), udp_fd_(-1) {
    udp_fd_ = socket(AF_INET, SOCK_DGRAM, 0);
    if (udp_fd_ < 0) {
        throw std::runtime_error("Failed to create UDP socket for QUIC");
    }

    sockaddr_in addr{};
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = INADDR_ANY;
    addr.sin_port = htons(port_);

    if (bind(udp_fd_, (struct sockaddr*)&addr, sizeof(addr)) < 0) {
        throw std::runtime_error("Failed to bind UDP socket");
    }
}

QuicServer::~QuicServer() {
    Stop();
    if (udp_fd_ != -1) ::close(udp_fd_);
}

void QuicServer::SetHandler(FrameHandler handler) {
    handler_ = std::move(handler);
}

void QuicServer::Run() {
    running_ = true;
    std::cout << "[EXPERIMENTAL] Aireya QUIC Server listening on UDP port " << port_ << std::endl;

    // A real QUIC implementation would hook this UDP socket into an Event Loop
    // and pass the datagrams into the QUIC state machine (e.g., quiche_conn_recv).
    // The QUIC state machine would then fire "Stream Data Readable" events,
    // which we would decode into Aireya Frames.

    while (running_) {
        HandleUdpRead();
    }
}

void QuicServer::Stop() {
    running_ = false;
}

void QuicServer::HandleUdpRead() {
    char buffer[65535];
    sockaddr_in client_addr{};
    socklen_t client_len = sizeof(client_addr);

    ssize_t n = recvfrom(udp_fd_, buffer, sizeof(buffer), 0, (struct sockaddr*)&client_addr, &client_len);
    if (n > 0) {
        // std::cout << "Received " << n << " bytes of UDP QUIC Datagrams." << std::endl;
        // TO DO: Pass `buffer` to QUIC TLS/crypto engine to decrypt stream data.
        // Once decrypted, parse as Frame::Decode(...)
    }
}

} // namespace quic
} // namespace aireya
