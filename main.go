package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var Grid [][]rune
var osArgs []string

type Command func()

type coordinates struct {
	x int
	y int
}

var ProgramRunning = false
var isString = false
var cursor_position = coordinates{x: -1, y: 0}
var direction = coordinates{x: 1, y: 0}
var runtime_values [][]int

var value_to_append int

var is_the_value_removed_from_stack = false
var removed_value_from_stack = 0

var dimension_working_array = 0

var help = []string{"", "rupi run [Filename]     -run file", "rupi --help             -get all the commands available"}

func readFile(file_path string) error {
	content, err := ioutil.ReadFile(file_path)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(content), "\n")

	Grid = make([][]rune, len(lines))
	var lineRune []rune
	for i, line := range lines {
		lineRune = []rune(line)
		if lineRune[len(lineRune)-1] == 13 {
			Grid[i] = lineRune[:len(lineRune)-1]
		} else {
			Grid[i] = lineRune
		}
		// if result1 := strings.IndexRune(string(lineRune), 2325); result1 != -1 {
		// 	cursor_position = coordinates{x: result1, y: i}
		// 	ProgramRunning = true
		// }
		ProgramRunning = true
	}
	return nil
}

func move() rune {
	if direction.x == 1 {
		if cursor_position.x < (len(Grid[cursor_position.y]) - 1) {
			cursor_position.x++
		} else {
			cursor_position.x = 0
		}
	} else if direction.x == -1 {
		if cursor_position.x == 0 {
			cursor_position.x = len(Grid[cursor_position.y]) - 1
		} else {
			cursor_position.x--
		}
	} else if direction.y == 1 {
		for {
			if cursor_position.y < len(Grid)-1 {
				cursor_position.y++
			} else {
				cursor_position.y = 0
			}
			if cursor_position.x < len(Grid[cursor_position.y]) {
				break
			}
		}
	} else if direction.y == -1 {
		for {
			if cursor_position.y == 0 {
				cursor_position.y = len(Grid) - 1
			} else {
				cursor_position.y--
			}
			if cursor_position.x < len(Grid[cursor_position.y]) {
				break
			}
		}
	}
	return Grid[cursor_position.y][cursor_position.x]
}

func apendString(programString rune) {
	runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], int(programString))
}

func add_dimension() {
	if len(runtime_values)-1 <= dimension_working_array {
		runtime_values = append(runtime_values, nil)
	}
}

func remove_dimension() {
	if len(runtime_values)-1 >= dimension_working_array {
		runtime_values[dimension_working_array] = nil
	}
}

