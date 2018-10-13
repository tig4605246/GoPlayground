package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	FileName = "Acsv.csv"
)

func main() {
	records := [][]string{
		{"first_name", "last_name", "username"},
		{"Rob", "Pike", "rob"},
		{"Ken", "Thompson", "ken"},
		{"Robert", "Griesemer", "gri"},
	}
	file, err := os.Create(FileName)
	//Is it opened ?
	if err != nil {
		fmt.Println("Fail to open file", err)
	}
	//Close file
	defer file.Close()
	w := csv.NewWriter(file)
	// Write any buffered data to the underlying writer (standard output).
	defer w.Flush()

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	//
	startPoint := time.Now().Unix() - 86400 - 28800
	fmt.Println(now)
}
