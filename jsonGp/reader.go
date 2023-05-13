package jsongp

type JsongpReader struct {
	content []rune
	flag    int
}

func NewJsongpReader(jsonStr string) *JsongpReader {
	return &JsongpReader{
		content: []rune(jsonStr),
		flag:    0,
	}
}

func (jr *JsongpReader) Next() rune {
	if jr.flag == len(jr.content) {
		return -1
	}
	ans := jr.content[jr.flag]
	jr.flag += 1
	return ans
}

func (jr *JsongpReader) HasNext() bool {
	return !(jr.flag == len(jr.content))
}

func (jr *JsongpReader) Peek() rune {
	if jr.flag == len(jr.content) {
		return -1
	}
	return jr.content[jr.flag]
}

func (jr *JsongpReader) Size() int {
	return len(jr.content)
}
