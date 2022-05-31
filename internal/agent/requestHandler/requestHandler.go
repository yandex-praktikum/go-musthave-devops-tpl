package requestHandler

import (
	"fmt"
	"metrics/internal/agent/statsReader"
)

func MemoryStatsUpload(memoryStats statsReader.MemoryStatsDump) error {
	fmt.Println(memoryStats)
	return nil
}
