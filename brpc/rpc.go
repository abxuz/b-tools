package brpc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/binary"

	"github.com/abxuz/b-tools/bcrypt"
)

func encrypt(privKey *bcrypt.NoisePrivateKey, pubKey *bcrypt.NoisePublicKey, t int64, data []byte) ([]byte, error) {
	ss := privKey.SharedSecret(pubKey)
	block, err := aes.NewCipher(ss[:])
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	var nonce [12]byte
	binary.BigEndian.PutUint64(nonce[:], uint64(t))
	data = aead.Seal(data[:0], nonce[:], data, nil)
	return data, nil
}

func decrypt(privKey *bcrypt.NoisePrivateKey, pubKey *bcrypt.NoisePublicKey, t int64, data []byte) ([]byte, error) {
	ss := privKey.SharedSecret(pubKey)
	block, err := aes.NewCipher(ss[:])
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	var nonce [12]byte
	binary.BigEndian.PutUint64(nonce[:], uint64(t))
	return aead.Open(data[:0], nonce[:], data, nil)
}

func hash(data []byte) []byte {
	h := md5.Sum(data)
	return h[4:12]
}
