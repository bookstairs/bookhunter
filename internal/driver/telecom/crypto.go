package telecom

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"sort"
	"strings"
)

const (
	b64map = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	biRm   = "0123456789abcdefghijklmnopqrstuvwxyz"
)

// base64Encode base64 encoded.
func base64Encode(raw []byte) []byte {
	var encoded bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &encoded)
	_, _ = encoder.Write(raw)
	_ = encoder.Close()
	return encoded.Bytes()
}

// rsaEncrypt RSA encrypt.
func rsaEncrypt(publicKey, origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// signatureOfMd5 MD5 Signature for telecom.
func signatureOfMd5(params map[string]string) string {
	var keys []string
	for k, v := range params {
		keys = append(keys, k+"="+v)
	}

	// sort
	sort.Strings(keys)

	signStr := strings.Join(keys, "&")

	h := md5.New()
	h.Write([]byte(signStr))
	return hex.EncodeToString(h.Sum(nil))
}

// b64toHex 将base64字符串转换成HEX十六进制字符串
func b64toHex(b64 []byte) string {
	b64str := string(b64)
	sb := strings.Builder{}
	e := 0
	c := 0
	for _, r := range b64str {
		if r != '=' {
			v := strings.Index(b64map, string(r))
			if e == 0 {
				e = 1
				sb.WriteByte(int2char(v >> 2))
				c = 3 & v
			} else if e == 1 {
				e = 2
				sb.WriteByte(int2char(c<<2 | v>>4))
				c = 15 & v
			} else if e == 2 {
				e = 3
				sb.WriteByte(int2char(c))
				sb.WriteByte(int2char(v >> 2))
				c = 3 & v
			} else {
				e = 0
				sb.WriteByte(int2char(c<<2 | v>>4))
				sb.WriteByte(int2char(15 & v))
			}
		}
	}
	if e == 1 {
		sb.WriteByte(int2char(c << 2))
	}
	return sb.String()
}

func int2char(i int) (r byte) {
	return biRm[i]
}
