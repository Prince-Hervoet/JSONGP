package jsong

import "strings"

const (
	EMPTY        = 10
	OBJECT_BEGIN = 100
	OBJECT_END   = -100
	ARRAY_BEGIN  = 200
	ARRAY_END    = -200
	STRING       = 300
	BOOLEAN      = 400
	NUMBER       = 500
	NULL         = 600
	COLON        = 700
	COMMA        = 800
	ERROR        = -1

	TRUE     = "true"
	FALSE    = "false"
	NULL_STR = "null"
)

func judgeTokenType(ch rune) int {
	if isEmptyChar(ch) {
		return EMPTY
	}
	switch ch {
	case '{':
		return OBJECT_BEGIN
	case '}':
		return OBJECT_END
	case '[':
		return ARRAY_BEGIN
	case ']':
		return ARRAY_END
	case ':':
		return COLON
	case ',':
		return COMMA
	case '"':
		return STRING
	case 't':
		return BOOLEAN
	case 'f':
		return BOOLEAN
	case 'n':
		return NULL
	case '-':
		return NUMBER
	default:
		if isCharNumber(ch) {
			return NUMBER
		} else {
			return ERROR
		}
	}
}

func judgeBooleanStr(source []rune, target bool, start int) bool {
	if len(source) == 0 {
		return false
	}
	targetStr := FALSE
	if target {
		targetStr = TRUE
	}
	if len(source)-start+1 < len(targetStr) {
		return false
	}
	k := 0
	for i := start; i < len(source); i++ {
		if source[i] != rune(targetStr[k]) {
			return false
		}
		k += 1
	}
	return true
}

func judgeNullStr(source []rune, start int) bool {
	if len(source) == 0 || len(source)-start+1 < 4 {
		return false
	}
	k := 0
	for i := start; i < len(source); i++ {
		if source[i] != rune(NULL_STR[k]) {
			return false
		}
	}
	return true
}

func judgeString(source []rune, start int) string {
	if len(source) == 0 || source[0] != '"' {
		return ""
	}
	sb := &strings.Builder{}
	for i := start + 1; i < len(source); i++ {
		switch source[i] {
		case '\\':
			if i+1 < len(source) && isEscape(source[i+1]) {
				sb.WriteRune(source[i])
				sb.WriteRune(source[i+1])
				if source[i] == 'u' {
					for k := i + 1; k <= i+5; k++ {
						if k == len(source) {
							return ""
						}
						if isHex(source[k]) {
							sb.WriteRune(source[k])
						} else {
							return ""
						}
					}
				}
			} else {
				return ""
			}
		}
	}
}

func isEmptyChar(ch rune) bool {
	if ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r' {
		return true
	}
	return false
}

func isCharNumber(ch rune) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}
	return false
}

func isEscape(ch rune) bool {
	return ch == '"' || ch == '\\' || ch == 'n' || ch == 't' || ch == 'r' || ch == 'u' || ch == 'f'
}

func isHex(ch rune) bool {
	return (ch >= '0' && ch <= '9') || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}
