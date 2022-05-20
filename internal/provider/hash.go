package provider

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

func GenerateHashes(filename string) (string, string, string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", "", "", fmt.Errorf("could not compute file '%s' checksum: %s", filename, err)
	}
	h := sha1.New()
	h.Write([]byte(data))
	sha1 := hex.EncodeToString(h.Sum(nil))

	h256 := sha256.New()
	h256.Write([]byte(data))
	shaSum := h256.Sum(nil)
	sha256base64 := base64.StdEncoding.EncodeToString(shaSum[:])

	md5 := md5.New()
	md5.Write([]byte(data))
	md5Sum := hex.EncodeToString(md5.Sum(nil))

	return sha1, sha256base64, md5Sum, nil
}
