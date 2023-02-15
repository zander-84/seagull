package boot

import "strings"

// 大驼峰
func upperCamelCase(word string) string {
	names := strings.Split(word, "_")
	out := ""
	for _, v := range names {
		out += strFirstToUpper(v)
	}
	return out
}

// 小驼峰
func lowerCamelCase(word string) string {
	names := strings.Split(word, "_")
	out := ""
	for k, v := range names {
		if k == 0 {
			out += strFirstToLower(v)
		} else {
			out += strFirstToUpper(v)

		}
	}
	return out
}

func strFirstToUpper(str string) string {
	if str == "" {
		return ""
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

func strFirstToLower(str string) string {
	if str == "" {
		return ""
	}
	return strings.ToLower(str[:1]) + str[1:]
}
