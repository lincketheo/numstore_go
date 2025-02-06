package utils

func IsAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func IsDigit(c byte) bool {
	return c <= '9' && c >= '0'
}

func IsAlphaNum(c byte) bool {
	return IsDigit(c) || IsAlpha(c)
}
