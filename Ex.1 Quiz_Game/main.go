package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Привет! Ты попал на викторину. Я буду отправлять тебе примеры, ты должен дать правильный ответ")
	records := read_from_csv("problems.csv")
	verify_otvet(records)
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

func verify_otvet(records [][]string) {
	var otvet int64
	var succes_otvet int64
	var failed_otvet int64
	for i, record := range records {
		fmt.Printf("Вопрос номер %d\n", i+1)
		fmt.Printf("%s = ", record[0])
		fmt.Scanf("%d", &otvet)
		num, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			fmt.Println("ошибка: ", err)
			panic(err)
		}
		if otvet == num {
			println("Верно! Следующий вопрос")
			succes_otvet++
		} else {
			println("Неправильно! Следующий вопрос")
			failed_otvet++
		}
	}
	fmt.Printf("Количество верных ответов: %d\nКоличество неверных ответов: %d\nВсего вопросов: %d\n", succes_otvet, failed_otvet, len(records))

}