var commands = map[string]func(){
	// "क": func() { // This is where the program starts
	// 	direction.x = 1
	// 	direction.y = 0
	// },
	">": func() { // Move The cursor Right
		direction.x = 1
		direction.y = 0
	},
	"<": func() { // Move The cursor Left
		direction.x = -1
		direction.y = 0
	},
	"^": func() { // Move The cursor Up
		direction.x = 0
		direction.y = -1
	},
	"v": func() { // Move The cursor Down
		direction.x = 0
		direction.y = 1
	},
	"V": func() { // Move The cursor Down
		direction.x = 0
		direction.y = 1
	},
	"/": func() { // Move The cursor Down
		if direction.x == 1 {
			direction.y = -1
			direction.x = 0
		} else if direction.x == -1 {
			direction.y = 1
			direction.x = 0
		} else if direction.y == 1 {
			direction.x = -1
			direction.y = 0
		} else if direction.y == -1 {
			direction.x = 1
			direction.y = 0
		}
	},
	"\\": func() { // Move The cursor Down
		if direction.x == 1 {
			direction.y = 1
			direction.x = 0
		} else if direction.x == -1 {
			direction.y = -1
			direction.x = 0
		} else if direction.y == 1 {
			direction.x = 1
			direction.y = 0
		} else if direction.y == -1 {
			direction.x = -1
			direction.y = 0
		}
	},
	"|": func() { // Move The cursor Down
		switch direction.x {
		case 1:
			direction.y = 0
			direction.x = -1
		case -1:
			direction.y = 0
			direction.x = 1
		}
	},
	"_": func() { // Move The cursor Down
		switch direction.y {
		case 1:
			direction.y = -1
			direction.x = 0
		case -1:
			direction.y = 1
			direction.x = 0
		}
	},
	"#": func() { // Move The cursor Down
		switch direction.x {
		case 1:
			direction.y = 0
			direction.x = -1
		case -1:
			direction.y = 0
			direction.x = 1
		}
		switch direction.y {
		case 1:
			direction.x = 0
			direction.y = -1
		case -1:
			direction.x = 0
			direction.y = 1
		}
	},
	"!": func() { // Skip the next command
		move()
	},
	".": func() { // Change the cursor position
		cursor_position.y = runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1]
		cursor_position.x = runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-2]
		runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-2]
	},
	"?": func() { // Skip the next command if value in the current stack is 0
		if runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1] == 0 {
			move()
		}
		runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-1]
	},
	"0": func() { // Append 9 To Stack
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], 0)
	},
	"1": func() { // Append 1 To Stack
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], 1)
	},
	"2": func() { // Append 2 To Stack
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], 2)
	},
	"3": func() { // Append 3 To Stack
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], 3)
	},
	"4": func() { // Append 4 To Stack
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], 4)
	},
	"5": func() { // Append 5 To Stack
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], 5)
	},
	"6": func() { // Append 6 To Stack
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], 6)
	},
	"7": func() { // Append 7 To Stack
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], 7)
	},
	"8": func() { // Append 8 To Stack
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], 8)
	},
	"9": func() { // Append 9 To Stack
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], 9)
	},
	"a": func() { // Append 10 To Stack
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], 10)
	},
	"b": func() { // Append 11 To Stack
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], 11)
	},
	"c": func() { // Append 12 To Stack
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], 12)
	},
	"d": func() { // Append 13 To Stack
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], 13)
	},
	"e": func() { // Append 14 To Stack
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], 14)
	},
	"f": func() { // Append 15 To Stack
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], 15)
	},
	"'": func() { // Make next Values Strings
		isString = !isString
	},
	"\"": func() { // Make next Values Strings
		isString = !isString
	},
	"s": func() { // Substract value
		if len(runtime_values[dimension_working_array]) > 0 {
			runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1] = runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1] - 1
		}
	},
	"P": func() { // Print each value in next line
		for i := len(runtime_values[dimension_working_array]) - 1; i >= 0; i-- {
			fmt.Printf("%v \n", string(runtime_values[dimension_working_array][i]))
		}
	},
	"p": func() { // Print values in same line
		for i := len(runtime_values[dimension_working_array]) - 1; i >= 0; i-- {
			fmt.Printf("%v", string(runtime_values[dimension_working_array][i]))
		}
		fmt.Printf("\n")
	},
	"O": func() { // Print value of top in stack and remove it from stack
		if len(runtime_values[dimension_working_array]) > 0 {
			fmt.Println(string(runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1]))
			runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-1]
		}
	},
	"o": func() { // Print value of top in stack and remove it from stack
		if len(runtime_values[dimension_working_array]) > 0 {
			fmt.Print(string(runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1]))
			runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-1]
		}
	},
	"+": func() { // Add addition value of second last and last value to stack
		if len(runtime_values[dimension_working_array]) > 1 {
			value_to_append = runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-2] + runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1]
			runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-2]
			runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], value_to_append)
		}
	},
	"-": func() { // Add Substracted value of second last and last value to stack
		if len(runtime_values[dimension_working_array]) > 1 {
			value_to_append = runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-2] - runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1]
			runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-2]
			runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], value_to_append)
		}
	},
	"*": func() { // Add Multiplication value of second last and last value to stack
		if len(runtime_values[dimension_working_array]) > 1 {
			value_to_append = runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-2] * runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1]
			runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-2]
			runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], value_to_append)
		}
	},
	",": func() { // Add Divided value of second last and last value to stack
		if len(runtime_values[dimension_working_array]) > 1 {
			value_to_append = runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-2] / runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1]
			runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-2]
			runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], value_to_append)
		}
	},
	"%": func() { // Add Modules of second last and last value to stack
		if len(runtime_values[dimension_working_array]) > 1 {
			value_to_append = runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-2] % runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1]
			runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-2]
			runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], value_to_append)
		}
	},
	"@": func() { // Add Modules of second last and last value to stack
		if len(runtime_values[dimension_working_array]) >= 3 {
			values_array := []int{runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1], runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-3], runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-2]}
			runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-3]
			runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], values_array...)
		}
	},
	"{": func() {
		value_to_append = runtime_values[dimension_working_array][0]
		runtime_values[dimension_working_array] = runtime_values[dimension_working_array][1:]
		runtime_values[dimension_working_array] = append([]int{value_to_append}, runtime_values[dimension_working_array]...)

	},
	"}": func() {
		value_to_append = runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1]
		runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-1]
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], value_to_append)

	},
	"=": func() { // Print each value in next line
		if len(runtime_values[dimension_working_array]) > 1 {
			value_to_append = 0
			if runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-2] == runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1] {
				value_to_append = 1
			}
			runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-2]
			runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], value_to_append)
		}
	},
	")": func() { // Add Substracted value of second last and last value to stack
		if len(runtime_values[dimension_working_array]) > 1 {
			value_to_append = 0
			if runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-2] > runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1] {
				value_to_append = 1
			}
			runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-2]
			runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], value_to_append)
		}
	},
	"(": func() { // Add Substracted value of second last and last value to stack
		if len(runtime_values[dimension_working_array]) > 1 {
			value_to_append = 0
			if runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-2] < runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1] {
				value_to_append = 1
			}
			runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-2]
			runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], value_to_append)
		}
	},
	":": func() { // Duplicate top value from stack
		if len(runtime_values[dimension_working_array]) > 0 {
			runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1])
		}
	},
	"~": func() { // Remove top value from stack
		if len(runtime_values[dimension_working_array]) > 0 {
			runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-1]
		}
	},
	"$": func() { // Remove top value from stack
		if len(runtime_values[dimension_working_array]) > 1 {
			value_to_append = runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1]
			runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1] = runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-2]
			runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-2] = value_to_append
		}
	},
	"&": func() {
		if !is_the_value_removed_from_stack {
			removed_value_from_stack = runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1]
			runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-1]
		} else {
			runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], removed_value_from_stack)
		}
		is_the_value_removed_from_stack = !is_the_value_removed_from_stack
	},
	"r": func() { // Reverse Value in Stack
		for i, j := 0, len(runtime_values[dimension_working_array])-1; i < j; i, j = i+1, j-1 {
			runtime_values[dimension_working_array][i], runtime_values[dimension_working_array][j] = runtime_values[dimension_working_array][j], runtime_values[dimension_working_array][i]
		}
	},
	"l": func() { // Append the length of stack in the stack
		runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array], len(runtime_values[dimension_working_array]))
	},
	"n": func() { // Print each value in next line
		if len(runtime_values[dimension_working_array]) > 0 {
			fmt.Println(runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1])
			runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-1]
		}
	},
	"[": func() {
		if len(runtime_values[dimension_working_array]) > 0 {
			add_dimension()
			value_to_append = runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1]
			runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-1]
			for i := 0; i < value_to_append; i++ {
				new_value_to_append := runtime_values[dimension_working_array][len(runtime_values[dimension_working_array])-1]
				runtime_values[dimension_working_array+1] = append(runtime_values[dimension_working_array+1], new_value_to_append)
				runtime_values[dimension_working_array] = runtime_values[dimension_working_array][:len(runtime_values[dimension_working_array])-1]
			}

			dimension_working_array += 1
		}
	},
	"]": func() {
		if len(runtime_values[dimension_working_array])-1 >= dimension_working_array {
			fmt.Println()
			runtime_values[dimension_working_array] = append(runtime_values[dimension_working_array-1], runtime_values[dimension_working_array]...)
			remove_dimension()
			dimension_working_array -= 1
		}
	},
	"र": func() { // End the Program
		ProgramRunning = false
	},
}

func runProgram() {
	var val rune
	add_dimension()
	for ProgramRunning {
		val = move()
		if !isString || val == 39 || val == 34 {
			if command, ok := commands[string(val)]; ok {
				// fmt.Println(runtime_values)
				command()
			}
		} else {
			apendString(val)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func main() {
	osArgs = os.Args
	switch lenArgs := len(osArgs); lenArgs {
	case 1:
		fmt.Println("Please use --help")
	case 2:
		switch secondValue := osArgs[1]; secondValue {
		case "--help":
			for _, v := range help {
				fmt.Println(v)
			}
		case "run":
			fmt.Println()
			fmt.Println("Define a filename -- eg: rupi run [filename]")
		}
	case 3:
		switch secondValue := osArgs[1]; secondValue {
		case "run":
			path, err := filepath.Abs(osArgs[2])

			if err != nil {
				panic("Error Reading the file")
			}
			error := readFile(path)
			if error != nil {
				panic("Deciphering the file")
			}

			runProgram()
		}
	}

}
