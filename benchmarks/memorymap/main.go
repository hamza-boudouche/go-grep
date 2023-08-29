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
	readFile("../sharedData/input2.html")
	tik := time.Now()
	fmt.Printf("benchmark took : %s\n", tik.Sub(tok))
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func readFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

    // use file stats to get file size
	fInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
    fSize := int(fInfo.Size())

    // get system page Size
    pageSize := os.Getpagesize()

    // number of concurrent pages to read
	numberOfPages := 4

    // iterate the number of times the reading chunk fits in the file
    // the reading chunk being `numberOfPages*pageSize`
	for i := 0; i < fSize/(numberOfPages*pageSize); i++ {
        // we use an anonymous function to execute a deferred statement on each iteration of the enclosing loop
		func() {
			mmap, err := mmap.MapRegion(
				file,
				numberOfPages*pageSize,
				mmap.RDONLY,
				0,
				int64(numberOfPages*pageSize*i),
			)
			if err != nil {
				panic(err)
			}
			defer mmap.Unmap()

			fmt.Println(string(mmap))
		}()
	}

    // this is the last value of i in the previous for loop
    // which means it's the index of the last whole chunk of data that was read from the file
    var last int
    if fSize/(numberOfPages*pageSize) == int(fSize/(numberOfPages*pageSize)) {
        last = fSize/(numberOfPages*pageSize) - 1
    } else {
        last = int(fSize/(numberOfPages*pageSize))
    }

    // here we use the variable `last` to get the remaining data which size is smaller than the reading chunk
	mmap, err := mmap.MapRegion(
		file,
		int(fSize)-numberOfPages*pageSize*(last+1),
		mmap.RDONLY,
		0,
		int64(numberOfPages*pageSize*int(fSize/(numberOfPages*pageSize))),
	)
	if err != nil {
		panic(err)
	}
	defer mmap.Unmap()

	fmt.Println(string(mmap))
}
