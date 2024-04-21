package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

// creating struct to be able to close file safely
type problemsReader struct {
	csvReader *csv.Reader
	file      *os.File
}

type problem struct {
	question string
	answer   string
}

func main() {
	var correct, total int

	fileName := readArguments()
	problems := readCSV(fileName)
	defer problems.closeFile()

	//opening statements
	fmt.Println("Welcome to the Quiz!")
	fmt.Println("We will present to you your questions now.")
	time.Sleep(time.Second)
	fmt.Println("Answer wisely!")
	time.Sleep(time.Second)
	fmt.Println()

	for {
		problem, err := problems.readProblem()
		if err == io.EOF {
			break
		}
		checkErr(err, "Error parsing CSV file. Make sure it is a properly formatted CSV file!")

		total++

		var answer string
		fmt.Println(problem.question)
		fmt.Scanln(&answer)
		fmt.Println()

		if answer == problem.answer {
			correct++
		}
	}

	fmt.Println("Finished!")
	time.Sleep(time.Second)
	fmt.Printf("Your score: %v/%v. Well played!\n", correct, total)
}

// to open problems file and attach csv reader to read from file
func readCSV(f string) problemsReader {
	file, err := os.Open(f)
	checkErr(err, "Error reading file.")
	return problemsReader{
		csvReader: csv.NewReader(file),
		file:      file,
	}
}

// to read a problem using the csv reader on the opened problems file
func (p problemsReader) readProblem() (problem, error) {
	r, e := p.csvReader.Read()
	if r == nil {
		return problem{}, e
	} else {
		return problem{
			question: r[0],
			answer:   r[1],
		}, e
	}
}

// to close file opened to read problems
func (p problemsReader) closeFile() {
	p.file.Close()
}

// read the flag value for filename
func readArguments() string {
	fileName := flag.String(
		"filename",
		"problems.csv",
		"Name of the CSV file with the problems and the solutions",
	)

	flag.Parse()

	return *fileName
}

// error handling
func checkErr(e error, msg string) {
	if e != nil {
		fmt.Println(msg)
		os.Exit(1)
	}
}
