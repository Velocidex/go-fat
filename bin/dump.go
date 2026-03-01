package main

import (
	"encoding/json"
	"fmt"
)

func Dump(v interface{}) {
	serialized, _ := json.MarshalIndent(v, " ", " ")
	fmt.Printf("%v", string(serialized))
}
