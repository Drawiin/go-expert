package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	i := 0
	for {
		f, err := os.Create(fmt.Sprintf("./tmp/file%d.txt", i))
		if err !=nil {
			panic(err)
		}
		defer f.Close()
		f.WriteString("Hello, World!")
		time.Sleep(1 * time.Second)
		i++
	}
}