package tests

import (
	"github.com/stretchr/testify/assert"
	"metrics/internal/agent/statsReader"
	"testing"
)

func TestRefresh(t *testing.T) {
	var memStatistics statsReader.MemoryStatsDump
	memStatistics.Refresh()
	memStatistics.Refresh()
	memStatistics.Refresh()

	assert.Equal(t, 3, int(memStatistics.PollCount))
}
