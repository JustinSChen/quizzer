package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type Problem struct {
	question string
	answer string
}

// Parses problems from file. Shuffles them if the optional flag is provided
func parseProblems(problems [][]string, shuffle bool) []Problem {
	parsed := make([]Problem, len(problems))

	if shuffle {
		rand.Seed(time.Now().UTC().UnixNano())
		list := rand.Perm(len(problems))
		for i, p := range list {
			parsed[p] = Problem {
				question: problems[i][0],
				answer: problems[i][1],
			}
		}
	} else {
		for i, problem := range problems {
			parsed[i] = Problem {
				question: problem[0],
				answer: problem[1],
			}
		}
	}

	return parsed
}

func main() {
	// Setup flag options 
	file := flag.String("file", "problems.csv", "problems in 'question,answer' format. Default is \"problems.csv\"")
	limit := flag.Int("limit", 30, "the time limit for the quiz, in seconds. Default is 30 seconds")
	shuffle := flag.Bool("shuffle", false, "whether to shuffle the order of the quiz. Default is false")
	flag.Parse()

	// Open File
	f, err := os.Open(*file);
	if err != nil {
		log.Fatalf("Error while opening file: %v\n", err)
	}
	defer f.Close();

	// Read problems from file
	problems, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatalf("Error while reading problems: %v\n", err)
	}

	correct := 0
	timer := time.NewTimer(time.Duration(*limit) * time.Second)
	parsedProblems := parseProblems(problems, *shuffle)
	answerCh := make(chan string)

loop:
	for i, problem := range parsedProblems {
		fmt.Printf("Problem #%d: %s = \n", i + 1, problem.question)

		// Get user answer input
		go func() {
			answer := ""
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		// When answer is received, check if correct
		case a := <-answerCh:
			if a == problem.answer {
				correct++
				fmt.Println("Correct")
			} else {
				fmt.Println("Incorrect")
			}
		// Once time limit is reached, end quiz loop
		case <-timer.C:
			break loop
		}
	}
	close(answerCh)

	fmt.Printf("Your score is: %v out of %v\n", correct, len(parsedProblems))
}
