package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/gdamore/tcell"
)

type eDirecton int32

const (
	eDirecton_LEFT  eDirecton = 0
	eDirecton_RIGHT eDirecton = 1
	eDirecton_DOWN  eDirecton = 2
	eDirecton_UP    eDirecton = 3
	VIEW_OFFSET   int = 3
	HEIGHT_OFFSET   int = 3
	WIDTH_OFFSET   int = 10
)

type point struct {
	x int
	y int
}

type gameState struct {
	_height       int
	_width        int
	_currentRound int
	_score        int
	_gameOver     bool
	snake         []point
	_food         point
	dir           eDirecton
}

 
func init()  {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func main() {
	lastState := Setup()
	
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	defer screen.Fini() 
	initError := screen.Init()
	if initError != nil {
		panic(initError)
	}

	var directionChan = make(chan eDirecton)              
	var isGameOverChan = make(chan bool)                  
	go getinput(screen, directionChan, isGameOverChan)    
	defer close(directionChan)                            
	defer close(isGameOverChan)                           

	PrintGameView(screen,lastState)
	for !lastState._gameOver {
		select {                                     
		case DIR := <-directionChan:
		lastState.dir = DIR
		default:
		}
		lastState = updateGameState(lastState)
		PrintGameView(screen,lastState)
		time.Sleep(time.Duration(200) * time.Millisecond)
	}
	tc_println(screen,0,1,"Game Over, thanks :)  press ESC key to exit.")
	screen.Show()

	<-isGameOverChan                                 
}

func Setup() gameState {
	var initState gameState
	validInput := false
	initState._gameOver = false

	fmt.Println("Welcom! Please input <width_number, height_number> pair(range[10~100]) to set board size: ")
	
	for !validInput {
		cmd := getBoardSize()
		width, height, err := ParseCmd(cmd)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("Please input <width_number, height_number> pair(range[10~100]) to set board size: ")
			continue
		}
		initState._height = height
		initState._width = width
		validInput = true
	}
	head := point{initState._width / 2, initState._height / 2}
	initState.snake = append(initState.snake, head)
	foodX := rand.Intn(initState._width -1) + WIDTH_OFFSET + 1
	foodY := rand.Intn(initState._height -1)  + HEIGHT_OFFSET+VIEW_OFFSET + 1
	initState._food = point{foodX, foodY}
	return initState
}
 
func GetHeadSymbol(directon eDirecton) rune {
	var symbol rune
	switch directon{
	case eDirecton_DOWN:
		symbol = '^'
	case eDirecton_LEFT:
		symbol = '>'
	case eDirecton_RIGHT:
		symbol = '<'
	case eDirecton_UP:
		symbol = 'v'
	}

	return symbol
}

func PrintGameView(screen tcell.Screen, currentstate gameState) {
	screen.Clear()
	height := currentstate._height
	width := currentstate._width
	theRightWall := WIDTH_OFFSET + width 
	headX := currentstate.snake[0].x
	headY := currentstate.snake[0].y
	foodX := currentstate._food.x
	foodY := currentstate._food.y
	snakeLength := len(currentstate.snake)

	tc_println(screen,0,VIEW_OFFSET,"==============Game View:================================")
	for j := WIDTH_OFFSET; j < width+WIDTH_OFFSET+1; j++ {
		tc_printcell(screen,j, HEIGHT_OFFSET+VIEW_OFFSET, '#')
	}

	for i := HEIGHT_OFFSET+VIEW_OFFSET; i < height+HEIGHT_OFFSET+VIEW_OFFSET; i++ {
		for j := WIDTH_OFFSET; j < width + WIDTH_OFFSET + 1; j++ {
			if j == WIDTH_OFFSET {
				tc_printcell(screen,j, i, '#')
			}
			if i == headY && j == headX {
				tc_printcell(screen,j, i, GetHeadSymbol(currentstate.dir))
			} else if i == foodY && j == foodX {
				tc_printcell(screen,j, i, '$')
			} else {
				for k := 1; k < snakeLength; k++ {
					if currentstate.snake[k].x == j && currentstate.snake[k].y == i {
						tc_printcell(screen,j, i, 'o')
					}
				}
				
			}
			if j == theRightWall {
				tc_printcell(screen,j, i, '#')
			}
		}
	}
	for j := WIDTH_OFFSET; j < width+WIDTH_OFFSET+1; j++ {
		tc_printcell(screen,j, height + HEIGHT_OFFSET+VIEW_OFFSET, '#')
	}

	ViewBottomOffset := VIEW_OFFSET + HEIGHT_OFFSET + height + HEIGHT_OFFSET
	tc_println(screen,0,ViewBottomOffset,"==============Game View:================================")
	tc_println(screen,0,ViewBottomOffset + 3, fmt.Sprintf("Score: %v", currentstate._score))
	tc_println(screen,0,ViewBottomOffset + 5, fmt.Sprintf("Snake length:  %d", len(currentstate.snake)))
	tc_println(screen,0,ViewBottomOffset + 7, fmt.Sprintf("Current round:  %d", currentstate._currentRound))
	tc_println(screen,0,ViewBottomOffset + 9, fmt.Sprintf("Current board size:  %d, %d", currentstate._width, currentstate._height))
	tc_println(screen,0,ViewBottomOffset + 11, "use key A for left, key D for right, key W for up, key s for down")
	screen.Show()
}

