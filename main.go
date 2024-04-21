package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	fileName := readArguments()
	file, err := os.Open(fileName)
	checkErr(err)
	defer file.Close()
	problems := csv.NewReader(file)
	// correct := 0

	fmt.Println("Welcome to the Quiz!")
	fmt.Println("We will present to you your questions now.")
	time.Sleep(time.Second)
	fmt.Println("Answer wisely!")
	time.Sleep(time.Second)

	for {
		problem, err := problems.Read()
		if err == io.EOF {
			break
		}
		checkErr(err)
		// evaluateProblem(problem)
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

// func evaluateProblem(p []string) bool {

// }

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
