package namegenerators

import "crypto/sha256"

type SHA256Encoder struct {
}

func (s *SHA256Encoder) Encode(contents []byte) []byte {
	arr := sha256.Sum256(contents)
	return arr[:]
}
