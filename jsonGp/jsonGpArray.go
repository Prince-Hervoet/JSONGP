package jsongp

import "strings"

type JsongpObjArray struct {
	objs []*JsongpObj
}

func GetJsongpObjArray() *JsongpObjArray {
	return &JsongpObjArray{
		objs: make([]*JsongpObj, 0),
	}
}

func (ja *JsongpObjArray) Get(index int) *JsongpObj {
	if index < 0 || index > len(ja.objs) {
		return nil
	}
	return ja.objs[index]
}

func (ja *JsongpObjArray) Add(jg *JsongpObj) {
	ja.objs = append(ja.objs, jg)
}

func (ja *JsongpObjArray) Stringify() string {
	ans := &strings.Builder{}
	ans.WriteRune('[')
	for i := 0; i < len(ja.objs); i++ {
		jg := ja.objs[i]
		str := jg.Stringify()
		ans.WriteString(str)
		if i < len(ja.objs)-1 {
			ans.WriteRune(',')
		}
	}
	ans.WriteRune(']')
	return ans.String()
}