func getBoardSize() string {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	result := input.Text()
	return result
}

func ParseCmd(cmd string) (int, int, error) {
	var err error
	var width int
	var height int

	if len(cmd) < 1 {
		err = errors.New("bad format--not an ints pair, no comma")
		return width, height, err
	}
	if !strings.Contains(cmd, ",") {
		err = errors.New("bad format--not an ints pair, no comma")
		return width, height, err
	}
	temp := strings.Split(cmd, ",")
	if len(temp) != 2 {
		err = errors.New("bad format--please input <width_number, height_number> pair")
		return width, height, err
	}
	width, er := strconv.Atoi(temp[0])
	if er != nil {
		err = errors.New("bad format--width_number")
		return width, height, err
	}
	height, er = strconv.Atoi(temp[1])
	if er != nil {
		err = errors.New("bad format--height_number")
		return width, height, err
	}
	if width < 10 || width > 100 {
		err = errors.New("bad input--width number must be within range[10~100]")
		return width, height, err
	}
	if height < 10 || height > 100 {
		err = errors.New("bad input--height number must be within range[10~100]")
		return width, height, err
	}

	return width, height, err
}

func updateGameState(newState gameState) gameState {
	snakeHead := newState.snake[0]
	headX := snakeHead.x
	headY := snakeHead.y
	width := newState._width
	height := newState._height
	TopEnd := HEIGHT_OFFSET+VIEW_OFFSET
	BottomEnd := HEIGHT_OFFSET+VIEW_OFFSET + height
	LeftEnd := WIDTH_OFFSET
	RightEnd := WIDTH_OFFSET + width

	switch newState.dir {
	case eDirecton_LEFT:
		headX--
	case eDirecton_RIGHT:
		headX++
	case eDirecton_UP:
		headY--
	case eDirecton_DOWN:
		headY++
	default:
	}

	if headX >= RightEnd || headX <= LeftEnd || headY >=  BottomEnd || headY <= TopEnd {
		newState._gameOver = true
		return newState
	}
	tailee := newState.snake
	if len(tailee) > 1 {
		if tailee[1].x == headX && tailee[1].y == headY {
			return newState
		}
	}
	for i := 1; i < len(tailee); i++ {
		if headX == tailee[i].x && headY == tailee[i].y {
			newState._gameOver = true
			return newState
		}
	}
	if headX == newState._food.x && headY == newState._food.y {
		newState._score++
		generateFood:
		foodX := rand.Intn(width-1) + WIDTH_OFFSET + 1
		foodY := rand.Intn(height-1) + HEIGHT_OFFSET+VIEW_OFFSET + 1
		for _,loc := range newState.snake{
			if loc.x == foodX && loc.y == foodY{
				goto generateFood
			}
		}

		newState._food.x = foodX
		newState._food.y = foodY
		newPoint := point{headX, headY}
		var newList []point
		newList = append(newList, newPoint)
		newList = append(newList, tailee...)
		newState.snake = newList
	} else {
		newSnake := []point{}
		newSnake = append(newSnake, point{headX, headY})
		newState.snake = append(newSnake, newState.snake...)
		if len(newState.snake) > 0 {
			newState.snake = newState.snake[:len(newState.snake)-1]
		}
	}
	newState._currentRound++

	return newState
}

func tc_println(screen tcell.Screen, x, y int, msg string) {
	str := []rune(msg)
	screen.SetContent(x, y, ' ', str, tcell.StyleDefault)
}

func tc_printcell(screen tcell.Screen, x,y int, ch rune){
	screen.SetCell(x, y, tcell.StyleDefault, ch,)
}

func getinput(screen tcell.Screen, directionChan chan eDirecton, isGameOverChan chan bool){    
	loop:                             
	for  {
		switch ev := screen.PollEvent(); et := ev.(type) {
		case *tcell.EventKey:
			if(et.Key() == tcell.KeyEsc){
				isGameOverChan <- true
				break loop
			}
			switch et.Rune() {
			case 97,65:                   // key A
				directionChan <- eDirecton_LEFT
			case 115,83:
				directionChan <- eDirecton_DOWN
			case 100,68:                   // key D
				directionChan <- eDirecton_RIGHT
			case 119,87:
				directionChan <- eDirecton_UP
			}
			 
		}
	}
}