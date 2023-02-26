package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/hashicorp/go-set"
)

type FoundLine struct {
	lineNo int
	line string
	filePath string
}

func (foundLine *FoundLine) Print() {
	fmt.Printf("%+v\n", foundLine)
}

var (
	wgWorkers = sync.WaitGroup{}
	wgListeners = sync.WaitGroup{}

)

func main() {
	// setting and parsing flags
	force := flag.Bool("f", false, "ignore any error messages and warnings")
	query := flag.String("q", "", "your query, aka what you're searching for, can be one or multiple regexps")
	location := flag.String("l", ".", "the files you want to search through, can be one or multiple regexps")
	name := flag.String("n", ".*", "the files you want to search through, can be one or multiple regexps")
	flag.Parse()

	fmt.Println("-----------------------------debugging")
	fmt.Println(*query)
	fmt.Println(*location)
	fmt.Println("-----------------------------debugging")

	// check params
	if !*force && len(*query) == 0 {
		fmt.Println("invalid parameters - parameter `-query` is required")
		flag.PrintDefaults()
		fmt.Println("if you'd like to ignore this error message and proceed with the parameters you gave, use the flag `-force`")
	}

	// preprocessing flags
	locations := strings.Fields(*location)
	names_re := []*regexp.Regexp{}
	for _, n := range strings.Fields(*name) {
		re, err := regexp.Compile(n)
		if err != nil {
			panic(n + " can't be compiled into a regexp")
		}
		names_re = append(names_re, re)
	}
	queries_re := []*regexp.Regexp{}
	for _, q := range strings.Fields(*query) {
		re, err := regexp.Compile(q)
		if err != nil {
			panic(q + " can't be compiled into a regexp")
		}
		queries_re = append(queries_re, re)
	}

	// FIXME this might be useless because the file system insures that there will be no duplicates
	file_set := set.New[string](10)

	ch := make(chan FoundLine, 100)

	wgListeners.Add(1)
	go func(ch <-chan FoundLine) {
		defer wgListeners.Done()
		for {
			if foundLine, ok := <- ch; ok {
				foundLine.Print()
			} else {
				break
			}
		}
	}(ch)

	for _, filepath := range locations {
		handle_files(filepath, file_set, &names_re)
		for _, file := range file_set.Slice() {
			fmt.Println(file)
			wgWorkers.Add(1)
			go read_file(file, &queries_re, ch)
		}
	}
	wgWorkers.Wait()
	close(ch)
	wgListeners.Wait()
}

// recursively search for files with matching names
func handle_files(basePath string, s *set.Set[string], names_re *[]*regexp.Regexp) error {
	files, err := os.ReadDir(basePath)
	if err != nil {
		panic("invalid filepath " + basePath)
	}
	for _, file := range files {
		if file.IsDir() {
			err := handle_files(filepath.Join(basePath, file.Name()), s, names_re)
			if err != nil {
				return err
			}
		} else if fileInfo, err := os.Lstat(filepath.Join(basePath, file.Name())); err == nil && (fileInfo.Mode()&os.ModeSymlink == 0) {
			// this above condition checks if the file is not a symlink, in fact this program ignores symlinks
			filename := filepath.Join(basePath, file.Name())
			match := false
			for _, re := range *names_re {
				if re.Match([]byte(filename)) {
					match = true
					break
				}
			}
			if match {
				s.Insert(filename)
			}
		}
	}
	return nil
}

func read_file(filePath string, queries_re *[]*regexp.Regexp, ch chan<- FoundLine) {
	defer wgWorkers.Done()

	check := func(e error) {
		if e != nil {
			panic(e)
		}
	}

	file, err := os.Open(filePath)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNo := 0
	// while loop: read while scanner not finished
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		match := false
		for _, re := range *queries_re {
			if re.Match([]byte(scanner.Text())) {
				match = true
				break
			}
		}
		if match {
			fmt.Println(lineNo, line)
			ch <- FoundLine{
				lineNo: lineNo,
				line: line,
				filePath: filePath,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		check(err)
	}
}

