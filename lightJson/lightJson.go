package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type void struct{}

var voidVar void

const TRUE_STR = "true"
const FALSE_STR = "false"
const NULL_STR = "null"
const OPENING_BRACE int16 = 1
const CLOSING_BRACE int16 = 1 << 1
const OPENING_BRACKET int16 = 1 << 2
const CLOSING_BRACKET int16 = 1 << 3
const DOUBLE_QUOTATION_MARKS int16 = 1 << 4
const COLON int16 = 1 << 5
const COMMA int16 = 1 << 6
const FULL_STOP int16 = 1 << 7
const MINUS int16 = 1 << 8
const NUMBER int16 = 1 << 9
const WHITE_SPACE int16 = 1 << 10
const TRUE = 1 << 11
const FALSE = 1 << 12
const NULL = 1 << 13

const ALL = NUMBER | MINUS | DOUBLE_QUOTATION_MARKS | TRUE | FALSE | NULL | OPENING_BRACE | OPENING_BRACKET | CLOSING_BRACKET | WHITE_SPACE

// DOUBLE_QUOTATION_MARKS | NUMBER | MINUS | TRUE | FALSE | NULL | OPENING_BRACE | OPENING_BRACKET | WHITE_SPACE
var expectMap = map[rune]int16{
	'{': OPENING_BRACE,
	'}': CLOSING_BRACE,
	'[': OPENING_BRACKET,
	']': CLOSING_BRACKET,
	'"': DOUBLE_QUOTATION_MARKS,
	':': COLON,
	',': COMMA,
	'.': FULL_STOP,
	'-': MINUS,
	'0': NUMBER,
	'1': NUMBER,
	'2': NUMBER,
	'3': NUMBER,
	'4': NUMBER,
	'5': NUMBER,
	'6': NUMBER,
	'7': NUMBER,
	'8': NUMBER,
	'9': NUMBER,
	' ': WHITE_SPACE,
	't': TRUE,
	'f': FALSE,
	'n': NULL,
}

var hexDigitSet = map[rune]struct{}{
	'0': voidVar,
	'1': voidVar,
	'2': voidVar,
	'3': voidVar,
	'4': voidVar,
	'5': voidVar,
	'6': voidVar,
	'7': voidVar,
	'8': voidVar,
	'9': voidVar,
	'A': voidVar,
	'B': voidVar,
	'C': voidVar,
	'D': voidVar,
	'E': voidVar,
	'F': voidVar,
	'a': voidVar,
	'b': voidVar,
	'c': voidVar,
	'd': voidVar,
	'e': voidVar,
	'f': voidVar,
}

var escapeSet = map[rune]struct{}{
	'"':  voidVar,
	'\\': voidVar,
	'/':  voidVar,
	'b':  voidVar,
	'f':  voidVar,
	'n':  voidVar,
	'r':  voidVar,
	't':  voidVar,
	'u':  voidVar,
}

type processStr struct {
	str    string
	runes  []rune
	cursor int
	expect int16
}

func (ps *processStr) next() bool {
	if ps.cursor+1 >= len(ps.runes) {
		ps.cursor = len(ps.runes)
		return false
	}
	ps.cursor += 1
	return true
}

func (ps *processStr) current() (rune, error) {
	if ps.cursor >= len(ps.runes) {
		return 0, errors.New("invalid index")
	}
	return ps.runes[ps.cursor], nil
}

func (ps *processStr) isExpect(current int16) bool {
	return (current & ps.expect) != 0
}

func (ps *processStr) setExpect(value int16) {
	ps.expect = value
}

func solveString(ps *processStr) {
	c, err := ps.current()
	if err != nil {
		return
	}
	currentExpect, has := expectMap[c]
	if !has || !ps.isExpect(currentExpect) {
		return
	}
	if c == '{' {
		parseObject(ps)
	} else if c == '[' {
		ans, err := parseArray(ps)
		fmt.Println(ans)
		fmt.Println(err)
	}
}

