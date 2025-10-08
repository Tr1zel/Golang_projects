package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
)

/*
код писать в этом файле
наверняка у вас будут какие-то структуры с методами, глобальные переменные ( тут можно ), функции
*/
var rooms = []string{"Комната", "Улица", "Кухня", "Коридор", "Домой"}

type Player struct {
	cur_room  *Room
	inventory []string
	bacpack   bool
}

type Room struct {
	name            string
	description     string
	objects_in_room []string
	availible_move  map[string]string
	move_order      []string // порядок направлений
	is_locked       bool
}

type World struct {
	all_rooms map[string]*Room
	Player    *Player
}

var gameWorld World
var gamePlayer Player

func main() {

	initGame()
	var command string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		command = scanner.Text()
		res := handleCommand(command)
		if res == "Неизвестная команда" {
			fmt.Println("Неизвестная команда!")
		}
		fmt.Println(res)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка чтения = %s", err)
	}
	/*
		в этой функции можно ничего не писать,
		но тогда у вас не будет работать через go run main.go
		очень круто будет сделать построчный ввод команд тут, хотя это и не требуется по заданию
	*/
}

func initGame() {
	gameWorld.all_rooms = make(map[string]*Room)

	kitchen := &Room{
		name:            "Кухня",
		description:     "ты находишься на кухне,",
		objects_in_room: []string{"на столе: чай, надо собрать рюкзак и идти в универ."},
		availible_move:  make(map[string]string),
		is_locked:       false,
	}
	coridor := &Room{
		name:            "Коридор",
		description:     "ничего интересного.",
		objects_in_room: []string{},
		availible_move:  make(map[string]string),
		is_locked:       false,
	}
	comnata := &Room{
		name:            "Комната",
		description:     "ты в своей комнате.",
		objects_in_room: []string{"на столе: ключи, конспекты", "на стуле: рюкзак."},
		availible_move:  make(map[string]string),
		is_locked:       false,
	}
	street := &Room{
		name:            "Улица",
		description:     "на улице весна.",
		objects_in_room: []string{},
		availible_move:  make(map[string]string),
		is_locked:       true,
	}

	kitchen.availible_move["коридор"] = "коридор"
	coridor.availible_move["улица"] = "улица"
	coridor.availible_move["кухня"] = "кухня"
	coridor.availible_move["комната"] = "комната"
	coridor.move_order = []string{"кухня", "комната", "улица"}
	comnata.availible_move["коридор"] = "коридор"
	street.availible_move["домой"] = "коридор"

	gameWorld.all_rooms["кухня"] = kitchen
	gameWorld.all_rooms["коридор"] = coridor
	gameWorld.all_rooms["комната"] = comnata
	gameWorld.all_rooms["улица"] = street

	gamePlayer = Player{
		cur_room:  gameWorld.all_rooms["кухня"],
		inventory: []string{},
		bacpack:   false,
	}
	/*
		эта функция инициализирует игровой мир - все комнаты
		если что-то было - оно корректно перезатирается
	*/
}

