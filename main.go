package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/eiannone/keyboard"
)

const maxX int = 25
const maxY int = 15

var board [maxX][maxY]Field
var direction int = 0 //0 - lewo, 2 prawo, 1 góra, 3 dół
var isAlive = true

var points int = 0
var record int

type Field struct {
	char string
}

type Position struct {
	x, y int
}

type Snake struct {
	length, pivot int
	parts         []Position
}

func clearTerminal() {
	for i := 0; i < 50; i++ {
		for j := 0; j < 120; j++ {
			fmt.Printf("\x1b[%d;%df ", i+1, j+1)
		}
	}
}

func drawBoard() {
	readRecordFromFile()
	for i := 0; i < maxX; i++ {
		for j := 0; j < maxY; j++ {
			if i == 0 || j == 0 || i == maxX-1 || j == maxY-1 {
				board[i][j].char = "#"
				putxy(i, j, "#")
			} else {
				board[i][j].char = " "
				putxy(i, j, " ")
			}
		}
	}
	fmt.Printf("\x1b[%d;%df%s", maxY, maxX+25, "Points: 0")
	fmt.Printf("\x1b[%d;%df%s%d", maxY-1, maxX+25, "Record: ", record)
}

func drawMenu() {
	putxy(21, 0, "SNAKE")
	putxy(20, 4, "New Game")
	putxy(20, 7, "Options")
	putxy(21, 10, "Exit")
}

func putxy(x, y int, s string) {
	fmt.Printf("\x1b[%d;%df%s", y+1, x+21, s)
	if x > 0 && x < maxX && y > 0 && y < maxY {
		board[x][y].char = s
	}

}
func putPosition(p Position, s string) {
	x := p.x
	y := p.y
	fmt.Printf("\x1b[%d;%df%s", y+1, x+21, s)
	if x > 0 && x < maxX && y > 0 && y < maxY {
		board[x][y].char = s
	}
}

func putSnake(s Snake) { //przełożenie snake'a na board
	for i := s.length - 1; i >= 0; i-- {
		putPosition(s.parts[i], "S")
	}
	if direction == 0 {
		putPosition(s.parts[s.pivot], "<")

	} else if direction == 2 {
		putPosition(s.parts[s.pivot], ">")

	} else if direction == 1 {
		putPosition(s.parts[s.pivot], "^")

	} else if direction == 3 {
		putPosition(s.parts[s.pivot], "V")
	}
	putxy(maxX+10, maxY+10, " ")
}

func clearSnake(s Snake) {
	for i := 0; i < s.length; i++ {
		putPosition(s.parts[i], " ")
	}
}

func setSnakeBegin() Snake { //stworzenie snake'a początkowego o długości 3
	var s Snake
	s.length = 3
	s.pivot = 0
	var p1 Position
	p1.x = maxX/2 - 1
	p1.y = maxY / 2
	var p2 Position
	p2.x = maxX / 2
	p2.y = maxY / 2
	var p3 Position
	p3.x = maxX/2 + 1
	p3.y = maxY / 2
	s.parts = make([]Position, 3)
	s.parts[0] = p1
	s.parts[1] = p2
	s.parts[2] = p3
	return s
}

func goInDirection(s *Snake) {
	var tempPositon Position = s.parts[s.pivot]
	var t Position = s.parts[s.pivot]
	var char string
	s.pivot--
	if s.pivot < 0 {
		s.pivot = s.length - 1
	}

	if direction == 0 { //x--
		tempPositon.x--
		char = board[tempPositon.x][tempPositon.y].char
		clearSnake(*s)
		s.parts[s.pivot] = tempPositon

	} else if direction == 1 {
		tempPositon.y--
		char = board[tempPositon.x][tempPositon.y].char
		clearSnake(*s)
		s.parts[s.pivot] = tempPositon

	} else if direction == 2 {
		tempPositon.x++
		char = board[tempPositon.x][tempPositon.y].char
		clearSnake(*s)
		s.parts[s.pivot] = tempPositon

	} else if direction == 3 {
		tempPositon.y++
		char = board[tempPositon.x][tempPositon.y].char
		clearSnake(*s)
		s.parts[s.pivot] = tempPositon
	}
	if char == " " {
		putSnake(*s)

	} else if char == "X" { //przypadek gdy trafimy na jedzenie
		s.length++
		s.parts = append(s.parts, t)
		putSnake(*s)
		generateFood()
		givePoints(10)

	} else if char == "#" { //trafienie na sciane
		putSnake(*s)
		isAlive = false

	} else if char == "S" { //trafienie na swój ogon
		putSnake(*s)
		isAlive = false
	}

}

