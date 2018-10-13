package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// our main function
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/test", getPostWithID).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getPostWithID(w http.ResponseWriter, r *http.Request) {
	var x map[string]interface{}
	resp, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal([]byte(resp), &x)
	//out, _ := json.Marshal(x)
	b, err := json.MarshalIndent(&x, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	//fmt.Println(b)
	//fmt.Println(x)
}
