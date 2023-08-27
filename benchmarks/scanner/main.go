package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
    fmt.Println("executing the scanner benchmark")
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

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        fmt.Println(line)
    }
}

