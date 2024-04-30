package strings

import "encoding/base64"

// Base64Decode base64解析
func Base64Decode(str string) string {
	switch len(str) % 4 {
	case 2:
		str += "=="
	case 3:
		str += "="
	}

	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}

	return BytesToStr(b)
}

// Base64Encode base64编码
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString(StrToBytes(str))
}
