package main

import (
	"fmt"
	"io/ioutil"
	//    "log"
	"bytes"
	//"encoding/json"
	//"math/rand"
	"net/http"
	//"strings"
	//"time"
)

func main() {

	dat, err := ioutil.ReadFile("./registration.json")
	fmt.Print(string(dat))
	//jsonVal, err := json.Marshal(new)
	//var prettyJSON bytes.Buffer
	//err = json.Indent(&prettyJSON, jsonVal, "", "\t")
	//fmt.Println("json:\n", string(prettyJSON.Bytes()))
	res, err := http.Post("https://140.118.123.214:5000/registration", "application/json", bytes.NewBuffer(dat))
	if err != nil {
		fmt.Println("Post failed")
		fmt.Println(err.Error())
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("Post return:\n", string(body))

}
