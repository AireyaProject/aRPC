#include "aireya/server.h"
#include <sys/epoll.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <unistd.h>
#include <fcntl.h>
#include <stdexcept>
#include <iostream>
#include <vector>

namespace aireya {

class TcpConnection : public Connection {
public:
    TcpConnection(int fd) : fd_(fd) {}
    ~TcpConnection() { Close(); }

    void SendFrame(const Frame& frame) override {
        std::vector<uint8_t> data = frame.Encode();
        // In a real hardcore implementation, we'd handle partial writes and EAGAIN
        // and register for EPOLLOUT. For now, blocking write for simplicity of skeleton.
        ssize_t ret = ::write(fd_, data.data(), data.size());
        (void)ret; // Suppress unused warning for the skeleton
    }

    void Close() override {
        if (fd_ != -1) {
            ::close(fd_);
            fd_ = -1;
        }
    }

private:
    int fd_;
};

Server::Server(uint16_t port) : port_(port), listen_fd_(-1), epoll_fd_(-1) {
    listen_fd_ = socket(AF_INET, SOCK_STREAM, 0);
    if (listen_fd_ < 0) {
        throw std::runtime_error("Failed to create socket");
    }

    int opt = 1;
    setsockopt(listen_fd_, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));

    // Set non-blocking
    int flags = fcntl(listen_fd_, F_GETFL, 0);
    fcntl(listen_fd_, F_SETFL, flags | O_NONBLOCK);

    sockaddr_in addr{};
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = INADDR_ANY;
    addr.sin_port = htons(port_);

    if (bind(listen_fd_, (struct sockaddr*)&addr, sizeof(addr)) < 0) {
        throw std::runtime_error("Failed to bind socket");
    }

    if (listen(listen_fd_, SOMAXCONN) < 0) {
        throw std::runtime_error("Failed to listen on socket");
    }

    epoll_fd_ = epoll_create1(0);
    if (epoll_fd_ < 0) {
        throw std::runtime_error("Failed to create epoll fd");
    }

    struct epoll_event ev;
    ev.events = EPOLLIN | EPOLLET; // Edge-triggered
    ev.data.fd = listen_fd_;
    if (epoll_ctl(epoll_fd_, EPOLL_CTL_ADD, listen_fd_, &ev) < 0) {
        throw std::runtime_error("Failed to add listen fd to epoll");
    }
}

Server::~Server() {
    Stop();
    if (listen_fd_ != -1) ::close(listen_fd_);
    if (epoll_fd_ != -1) ::close(epoll_fd_);
}

void Server::SetHandler(FrameHandler handler) {
    handler_ = std::move(handler);
}

void Server::Run() {
    running_ = true;
    const int MAX_EVENTS = 64;
    struct epoll_event events[MAX_EVENTS];

    std::cout << "Aireya Server listening on port " << port_ << " (epoll)" << std::endl;

    while (running_) {
        int n = epoll_wait(epoll_fd_, events, MAX_EVENTS, 100);
        if (n < 0) {
            if (errno == EINTR) continue;
            break;
        }

        for (int i = 0; i < n; ++i) {
            if (events[i].data.fd == listen_fd_) {
                AcceptConnection();
            } else {
                HandleRead(events[i].data.fd);
            }
        }
    }
}

void Server::Stop() {
    running_ = false;
}

void Server::AcceptConnection() {
    while (true) {
        sockaddr_in client_addr{};
        socklen_t client_len = sizeof(client_addr);
        int client_fd = accept(listen_fd_, (struct sockaddr*)&client_addr, &client_len);

        if (client_fd < 0) {
            if (errno == EAGAIN || errno == EWOULDBLOCK) {
                break; // No more connections to accept
            }
            continue;
        }

        // Set non-blocking
        int flags = fcntl(client_fd, F_GETFL, 0);
        fcntl(client_fd, F_SETFL, flags | O_NONBLOCK);

        struct epoll_event ev;
        ev.events = EPOLLIN | EPOLLET | EPOLLRDHUP;
        ev.data.fd = client_fd;
        epoll_ctl(epoll_fd_, EPOLL_CTL_ADD, client_fd, &ev);
    }
}

void Server::HandleRead(int client_fd) {
    std::vector<uint8_t> buffer(4096);
    ssize_t n = read(client_fd, buffer.data(), buffer.size());

    if (n <= 0) {
        epoll_ctl(epoll_fd_, EPOLL_CTL_DEL, client_fd, nullptr);
        ::close(client_fd);
        return;
    }

    // Attempt to decode frame
    Frame frame;
    int32_t consumed = Frame::Decode(buffer.data(), n, frame);
    
    if (consumed > 0 && handler_) {
        auto conn = std::make_shared<TcpConnection>(client_fd);
        handler_(conn, frame);
    }
}

} // namespace aireya
