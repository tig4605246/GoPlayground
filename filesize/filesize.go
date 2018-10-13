package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("./test")
	_, err = file.WriteString("1234")
	if err != nil {
		// handle the error here
		fmt.Println("dd")
		return
	}
	defer file.Close()

	// get the file size
	stat, err := file.Stat()
	if err != nil {
		return
	}

	fmt.Println("File size is ", stat.Size())

	//fd := file.Fd()
	file.Truncate(0)
	stat, err = file.Stat()
	fmt.Println("File size is ", stat.Size())

}
