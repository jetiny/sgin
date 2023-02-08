package utils

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
)

func MakeUploadHash(uploadId, targetId int64) (string, error) {
	str := fmt.Sprintf("%d_$_%d", uploadId, targetId)
	h := md5.New()
	_, err := h.Write([]byte(str))
	if err != nil {
		return "", err
	}
	return string(base64.RawStdEncoding.EncodeToString(h.Sum([]byte("")))), nil
}

func CheckUploadHash(uploadId, targetId int64, hash string) (bool, error) {
	str, err := MakeUploadHash(uploadId, targetId)
	if err != nil {
		return false, err
	}
	return str == hash, nil
}

func CreateUploadToken(account string, ts int64, id int64, mod int8, path string) (string, error) {
	text := fmt.Sprintf("%d_%d_%d_%s", ts, id, mod, path)
	return GenerateHMAC(text, account)
}

func CheckUploadToken(account string, ts int64, id int64, mod int8, path string, token string) (bool, error) {
	text := fmt.Sprintf("%d_%d_%d_%s", ts, id, mod, path)
	return VerifyHMAC(token, text, account)
}
