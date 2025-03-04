package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	quiz    = flag.String("quiz", "problems.csv", "Questions and answers file")
	limit   = flag.Int("limit", 30, "Time given to answer each question")
	shuffle = flag.Bool("shuffle", false, "shuffle the quiz questions")
)

func loadFile(string) ([][]string, error) {
	file, err := os.Open(*quiz)
	if err != nil {
		return nil, fmt.Errorf("failed to open the file: %w", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	questions, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read the file: %w", err)
	}
	return questions, nil
}

func main() {
	flag.Parse()
	qs, err := loadFile(*quiz)
	if err != nil {
		log.Println(err)
		return
	}

	if *shuffle {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(qs), func(i int, j int) {
			qs[i], qs[j] = qs[j], qs[i]
		})
	}

	fmt.Println("Press Enter to start the quiz: ")
	fmt.Scanln()

	s := bufio.NewScanner(os.Stdin)
	timer := time.NewTimer(time.Duration(*limit) * time.Second)
	done := make(chan bool)
	correct := 0

	go func() {
		for _, q := range qs {
			fmt.Printf("How much is %s? ", q[0])
			s.Scan()
			i := s.Text()
			i = strings.Trim(i, " /,.!?-=")
			i = strings.ToLower(i)
			if i == q[1] {
				correct++
			}
		}
		done <- true
	}()
	select {
	case <-timer.C:
		fmt.Println("\nTime's up!")
	case <-done:
	}
	fmt.Printf("Your total score is %d/%d\n", correct, len(qs))
}
