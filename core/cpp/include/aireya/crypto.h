#pragma once

#include <cstdint>
#include <vector>
#include <memory>
#include <stdexcept>

namespace aireya {
namespace crypto {

// The 10 Ultra-Strong Cryptographic Suites
enum class CipherSuite : uint8_t {
    NONE = 0x00,
    AES_256_GCM = 0x01,             // Industry standard, hardware accelerated
    CHACHA20_POLY1305 = 0x02,       // Extremely fast in software, mobile friendly
    XCHACHA20_POLY1305 = 0x03,      // Extended nonce for extreme security
    SALSA20_MAC = 0x04,             // High-speed stream cipher
    CAMELLIA_256_GCM = 0x05,        // Strong alternative block cipher
    SERPENT_256_GCM = 0x06,         // Highest security margin block cipher
    TWOFISH_256_GCM = 0x07,         // Complex key schedule, Bruce Schneier design
    KYBER_HYBRID = 0x08,            // Post-Quantum Cryptography Hybrid
    AIREYA_CUSTOM_A = 0x09,         // Self-developed proprietary cipher A
    AIREYA_CUSTOM_B = 0x0A          // Self-developed proprietary cipher B
};

// 10 Custom Proprietary Encryption Headers (Bitmask/Types)
enum class CryptoHeaderType : uint8_t {
    STANDARD_IV = 0x00,           // Standard Initialization Vector
    ROLLING_KEY_EPOCH = 0x01,     // Embeds an epoch for key rotation every N frames
    ONION_ROUTED = 0x02,          // Header designed for multi-hop Tor-like routing
    POST_QUANTUM_NONCE = 0x03,    // Extended header for Quantum resistance
    STEGANOGRAPHY_PADDING = 0x04, // Header disguises traffic size (adds random noise padding)
    DECENTRALIZED_SIG = 0x05,     // Embeds an ECDSA/Ed25519 signature of the sender
    ZERO_KNOWLEDGE_PROOF = 0x06,  // Embeds a zk-SNARK proof of authorization
    TIME_LOCKED = 0x07,           // Payload cannot be decrypted before a certain timestamp
    EPHEMERAL_DH = 0x08,          // Diffie-Hellman public key for perfect forward secrecy per-frame
    AIREYA_PROPRIETARY = 0x09     // Fully custom obscured header layout
};

// The generic interface for all our cipher suites
class ICipher {
public:
    virtual ~ICipher() = default;
    virtual std::vector<uint8_t> Encrypt(const std::vector<uint8_t>& plaintext, const std::vector<uint8_t>& key) = 0;
    virtual std::vector<uint8_t> Decrypt(const std::vector<uint8_t>& ciphertext, const std::vector<uint8_t>& key) = 0;
};

// Factory to spawn the selected cryptography engine
class CryptoEngine {
public:
    static std::unique_ptr<ICipher> CreateCipher(CipherSuite suite);

    // Embeds the 10 custom headers into the ciphertext
    static std::vector<uint8_t> WrapWithHeader(const std::vector<uint8_t>& encrypted_payload, CryptoHeaderType header_type);
    
    // Strips and verifies the custom header before decryption
    static std::vector<uint8_t> UnwrapHeader(const std::vector<uint8_t>& framed_payload, CryptoHeaderType& out_header_type);
};

} // namespace crypto
} // namespace aireya
