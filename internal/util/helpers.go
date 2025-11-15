package util

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"math/rand"
	fs "mmplat/internal/filesystem"
	"regexp"
	"strconv"
	"time"
	"unsafe"
)

func LogStringToLevel(str string) logrus.Level {
	switch str {
		case "info":
			return logrus.InfoLevel
		case "debug":
			return logrus.DebugLevel
		default:
			return logrus.ErrorLevel
	}
}

// PrepareTemplateItem templ[name,size,type]
func PrepareTemplateItem(id fs.Id, item fs.IItem) map[string]string {
	templateItem := make(map[string]string)
	templateItem["name"] = item.Name()
	templateItem["size"] = strconv.Itoa(int(item.Size()))
	templateItem["type"] = ExtToMetadata(item)
	templateItem["id"] = strconv.Itoa(int(id))
	// TODO obfuscate actual file behind exposed API interface
	templateItem["path"] = item.Path()
	return templateItem
}

func ValidateUserInput(input string) (validated string) {
	var pattern = `~\.|\\|\(|\)|\{|\}~`
	reg, err := regexp.Compile(pattern)
	if err != nil {
		panic("regex compilation error")
	}
	return reg.ReplaceAllString(input, pattern)
}

func CheckRange(ctx *fasthttp.RequestCtx) int {
	bytes := ctx.Request.Header.Peek("Range")
	if len(bytes) > 0 {
		res, _ := strconv.Atoi(string(bytes[6 : len(bytes)-1]))
		return res
	}
	return -1
}

func Empty(v any) bool {
	switch v.(type) {
	case string:
		return len(v.(string)) == 0
	case []string:
		return len(v.([]string)) == 0
	case []byte:
		return len(v.([]byte)) == 0
	case nil:
		return true
	default:
		panic("helpers: error: type unspecified")
	}
}

func Not(v bool) bool {
	return !v
}

func NotEmpty(v any) bool {
	return !Empty(v)
}

func Equals(s1, s2 string) bool {
	return s1 == s2
}

func SetTitle(ctx *fasthttp.RequestCtx, str string) {
	fmt.Fprintf(ctx, "<head><title>%s</title></head>", str)
}

// RandStr length.(int)
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandStr(n int) string {
	if n == 0 {
		n = DefRndLen
	}
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	var src = rand.NewSource(time.Now().UnixNano())
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
