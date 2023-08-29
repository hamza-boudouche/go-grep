# Benchmarks

You will find here the file reading strategies that were ran and compared using
a 20 MB text file (approximately 32K lines), including:

- Using a scanner (`bufio.NewScanner`) to read the files line by line. the
advantages of this approache is that it's fairly simple to implement and
it's guaranteed to not consume too much physical memory (because go's garbage
collector will eventually clean up the lines that were previously read in
previous iterations of a for loop). Its main disadvantage is that it's
relatively slow compared to the next approache.
- Using memory mapping to map chunks of the file into virtual memory (which is
mapped to physical memory by the kernel) which avoids the time penalties
incurred by reading each line directly from disk (as reading from physical
memory is orders of magnitude faster).

Tests executed using the 20 MB file mentionned above resulted in it being read
using the scanner aproache in 6.5s on average, and in 3.3 seconds on average
using memory mapping, almost a 50% decrease in reading times!

The memory mapping approache uses a sliding window mechanism. The size of the
sliding window is a multiple of the system's page size, which makes the
implementation easier as the memory mapping syscall require the offset used to
map a chunk of a file into virtual memory to be a multiple of the system's page
size.

The sliding window's size is determined by the memory usage constraints of the
user (the garbage collection cycles make it so that go-grep keeps a reasonable
memory usage). The tests mentionned above were executed using chunks of 250
pages (the page size on the machine that was used is 4 KB).

