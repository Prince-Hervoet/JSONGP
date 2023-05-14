package main

import (
	"fmt"
	jsongp "json_g/jsonGp"
	"time"
)

func main() {
	start := time.Now().UnixMicro()
	t := jsongp.ParseTokens("{\"name\":123,\"huhu\":true}")
	jp := jsongp.ParseGrammar(t)
	end := time.Now().UnixMicro()
	fmt.Print("耗时: ")
	fmt.Print(end - start)
	fmt.Println(" ns")
	fmt.Println(jp.Get("huhu"))
}
