package passwordhashing

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

var b64 = base64.RawStdEncoding.Strict()

type Argon2di struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	Keylen  uint32
}

type Argon2diHash struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	Keylen  uint32
	Salt    []byte
	Hash    []byte
}

func (h Argon2diHash) string() string {
	return fmt.Sprintf(
		"%d %d %d %d %s %s",
		h.Time,
		h.Memory,
		h.Threads,
		h.Keylen,
		b64.EncodeToString(h.Salt),
		b64.EncodeToString(h.Hash),
	)
}

func decodeHash(hs string) (*Argon2diHash, error) {
	var memory, time, keylen uint32
	var threads uint8
	var salt, hash string

	if _, err := fmt.Sscanf(
		hs,
		"%d %d %d %d %s %s",
		&time,
		&memory,
		&threads,
		&keylen,
		&salt,
		&hash,
	); err != nil {
		return nil, err
	}

	h := Argon2diHash{}

	h.Time = time
	h.Memory = memory
	h.Threads = threads
	h.Keylen = keylen

	hh, err := b64.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	h.Hash = hh

	ss, err := b64.DecodeString(salt)
	if err != nil {
		return nil, err
	}

	h.Salt = ss

	return &h, nil
}

// Compile time check that providers implement interface
var _ IPasswordHashing = Argon2di{}

func NewArgon2di() Argon2di {
	return Argon2di{
		Time:    3,
		Memory:  46 * 1024,
		Threads: 1,
		Keylen:  32,
	}
}

func (a Argon2di) Compare(hash string, plaintext string) bool {
	h, err := decodeHash(hash)
	if err != nil {
		return false
	}

	hh := argon2.IDKey([]byte(plaintext), h.Salt, h.Time, h.Memory, h.Threads, 32)

	if subtle.ConstantTimeEq(int32(h.Keylen), int32(len(hh))) == 0 {
		return false
	}

	if subtle.ConstantTimeCompare(h.Hash, hh) == 0 {
		return false
	}

	return true
}

func (a Argon2di) Hash(plaintext string) string {
	salt := make([]byte, 16)
	rand.Read(salt)

	hash := argon2.IDKey([]byte(plaintext), salt, a.Time, a.Memory, a.Threads, a.Keylen)

	return Argon2diHash{
		Salt:    salt,
		Hash:    hash,
		Time:    a.Time,
		Memory:  a.Memory,
		Threads: a.Threads,
		Keylen:  a.Keylen,
	}.string()
}
