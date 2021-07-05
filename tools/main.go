package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/minus5/svckit/nsq"
)

func main() {
	pub := nsq.Pub("")
	defer pub.Close()
	for bytes := range readFileLines("./ponuda_short.json") {
		line := string(bytes)
		fmt.Println(line)
		pub.PublishTo("ponuda.req", bytes)
	}
}

func readFileLines(fileName string) chan []byte {
	lines := make(chan []byte, 64)
	go func() {
		defer close(lines)
		f, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Fatal("no such file")
		}
		defer f.Close()
		rd := bufio.NewReader(f)
		for {
			line, err := rd.ReadBytes('\n')
			if err != nil {
				return
			}
			lines <- line[:len(line)-1] // trim newline, last char
		}
	}()
	return lines
}
