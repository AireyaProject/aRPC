#pragma once

#include <string>
#include <functional>
#include <memory>
#include <thread>
#include <atomic>
#include "aireya/frame.h"

namespace aireya {

class Connection {
public:
    virtual ~Connection() = default;
    virtual void SendFrame(const Frame& frame) = 0;
    virtual void Close() = 0;
};

// Callback type for incoming RPC frames
using FrameHandler = std::function<void(std::shared_ptr<Connection> conn, const Frame& frame)>;

class Server {
public:
    Server(uint16_t port);
    ~Server();

    // Set the handler for incoming frames
    void SetHandler(FrameHandler handler);

    // Start the event loop (blocking)
    void Run();

    // Stop the event loop
    void Stop();

private:
    uint16_t port_;
    int listen_fd_;
    int epoll_fd_;
    std::atomic<bool> running_{false};
    FrameHandler handler_;

    void AcceptConnection();
    void HandleRead(int client_fd);
};

} // namespace aireya
