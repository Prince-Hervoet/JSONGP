package jsongp

import "strings"

const (
	_OBJECT_BEGIN = 10
	_OBJECT_END   = 20
	_ARRAY_BEGIN  = 30
	_ARRAY_END    = 40
	_STRING       = 50
	_COLON        = 60
	_COMMA        = 70
	_BOOL         = 80
	_NULL         = 90
)

type token struct {
	data      string
	tokenType int
}

func parseJsonStr(jsonStr string) {
	jr := NewJsongpReader(jsonStr)
	tokens := make([]*token, 0)
	for jr.HasNext() {
		t := judgeTokenType(jr)
		if t == nil {
			return
		}
		tokens = append(tokens, t)
	}
}

func judgeTokenType(reader *JsongpReader) *token {
	ru := reader.Peek()
	switch ru {
	case '[':
		return &token{data: "[", tokenType: _ARRAY_BEGIN}
	case ']':
		return &token{data: "]", tokenType: _ARRAY_END}
	case '{':
		return &token{data: "{", tokenType: _OBJECT_BEGIN}
	case '}':
		return &token{data: "}", tokenType: _OBJECT_END}
	case '"':
		return judgeString(reader)
	case ':':
		return &token{data: ":", tokenType: _COLON}
	case ',':
		return &token{data: ",", tokenType: _COMMA}
	case 't':
		return judgeBool(reader)
	case 'f':
		return judgeBool(reader)
	case 'n':
		return judgeNull(reader)
	case '-':

	}
	return nil
}

func judgeBool(reader *JsongpReader) *token {
	if !reader.HasNext() {
		return nil
	}
	ru := reader.Next()
	target := ""
	if ru == 't' {
		target = "ture"
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

func judgeNull(reader *JsongpReader) *token {
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

func judgeNumber(reader *JsongpReader) *token {
	if !reader.HasNext() {
		return nil
	}
	ru := reader.Next()
	if ru != '-' && (ru <= '0' || ru > '9') {
		return nil
	}
	isPoint := true
	sbu := &strings.Builder{}
	sbu.WriteRune(ru)
	for reader.HasNext() {
		ru := reader.Peek()
		if isEmptyCh(ru) {
			break
		} else if ru >= '0' && ru <= '9' {
			sbu.WriteRune(ru)
		} else if ru == '.' {
			if isPoint && sbu.Len() > 1 {
				sbu.WriteRune(ru)
				isPoint = false
			} else {
				return nil
			}
		}
	}
	if ru == '-' && sbu.Len() < 2 {
		return nil
	}
	return nil
}

func judgeString(reader *JsongpReader) *token {
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
		if ru == '"' {
			isStr = true
			break
		} else if ru == '\\' {
			rux := reader.Peek()
			if !isEscapeCh(rux) {
				return nil
			}

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
