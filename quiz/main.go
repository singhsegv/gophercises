package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

type quizConfig struct {
	filename string
	shuffle  bool
	time     int
}

type problem struct {
	question, answer string
}

type quiz struct {
	problems         []problem
	totalProblems    int
	correctSolutions int
	shouldShuffle    bool
	timeout          int
}

var config quizConfig

func init() {
	flag.StringVar(&config.filename, "csv", "problems.csv", "path for a csv file in the format 'question,answer'")
	flag.BoolVar(&config.shuffle, "shuffle", false, "true/false flag to shuffle order of questions")
	flag.IntVar(&config.time, "time", 30, "time in seconds after which the quiz ends")
	flag.Parse()
}

func NewQuiz(input io.Reader, shouldShuffle bool, timeout int) (*quiz, error) {
	reader := csv.NewReader(input)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	problems := make([]problem, len(records))
	for idx, record := range records {
		var problem problem
		problem.question = record[0]
		problem.answer = strings.TrimSpace(record[1])

		problems[idx] = problem
	}

	q := quiz{
		problems:         problems,
		totalProblems:    len(problems),
		correctSolutions: 0,
		shouldShuffle:    shouldShuffle,
		timeout:          timeout,
	}

	return &q, nil
}

// TODO: Add a debugger and check how the loop and deletion
// works fine together
func (q *quiz) shuffle() {
	reorderedProblems := make([]problem, q.totalProblems)

	for i := 0; i < q.totalProblems; i++ {
		idx := rand.Intn(len(q.problems))
		reorderedProblems[i] = q.problems[idx]
		q.problems = append(q.problems[:idx], q.problems[idx+1:]...)
	}

	q.problems = reorderedProblems
}

func (q *quiz) startTimer() {
	defer wg.Done()

	time.Sleep(time.Duration(q.timeout) * time.Second)
}

func (q *quiz) Run() {
	defer wg.Done()

	if q.shouldShuffle {
		q.shuffle()
	}

	go q.startTimer()

	for i := 0; i < q.totalProblems; i++ {
		var answer string

		fmt.Printf("Problem #%d: %s\n", i+1, q.problems[i].question)
		fmt.Scanf("%s\n", &answer)

		if answer == q.problems[i].answer {
			q.correctSolutions += 1
		}
	}
	wg.Done()
}

func (q *quiz) PrintResults() {
	fmt.Printf("You score %d out of %d\n", q.correctSolutions, q.totalProblems)
}

func main() {
	quizFile, err := os.Open(config.filename)
	if err != nil {
		log.Fatalf("Failed to open the csv file: %s\nError: %s\n", config.filename, err)
	}
	defer quizFile.Close()

	quiz, err := NewQuiz(quizFile, config.shuffle, config.time)
	if err != nil {
		log.Fatalf("Failed to parse file.\nError: %s\n", err)
	}

	fmt.Println("Press enter to start the quiz")
	fmt.Scanf("%s")

	wg.Add(1)
	go quiz.Run()

	wg.Wait()
	quiz.PrintResults()
}
