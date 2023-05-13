package jsongp

import (
	"reflect"
	"strconv"
	"strings"
)

const (
	_TYPE_UNKONWN = ""
	_TYPE_STRING  = "string"
	_TYPE_NUMBER  = "number"
	_TYPE_BOOLEAN = "bool"
	_TYPE_NULL    = "null"
	_TYPE_JSONOBJ = "jsonObj"
	_TYPE_JSONARR = "jsonObjArray"
)

type JsongpObj struct {
	keyToNode map[string]*jsonGpNode
	flag      *jsonGpNode
	head      *jsonGpNode
	tail      *jsonGpNode
}

type jsonGpNode struct {
	key   string
	value any
	next  *jsonGpNode
	prev  *jsonGpNode
}

func GetJsonGpObj() *JsongpObj {
	return &JsongpObj{
		keyToNode: make(map[string]*jsonGpNode),
		flag:      nil,
		head:      nil,
		tail:      nil,
	}
}

func (jg *JsongpObj) Set(key string, value any) {
	if key == "" {
		return
	}
	if value != nil {
		tName := reflect.TypeOf(value).Name()
		if tName == "" {
			eName := reflect.TypeOf(value).Elem().Name()
			if eName != "JsongpObj" && eName != "JsongpObjArray" {
				return
			}
		}
	}
	if _, has := jg.keyToNode[key]; has {
		jg.keyToNode[key].value = value
		return
	}
	n := &jsonGpNode{
		key:   key,
		value: value,
		next:  nil,
		prev:  nil,
	}
	if jg.head == nil {
		jg.head = n
		jg.tail = jg.head
		jg.flag = n
	} else {
		jg.tail.next = n
		n.prev = jg.tail
		jg.tail = n
	}
	jg.keyToNode[key] = n
}

func (jg *JsongpObj) Get(key string) any {
	if _, has := jg.keyToNode[key]; !has {
		return nil
	}
	n := jg.keyToNode[key]
	return n.value
}

func (jg *JsongpObj) Remove(key string) {
	if _, has := jg.keyToNode[key]; !has {
		return
	}
	target := jg.keyToNode[key]
	if target == jg.head {
		temp := target.next
		target.next = nil
		if temp != nil {
			temp.prev = nil
		}
		jg.head = temp
		if jg.head == nil {
			jg.tail = jg.head
		}
	} else if target == jg.tail {
		temp := target.prev
		temp.next = nil
		target.prev = nil
		jg.tail = temp
	} else {
		p1 := target.prev
		p2 := target.next
		p1.next = p2
		p2.prev = p1
		target.next = nil
		target.prev = nil
	}
	delete(jg.keyToNode, key)
}

func (jg *JsongpObj) ResetNext() {
	jg.flag = jg.head
}

func (jg *JsongpObj) Stringify() string {
	run := jg.head
	ans := &strings.Builder{}
	ans.WriteRune('{')
	for run != nil {
		key := run.key
		value := run.value
		ans.WriteRune('"')
		ans.WriteString(key)
		ans.WriteRune('"')
		ans.WriteRune(':')
		if run.value == nil {
			ans.WriteString("null")
			run = run.next
			continue
		}
		if reflect.TypeOf(run.value).Name() != "" {
			switch reflect.TypeOf(run.value).Name() {
			case "string":
				ans.WriteRune('"')
				ans.WriteString(value.(string))
				ans.WriteRune('"')
			case "bool":
				boo := value.(bool)
				if boo {
					ans.WriteString("true")
				} else {
					ans.WriteString("false")
				}
			case "int":
				temp := run.value.(int)
				ans.WriteString(strconv.FormatInt(int64(temp), 10))
			case "int8":
				temp := run.value.(int8)
				ans.WriteString(strconv.FormatInt(int64(temp), 10))
			case "int16":
				temp := run.value.(int16)
				ans.WriteString(strconv.FormatInt(int64(temp), 10))
			case "int32":
				temp := run.value.(int32)
				ans.WriteString(strconv.FormatInt(int64(temp), 10))
			case "int64":
				temp := run.value.(int64)
				ans.WriteString(strconv.FormatInt(temp, 10))
			case "float32":
				temp := run.value.(float32)
				ans.WriteString(strconv.FormatFloat(float64(temp), 'f', 5, 64))
			case "float64":
				temp := run.value.(float64)
				ans.WriteString(strconv.FormatFloat(temp, 'f', 5, 64))
			}
		} else {
			tName := reflect.TypeOf(run.value).Elem().Name()
			if tName == reflect.TypeOf(jg).Elem().Name() {
				jg := run.value.(*JsongpObj)
				str := jg.Stringify()
				ans.WriteString(str)
			} else if tName == "JsongpObjArray" {
				ja := run.value.(*JsongpObjArray)
				str := ja.Stringify()
				ans.WriteString(str)
			}
		}
		if run != jg.tail {
			ans.WriteRune(',')
		}
		run = run.next
	}
	ans.WriteRune('}')
	return ans.String()
}

func (jg *JsongpObj) Size() int {
	return len(jg.keyToNode)
}

func isTypeInt(typeName string) bool {
	return typeName == "int" || typeName == "int32" || typeName == "int64" || typeName == "int8" || typeName == "uint8" || typeName == "uint32" || typeName == "uint64"
}

func isTypeFloat(typeName string) bool {
	return typeName == "float32" || typeName == "float64"
}
