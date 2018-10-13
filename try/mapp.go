package his

import (
	"fmt"
)

func main() {

	a := []string{"apple", "applause"}
	b := []string{"banana","balista"}
	c := []string{"catherine","catal#yst"}

	x := map[string][]string{
		"a": a,
		"b": b,
		"c": c,
	}

	if test,ok := x["b"]; ok{
		fmt.Println("WOw ",test)
	}

	for key, value := range x {
		
		for i := 0 ; i < len(value) ; i++{
			fmt.Println(key, " and ", value[i])

			switch value[i]{
			case "apple":
				fmt.Println("It's an apple!")
				break
			default:
				fmt.Println(value[i]," is not a type for checking")
			}
			

		}
			
	}

}
