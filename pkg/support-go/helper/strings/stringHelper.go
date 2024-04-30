package strings

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/leeqvip/gophp/serialize"
	"reflect"
	"unsafe"
)

func Empty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// StrToBytes string转bytes
func StrToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// BytesToStr bytes转string
func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func SubStr(str string, offset int, length int) string {
	runeData := []rune(str)
	runeLen := len(runeData)

	if length < 0 {
		return ""
	}
	if (offset == 0 && length == runeLen) || runeLen <= length {
		return str
	}
	if offset >= 0 && length > 0 {
		limit := length + offset
		return string(runeData[offset:limit])
	}

	return ""
}

// UnSerialize 反序列化
func UnSerialize(str string) interface{} {
	out, _ := serialize.UnMarshal([]byte(str))
	if out == nil {
		return ""
	}
	return out
}

// Serialize 序列化
func Serialize(data interface{}) ([]byte, error) {
	out, err := serialize.Marshal(data)
	if out == nil {
		return nil, err
	}
	return out, err
}

func Md5(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))

	return hex.EncodeToString(hash.Sum(nil))
}
