#include "aireya/resolver.h"
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>
#include <stdexcept>
#include <iostream>
#include <cstring>
#include <sstream>

namespace aireya {

Endpoint Resolver::Resolve(const std::string& service_name) {
    int sock = socket(AF_UNIX, SOCK_STREAM, 0);
    if (sock < 0) {
        throw std::runtime_error("Failed to create UNIX domain socket");
    }

    struct sockaddr_un addr{};
    addr.sun_family = AF_UNIX;
    strncpy(addr.sun_path, "/tmp/aireya.sock", sizeof(addr.sun_path) - 1);

    if (connect(sock, (struct sockaddr*)&addr, sizeof(addr)) < 0) {
        ::close(sock);
        throw std::runtime_error("Failed to connect to Aireya Node Agent (ANA). Is it running?");
    }

    // Send JSON request
    std::string request = "{\"service\": \"" + service_name + "\"}";
    ssize_t wret = ::write(sock, request.c_str(), request.length());
    (void)wret;

    // Read JSON response
    char buffer[1024];
    ssize_t n = ::read(sock, buffer, sizeof(buffer) - 1);
    ::close(sock);

    if (n <= 0) {
        throw std::runtime_error("Failed to read from ANA");
    }
    buffer[n] = '\0';
    std::string response(buffer);

    // Ultra-basic JSON parsing for the skeleton (avoids pulling in nlohmann/json for now)
    Endpoint ep;
    ep.port = 0;
    
    // Look for "host":"xxx"
    size_t host_pos = response.find("\"host\":\"");
    if (host_pos != std::string::npos) {
        host_pos += 8;
        size_t host_end = response.find("\"", host_pos);
        if (host_end != std::string::npos) {
            ep.host = response.substr(host_pos, host_end - host_pos);
        }
    }

    // Look for "port":xxx
    size_t port_pos = response.find("\"port\":");
    if (port_pos != std::string::npos) {
        port_pos += 7;
        size_t port_end = response.find_first_of(",}", port_pos);
        if (port_end != std::string::npos) {
            ep.port = std::stoi(response.substr(port_pos, port_end - port_pos));
        }
    }

    if (ep.host.empty() || ep.port == 0) {
        throw std::runtime_error("Service not found or invalid response from ANA: " + response);
    }

    return ep;
}

} // namespace aireya
