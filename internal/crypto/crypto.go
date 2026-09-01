package crypto

import (
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
)

type CipherService struct{
	aead cipher.AEAD
}

func New(key []byte) (*CipherService, error) { 
	aead, err := chacha20poly1305.New(key)
	if err != nil { 
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	return &CipherService{
		aead: aead,
	}, nil
}

func (s *CipherService) Encrypt(text []byte) ([]byte, error) { 
	nonce := make([]byte, s.aead.NonceSize())

	if _, err := rand.Read(nonce); err != nil { 
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	cipheredText := s.aead.Seal(nil, nonce, text, nil)
	return append(nonce, cipheredText...), nil
}

func (s *CipherService) Decrypt(packet []byte) ([]byte, error) { 
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