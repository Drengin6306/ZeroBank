package account

import "regexp"

func IsValidEmail(email string) bool {
	// 使用正则表达式检查邮箱格式
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
