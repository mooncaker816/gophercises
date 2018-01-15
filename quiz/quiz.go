package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "格式：问题，答案")
	flag.Parse()
	//_ = csvFilename
	file, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Printf("can not open %s\n", *csvFilename)
	}
	defer file.Close()
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	problems := parseLines(lines)
	rand.Seed(time.Now().UnixNano())
	Shuffle(problems)
	correct := 0
	ansCh := make(chan string)
	go func() {
		for {
			var answer string
			fmt.Scanf("%s\n", &answer)
			ansCh <- answer
		}
	}()
	for i, v := range problems {
		fmt.Printf("Problem # %d: %s = ", i+1, v.q)
		subTimer := time.NewTimer(5 * time.Second)
		select {
		case <-subTimer.C:
			fmt.Printf("\ntime for this problem is exhausted!\n")
		case answer := <-ansCh:
			if strings.ToLower(strings.TrimSpace(answer)) == v.a {
				correct++
			}
		}
	}
	fmt.Printf("Quiz completed! You've got %d/%d!\n", correct, len(problems))
}

type problem struct {
	q string
	a string
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, v := range lines {
		ret[i].q = v[0]
		ret[i].a = strings.ToLower(strings.TrimSpace(v[1]))
	}
	return ret
}

// func Shuffle(array []interface{}, source rand.Source) {
// 	random := rand.New(source)
// 	for i := len(array) - 1; i > 0; i-- {
// 		j := random.Intn(i + 1)
// 		array[i], array[j] = array[j], array[i]
// 	}
// }

// [Min] common method to shuffle any slice
func Shuffle(slice interface{}) {
	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	length := rv.Len()
	for i := length - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		swap(i, j)
	}
}
