package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const BUFSIZE = 5000

var mem [BUFSIZE]int

func print_usage() {
	fmt.Println("usage: go run main.go <file>")
}

type State struct {
	index       int
	loops       []string
	loopIndexes []int
}

var state = State{
	index: 0,
}

func (s *State) Next() {
	s.index = (s.index + 1) % BUFSIZE
}

func (s *State) Prev() {
	s.index = (s.index - 1) % BUFSIZE
}

func (s *State) Increment() {
	mem[s.index] = mem[s.index] + 1
}

func (s *State) Decrement() {
	mem[s.index]--
}

func (s *State) Print() {
	fmt.Printf("%c", mem[s.index])
}

func (s *State) Save(char rune) {
	mem[s.index] = int(char)
}

func (s *State) registerLoop(char rune) {
	if len(s.loops) <= 0 {
		return
	}
	s.loops[len(s.loops)-1] = s.loops[len(s.loops)-1] + string(char)
}

func (s *State) CreateLoop() {
	s.loops = append(s.loops, "")
	s.loopIndexes = append(s.loopIndexes, s.index)
}

func (s *State) Repeat() {
	str := s.loops[len(s.loops)-1]
	index := s.loopIndexes[len(s.loopIndexes)-1]
	for mem[index] > 0 {
		interpret(bufio.NewReader(strings.NewReader(str)))
	}
	s.loops = s.loops[:len(s.loops)-1]
}

func main() {
	if len(os.Args) != 2 {
		print_usage()
		return
	}
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := bufio.NewReader(file)
	interpret(r)
}

func interpret(r *bufio.Reader) {
	for {
		if char, _, err := r.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		} else {
			switch char {
			case '>':
				state.Next()
			case '<':
				state.Prev()
			case '+':
				state.Increment()
			case '-':
				state.Decrement()
			case '.':
				state.Print()
			case ',':
				input := ' '
				fmt.Scanf("%c", &input)
				state.Save(input)
			case '[':
				state.CreateLoop()
			case ']':
				state.Repeat()
			case '\n', '\r', ' ':
				continue
			default:
				panic("Unknown character")
			}
			state.registerLoop(char)
		}
	}
}
