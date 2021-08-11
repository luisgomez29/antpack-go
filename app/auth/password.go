package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Errors
var (
	errInvalidHash         = errors.New("the encoded hash is not in the correct format")
	errIncompatibleVersion = errors.New("incompatible version of argon2")
)

// passwordConfig is the configuration to generate the password hash using the Argon2id variant.
type passwordConfig struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint8
	keyLength   uint32
}

// NewPasswordConfig establish the parameters to use for Argon2.
func NewPasswordConfig() *passwordConfig {
	return &passwordConfig{
		memory:      102400,
		iterations:  2,
		parallelism: 8,
		saltLength:  16,
		keyLength:   16,
	}
}

// generateFromPassword generate a hash of the password using the Argon2id variant.
func generateFromPassword(c *passwordConfig, password string) (string, error) {
	// Generate a cryptographically secure random salt.
	salt, err := generateRandomBytes(c.saltLength)
	if err != nil {
		return "", err
	}

	// Pass the plaintext password, salt and parameters to the argon2.IDKey function.
	hash := argon2.IDKey([]byte(password), salt, c.iterations, c.memory, c.parallelism, c.keyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b65Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodeHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, c.memory, c.iterations, c.parallelism, b64Salt, b65Hash,
	)
	return encodeHash, err
}

// comparePasswordAndHash check that the contents of the hashed passwords are identical. Note that we are using the
// subtle.ConstantTimeCompare() function for this to help prevent timing attacks.
func comparePasswordAndHash(password, encodedHash string) (bool, error) {
	c, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, c.iterations, c.memory, c.parallelism, c.keyLength)
	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}

// generateRandomBytes generate byte slice used as a cryptographically secure random salt.
func generateRandomBytes(n uint8) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

// decodeHash extract the parameters, salt and derived key from the encoded password hash.
func decodeHash(encodedHash string) (c *passwordConfig, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errInvalidHash
	}

	var version int
	if _, err = fmt.Sscanf(vals[2], "v=%d", &version); err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errIncompatibleVersion
	}

	c = &passwordConfig{}
	if _, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &c.memory, &c.iterations, &c.parallelism); err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	c.saltLength = uint8(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	c.keyLength = uint32(len(hash))

	return c, salt, hash, nil
}
