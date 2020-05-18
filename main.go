package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var csvFile = flag.String("csv", "problems.csv", "File containing problems file in CSV file.")
var limit = flag.Int("limit", 30, "the time limit for quiz in seconds")

type problem struct {
	question string
	answer string
}

func convert(records [][]string) []problem {
	var result []problem
	for _, r := range records {
		if len(r) > 1 {
			result = append(result, problem{r[0], r[1]})
		}
	}
	return result
}

func main() {

	flag.Parse()
	file, err := os.Open(*csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	problems := convert(records)
  fmt.Println("Press enter to start")
  fmt.Scanf("\n")

	ansCh := make(chan string)
	timer := time.NewTimer(time.Duration(*limit) * time.Second)
	elapsed := false
	correct := 0
	for i, p := range problems {
		if elapsed {
			break
		}
		fmt.Printf("Problem %02d: %s\n", i+1, p.question)
		go func() {
			var ans string
			fmt.Scanf("%s", &ans)
			ansCh <- ans
		}()

		select {
		case <- timer.C:
			elapsed = true
		case ans := <- ansCh:
			if ans == p.answer {
				correct++
			}
		}
	}
	fmt.Printf("You got %d of %d correct\n", correct, len(problems))
}
