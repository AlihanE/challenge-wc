package main

import (
	"bufio"
	"fmt"
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
		f := openFile(os.Args[2])
		defer f.Close()

		switch os.Args[1] {
		case "-c":
			fmt.Println(c(f), os.Args[2])
		case "-l":
			resCh, inCh := l()
			go readFile(f, inCh)
			fmt.Println(<-resCh, os.Args[2])
		case "-w":
			resCh, inCh := w()
			go readFile(f, inCh)
			fmt.Println(<-resCh, os.Args[2])
		case "-m":
			fmt.Println(m(os.Args[2]), os.Args[2])
		}
	} else {
		f := openFile(os.Args[1])
		defer f.Close()
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

func c(f *os.File) int64 {
	s, err := f.Stat()
	if err != nil {
		errorExit(errFile, err)
	}

	return s.Size()
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

func openFile(fileName string) (file *os.File) {
	file, err := os.Open(fileName)
	if err != nil {
		errorExit(errFile, err)
	}

	return
}

func readFile(f *os.File, chans ...chan<- string) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		for _, ch := range chans {
			ch <- scanner.Text()
		}
	}

	for _, ch := range chans {
		close(ch)
	}
}
