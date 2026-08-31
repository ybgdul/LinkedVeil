package crypto

import (
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
)

type CryptService struct{
	aead cipher.AEAD
}

func New(key []byte) (*CryptService, error) { 
	aead, err := chacha20poly1305.New(key)
	if err != nil { 
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	return &CryptService{
		aead: aead,
	}, nil
}

func (s *CryptService) Encrypt(text []byte) ([]byte, error) { 
	nonce := make([]byte, s.aead.NonceSize())

	if _, err := rand.Read(nonce); err != nil { 
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	cipheredText := s.aead.Seal(nil, nonce, text, nil)
	return append(nonce, cipheredText...), nil
}

func (s *CryptService) Decrypt(packet []byte) ([]byte, error) { 
	nonceSize := s.aead.NonceSize()

	if len(packet) < nonceSize { 
		return nil, fmt.Errorf("packet too short")
	}

	nonce:= packet[:nonceSize]
	cipheredText := packet[nonceSize:]

	plainText, err := s.aead.Open(nil, nonce, cipheredText, nil)
	if err != nil { 
		return nil, fmt.Errorf("decrypting packet: %w", err)
	}

	return plainText, nil
}