func chooseDirection() { //sterowanie
	for {
		_, key, _ := keyboard.GetKey()
		if key == keyboard.KeyArrowLeft && direction != 2 {
			direction = 0
		} else if key == keyboard.KeyArrowRight && direction != 0 {
			direction = 2
		} else if key == keyboard.KeyArrowUp && direction != 3 {
			direction = 1
		} else if key == keyboard.KeyArrowDown && direction != 1 {
			direction = 3
		}
	}
}

func clearChooseMenu() {
	putxy(18, 4, "  ")
	putxy(28, 4, "  ")
	putxy(18, 7, "  ")
	putxy(27, 7, "  ")
	putxy(19, 10, "  ")
	putxy(25, 10, "  ")

}

func chooseMenu() {

	for {
		clearChooseMenu()
		if direction == 0 || direction > 2 {
			direction = 0
			putxy(18, 4, "[_")
			putxy(28, 4, "_]")
		} else if direction == 1 {
			putxy(18, 7, "[_")
			putxy(27, 7, "_]")
		} else if direction == 2 || direction < 0 {
			direction = 2
			putxy(19, 10, "[_")
			putxy(25, 10, "_]")
		}
		putxy(maxX+10, maxY+10, " ")
		_, key, _ := keyboard.GetKey()
		if key == keyboard.KeyArrowUp {
			direction--
		} else if key == keyboard.KeyArrowDown {
			direction++
		} else if key == keyboard.KeyEnter {
			if direction == 0 {
				break

			} else if direction == 1 {
				//options

			} else if direction == 2 {
				//exit
			}

		}

	}
}

func generateFood() {
	x := rand.Intn(maxX-2) + 2
	y := rand.Intn(maxY-2) + 2
	if board[x][y].char == " " {
		putxy(x, y, "X")
	} else {
		generateFood()
	}
}

func givePoints(p int) {
	points += p
	fmt.Printf("\x1b[%d;%df%d", maxY, maxX+33, points)
	if points > record {
		record = points
		fmt.Printf("\x1b[%d;%df%d", maxY-1, maxX+33, record)
	}
}

func readRecordFromFile() {
	content, _ := ioutil.ReadFile("data.dat")
	record, _ = strconv.Atoi(string(content))

}

func setNewRecord() {
	F, err := os.Create("data.dat")
	if err != nil {
		fmt.Println(err)
	}
	defer F.Close()
	s := strconv.Itoa(record)
	F.WriteString(s)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var speed uint = 200
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	var snake Snake
	snake = setSnakeBegin()

	clearTerminal() //oczyszczenie terminala i stworzenie srodowiska do gry
	drawMenu()
	chooseMenu()
	clearTerminal()
	drawBoard()

	putSnake(snake) //przygotowanie rozgrywki
	generateFood()
	var tempRecord int = record

	go chooseDirection() //poboczny keyboard listener
	for isAlive {        //dopóki wąż jest żywy to program działa
		time.Sleep(time.Millisecond * time.Duration(speed))
		goInDirection(&snake)
	}
	//skonczona gra
	if record > tempRecord { //ustanowiony nowy rekord
		setNewRecord()
		putxy(maxX/2-7, maxY/2, " NOWY REKORD! ")
	} else {
		putxy(maxX/2-6, maxY/2, " KONIEC GRY ")
	}
	s := strconv.Itoa(points)
	putxy(maxX/2-8, maxY/2+1, " Zdobyte pkt: "+s+" ")

}
