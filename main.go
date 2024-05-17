package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Data struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	//Пример запуска скрипта: go run main.go tz.json 4, где "4" это кол-во гоурутин
	if len(os.Args) < 3 {
		fmt.Println("error with usage")
		return
	}

	filename := os.Args[1]
	numGoroutines, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("error with number of goroutines:", err)
		return
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("error with file:", err)
		return
	}

	var arrays []Data
	err = json.Unmarshal(data, &arrays)
	if err != nil {
		fmt.Println("error with unmarshalling:", err)
		return
	}

	resultChan := make(chan int, numGoroutines)
	var wg sync.WaitGroup

	size := (len(arrays) + numGoroutines - 1) / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		start := i * size
		end := start + size
		if end > len(arrays) {
			end = len(arrays)
		}

		if start < end {
			wg.Add(1)
			go sum(arrays[start:end], &wg, resultChan)
		}
	}

	wg.Wait()
	close(resultChan)

	totalSum := 0
	for s := range resultChan {
		totalSum += s
	}

	fmt.Println("Total Sum:", totalSum)
}

func sum(data []Data, wg *sync.WaitGroup, resultChan chan int) {
	defer wg.Done()
	total := 0
	for _, d := range data {
		total += d.A + d.B
	}
	resultChan <- total
}
