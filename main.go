package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var csvFile = flag.String("csv", "problems.csv", "File containing problems file in CSV file.")
var limit = flag.Int("limit", 30, "the time limit for quiz in seconds")

func askQuestion(counter int, rec []string, ch chan bool, waitTime time.Duration) int{
	fmt.Printf("Problem %02d: %s\n", counter, rec[0])
	var result string
  go func() {
        fmt.Scanf("%s", &result)
        if rec[1] == result {
      		ch <- true
      		return
      	}
      	ch <- false
    }()
    select {
    case isCorrect := <- ch:
      if isCorrect {
        return 1
      }
      return 0
    case <-time.After(waitTime):
      return 0
    }
}

func main() {

	flag.Parse()

	file, err := os.Open(*csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rec, err := reader.Read()

  fmt.Println("Press enter to start")
  fmt.Scanf("\n")
	ch := make(chan bool)
	total, correct := 0, 0
	start := time.Now()
	for err == nil {
		if len(rec) != 2 {
			log.Print("WARNING: Skipping invalid entry %v", rec)
			continue
		}

		total++
    timeWait := time.Duration(*limit) *time.Second - time.Now().Sub(start)
    if timeWait >0 {
      correct = correct + askQuestion(total, rec, ch, timeWait)
    }
    rec, err = reader.Read()
	}
	if err != io.EOF {
		log.Fatal(err)
	}
	if len(rec) > 0 {
		total++
    timeWait := time.Duration(*limit) *time.Second - time.Now().Sub(start)
    if timeWait >0 {
      correct = correct + askQuestion(total, rec, ch, timeWait)
    }
	}

	fmt.Printf("You got %d of %d correct\n", correct, total)
}
