package jsongp

type pair struct {
	key   string
	value any
}

func ParseGrammar(tokens []*token) *JsongpObj {
	if len(tokens) == 0 {
		return nil
	}
	t := tokens[0].tokenType
	if t == _OBJECT_BEGIN {
		return action(tokens, 1)
	} else if t == _ARRAY_BEGIN {
	}
	return nil
}

func action(tokens []*token, index int) *JsongpObj {
	expect := _STRING
	jg := GetJsonGpObj()
	key := ""
	var value any
	for i := index; i < len(tokens); i++ {
		token := tokens[i]
		tt := token.tokenType
		switch tt {
		case _OBJECT_BEGIN:
			if checkStatus(tt, expect) {
				value = action(tokens, i+1)
				expect = _COMMA | _OBJECT_END
				jg.Set(key, value)
				key = ""
			}
		case _OBJECT_END:
			if checkStatus(tt, expect) {
				return jg
			}
		case _ARRAY_BEGIN:
			if checkStatus(tt, expect) {
				expect = _STRING
			}
		case _ARRAY_END:
			if checkStatus(tt, expect) {
				expect = _STRING
			}
		case _STRING:
			if checkStatus(tt, expect) {
				if key == "" {
					key = token.data
					expect = _COLON
				} else {
					value = token.data
					expect = _COMMA | _OBJECT_END
					jg.Set(key, value)
					key = ""
				}
			} else {
				return nil
			}
		case _COLON:
			if checkStatus(tt, expect) {
				expect = _STRING | _NUMBER | _ARRAY_BEGIN | _OBJECT_BEGIN | _BOOL | _NULL
			} else {
				return nil
			}
		case _COMMA:
			if checkStatus(tt, expect) {
				expect = _STRING
			} else {
				return nil
			}
		case _BOOL:
			if checkStatus(tt, expect) {
				value = token.data
				expect = _COMMA | _OBJECT_END
				jg.Set(key, value)
				key = ""
			} else {
				return nil
			}
		case _NULL:
			if checkStatus(tt, expect) {
				value = token.data
				expect = _COMMA | _OBJECT_END
				jg.Set(key, value)
				key = ""
			} else {
				return nil
			}
		case _NUMBER:
			if checkStatus(tt, expect) {
				value = token.data
				expect = _COMMA | _OBJECT_END
				jg.Set(key, value)
				key = ""
			} else {
				return nil
			}
		default:
			return nil
		}
	}
	return nil
}
