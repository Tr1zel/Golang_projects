package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

// waitgroup
func main() {
	fmt.Println("Привет! Ты попал на викторину. Я буду отправлять тебе примеры, ты должен дать правильный ответ")
	fmt.Println("У тебя есть 10 секунд для ответа на все вопросы")
	records := read_from_csv("problems.csv")

	timerCh := time.After(30 * time.Second) // канал для таймера
	resultCh := make(chan string)           // канал для флага что вопросы закончились

	go func() {
		result := verify_otvet(records, timerCh)
		resultCh <- result
	}()

	select {
	case <-timerCh:
		fmt.Printf("\nВремя вышло! ")

	case result := <-resultCh:
		fmt.Println(result)

		fmt.Printf("Все вопросы пройдены!%s\n", result)
	}

}

func read_from_csv(filename string) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	return records
}

func verify_otvet(records [][]string, timerCh <-chan time.Time) string {
	succes_otvet := 0
	failed_otvet := 0
	otvet := 0
	for i, record := range records {
		fmt.Printf("Вопрос номер %d: %s = ", i+1, record[0])
		fmt.Scanf("%d", &otvet)
		res, err := strconv.Atoi(record[1])
		if err != nil {
			panic(err)
		}
		if otvet == res {
			succes_otvet++
			fmt.Println("Верно! Следующий вопрос")
		} else {
			failed_otvet++
			fmt.Println("Неверно! Следующий вопрос")
		}

	}
	return fmt.Sprintf("\nРезультат: %d/%d правильных ответов\n", succes_otvet, len(records))

}
