package main

import (
	"fmt"
)

import (
	"metrics/internal/agent/statsReader"
)

func main() {
	var MemStatistics statsReader.MemoryStatsDump
	MemStatistics.Refresh()

	fmt.Println("Hi!")
	fmt.Println(MemStatistics.Alloc)
}
