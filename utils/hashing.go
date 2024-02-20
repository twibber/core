package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

// ArgonParams defines the parameters for the Argon2 hashing algorithm.
type ArgonParams struct {
	memory      uint32 // Memory usage
	iterations  uint32 // Number of iterations
	parallelism uint8  // Number of threads and lanes
	saltLength  uint32 // Length of the salt
	keyLength   uint32 // Length of the generated key
}

// DefaultArgonParams provides default settings for Argon2 parameters.
var DefaultArgonParams = ArgonParams{
	memory:      64 * 1024, // 64 MB
	iterations:  3,
	parallelism: 2,
	saltLength:  16, // 16 bytes
	keyLength:   32, // 32 bytes
}

// Predefined errors for hash validation and processing.
var (
	ErrInvalidHash         = fmt.Errorf("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = fmt.Errorf("incompatible version of argon2")
)

// CreateHash generates a hash for a given password using Argon2.
func CreateHash(password string) (encodedHash string, err error) {
	// Generate a cryptographically secure random salt.
	salt, err := GenerateRandomBytes(DefaultArgonParams.saltLength)
	if err != nil {
		return "", err
	}

	// Generate the hash using Argon2id with the provided password and salt.
	hash := argon2.IDKey([]byte(password), salt, DefaultArgonParams.iterations, DefaultArgonParams.memory, DefaultArgonParams.parallelism, DefaultArgonParams.keyLength)

	// Encode the salt and hash in base64 for storage.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format the final encoded hash string with all parameters for verification later.
	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, DefaultArgonParams.memory, DefaultArgonParams.iterations, DefaultArgonParams.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// GenerateRandomBytes creates a slice of random bytes of specified length.
func GenerateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// CompareHash checks if a password matches the hash.
func CompareHash(password, encodedHash string) (match bool, err error) {
	// Decode the hash into its components.
	params, salt, hash, err := DecodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Hash the password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, params.iterations, params.memory, params.parallelism, params.keyLength)

	// Use constant time comparison to mitigate timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

// DecodeHash extracts the parameters, salt, and hash from an encoded hash string.
func DecodeHash(encodedHash string) (p *ArgonParams, salt, hash []byte, err error) {
	// Split the encoded hash into its components.
	values := strings.Split(encodedHash, "$")
	if len(values) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	// Check the Argon2 version.
	var version int
	_, err = fmt.Sscanf(values[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	// Extract the Argon2 parameters.
	p = &ArgonParams{}
	_, err = fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	// Decode the salt.
	salt, err = base64.RawStdEncoding.Strict().DecodeString(values[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	// Decode the hash.
	hash, err = base64.RawStdEncoding.Strict().DecodeString(values[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
