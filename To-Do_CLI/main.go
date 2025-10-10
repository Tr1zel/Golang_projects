package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	menu = `В данный момент доступны команды:
	add - добавить задачу;
	list - посмотреть все задачи
	exit - выйти из приложения 
`
)

type Task struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Done        bool
}

func main() {
	if len(os.Args) < 2 {
		panic("Недостаточно аргументов, введите имя файла в аргументах")
	}
	filename := os.Args[1]
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("Файла не сущеcтвует %s\n", filename)
		_, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Привет! Твой выбранный файл: ", filename)
	fmt.Printf(menu)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		command := scanner.Text()
		split_command := strings.Split(command, " ")
		switch split_command[0] {
		case "exit":
			fmt.Println("Спасибо за использование нашей программы")
			return
		case "add":
			addTask(filename)
		case "list":
			readFile(filename)
		case "done":
			doneTask(filename)
		case "delete":
			deleteTask(filename)
		default:
			fmt.Println("Неизвестная команда")
			fmt.Printf(menu)
		}
	}
}

func addTask(filename string) {
	_, err := os.OpenFile(filename, 0, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	var tempTask Task
	scanner := bufio.NewScanner(os.Stdin)

	// Читаем существующие задачи чтобы узнать максимальный ID
	data, _ := os.ReadFile(filename)
	var existingTasks []Task
	if len(data) > 0 {
		json.Unmarshal(data, &existingTasks)
	}
	// Находим максимальный ID
	maxID := 0
	for _, task := range existingTasks {
		if task.Id > maxID {
			maxID = task.Id
		}
	}
	tempTask.Id = maxID + 1

	fmt.Println("Введите имя задачи:")
	scanner.Scan()
	tempTask.Name = scanner.Text()
	fmt.Println("Введите описание задачи в одну строку:")
	scanner.Scan()
	tempTask.Description = scanner.Text()
	fmt.Println("Введите дедлайн по задаче в формате DD/MM/YY:")
	scanner.Scan()
	tempTask.Deadline = parseDate(scanner.Text())
	tempTask.Done = false
	err = appendToJson(filename, tempTask)
	if err != nil {
		panic(err)
	}
	fmt.Println("Задача добавлена в файл!")
}

func appendToJson(filename string, newitem Task) error {
	var items []Task

	data, err := os.ReadFile(filename)
	if err == nil && len(data) > 0 {
		json.Unmarshal(data, &items)
	} else if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	items = append(items, newitem)

	newData, err := json.MarshalIndent(items, "", "	")
	if err != nil {
		panic(err)
	}

	return os.WriteFile(filename, newData, 0644)

}

func parseDate(dateStr string) time.Time {
	var ResTime time.Time
	dateSplit := strings.Split(dateStr, "/")
	day, _ := strconv.Atoi(dateSplit[0])
	month, _ := strconv.Atoi(dateSplit[1])
	year, _ := strconv.Atoi(dateSplit[2])

	ResTime = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	fmt.Println(ResTime)
	return ResTime
}

func readFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var tasks []Task
	if len(data) > 0 {
		err = json.Unmarshal(data, &tasks)
		if err != nil {
			panic(err)
		}
	}

	if len(tasks) == 0 {
		fmt.Println("В данный момент у вас нет задач")
	} else {
		for i, task := range tasks {
			fmt.Printf("Задача %d:\n", i+1)
			fmt.Printf("  ID: %d\n", task.Id)
			fmt.Printf("  Название: %s\n", task.Name)
			fmt.Printf("  Описание: %s\n", task.Description)
			fmt.Printf("  Дедлайн: %s\n", task.Deadline.Format("02/01/2006"))
			status := "НЕТ"
			if task.Done {
				status = "ДА"
			}
			fmt.Printf("  Выполнена: %s\n\n", status)
		}
	}
}

func doneTask(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		panic(err)
	}

	fmt.Println("Введите айди задачи которую выполнили:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	taskId := scanner.Text()

	id, _ := strconv.Atoi(taskId)
	for i := range tasks {
		if tasks[i].Id == id {
			tasks[i].Done = true
			break
		}
	}
	newData, _ := json.MarshalIndent(tasks, "", "\t")
	os.WriteFile(filename, newData, 0644)
	fmt.Printf("Задача с айди %d отмечена как выполненная!\n", id)
}

func deleteTask(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		panic(err)
	}

	fmt.Println("Введите айди задачи которую хотите удалить:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	taskId := scanner.Text()

	tasks = slices.DeleteFunc(tasks,
		func(item Task) bool {
			id, _ := strconv.Atoi(taskId)
			return item.Id == id
		})
	newData, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		panic(err)
	}
	os.WriteFile(filename, newData, 0644)
	fmt.Printf("Задача с айди %s успешно удалена!\n", taskId)
}
