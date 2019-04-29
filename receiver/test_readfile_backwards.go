package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func minInt64(first int64, second int64) int64 {
	if first < second {
		return first
	}
	return second
}

func maxInt64(first int64, second int64) int64 {
	if first > second {
		return first
	}
	return second
}

func getLastLines(fname string, numLines int) []string {
	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}
	var superArr []byte
	buf := make([]byte, minInt64(fi.Size(), 100000))

	curStart := fi.Size()
	for curStart > 0 {
		if curStart < int64(len(buf)) {
			buf = make([]byte, curStart)
			curStart = int64(0)
		} else {
			curStart = curStart - int64(len(buf))
		}

		n, err := file.ReadAt(buf, curStart)
		if err != nil {
			fmt.Println(err)
		}
		buf = buf[:n]
		bufcpy := make([]byte, len(buf))
		copy(bufcpy, buf)
		superArr = append(bufcpy, superArr...)

		if countChars(superArr, '\n') > numLines+12 {
			fmt.Println("Breaking Early")
			break
		}
	}

	lastLines := strings.Split(fmt.Sprintf("%s", superArr), "\n")
	// remove trailing newline-line
	if lastLines[len(lastLines)-1] == "" {
		lastLines = lastLines[:len(lastLines)-1]
	}

	if len(lastLines) <= numLines {
		return lastLines
	}
	return lastLines[len(lastLines)-(numLines):]

}

func main_blub() {
	var numLines int = 10

	if len(os.Args) > 1 {
		lin, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("%v", err)
			return
		} else {
			numLines = lin
		}
	} else {
		numLines = 10
	}

	lastLines := getLastLines("./test/large", numLines)
	for i := 0; i < len(lastLines); i++ {
		fmt.Printf("Line %d: %s\n", i, lastLines[i])
	}

}

func countChars(bytes []byte, char byte) int {
	count := 0
	for i := range bytes {
		if bytes[i] == char {
			count++
		}
	}
	return count
}
