package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/jetiny/gdash/passwd"
	uuid "github.com/satori/go.uuid"
)

func Uuid() string {
	id := uuid.NewV4().String()
	return id
}

func ShortUuid() string {
	id := uuid.NewV4().String()
	return id[:8]
}

func HashUuid() string {
	id := uuid.NewV4().String()
	return strings.Replace(id, "-", "", -1)
}

var flakeNode *snowflake.Node

func InitSnowflake(node int64) error {
	flake, err := snowflake.NewNode(node)
	if err != nil {
		return err
	}
	flakeNode = flake
	return nil
}

func init() {
	InitSnowflake(1)
}

func SnowId() int64 {
	return flakeNode.Generate().Int64()
}

func SnowCode() string {
	return flakeNode.Generate().String()
}

func SnowWithPrefix(prefix string, id *int64) string {
	n := int64(0)
	if id != nil {
		n = *id
	} else {
		n = SnowId()
	}
	return prefix + strconv.FormatInt(n, 36)
}

func CreatePassword(password string) (string, error) {
	str, err := passwd.Create(password)
	return str, err
}

func CheckPassword(origin, password string) (isValid bool, err error) {
	return passwd.Verify(origin, password)
}

func ValueTo[T any](value T) *T {
	return &value
}

func ValueOf[T any](value *T) T {
	return *value
}

func CreateCaptcha(num int) string {
	str := "1"
	for i := 0; i < num; i++ {
		str += strconv.Itoa(0)
	}
	str10 := str
	int10, err := strconv.ParseInt(str10, 10, 32)
	if err != nil {
		fmt.Println(err)
		return ""
	} else {
		j := int32(int10)
		return fmt.Sprintf("%0"+strconv.Itoa(num)+"v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(j))
	}
}

func VerifyHMAC(HMAC string, text string, key string) (bool, error) {
	nowHMAC, err := GenerateHMAC(text, key)
	if err != nil {
		return false, nil
	}
	return hmac.Equal([]byte(HMAC), []byte(nowHMAC)), nil
}

func GenerateHMAC(text string, key string) (string, error) {
	textBytes := []byte(text)
	keyBytes := []byte(key)
	hash := hmac.New(sha256.New, keyBytes)
	_, err := hash.Write(textBytes)
	if err != nil {
		return "", err
	}
	result := hash.Sum(nil)
	return string(result), nil
}

func FormatPhone(phone string) string {
	length := len(phone)
	if length > 0 { //189****8822
		start := 3
		end := start + int(length-start)/2
		if length < 6 {
			start = 1
			end = length - 1
		} else if length < 11 {
			start = 2
			end = length - 2
		}
		newStr := phone[:start] + strings.Repeat("*", end-start) + phone[end:]
		return newStr
	}
	return phone
}

func FormatEmail(email string) string {
	length := strings.Index(email, "@")
	if length > 0 { //8******4@qq.com
		start := 1
		end := length - 1
		if length > 7 {
			start = 2
			end = length - 2
		}
		newStr := email[:start] + strings.Repeat("*", end-start) + email[end:]
		return newStr
	}
	return email
}
