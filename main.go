package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {
	fileName, timeout := readArguments()
	problems, total := readCSV(fileName)

	//opening statements
	fmt.Println("Welcome to the Quiz!")
	fmt.Println("We will present to you your questions now.")
	time.Sleep(time.Second)
	fmt.Println("Answer wisely!")
	time.Sleep(time.Second)
	fmt.Println()

	timer := time.NewTimer(time.Duration(timeout) * time.Second)
	correct := evaluate(problems, timer)

	time.Sleep(time.Second)
	fmt.Printf("Your score: %v/%v. Well played!\n", correct, total)
}

func evaluate(problems []problem, timer *time.Timer) int {
	correct := 0

	for i, p := range problems {
		fmt.Printf("Question %v: %v:\n", i+1, p.question)

		answerChannel := make(chan string)

		/*
			running a goroutine here, because if we run a goroutine of the entire
			evaluate function, and use select statement in main() to race the
			evaluate function	against the timer, it will continue running the
			goroutine of the entire	loop even after the timer runs out (remember, the
			select statement will	not close the goroutines even after they've been
			exited from) Thus the evaluate goroutine will keep running in the
			background, waiting to run when there's a blocking call on the main
			routine, in which case it will continue to take input and prcoess that
			input. And we do have a blocking call in the form of the sleep call kept
			for presentation
		*/
		go func() {
			var answer string
			fmt.Scanln(&answer)
			fmt.Println()
			answerChannel <- answer
		}()

		/*
			passing the answer in the channel instead of updating the correct counter
			in the gorooutine itself, is done to avoid any extra inputs after timer
			timeout to be processed within the goroutine and affect the final count
			of correct answers. the answer should be processed synchronously, not asynchronously
		*/

		select {
		case <-timer.C:
			fmt.Println("\nTime's up!")
			return correct
		case answer := <-answerChannel:
			if answer == p.answer {
				correct++
			}
		}
	}
	return correct
}

// read csv file or problems and return slice of problem structs
func readCSV(f string) ([]problem, int) {
	file, err := os.Open(f)
	checkErr(err, "Error reading file")
	defer file.Close()

	csvReader := csv.NewReader(file)
	lines, err := csvReader.ReadAll()
	checkErr(err, "Error while parsing CSV file. Make sure it is properly formatted!")
	return parseProblems(lines), len(lines)
}

// parse array of records from CSV file to array of problem types
func parseProblems(lines [][]string) []problem {
	problems := []problem{}
	for _, line := range lines {
		p := problem{
			question: line[0],
			answer:   line[1],
		}
		problems = append(problems, p)
	}
	return problems
}

// read the flag value for filename
func readArguments() (string, int) {
	fileName := flag.String(
		"filename",
		"problems.csv",
		"Name of the CSV file with the problems and the solutions",
	)

	timer := flag.Int(
		"timer",
		30,
		"Time limit for each question (in seconds)",
	)

	flag.Parse()

	return *fileName, *timer
}

// error handling
func checkErr(e error, msg string) {
	if e != nil {
		fmt.Println(msg)
		os.Exit(1)
	}
}
