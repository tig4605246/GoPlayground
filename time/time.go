package main

import (
	"fmt"
	"time"
)

func main() {
	catchTime, _ := time.Parse("2006-01-02 15:04:05", "2018-08-28 10:23:05")
	//timeString := catchTime.Format("2006-01-02 15:04:05")
	timeUnix := catchTime.Unix()
	fmt.Println(timeUnix - 28800)
}
