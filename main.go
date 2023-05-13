package main

import (
	"fmt"
	jsongp "json_g/jsonGp"
)

func main() {
	jgo := jsongp.GetJsonGpObj()
	jgo2 := jsongp.GetJsonGpObj()
	jgo3 := jsongp.GetJsonGpObj()
	jgo.Set("asdf", "asdfasdfasdf331")
	jgo.Set("123", 22134)
	jgo.Set("123as", false)
	jgo2.Set("asdlkfjasildfhj", 12312421345.56585858)
	jgo2.Set("asdfa23", "asdfoiq")
	jgo3.Set("opppp", true)
	jgo.Set("jgo2", jgo2)
	jgo2.Set("jgo3", jgo3)
	jgo.Set("jgo3", jgo3)
	s := jgo.Stringify()
	// var a any
	// a = jgo
	fmt.Println(s)
}