func parseObject(ps *processStr) (map[string]any, error) {
	ans := make(map[string]any)
	ps.setExpect(DOUBLE_QUOTATION_MARKS | WHITE_SPACE | CLOSING_BRACE)
	ps.next()
	var c rune
	var err error
	var key string = ""
	for {
		c, err = ps.current()
		if err != nil {
			break
		}
		currentExpect, has := expectMap[c]
		if !has || !ps.isExpect(currentExpect) {
			err = errors.New("invalid rune")
			break
		}
		if c == '"' {
			str, err := parseString(ps)
			if err != nil {
				break
			}
			if key == "" {
				key = str
				ps.setExpect(COLON | WHITE_SPACE)
			} else {
				ans[key] = str
				ps.setExpect(COMMA | CLOSING_BRACE | WHITE_SPACE)
				key = ""
			}
		} else if c == ':' {
			ps.setExpect(DOUBLE_QUOTATION_MARKS | NUMBER | MINUS | TRUE | FALSE | NULL | OPENING_BRACE | OPENING_BRACKET | WHITE_SPACE)
			ps.next()
		} else if c == '-' || (c >= '0' && c <= '9') {
			isPos := true
			if c == '-' {
				ps.next()
				isPos = false
			}
			ansInt, ansFloat, isFloat, err := parseNumber(ps)
			if err != nil {
				break
			}
			if key == "" {
				err = errors.New("invalid json")
				break
			}
			if !isPos {
				ansInt = -ansInt
				ansFloat = -ansFloat
			}
			if isFloat {
				ans[key] = ansFloat
			} else {
				ans[key] = ansInt
			}
			ps.setExpect(COMMA | CLOSING_BRACE | WHITE_SPACE)
			key = ""
		} else if c == 't' {
			if !parseTrue(ps) {
				err = errors.New("invalid rune")
				break
			}
			if key == "" {
				err = errors.New("invalid json")
				break
			}
			ans[key] = true
			ps.setExpect(COMMA | CLOSING_BRACE | WHITE_SPACE)
			key = ""
		} else if c == 'f' {
			if !parseFalse(ps) {
				err = errors.New("invalid rune")
				break
			}
			if key == "" {
				err = errors.New("invalid json")
				break
			}
			ans[key] = false
			ps.setExpect(COMMA | CLOSING_BRACE | WHITE_SPACE)
			key = ""
		} else if c == 'n' {
			if !parseNull(ps) {
				err = errors.New("invalid rune")
				break
			}
			if key == "" {
				err = errors.New("invalid json")
				break
			}
			ans[key] = nil
			ps.setExpect(COMMA | CLOSING_BRACE | WHITE_SPACE)
			key = ""
		} else if c == ' ' {
			ps.next()
		} else if c == ',' {
			ps.setExpect(DOUBLE_QUOTATION_MARKS | WHITE_SPACE)
			ps.next()
		} else if c == '}' {
			break
		} else if c == '{' {
			value, err := parseObject(ps)
			if err != nil {
				break
			}
			if key == "" {
				err = errors.New("invalid json")
				break
			}
			ans[key] = value
			ps.setExpect(COMMA | CLOSING_BRACE | WHITE_SPACE)
			key = ""
		}
	}

	if err != nil {
		return nil, err
	}
	return ans, nil
}

func parseArray(ps *processStr) ([]any, error) {
	ans := make([]any, 0)
	ps.setExpect(NUMBER | MINUS | DOUBLE_QUOTATION_MARKS | TRUE | FALSE | NULL | OPENING_BRACE | OPENING_BRACKET | CLOSING_BRACKET | WHITE_SPACE)
	ps.next()
	var c rune
	var err error
	for {
		c, err = ps.current()
		if err != nil {
			return nil, err
		}
		currentExpect, has := expectMap[c]
		if !has || !ps.isExpect(currentExpect) {
			err = errors.New("invalid rune")
			break
		}
		if c == '"' {
			str, err := parseString(ps)
			if err != nil {
				break
			}
			ans = append(ans, str)
			ps.setExpect(COMMA | CLOSING_BRACKET | WHITE_SPACE)
		} else if c == '-' || (c >= '0' && c <= '9') {
			isPos := true
			if c == '-' {
				ps.next()
				isPos = false
			}
			ansInt, ansFloat, isFloat, err := parseNumber(ps)
			if err != nil {
				break
			}
			if !isPos {
				ansInt = -ansInt
				ansFloat = -ansFloat
			}
			if isFloat {
				ans = append(ans, ansFloat)
			} else {
				ans = append(ans, ansInt)
			}
			ps.setExpect(COMMA | CLOSING_BRACKET | WHITE_SPACE)
		} else if c == 't' {
			if !parseTrue(ps) {
				err = errors.New("invalid rune")
				break
			}
			ans = append(ans, true)
			ps.setExpect(COMMA | CLOSING_BRACKET | WHITE_SPACE)
		} else if c == 'f' {
			if !parseFalse(ps) {
				err = errors.New("invalid rune")
				break
			}
			ans = append(ans, false)
			ps.setExpect(COMMA | CLOSING_BRACKET | WHITE_SPACE)
		} else if c == 'n' {
			if !parseNull(ps) {
				err = errors.New("invalid rune")
				break
			}
			ans = append(ans, nil)
			ps.setExpect(COMMA | CLOSING_BRACKET | WHITE_SPACE)
		} else if c == ',' {
			ps.setExpect(NUMBER | MINUS | DOUBLE_QUOTATION_MARKS | TRUE | FALSE | NULL | OPENING_BRACE | OPENING_BRACKET | CLOSING_BRACKET | WHITE_SPACE)
			ps.next()
		} else if c == '{' {
			value, err := parseObject(ps)
			if err != nil {
				break
			}
			ans = append(ans, value)
			ps.setExpect(COMMA | CLOSING_BRACKET | WHITE_SPACE)
			ps.next()
		} else if c == ']' {
			break
		} else if c == '[' {
			value, err := parseArray(ps)
			if err != nil {
				break
			}
			ans = append(ans, value)
			ps.setExpect(COMMA | CLOSING_BRACKET | WHITE_SPACE)
			ps.next()
		} else if c == ' ' {
			ps.next()
		}
	}
	if err != nil {
		return nil, err
	}
	return ans, nil
}

