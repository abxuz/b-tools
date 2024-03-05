package bcrypt

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/curve25519"
)

const (
	NoisePublicKeySize  = 32
	NoisePrivateKeySize = 32
)

var (
	ErrInvalidPrivateKey = errors.New("invalid private key")
	ErrInvalidPublicKey  = errors.New("invalid public key")
)

type (
	NoisePublicKey  [NoisePublicKeySize]byte
	NoisePrivateKey [NoisePrivateKeySize]byte
)

func NewPrivateKey() (sk NoisePrivateKey, err error) {
	_, err = rand.Read(sk[:])
	sk.Clamp()
	return
}

func NewPrivateKeyFromString(s string) (sk NoisePrivateKey, err error) {
	err = sk.FromString(s)
	return
}

func NewPublicKeyFromString(s string) (sk NoisePublicKey, err error) {
	err = sk.FromString(s)
	return
}

func NewPrivateKeyFromData(data []byte) (sk NoisePrivateKey, err error) {
	if len(data) != NoisePrivateKeySize {
		err = ErrInvalidPrivateKey
		return
	}

	copy(sk[:], data)
	return
}

func NewPublicKeyFromData(data []byte) (sk NoisePublicKey, err error) {
	if len(data) != NoisePublicKeySize {
		err = ErrInvalidPublicKey
		return
	}

	copy(sk[:], data)
	return
}

func (sk *NoisePrivateKey) FromString(s string) error {
	n, err := base64.StdEncoding.Decode(sk[:], []byte(s))
	if err != nil {
		return err
	}
	if n != NoisePrivateKeySize {
		return ErrInvalidPrivateKey
	}
	return nil
}

func (sk *NoisePrivateKey) String() string {
	return base64.StdEncoding.EncodeToString(sk[:])
}

func (sk *NoisePrivateKey) PublicKey() (pk NoisePublicKey) {
	apk := (*[NoisePublicKeySize]byte)(&pk)
	ask := (*[NoisePrivateKeySize]byte)(sk)
	curve25519.ScalarBaseMult(apk, ask)
	return
}

func (sk *NoisePrivateKey) SharedSecret(pk *NoisePublicKey) (ss [NoisePublicKeySize]byte) {
	apk := (*[NoisePublicKeySize]byte)(pk)
	ask := (*[NoisePrivateKeySize]byte)(sk)
	curve25519.ScalarMult(&ss, ask, apk)
	return
}

func (key NoisePrivateKey) IsZero() bool {
	var zero NoisePrivateKey
	return key.Equals(zero)
}

func (key NoisePrivateKey) Equals(tar NoisePrivateKey) bool {
	return Equals(key[:], tar[:])
}

func (sk *NoisePrivateKey) Clamp() {
	sk[0] &= 248
	sk[31] = (sk[31] & 127) | 64
}

func (sk *NoisePublicKey) FromString(s string) error {
	n, err := base64.StdEncoding.Decode(sk[:], []byte(s))
	if err != nil {
		return err
	}
	if n != NoisePublicKeySize {
		return ErrInvalidPublicKey
	}
	return nil
}

func (sk *NoisePublicKey) String() string {
	return base64.StdEncoding.EncodeToString(sk[:])
}

func (key NoisePublicKey) IsZero() bool {
	var zero NoisePublicKey
	return key.Equals(zero)
}

func (key NoisePublicKey) Equals(tar NoisePublicKey) bool {
	return Equals(key[:], tar[:])
}

func Equals(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}
