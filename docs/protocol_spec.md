# Aireya-RPC Protocol Specification (v1)

## 1. Frame Layout
All data transmitted over TCP or QUIC is broken into Frames.

| Offset | Length (Bytes) | Field Name | Description |
| :--- | :--- | :--- | :--- |
| 0 | 4 | `Frame Length` | Total length of the frame (excludes this 4-byte length header) |
| 4 | 8 | `Magic & Version`| String identifying protocol and version, e.g., `"Aireyav1"` or `"Aireya RPC v1\0\0\0"` (up to 16 bytes depending on preference. Let's use 8 bytes: `"Aireyav1"`) |
| 12 | 1 | `Frame Type` | 0x00: Request, 0x01: Response, 0x02: Error, 0x03: Ping, 0x04: Pong |
| 13 | 1 | `Flags` | Bit 0: End of Stream, Bit 1: Compressed (Snappy/LZ4) |
| 14 | 4 | `Stream ID` | Unique ID for multiplexing, big-endian |
| 18 | 2 | `Metadata Len` | Length of the metadata section, big-endian |
| 20 | Var | `Metadata` | Key-Value pairs, UTF-8. Encoded as: `[KeyLen(1 byte)] [Key] [ValLen(2 bytes)] [Val] ...` |
| Var | Var | `Payload` | Serialized IDL struct payload |

## 2. IDL Syntax (.aya)

```text
// .aya file syntax
namespace <identifier>.<identifier>;

enum <identifier> {
    <IDENTIFIER> = <NUMBER>;
}

struct <identifier> {
    <number>: <type> <identifier>;
}

service <identifier> {
    rpc <identifier>(<type>) -> (<type>);
    rpc <identifier>(stream <type>) -> stream(<type>);
}
```

### Supported Types
- `int32`, `int64`, `uint32`, `uint64`
- `float`, `double`
- `string`, `bytes`, `bool`
- `list<T>`
- Custom Structs and Enums
