package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)
	go func() {
		defer f.Close()
		defer close(out)

		var buf strings.Builder
		data := make([]byte, 8)

		for {
			n, err := f.Read(data)
			if n > 0 {
				for _, b := range data[:n] {
					if b == '\n' {
						out <- buf.String()
						buf.Reset()
					} else {
						buf.WriteByte(b)
					}
				}
			}
			if err != nil {
				break
			}
		}

		if buf.Len() > 0 {
			out <- buf.String()
		}
	}()
	return out
}

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	lines := getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}
