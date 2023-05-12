package jsong

import "errors"

type PackStringReader struct {
	jsonStr []rune
	flag    int
}

func NewJsonStrReader(source string) *PackStringReader {
	return &PackStringReader{
		jsonStr: []rune(source),
		flag:    0,
	}
}

func (jr *PackStringReader) Next() (rune, error) {
	if jr.flag == len(jr.jsonStr) {
		return 0, errors.New("data empty")
	}
	ans := jr.jsonStr[jr.flag]
	jr.flag += 1
	return ans, nil
}

func (jr *PackStringReader) HasNext() bool {
	return jr.flag < len(jr.jsonStr)
}

func (jr *PackStringReader) PeekNext() (rune, error) {
	if jr.flag == len(jr.jsonStr) {
		return 0, errors.New("data empty")
	}
	return jr.jsonStr[jr.flag], nil
}
