package jsongp

import (
	"strings"
)

const (
	// 对象开始
	_OBJECT_BEGIN = 10
	// 对象结束
	_OBJECT_END = 20
	// 数组开始
	_ARRAY_BEGIN = 30
	// 数组结束
	_ARRAY_END = 40
	_STRING    = 50
	_COLON     = 60
	_COMMA     = 70
	_BOOL      = 80
	_NULL      = 90
	_NUMBER    = 100
)

const (
	// 前导零
	_NUMBER_FIRST_ZERO = 1
	// 正常零
	_NUMBER_MULTI_ZERO = 2
	// 中划线
	_NUMBER_MINUS = 4
	// 小数点
	_NUMBER_POINT = 8
	// 正常数字，不包含0
	_NUMBER_COMMON      = 16
	_NUMBER_POINT_FALSE = 32
)

type token struct {
	data      string
	tokenType int
}

func ParseTokens(jsonStr string) []*token {
	jr := NewJsongpReader(jsonStr)
	tokens := make([]*token, 0)
	for jr.HasNext() {
		cur := jr.Peek()
		if isEmptyCh(cur) {
			jr.Next()
			continue
		}
		t := judgeTokenType(jr)
		if t == nil {
			return nil
		}
		tokens = append(tokens, t)
	}
	return tokens
}

func judgeTokenType(reader *JsongpReader) *token {
	ru := reader.Peek()
	switch ru {
	case '[':
		reader.Next()
		return &token{data: "[", tokenType: _ARRAY_BEGIN}
	case ']':
		reader.Next()
		return &token{data: "]", tokenType: _ARRAY_END}
	case '{':
		reader.Next()
		return &token{data: "{", tokenType: _OBJECT_BEGIN}
	case '}':
		reader.Next()
		return &token{data: "}", tokenType: _OBJECT_END}
	case '"':
		return JudgeString(reader)
	case ':':
		reader.Next()
		return &token{data: ":", tokenType: _COLON}
	case ',':
		reader.Next()
		return &token{data: ",", tokenType: _COMMA}
	case 't':
		return JudgeBool(reader)
	case 'f':
		return JudgeBool(reader)
	case 'n':
		return JudgeNull(reader)
	case '-':
		return JudgeNumber(reader)
	default:
		if isChNumber(ru) {
			return JudgeNumber(reader)
		} else {
			return nil
		}
	}
}

func JudgeBool(reader *JsongpReader) *token {
	if !reader.HasNext() {
		return nil
	}
	ru := reader.Next()
	target := ""
	if ru == 't' {
		target = "true"
	} else if ru == 'f' {
		target = "false"
	} else {
		return nil
	}
	for i := 1; i < len(target); i++ {
		if !reader.HasNext() {
			return nil
		}
		ru := reader.Next()
		if ru != rune(target[i]) {
			return nil
		}
	}
	return &token{data: target, tokenType: _BOOL}
}

func JudgeNull(reader *JsongpReader) *token {
	if !reader.HasNext() {
		return nil
	}
	ru := reader.Next()
	if ru != 'n' {
		return nil
	}
	target := "null"
	for i := 1; i < len(target); i++ {
		if !reader.HasNext() {
			return nil
		}
		ru := reader.Next()
		if ru != rune(target[i]) {
			return nil
		}
	}
	return &token{data: target, tokenType: _NULL}
}

func JudgeNumber(reader *JsongpReader) *token {
	if !reader.HasNext() {
		return nil
	}
	sbu := &strings.Builder{}
	end := false
	expect := _NUMBER_MINUS | _NUMBER_COMMON | _NUMBER_FIRST_ZERO
	for reader.HasNext() {
		ru := reader.Peek()
		switch ru {
		case ',':
			end = true
		case '-':
			if checkStatus(_NUMBER_MINUS, expect) {
				sbu.WriteRune('-')
				expect = _NUMBER_COMMON | _NUMBER_FIRST_ZERO
			} else {
				return nil
			}
		case '0':
			if checkStatus(_NUMBER_FIRST_ZERO, expect) {
				sbu.WriteRune(ru)
				expect = _NUMBER_POINT
			} else if checkStatus(_NUMBER_MULTI_ZERO, expect) {
				sbu.WriteRune(ru)
				expect = _NUMBER_COMMON | _NUMBER_MULTI_ZERO
			} else {
				return nil
			}
		case '.':
			if checkStatus(_NUMBER_POINT, expect) {
				sbu.WriteRune(ru)
				expect = _NUMBER_COMMON | _NUMBER_MULTI_ZERO | _NUMBER_POINT_FALSE
			} else {
				return nil
			}
		default:
			if isChNumber(ru) && checkStatus(_NUMBER_COMMON, expect) {
				if checkStatus(_NUMBER_POINT_FALSE, expect) {
					expect = _NUMBER_COMMON | _NUMBER_MULTI_ZERO | _NUMBER_POINT_FALSE
				} else {
					expect = _NUMBER_COMMON | _NUMBER_MULTI_ZERO | _NUMBER_POINT
				}
				sbu.WriteRune(ru)
			} else {
				end = true
			}
		}
		if end {
			break
		}
		reader.Next()
	}
	ans := sbu.String()
	if !isChNumber(rune(ans[len(ans)-1])) {
		return nil
	}
	return &token{data: ans, tokenType: _NUMBER}
}

func JudgeString(reader *JsongpReader) *token {
	if !reader.HasNext() {
		return nil
	}
	ru := reader.Next()
	if ru != '"' {
		return nil
	}
	isStr := false
	sbu := &strings.Builder{}
	for reader.HasNext() {
		ru := reader.Next()
		if ru == '\\' {
			switch ru {
			case '"':
				sbu.WriteRune('"')
			case '\\':
				sbu.WriteRune('\\')
			case '/':
				sbu.WriteRune('/')
			case 'b':
				sbu.WriteRune('\b')
			case 'f':
				sbu.WriteRune('\f')
			case 'n':
				sbu.WriteRune('\n')
			case 'r':
				sbu.WriteRune('\r')
			case 't':
				sbu.WriteRune('\t')
			case 'u':
				u := 0
				for i := 0; i < 4; i++ {
					uch := reader.Next()
					if uch >= '0' && uch <= '9' {
						u = (u << 4) + (int(uch) - '0')
					} else if uch >= 'a' && uch <= 'f' {
						u = (u << 4) + (int(uch) - 'a') + 10
					} else if uch >= 'A' && uch <= 'F' {
						u = (u << 4) + (int(uch) - 'A') + 10
					} else {
						return nil
					}
				}
			default:
				return nil
			}
		} else if ru == '"' {
			isStr = true
			break
		} else if ru == '\r' || ru == '\n' {
			return nil
		} else {
			sbu.WriteRune(ru)
		}
	}
	if isStr {
		return &token{data: sbu.String(), tokenType: _STRING}
	}
	return nil
}

func isEscapeCh(ru rune) bool {
	return ru == 't' || ru == 'n' || ru == '"' || ru == 'u' || ru == '\\'
}

func isEmptyCh(ru rune) bool {
	return ru == ' ' || ru == '\n' || ru == '\t'
}

func isChNumber(ru rune) bool {
	return ru >= '0' && ru <= '9'
}

func checkStatus(current, target int) bool {
	return !((current & target) == 0)
}