func handleCommand(command string) string {
	/*
		данная функция принимает команду от "пользователя"
		и наверняка вызывает какой-то другой метод или функцию у "мира" - списка комнат
	*/

	split_command := strings.Split(command, " ")
	if split_command[0] == "осмотреться" {
		return view_around()
	} else if split_command[0] == "идти" {
		if value, exists := gamePlayer.cur_room.availible_move[split_command[1]]; exists {
			if value == "улица" && gameWorld.all_rooms[value].is_locked {
				return fmt.Sprintf("дверь закрыта")
			} else {
				gamePlayer.cur_room = gameWorld.all_rooms[value]
				return go_to()
			}
		} else {
			return fmt.Sprintf("нет пути в %s", split_command[1])
		}
	} else if split_command[0] == "надеть" {
		if !gamePlayer.bacpack {
			gamePlayer.bacpack = true
			gamePlayer.cur_room.objects_in_room = slices.DeleteFunc(
				gamePlayer.cur_room.objects_in_room,
				func(item string) bool {
					return strings.Contains(item, split_command[1])
				},
			)
			return fmt.Sprintf("вы надели: %s", split_command[1])
		}
		return "уже надето"
	} else if split_command[0] == "взять" {
		if !gamePlayer.bacpack {
			return "некуда класть"
		}
		found := false
		for _, item := range gamePlayer.cur_room.objects_in_room {
			if strings.Contains(item, split_command[1]) {
				found = true
				break
			}
		}
		if found {
			gamePlayer.inventory = append(gamePlayer.inventory, split_command[1])
			// Удаляем предмет из строки или всю строку если она станет пустой
			newObjects := []string{}
			for _, item := range gamePlayer.cur_room.objects_in_room {
				if strings.Contains(item, split_command[1]) {
					// Удаляем предмет из строки
					updated := strings.ReplaceAll(item, split_command[1]+", ", "")
					updated = strings.ReplaceAll(updated, ", "+split_command[1], "")
					updated = strings.ReplaceAll(updated, split_command[1], "")
					updated = strings.TrimSpace(updated)
					// Проверяем осталось ли что-то кроме "на столе:" или "на стуле:"
					if strings.HasSuffix(updated, ":") || updated == "" || strings.HasSuffix(updated, ": ") {
						// Если осталось только "на столе:" - не добавляем
						continue
					}
					newObjects = append(newObjects, updated)
				} else {
					newObjects = append(newObjects, item)
				}
			}
			gamePlayer.cur_room.objects_in_room = newObjects
			return fmt.Sprintf("предмет добавлен в инвентарь: %s", split_command[1])
		}
		return "нет такого"
	} else if split_command[0] == "применить" {
		return apply(split_command[1], split_command[2])
	} else {
		return "неизвестная команда"
	}
}
func view_around() string { // функция осмотреться
	var directions []string
	if len(gamePlayer.cur_room.move_order) > 0 {
		directions = gamePlayer.cur_room.move_order
	} else {
		directions = []string{}
		for dir := range gamePlayer.cur_room.availible_move {
			directions = append(directions, dir)
		}
		sort.Strings(directions)
	}
	directionsStr := strings.Join(directions, ", ")

	objectsStr := strings.Join(gamePlayer.cur_room.objects_in_room, ", ")

	// Для кухни меняем текст если рюкзак уже надет
	if gamePlayer.cur_room.name == "Кухня" {
		if gamePlayer.bacpack {
			objectsStr = strings.ReplaceAll(objectsStr, "надо собрать рюкзак и идти в универ.", "надо идти в универ.")
		}
		return fmt.Sprintf("%s %s можно пройти - %s", gamePlayer.cur_room.description, objectsStr, directionsStr)
	} else if len(gamePlayer.cur_room.objects_in_room) == 0 {
		return fmt.Sprintf("пустая комната. можно пройти - %s", directionsStr)
	} else {
		// Добавляем точку если её нет
		if !strings.HasSuffix(objectsStr, ".") {
			objectsStr += "."
		}
		return fmt.Sprintf("%s можно пройти - %s", objectsStr, directionsStr)
	}
}

func go_to() string { // функция идти
	var directions []string
	if len(gamePlayer.cur_room.move_order) > 0 {
		directions = gamePlayer.cur_room.move_order
	} else {
		directions = []string{}
		for dir := range gamePlayer.cur_room.availible_move {
			directions = append(directions, dir)
		}
		sort.Strings(directions)
	}
	directionsStr := strings.Join(directions, ", ")

	if gamePlayer.cur_room.name == "Кухня" {
		return fmt.Sprintf("кухня, ничего интересного. можно пройти - %s", directionsStr)
	}
	return fmt.Sprintf("%s можно пройти - %s", gamePlayer.cur_room.description, directionsStr)
}

func apply(object string, where string) string { // функция использовать
	if !slices.Contains(gamePlayer.inventory, object) {
		return fmt.Sprintf("нет предмета в инвентаре - %s", object)
	}
	if where == "дверь" {
		gameWorld.all_rooms["улица"].is_locked = false
		return "дверь открыта"
	} else {
		return "не к чему применить"
	}

}
