# Aireya aRPC ✨

**Aireya aRPC (Aireya Remote Procedure Call)** is a high-performance, decentralized RPC framework designed for distributed systems, edge computing, service mesh architectures, and Web3 federation.

aRPC provides a lightweight alternative to traditional HTTP/2-based RPC frameworks such as gRPC. It combines a custom binary wire protocol, the `.aya` Interface Definition Language, automatic C++ and Go code generation, proxyless service communication, peer-to-peer service discovery, and optional Layer 7 routing.

The project is designed for developers building high-performance microservices, distributed applications, edge infrastructure, decentralized networks, and low-latency backend systems.

## Why Aireya aRPC?

Traditional RPC and service mesh architectures often depend on HTTP/2, centralized control planes, or per-service sidecar proxies. Aireya aRPC explores a different architecture focused on low overhead, decentralization, and direct service-to-service communication.

Key features include:

* High-performance binary RPC protocol
* Custom `.aya` Interface Definition Language
* Aireya Compiler (`aireyac`)
* Automatic C++17+ and Go code generation
* Proxyless RPC communication
* Unix Domain Socket integration with Aireya Node Agent
* Decentralized service discovery
* libp2p and Kademlia DHT federation
* Rust and Tokio based Layer 7 Waypoint routing
* Edge computing support
* Service mesh architecture without traditional sidecars
* Streaming RPC support
* Multiple transport and encryption strategies
* Hybrid cryptography support
* Designed for low-latency and high-concurrency workloads

## Architecture

Aireya aRPC consists of several components:

### aRPC Core

The core RPC runtime implements the Aireya wire protocol, request routing, serialization, transport, connection management, and runtime behavior.

### `.aya` IDL

Aireya uses its own Interface Definition Language for defining structures and RPC services.

Example:

```aya
namespace aireya.example;

struct User {
    1: int64 id;
    2: string name;
}

service UserService {
    rpc GetUser(User) -> (User);
}
```

### Aireya Compiler

`aireyac` compiles `.aya` definitions into client and server bindings.

Current targets include:

* [x] C++17+
* [x] Go
* [x] Rust (Tokio async/await Native Support)
* [x] Java (CompletableFuture Async Hooks)
* [x] Node.js / TypeScript (V8 Optimized Promise API)

### ANA — Aireya Node Agent

ANA provides local service discovery, networking, federation, and node-level communication.

Applications can communicate with ANA through Unix Domain Sockets, allowing service discovery and networking logic to remain outside the application process without requiring a traditional sidecar proxy.

### Aireya Waypoint

Aireya Waypoint provides optional Layer 7 routing and traffic management for workloads that require advanced routing behavior.

The Waypoint implementation is based on Rust and Tokio.

### Decentralized Federation

Aireya aRPC can use libp2p and Kademlia DHT for decentralized node and service discovery, allowing multiple Aireya environments to participate in a federated network.

## Use Cases

Aireya aRPC is designed for scenarios such as:

* Microservice communication
* High-performance backend services
* Distributed systems
* Edge computing
* AI infrastructure
* Service mesh networking
* Internal cloud infrastructure
* P2P applications
* Web3 infrastructure
* Federated networks
* Low-latency APIs
* High-concurrency RPC services
* IoT and distributed device networks

## aRPC vs gRPC

Aireya aRPC is not intended to be a drop-in replacement for gRPC.

Instead, it explores an alternative RPC architecture focused on:

* Reduced protocol overhead
* Custom binary serialization
* Direct service communication
* Proxyless service mesh architecture
* Decentralized discovery
* Edge-friendly deployment
* P2P federation
* Independent protocol evolution

Developers familiar with Protocol Buffers and gRPC should find the `.aya` IDL and service model familiar while still being able to use Aireya-specific networking features.

## Project Status

Aireya aRPC is under active development.

APIs, protocols, wire formats, compiler behavior, and distributed networking components may change before stable releases.

Feedback, testing, issues, benchmarks, protocol reviews, and contributions are welcome.

## 📜 Credits, Ownership & Licensing

* **Authors**: Antigravity & [ChloeYuki](https://github.com/chloeyuki)
* **Software Ownership**: [ChloeYuki](https://github.com/chloeyuki)

### Dual-License Strategy
To protect the ecosystem while maximizing developer adoption, Aireya uses a dual-license model:
* **The SDKs & Generated Code**: Licensed under the **[MIT License](LICENSE-SDK)**. You can freely use, integrate, and compile the Aireya client SDKs into your proprietary, closed-source commercial applications without any restrictions.
* **The Core Engine, ANA & Waypoint**: Licensed under the **[AGPL-3.0 License](LICENSE)**. If you modify the core infrastructure or provide the Aireya network components as a managed cloud service, you must open-source your modifications.

## Open Source

Aireya aRPC is developed publicly by **AireyaProject**.

SDK components and generated client code are designed to remain easy to integrate into external applications, while the infrastructure components are developed as open-source software.

GitHub repository:
[AireyaProject/aRPC](https://github.com/AireyaProject/aRPC)

---
*Search keywords: Aireya RPC, Aireya aRPC, aRPC, RPC framework, remote procedure call, gRPC alternative, high performance RPC, decentralized RPC, distributed systems, service mesh, proxyless service mesh, libp2p RPC, Kademlia DHT, edge computing RPC, Web3 RPC framework, C++ RPC framework, Go RPC framework, Rust service mesh, custom RPC protocol, binary RPC protocol, microservices RPC.*