func parseTrue(ps *processStr) bool {
	for i := 0; i < 4; i++ {
		c, err := ps.current()
		if err != nil || c != rune(TRUE_STR[i]) {
			return false
		}
		ps.next()
	}
	return true
}

func parseFalse(ps *processStr) bool {
	for i := 0; i < 5; i++ {
		c, err := ps.current()
		if err != nil || c != rune(FALSE_STR[i]) {
			return false
		}
		ps.next()
	}
	return true
}

func parseNull(ps *processStr) bool {
	for i := 0; i < 4; i++ {
		c, err := ps.current()
		if err != nil || c != rune(NULL_STR[i]) {
			return false
		}
		ps.next()
	}
	return true
}

func parseString(ps *processStr) (string, error) {
	ps.next()
	var c rune
	var err error
	builder := strings.Builder{}
	for {
		c, err = ps.current()
		if err != nil || c == '"' {
			break
		}
		if c == '\\' {
			builder.WriteRune(c)
			ps.next()
			c, err = ps.current()
			if err != nil || !isEscape(c) {
				break
			}
			if c == 'u' {
				for i := 0; i < 4; i++ {
					ps.next()
					c, err = ps.current()
					if err != nil {
						return "", err
					}
					if !isHexDigit(c) {
						return "", errors.New("invalid rune")
					}
					builder.WriteRune(c)
				}
			}
			ps.next()
		} else {
			builder.WriteRune(c)
			ps.next()
		}
	}
	if err != nil {
		return "", err
	}
	ans := builder.String()
	ps.next()
	return ans, nil
}

func parseNumber(ps *processStr) (int64, float64, bool, error) {
	c, err := ps.current()
	if err != nil {
		return 0, 0, false, err
	}
	if c == '0' {
		ps.setExpect(FULL_STOP)
	} else if c >= '1' && c <= '9' {
		ps.setExpect(NUMBER | FULL_STOP)
	} else {
		return 0, 0, false, errors.New("invalid rune")
	}
	builder := strings.Builder{}
	isFloat := false
	builder.WriteRune(c)
	ps.next()

	for {
		c, err = ps.current()
		if err != nil || ((c < '0' || c > '9') && c != '.') {
			break
		}
		currentExpect, has := expectMap[c]
		if !has || !ps.isExpect(currentExpect) {
			err = errors.New("invalid rune")
			break
		}
		builder.WriteRune(c)
		if c == '.' {
			isFloat = true
			ps.setExpect(NUMBER)
		}
		ps.next()
	}
	if err != nil {
		return 0, 0, false, err
	}
	str := builder.String()
	if isFloat {
		ans, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0, 0, false, err
		}
		return 0, ans, true, nil
	}
	ans, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, 0, false, err
	}
	return ans, 0, false, nil
}

func isEscape(c rune) bool {
	_, has := escapeSet[c]
	return has
}

func isHexDigit(c rune) bool {
	_, has := hexDigitSet[c]
	return has
}

func Parse(str string) {
	ps := &processStr{
		str:    str,
		runes:  []rune(str),
		cursor: 0,
		expect: OPENING_BRACE | OPENING_BRACKET,
	}
	solveString(ps)
}

func main() {
	start := time.Now().UnixNano()
	testStr := "[ {\"name\":\"小红\" ,\"age\":30,\"goal\":100,\"stuId\":\"3120004242\"},{\"name\":\"小红\",\"age\":30,\"goal\":100,\"stuId\":\"3120004242\"},{\"name\":\"小红\",\"age\":30,\"goal\":100,\"stuId\":\"3120004242\"},{\"name\":\"小红\" ,\"age\":30,\"goal\":100,\"stuId\":\"3120004242\"},{\"name\":\"小红\" ,\"age\":30,\"goal\":100,\"stuId\":\"3120004242\"},{\"name\":\"小红\" ,\"age\":30,\"goal\":100,\"stuId\":\"3120004242\"}]"
	Parse(testStr)
	end := time.Now().UnixNano()
	fmt.Print("耗时: ")
	fmt.Print(float64(end-start) / float64(1000000000))
	fmt.Println("s")
}
