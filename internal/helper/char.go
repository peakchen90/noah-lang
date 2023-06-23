package helper

func IsIdentifierChar(ch rune, first bool) bool {
	if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || ch == '_' || ch == '$' {
		return true
	}
	if !first && (ch >= '0' && ch <= '9') {
		return true
	}
	return false
}
