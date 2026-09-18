# Aireya-RPC Benchmark Harness

To ensure transparency and reproducibility, this directory contains the exact methodology, hardware configurations, and scripts used to generate the performance metrics reported in the main README.

## 🖥️ Hardware & OS Environment

All benchmarks were executed on an isolated bare-metal server to prevent hypervisor-induced jitter.

* **CPU**: Dual Intel Xeon Platinum 8380 (80 Cores / 160 Threads total, locked at 2.30GHz).
* **RAM**: 256GB DDR4-3200 ECC.
* **Network**: Loopback interface (`lo`) for IPC testing, avoiding NIC driver bottlenecks.
* **OS**: Ubuntu 22.04 LTS (Kernel `5.15.0-100-generic`).

### Kernel Tuning (`sysctl`)
To achieve 200k+ QPS and sub-millisecond tail latencies, we applied the following strict network tuning:
```bash
# Disable Nagle's algorithm and TCP slow start
sysctl -w net.ipv4.tcp_nodelay=1
sysctl -w net.ipv4.tcp_slow_start_after_idle=0

# Increase ephemeral port range for 10k+ concurrent connections
sysctl -w net.ipv4.ip_local_port_range="1024 65535"
sysctl -w net.core.somaxconn=65535

# Increase socket buffer sizes
sysctl -w net.core.rmem_max=16777216
sysctl -w net.core.wmem_max=16777216
```

## 🧪 Methodology

* **Payload**: Exactly 1024 bytes (1KB) of randomized string/list data.
* **Concurrency**: 10,000 persistent connections.
* **Duration**: 5 minutes of sustained load (300 seconds), capturing the 99th percentile (P99) latency distribution to observe GC pauses (in Envoy/Go) and memory allocator spikes (in C++ Protobuf).
* **Tooling**: We use a custom C++ async load generator based on `epoll` for Aireya, and `ghz` / `wrk2` for the HTTP/2 based competitor frameworks.

## 🚀 Running the Benchmarks

To independently verify our results, compile the benchmark clients and run the automated suite:

```bash
# 1. Compile the Aireya benchmark client
cd ../core/cpp/build && make aireya_bench_client

# 2. Run the harness
./bench.sh --target=aireya --concurrency=10000 --payload=1kb --time=300s
```

*(Note: The implementations of "Framework A" and "Framework B" benchmark clients are provided as pre-compiled binaries in this suite to adhere to our naming anonymity policy. If you wish to inspect their source, please contact the maintainers).*
