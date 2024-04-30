package formats

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// CheckMobile 验证手机号格式
func CheckMobile(str string) bool {
	regRule := "^1[345789]{1}\\d{9}$"
	reg := regexp.MustCompile(regRule)
	return reg.MatchString(str)
}

// EncodePhone 手机号脱敏
func EncodePhone(phone string) string {
	strLen := len(phone)
	if strLen < 3 {
		return "****"
	}

	if strLen <= 7 {
		return phone[0:3] + "****"
	} else {
		return phone[0:3] + "****" + phone[strLen-4:strLen]
	}
}

// EncodeEmail 邮箱加密
func EncodeEmail(email string) string {
	if !strings.Contains(email, "@") {
		return email
	}
	str := strings.Split(email, "@")
	eStr := ""
	//防止中文乱码
	str2 := []rune(str[0])
	count := utf8.RuneCountInString(str[0])
	if count < 4 {
		for i := 0; i < len(str[0]); i++ {
			eStr += "*"
		}
	} else {
		eStr += string(str2[0:count-4]) + "****"
	}
	return eStr + "@" + str[1]
}

// CheckEmail 验证邮箱
func CheckEmail(str string) bool {
	regRule := "^[_a-z0-9-]+(\\.[_a-z0-9-]+)*@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,})$"
	reg := regexp.MustCompile(regRule)
	return reg.MatchString(str)
}
