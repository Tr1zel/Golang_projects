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

const (
	roomKitchen  = "кухня"
	roomCorridor = "коридор"
	roomRoom     = "комната"
	roomStreet   = "улица"
	roomHome     = "домой"
)

type Player struct {
	curRoom   *Room
	inventory []string
	bacpack   bool
}

type Room struct {
	name          string
	description   string
	objectsInRoom []string
	availibleMove map[string]string
	moveOrder     []string // порядок направлений
	isLocked      bool
}

type World struct {
	allRooms map[string]*Room
	Player   *Player
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
	gameWorld.allRooms = make(map[string]*Room)

	kitchen := &Room{
		name:          "Кухня",
		description:   "ты находишься на кухне,",
		objectsInRoom: []string{"на столе: чай, надо собрать рюкзак и идти в универ."},
		availibleMove: make(map[string]string),
		isLocked:      false,
	}
	coridor := &Room{
		name:          "Коридор",
		description:   "ничего интересного.",
		objectsInRoom: []string{},
		availibleMove: make(map[string]string),
		isLocked:      false,
	}
	comnata := &Room{
		name:          "Комната",
		description:   "ты в своей комнате.",
		objectsInRoom: []string{"на столе: ключи, конспекты", "на стуле: рюкзак."},
		availibleMove: make(map[string]string),
		isLocked:      false,
	}
	street := &Room{
		name:          "Улица",
		description:   "на улице весна.",
		objectsInRoom: []string{},
		availibleMove: make(map[string]string),
		isLocked:      true,
	}

	kitchen.availibleMove[roomCorridor] = roomCorridor
	coridor.availibleMove[roomStreet] = roomStreet
	coridor.availibleMove[roomKitchen] = roomKitchen
	coridor.availibleMove[roomRoom] = roomRoom
	coridor.moveOrder = []string{roomKitchen, roomRoom, roomStreet}
	comnata.availibleMove[roomCorridor] = roomCorridor
	street.availibleMove[roomHome] = roomCorridor

	gameWorld.allRooms[roomKitchen] = kitchen
	gameWorld.allRooms[roomCorridor] = coridor
	gameWorld.allRooms[roomRoom] = comnata
	gameWorld.allRooms[roomStreet] = street

	gamePlayer = Player{
		curRoom:   gameWorld.allRooms[roomKitchen],
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

	splitCommand := strings.Split(command, " ")
	switch splitCommand[0] {
	case "осмотреться":
		return ViewAround()
	case "идти":
		if value, exists := gamePlayer.curRoom.availibleMove[splitCommand[1]]; exists {
			if value == roomStreet && gameWorld.allRooms[value].isLocked {
				return "дверь закрыта"
			}
			gamePlayer.curRoom = gameWorld.allRooms[value]
			return GoTo()
		}
		return fmt.Sprintf("нет пути в %s", splitCommand[1])
	case "надеть":
		if !gamePlayer.bacpack {
			gamePlayer.bacpack = true
			gamePlayer.curRoom.objectsInRoom = slices.DeleteFunc(
				gamePlayer.curRoom.objectsInRoom,
				func(item string) bool {
					return strings.Contains(item, splitCommand[1])
				},
			)
			return fmt.Sprintf("вы надели: %s", splitCommand[1])
		}
		return "уже надето"
	case "взять":
		if !gamePlayer.bacpack {
			return "некуда класть"
		}
		found := false
		for _, item := range gamePlayer.curRoom.objectsInRoom {
			if strings.Contains(item, splitCommand[1]) {
				found = true
				break
			}
		}
		if found {
			gamePlayer.inventory = append(gamePlayer.inventory, splitCommand[1])
			// Удаляем предмет из строки или всю строку если она станет пустой
			newObjects := []string{}
			for _, item := range gamePlayer.curRoom.objectsInRoom {
				if strings.Contains(item, splitCommand[1]) {
					// Удаляем предмет из строки
					updated := strings.ReplaceAll(item, splitCommand[1]+", ", "")
					updated = strings.ReplaceAll(updated, ", "+splitCommand[1], "")
					updated = strings.ReplaceAll(updated, splitCommand[1], "")
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
			gamePlayer.curRoom.objectsInRoom = newObjects
			return fmt.Sprintf("предмет добавлен в инвентарь: %s", splitCommand[1])
		}
		return "нет такого"
	case "применить":
		return apply(splitCommand[1], splitCommand[2])
	default:
		return "неизвестная команда"
	}
}
func ViewAround() string { // функция осмотреться
	var directions []string
	if len(gamePlayer.curRoom.moveOrder) > 0 {
		directions = gamePlayer.curRoom.moveOrder
	} else {
		directions = []string{}
		for dir := range gamePlayer.curRoom.availibleMove {
			directions = append(directions, dir)
		}
		sort.Strings(directions)
	}
	directionsStr := strings.Join(directions, ", ")

	objectsStr := strings.Join(gamePlayer.curRoom.objectsInRoom, ", ")

	// Для кухни меняем текст если рюкзак уже надет
	switch {
	case gamePlayer.curRoom.name == "Кухня":
		if gamePlayer.bacpack {
			objectsStr = strings.ReplaceAll(objectsStr, "надо собрать рюкзак и идти в универ.", "надо идти в универ.")
		}
		return fmt.Sprintf("%s %s можно пройти - %s", gamePlayer.curRoom.description, objectsStr, directionsStr)
	case len(gamePlayer.curRoom.objectsInRoom) == 0:
		return fmt.Sprintf("пустая комната. можно пройти - %s", directionsStr)
	default:
		// Добавляем точку если её нет
		if !strings.HasSuffix(objectsStr, ".") {
			objectsStr += "."
		}
		return fmt.Sprintf("%s можно пройти - %s", objectsStr, directionsStr)
	}
}

func GoTo() string { // функция идти
	var directions []string
	if len(gamePlayer.curRoom.moveOrder) > 0 {
		directions = gamePlayer.curRoom.moveOrder
	} else {
		directions = []string{}
		for dir := range gamePlayer.curRoom.availibleMove {
			directions = append(directions, dir)
		}
		sort.Strings(directions)
	}
	directionsStr := strings.Join(directions, ", ")

	if gamePlayer.curRoom.name == "Кухня" {
		return fmt.Sprintf("кухня, ничего интересного. можно пройти - %s", directionsStr)
	}
	return fmt.Sprintf("%s можно пройти - %s", gamePlayer.curRoom.description, directionsStr)
}

func apply(object string, where string) string { // функция использовать
	if !slices.Contains(gamePlayer.inventory, object) {
		return fmt.Sprintf("нет предмета в инвентаре - %s", object)
	}
	if where == "дверь" {
		gameWorld.allRooms[roomStreet].isLocked = false
		return "дверь открыта"
	}
	return "не к чему применить"

}
