# go-grep

This is a clone of the popular grep command line utility that is implemented
in the go programming language. It makes use of go's simple and efficient
concurrency model to grep a file or set of files looking for a string that
matches a regular expression given as parameter.

While this was originally meant to be a simple exercice whose goal is the
simply apply newly learned go concepts, it later turned into an optimisation
quest, with the first challenge being to speed up the file reading times.

In the benchmarking directory, you'll find the file reading methods that were
implemented and tested against medium sized test files.


## TODO

- [ ] add overlapping secondary window to the memory mapping approache so that
strings that match the grepped regex and that span 2 pages are also found.
- [ ] use the memory mapping approache instead of the scanner in the CLI code
- [ ] use goroutines to match the regex on each file faster (currently only 1
goroutine is used for each file, the goal is to start mutliple goroutines per
file by relying on the memory mapped chunks already used in the algorithm).
