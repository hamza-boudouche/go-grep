package main

import (
	"fmt"
	"os"
	"time"
	"github.com/edsrzf/mmap-go"
)

func main() {
	fmt.Println("executing the memorymap benchmark")
	tok := time.Now()
	readFile("../sharedData/input.html")
	tik := time.Now()
	fmt.Printf("benchmark took : %s\n", tik.Sub(tok))
}

func readFile(filePath string) {
    file, err := os.Open(filePath)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    mmap, err := mmap.Map(file, mmap.RDONLY, 0)
    if err != nil {
        panic(err)
    }
    defer mmap.Unmap()

    fmt.Println(string(mmap))
    // fmt.Println(len(mmap))
    // fmt.Println(string(mmap[:20]))
}

