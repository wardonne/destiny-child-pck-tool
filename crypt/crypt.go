package crypt

import (
	"bytes"
	"crypto/aes"
)

var key = []byte{0x37, 0xea, 0x79, 0x85, 0x86, 0x29, 0xec, 0x94, 0x85, 0x20, 0x7c, 0x1a, 0x62, 0xc3, 0x72, 0x4f, 0x72, 0x75, 0x25, 0x0b, 0x99, 0x99, 0xbd, 0x7f, 0x0b, 0x24, 0x9a, 0x8d, 0x85, 0x38, 0x0e, 0x39}

func Decrypt(data []byte) ([]byte, error) {
	decipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	padding := 16 - (len(data) % 16)
	paddedBytes := append(data, bytes.Repeat([]byte{0x00}, padding)...)
	decrypted := make([]byte, len(paddedBytes))
	for bs, be := 0, decipher.BlockSize(); bs < len(paddedBytes); bs, be = bs+decipher.BlockSize(), be+decipher.BlockSize() {
		decipher.Decrypt(decrypted[bs:be], paddedBytes[bs:be])
	}

	return decrypted[:len(data)], nil
}
