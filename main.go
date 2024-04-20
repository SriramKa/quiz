package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	fileName := readArguments()
	problems := readCSV(fileName)

	for {
		problem, err := problems.Read()
		if err == io.EOF {
			break
		}
		checkErr(err)
		fmt.Println(problem)
	}
}

func readArguments() string {
	fileName := flag.String(
		"filename",
		"problems.csv",
		"Name of the CSV file with the problems and the solutions",
	)

	flag.Parse()

	return *fileName
}

func readCSV(f string) *csv.Reader {
	file, err := os.Open(f)
	checkErr(err)
	return csv.NewReader(file)
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
