package main

import(
	"fmt"
	"math/rand" 
)

const(
	BallNum = 100
	BoxList = 64
)

func main(){
	balls := BallNum
	var box  [64]int
	success := 0
	try := 0
	iMAX := 10
	rand.Seed(2)
	for balls > 0 && iMAX > 0{
		try++
		for i:=0 ; i < balls; i++{
			Boxnum := rand.Int() % (BoxList)
			box[Boxnum] = box[Boxnum] + 1
		} 
		for i:=0 ; i <  BoxList ; i++{
			if box[i] == 1{
				success = success + 1
				
			}
			box[i] = 0
		}
		//var rate float
		//rate = success/balls
		balls = balls - success
		fmt.Println("Trial ",try , " ",success, " balls success, ",balls, " balls failed", " iMAX = ",iMAX)
		success = 0
		iMAX = iMAX - 1
		//fmt.Println("balls ",balls)
	}
	fmt.Println("All balls success. Terminated")


}
