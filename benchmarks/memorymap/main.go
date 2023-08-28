package main

import (
	"fmt"
	"os"
	// "runtime"
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

	fInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}

	fmt.Printf("page size of this system : %d", os.Getpagesize())

	numberOfPages := 4

	for i := 0; i < int(fInfo.Size())/(numberOfPages*os.Getpagesize()); i++ {
		func() {
			mmap, err := mmap.MapRegion(
				file,
				numberOfPages*os.Getpagesize(),
				mmap.RDONLY,
				0,
				int64(numberOfPages*os.Getpagesize()*i),
			)
			if err != nil {
				panic(err)
			}
			defer mmap.Unmap()

			fmt.Println(string(mmap))
		}()
	}
    var last int
    if int(fInfo.Size())/(numberOfPages*os.Getpagesize()) == int(int(fInfo.Size())/(numberOfPages*os.Getpagesize())) {
        last = int(fInfo.Size())/(numberOfPages*os.Getpagesize()) - 1
    } else {
        last = int(int(fInfo.Size())/(numberOfPages*os.Getpagesize()))
    }
	mmap, err := mmap.MapRegion(
		file,
		int(fInfo.Size())-numberOfPages*os.Getpagesize()*(last+1),
		mmap.RDONLY,
		0,
		int64(numberOfPages*os.Getpagesize()*int(int(fInfo.Size())/(numberOfPages*os.Getpagesize()))),
	)
	if err != nil {
		panic(err)
	}
	defer mmap.Unmap()

	fmt.Println(string(mmap))
}
