package imp

import (
	"crypto"
	"crypto/sha256"
	"github.com/pkg/errors"
	"github.com/scryinfo/dot/dots/scrypto"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"
	"io"
)

type ende25519 struct{}

func EcdhDecoder25519() scrypto.EcdhDecoder {
	return &ende25519{}
}

func EcdhEncoder25519() scrypto.EcdhEncoder {
	return &ende25519{}
}

var (
	salt     []byte
	hash     = sha256.New
	info     = []byte("scry info")
	nonce    = make([]byte, chacha20poly1305.NonceSize)
	ecdh     = X25519()
	endeType = scrypto.EndeType_X25519
)

func (c *ende25519) EcdhDecode(privateKey crypto.PrivateKey, cipher scrypto.EndeData) (plain scrypto.EndeData, err error) {
	if !cipher.EnData {
		return cipher, nil
	}
	if cipher.EndeType != endeType {
		err = errors.New("the ende type is not " + endeType)
		return
	}

	var peersKey crypto.PublicKey = cipher.PublicKey

	cipher.Body, err = c._ecdhDecode(privateKey, peersKey, cipher.Body)
	if err != nil {
		return
	}
	cipher.EnData = false //decode data
	plain = cipher

	return
}

func (c *ende25519) EcdhEncode(privateKey crypto.PrivateKey, peersKey crypto.PublicKey, plain scrypto.EndeData) (cipher scrypto.EndeData, err error) {
	if cipher.EnData {
		return cipher, nil
	}
	plain.EndeType = endeType
	publicKey, err := ecdh.PublicKey(privateKey)
	if err != nil {
		return
	}
	plain.PublicKey, err = ecdh.PublicKeyToBytes(publicKey)
	if err != nil {
		return
	}

	plain.Body, err = c._ecdhEncode(privateKey, peersKey, plain.Body)
	if err != nil {
		return
	}
	plain.EnData = true
	cipher = plain
	return
}

// EcdhDecode
// privateKey x25519, peersKey x25519
//
func (c *ende25519) _ecdhDecode(privateKey crypto.PrivateKey, peersKey crypto.PublicKey, ciphertext []byte) (plaintext []byte, err error) {

	key, err := ecdh.ComputeSecret(privateKey, peersKey)
	if err != nil {
		return
	}
	dk := hkdf.New(hash, key, salt, info)
	wrappingKey := make([]byte, chacha20poly1305.KeySize)
	if _, err = io.ReadFull(dk, wrappingKey); err != nil {
		return
	}
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return
	}
	plaintext, err = aead.Open(nil, nonce, ciphertext, nil)
	return
}

// EcdhEncode
// privateKey x25519, peersKey x25519
func (c *ende25519) _ecdhEncode(privateKey crypto.PrivateKey, peersKey crypto.PublicKey, plaintext []byte) (ciphertext []byte, err error) {
	echg := X25519()
	key, err := echg.ComputeSecret(privateKey, peersKey)
	if err != nil {
		return
	}

	dk := hkdf.New(hash, key, salt, info)

	wrappingKey := make([]byte, chacha20poly1305.KeySize)
	if _, err = io.ReadFull(dk, wrappingKey); err != nil {
		return
	}

	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return
	}
	ciphertext = aead.Seal(nil, nonce, plaintext, nil)
	return
}
