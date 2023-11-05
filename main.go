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

	if len(os.Args) < 3 {
		errorExit(errC)
	}

	f := openFile(os.Args[2])
	defer f.Close()

	switch os.Args[1] {
	case "-c":
		c(f)
	case "-l":
		l(f)
	case "-w":
		w(f)
	case "-m":
		m(f)
	}
}

func errorExit(args ...any) {
	fmt.Println(args...)
	os.Exit(1)
}

func c(f *os.File) {
	s, err := f.Stat()
	if err != nil {
		errorExit(errFile, err)
	}

	fmt.Println(s.Size(), f.Name())
}

func l(f *os.File) {
	scanner := bufio.NewScanner(f)
	i := 0
	for scanner.Scan() {
		i++
	}

	fmt.Println(i, f.Name())
}

func w(f *os.File) {
	scanner := bufio.NewScanner(f)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		i += len(strings.Fields(line))
	}

	fmt.Println(i, f.Name())
}

func m(f *os.File) {
	scanner := bufio.NewScanner(f)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		i += strings.Count(line, "")
	}

	fmt.Println(i, f.Name())
}

func openFile(fileName string) (file *os.File) {
	file, err := os.Open(fileName)
	if err != nil {
		errorExit(errFile, err)
	}

	return
}
