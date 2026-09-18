# RFC-0001: Aireya Remote Procedure Call (aRPC) Core Specification

**Version:** 1.0.0
**Status:** DRAFT (Under Active Development)
**Authors:** Antigravity, ChloeYuki

---

## 1. Introduction
This document defines the formal specification for Aireya aRPC, separating the protocol design from any specific language implementation (C++, Go, Rust, etc.). All official client and server implementations MUST adhere to the rules defined herein.

---

## 2. Protocol Specification & Wire Format

### 2.1 The aRPC/1 Binary Frame
All data is transmitted in binary frames. The frame layout is designed for zero-copy parsing.

| Offset | Length | Field Name | Description |
| :--- | :--- | :--- | :--- |
| 0 | 4 | `Frame Length` | Total length of the frame (excludes this 4-byte header), Big-Endian. |
| 4 | 8 | `Magic & Version`| Fixed ASCII string: `"Aireyav1"`. Used for fast protocol sniffing. |
| 12 | 1 | `Frame Type` | `0x00`: Request, `0x01`: Response, `0x02`: Error, `0x03`: Ping, `0x04`: Pong, `0x05`: Handshake. |
| 13 | 1 | `Flags` | Bit 0 (`0x01`): `END_OF_STREAM`, Bit 1 (`0x02`): `COMPRESSED`, Bit 2 (`0x04`): `ENCRYPTED`. |
| 14 | 4 | `Stream ID` | Unique ID for multiplexing. Odd IDs initiated by client, Even by server. |
| 18 | 2 | `Metadata Len` | Length of the metadata section, Big-Endian. |
| 20 | Var | `Metadata` | Key-Value pairs. Encoded as: `[KeyLen(1 byte)] [Key] [ValLen(2 bytes)] [Val]` |
| Var | Var | `Payload` | Encrypted/Compressed IDL Payload (or Error Struct). |

### 2.2 Error Codes
Standardized 4-byte integer error codes returned in `0x02` Error Frames:
* `0`: OK
* `1`: CANCELLED
* `2`: UNKNOWN
* `3`: INVALID_ARGUMENT
* `4`: DEADLINE_EXCEEDED
* `5`: NOT_FOUND
* `7`: PERMISSION_DENIED
* `14`: UNAVAILABLE
* `16`: UNAUTHENTICATED

### 2.3 Streaming Behavior
aRPC supports fully bi-directional streaming multiplexed over a single TCP/QUIC connection.
* **Lifecycle**: A stream is opened implicitly when a frame with a new `Stream ID` is received.
* **Termination**: A stream is fully closed ONLY when both client and server have sent a frame with the `END_OF_STREAM (0x01)` flag set on that specific `Stream ID`.

---

## 3. Forward & Backward Compatibility Principles

To prevent catastrophic system failures during microservice upgrades, all aRPC implementations MUST follow these strict schema evolution rules:

1. **Unknown Fields**: If an old client receives a payload containing a new field tag it does not recognize, the decoder **MUST NOT crash**. It MUST calculate the length using the wire-type tag and safely skip the bytes.
2. **Field Reuse**: A field tag number (e.g., `1: string name`) that has been deleted **MUST NEVER** be reused for a different variable in the future. The compiler `aireyac` will enforce this via linting.
3. **Default Values**: Missing fields in a received payload evaluate to their semantic zero-value (`0`, `""`, `false`).
4. **Old Client -> New Server**: The server will treat missing new fields as zero-values and respond correctly.

---

## 4. Versioning Strategy

aRPC implements a three-tier versioning strategy to prevent "upgrade explosions":
1. **Transport Protocol Version (`aRPC/1`)**: Declared in the 8-byte Magic String (`"Aireyav1"`). If a future `aRPC/2` introduces a radical framing change, the server can instantly reject `v1` connections gracefully during the first 8 bytes.
2. **IDL Schema Version**: Handled gracefully via the Compatibility Principles (Section 3).
3. **Handshake Negotiation**: Upon establishing a TCP connection, the client sends a `0x05 Handshake` frame containing supported Crypto Suites and Compression algorithms. The server replies with the chosen suite.

---

## 5. Security Model

Security is negotiated during the `0x05 Handshake` frame.

* **Transport Encryption (QUIC)**: Relies on native TLS 1.3 mTLS.
* **Payload Encryption (TCP)**: Uses the **Noise Protocol Framework** (e.g., `Noise_XX_25519_AESGCM_SHA256`).
* **Node Identity**: Nodes in the decentralized ANA federation are identified by their `Ed25519` Public Keys.
* **Anti-Replay Protection**: All encrypted payloads include an incrementing 64-bit Nonce. The receiver strictly drops nonces that are $\le$ the highest seen nonce for that Stream ID.
* **Authentication/ACL**: Propagated via `Metadata` headers (e.g., `Authorization: Bearer <JWT>`).

---

## 6. Independent Benchmarking Verification

Performance metrics MUST NOT be accepted as fact without independent verification. 
The official aRPC repository maintains a strict, isolated benchmarking harness located in the `/benchmark/` directory.

* **Transparency**: The scripts (`bench.sh`) explicitly outline the OS kernel parameters (`sysctl`) and hardware used.
* **Fairness**: Competitor frameworks (gRPC, tRPC) are executed under identical loopback network conditions with identical 1KB payloads.
* All performance claims in the aRPC `README.md` must be reproducible via this open-source harness.
