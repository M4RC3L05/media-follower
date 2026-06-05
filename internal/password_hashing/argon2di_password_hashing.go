package passwordhashing

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2diPasswordHashing struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	Keylen  uint32
	Salt    *[]byte
}

type Argon2diHash struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	Salt    []byte
	Hash    []byte
}

func (h Argon2diHash) String() string {
	return fmt.Sprintf(
		"%d:%d:%d:%s:%s",
		h.Time,
		h.Memory,
		h.Threads,
		base64.RawStdEncoding.EncodeToString(h.Salt),
		base64.RawStdEncoding.EncodeToString(h.Hash),
	)
}

// Compile time check that providers implement interface
var _ IPasswordHashing = Argon2diPasswordHashing{}

func NewArgon2diPasswordHashing() Argon2diPasswordHashing {
	return Argon2diPasswordHashing{
		Time:    3,
		Memory:  46 * 1024,
		Threads: 1,
		Keylen:  32,
	}
}

func decodeHash(h string) (*Argon2diHash, error) {
	parts := strings.Split(h, ":")

	hh, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return nil, err
	}

	time, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return nil, err
	}

	memory, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return nil, err
	}

	threads, err := strconv.ParseUint(parts[2], 10, 8)
	if err != nil {
		return nil, err
	}

	return &Argon2diHash{
		Salt:    salt,
		Hash:    hh,
		Time:    uint32(time),
		Memory:  uint32(memory),
		Threads: uint8(threads),
	}, nil
}

func (a Argon2diPasswordHashing) Compare(hash string, plaintext string) bool {
	h, err := decodeHash(hash)
	if err != nil {
		return false
	}

	hh := pHash(plaintext, h.Salt, h.Time, h.Memory, h.Threads, 32)

	return bytes.Equal(h.Hash, hh)
}

func pHash(
	plaintext string,
	salt []byte,
	time uint32,
	memory uint32,
	threads uint8,
	keylen uint32,
) []byte {
	return argon2.IDKey([]byte(plaintext), salt, time, memory, threads, keylen)
}

func (a Argon2diPasswordHashing) Hash(plaintext string) string {
	salt := make([]byte, 16)
	rand.Read(salt)

	hash := pHash(plaintext, salt, a.Time, a.Memory, a.Threads, a.Keylen)

	return Argon2diHash{
		Salt:    salt,
		Hash:    hash,
		Time:    a.Time,
		Memory:  a.Memory,
		Threads: a.Threads,
	}.String()
}
