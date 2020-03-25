package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

type board struct {
	data [10][10]bool //change grid dimensions here
}

type snake struct {
	position []co_ordinate
}

type co_ordinate struct {
	xCoord int
	yCoord int
}

type path struct {
	length          int
	startCoordinate co_ordinate
}

type Node struct {
	coords   co_ordinate
	futility int //the higher the futility the worse the path is
}

const (
	EMPTY_SQUARE bool = false
	FOOD_SQUARE  bool = true

	//representations
	SNAKE_REP string = "#"
	HEAD_REP  string = ""
	EMPTY_REP string = " "
	FOOD_REP  string = "$"

	//fill colors if you want
	BLANK_COLOR string = ""
	FOOD_COLOR  string = ""
	SNAKE_COLOR string = ""
	RESET       string = ""

	//important constants
	DEATH int = 9999999999
)

var grid board
var snaky snake

var rootPointer *Node

var foodLocation co_ordinate
var foodGenerated bool
var totSquares int

func (self snake) len() int {
	return len(self.position)
}

func delay(seconds string) {
	delayTime, _ := time.ParseDuration(seconds)
	time.Sleep(delayTime)
}

func main() {
	initialize()
	x, y := snaky.generateFoodLocation(grid)

	for snaky.len() != totSquares {
		if !foodGenerated {
			x, y = snaky.generateFoodLocation(grid)
			(&grid).place(x, y, FOOD_SQUARE)

			foodLocation = co_ordinate{x, y}

			foodGenerated = true
		}
		nextCoordinate, futility := snaky.getNextCoordinate(grid)
		grid = (&snaky).moveTo(nextCoordinate, foodLocation, grid, futility)

		grid.display(snaky)
		delay("100ms")
	}
}

func initialize() {
	grid = board{}
	snaky = grid.initializeSnake(4)

	totSquares = grid.totalSquaresOfTheBoard()
}

func (self board) initializeSnake(length int) snake {
	yCenter := len(self.data) / 2
	xCenter := (len(self.data[0]) - length) / 2

	position := make([]co_ordinate, 0)

	for i := 0; i < length; i++ {
		position = append(position, co_ordinate{xCenter + i, yCenter})
	}

	return snake{position}

}

func (self board) totalSquaresOfTheBoard() int {
	return len(self.data) * len(self.data[0])
}

func print(matter, modifier string) {
	fmt.Print(modifier + matter + RESET)
}

func (self board) display(player snake) {
	clearScreen()

	for i := 0; i < len(self.data); i++ {
		for j := 0; j < len(self.data[0]); j++ {
			itIs := player.checkIfPositionIsOccupied(i, j)

			if itIs {
				print(SNAKE_REP, SNAKE_COLOR)
			} else {
				if isFood(self.data[i][j]) {
					print(FOOD_REP, FOOD_COLOR)
				} else {
					print(EMPTY_REP, BLANK_COLOR)
				}
			}
		}
		fmt.Println("")
	}
}

func clearScreen() {
	fmt.Println("\033[H\033[2J")
}

func (self snake) checkIfPositionIsOccupied(xCoord, yCoord int) bool {
	for i := 0; i < len(self.position); i++ {
		pos := self.position[i]

		if (pos.xCoord == xCoord) && (pos.yCoord == yCoord) {
			return true
		}
	}
	return false
}

func isFood(args bool) bool {
	return args == FOOD_SQUARE
}

func (self snake) generateFoodLocation(presentBoard board) (xCoord, yCoord int) {
	xCoord = rand.Intn(len(presentBoard.data))
	yCoord = rand.Intn(len(presentBoard.data[0]))

	for true {
		xCoord = rand.Intn(len(presentBoard.data))
		yCoord = rand.Intn(len(presentBoard.data[0]))

		if !self.checkIfPositionIsOccupied(xCoord, yCoord) {
			return xCoord, yCoord
		}
	}
	return xCoord, yCoord
}

func (presentBoard *board) place(xCoord, yCoord int, value bool) {
	presentBoard.data[xCoord][yCoord] = value
}

func (self snake) getNextCoordinate(presentBoard board) (nextCoords co_ordinate, futility int) {
	nextCoords, futility = self.generateBestPath(presentBoard, foodLocation)

	return
}

func (self snake) generateBestPath(presentBoard board, foodLocation co_ordinate) (co_ordinate, int) {
	headCoordinates := self.position[0]
	rootNode := Node{headCoordinates, self.getFutility(headCoordinates, foodLocation, presentBoard)}
	itrPointer, _ := &rootNode, &rootNode //rootPointer

	bestChild := itrPointer.getBestChild(presentBoard, foodLocation, self)

	return bestChild.coords, bestChild.futility
}

func (self snake) getFutility(myCoord, foodLocation co_ordinate, presentBoard board) int {
	deltaX := math.Abs(float64(myCoord.xCoord - foodLocation.xCoord))
	deltaY := math.Abs(float64(myCoord.yCoord - foodLocation.yCoord))

	occupied := self.checkIfPositionIsOccupied(myCoord.xCoord, myCoord.yCoord)

	if occupied {
		return DEATH
	} else if myCoord.exceedsLimitsOf(presentBoard) {
		return DEATH
	}
	return int(deltaY + deltaX)
}

func (self co_ordinate) exceedsLimitsOf(presentBoard board) bool {
	x, y := self.xCoord, self.yCoord

	return (x > (len(presentBoard.data) - 1)) || (y > (len(presentBoard.data[0]) - 1))
}

func (self *Node) getBestChild(presentBoard board, foodLocation co_ordinate, snaky snake) Node {
	my_coords := self.coords

	children := my_coords.allChildren(presentBoard)
	bestChild := snaky.getBestNode(children, foodLocation, presentBoard)

	return bestChild
}

func (self co_ordinate) allChildren(presentBoard board) (children []co_ordinate) {
	x, y := self.xCoord, self.yCoord

	children = append(children, co_ordinate{x - 1, y})
	children = append(children, co_ordinate{x + 1, y})
	children = append(children, co_ordinate{x, y - 1})
	children = append(children, co_ordinate{x, y + 1})

	return
}

func (self snake) getBestNode(children []co_ordinate, foodLocation co_ordinate, presentBoard board) (bestNode Node) {
	bestNode.futility = DEATH

	for i := 0; i < len(children); i++ {
		presentNode := Node{children[i], self.getFutility(children[i], foodLocation, presentBoard)}

		if bestNode.futility > presentNode.futility {
			bestNode = presentNode
		}
	}
	return
}

func (self *snake) moveTo(nextCoordinate, foodLocation co_ordinate, presentBoard board, futility int) board {
	if nextCoordinate == foodLocation {
		self.position = append(self.position, co_ordinate{})
		foodGenerated = false
		presentBoard.place(foodLocation.xCoord, foodLocation.yCoord, false)
	} else if futility == DEATH {
		fmt.Println("________GAME OVER________")
		fmt.Println("snake-length  ", len((*self).position))
		fmt.Println("food Eaten    ", len((*self).position)-4)
		os.Exit(0)
	}

	temp1, temp2 := nextCoordinate, co_ordinate{}

	for i := 0; i < len(self.position); i++ {
		temp2 = self.position[i]
		self.position[i] = temp1
		temp1 = temp2
	}

	return presentBoard
}
