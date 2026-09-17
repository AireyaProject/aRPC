#pragma once

#include <string>
#include <memory>
#include <cstdint>

namespace aireya {

struct Endpoint {
    std::string host;
    uint16_t port;
};

// Resolver asks the local Aireya Node Agent (ANA) for routing information.
// It uses a UNIX Domain Socket for microsecond-level local IPC latency.
class Resolver {
public:
    // Resolve a service name to an Endpoint
    static Endpoint Resolve(const std::string& service_name);
};

} // namespace aireya
