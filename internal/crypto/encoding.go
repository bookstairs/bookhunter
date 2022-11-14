package crypto

import (
	"bytes"
	"encoding/base64"
)

// Base64Encode base64 encoded.
func Base64Encode(raw []byte) []byte {
	var encoded bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &encoded)
	_, _ = encoder.Write(raw)
	_ = encoder.Close()
	return encoded.Bytes()
}
