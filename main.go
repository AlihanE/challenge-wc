package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

const errC = `Enter flag -c <file> to get bytes count`
const errFile = `Can't read file. Error:`

func main() {
	if len(os.Args) < 2 {
		errorExit(errC)
	}

	if strings.Contains(os.Args[1], "-") {
		var f []byte
		var err error
		var fileName string
		if len(os.Args) < 3 {
			f, err = io.ReadAll(os.Stdin)
			if err != nil {
				panic(err)
			}
			fileName = "stdin"
		} else {
			f = openFile(os.Args[2])
			fileName = os.Args[2]
		}

		switch os.Args[1] {
		case "-c":
			fmt.Println(c(f), fileName)
		case "-l":
			resCh, inCh := l()
			go readFile(f, inCh)
			fmt.Println(<-resCh, fileName)
		case "-w":
			resCh, inCh := w()
			go readFile(f, inCh)
			fmt.Println(<-resCh, fileName)
		case "-m":
			fmt.Println(m(fileName), fileName)
		}
	} else {
		f := openFile(os.Args[1])
		lResCh, lInCh := l()
		wResCh, wInCh := w()
		go readFile(f, lInCh, wInCh)
		lines := <-lResCh
		words := <-wResCh
		fmt.Println(lines, words, c(f), os.Args[1])
	}
}

func errorExit(args ...any) {
	fmt.Println(args...)
	os.Exit(1)
}

func c(b []byte) int {
	return len(b)
}

func l() (chan int, chan string) {
	i := 0
	resCh := make(chan int)
	inCh := make(chan string)

	go func() {
		for range inCh {
			i++
		}
		resCh <- i
	}()

	return resCh, inCh
}

func w() (chan int, chan string) {
	i := 0
	resCh := make(chan int)
	inCh := make(chan string)

	go func() {
		for line := range inCh {
			i += len(strings.Fields(line))
		}
		resCh <- i
	}()

	return resCh, inCh
}

func m(fileName string) int {
	b, err := os.ReadFile(fileName)
	if err != nil {
		errorExit(errFile, err)
	}

	return len([]rune(string(b)))
}

func openFile(fileName string) []byte {
	b, err := os.ReadFile(fileName)
	if err != nil {
		errorExit(errFile, err)
	}

	return b
}

func readFile(b []byte, chans ...chan<- string) {
	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		for _, ch := range chans {
			ch <- scanner.Text()
		}
	}

	for _, ch := range chans {
		close(ch)
	}
}
