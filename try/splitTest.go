package main

import (
	"fmt"
	"strings"
)

func main() {
	testString := "1234;2234"
	subString := strings.Split(testString, ";")
	fmt.Println(subString)
	fmt.Println("length is ", len(subString))
}
