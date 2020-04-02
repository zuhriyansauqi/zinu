package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/zuhriyan/zinu/zinu"
)

func main() {
	start := time.Now()

	z, err := zinu.Load("templates/idn/idn.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := z.Generate("output/output." + start.Format("2006.01.02")); err != nil {
		log.Fatal(err)
	}

	fmt.Println(time.Since(start))
	PrintMemUsage()
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
